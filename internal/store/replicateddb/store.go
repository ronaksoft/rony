package replicateddb

import (
	"context"
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/internal/store/replicateddb/raftwal"
	"github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"hash/crc64"

	"path/filepath"
	"sync"
	"time"
)

/*
   Creation Time: 2021 - May - 15
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. --go_out=paths=source_relative:. commands.proto
//go:generate protoc -I=. --gorony_out=paths=source_relative,option=no_edge_dep:. commands.proto

var crcTable = crc64.MakeTable(crc64.ECMA)

const (
	ElectionTick       = 10
	HeartbeatTick      = 1
	ProposeTimeout     = 5 * time.Second
	ProposeConfTimeout = 5 * time.Second
)

// Store is the finite state machine which will be used when Raft is enabled.
type Store struct {
	c     cluster.Cluster
	raft  raft.Node
	wal   *raftwal.DiskStorage
	db    *badger.DB
	state *raftpb.ConfState

	// configs
	conflictRetry         int
	conflictRetryInterval time.Duration
	openTxnMtx            sync.RWMutex
	openTxn               map[int64]*store.LTxn
}

func (fsm *Store) OnJoin(hash uint64) {
	for _, v := range fsm.state.Voters {
		if v == hash {
			return
		}
	}
	ctx, cf := context.WithTimeout(context.TODO(), ProposeConfTimeout)
	_ = fsm.raft.ProposeConfChange(ctx, raftpb.ConfChange{
		Type:   raftpb.ConfChangeAddNode,
		NodeID: hash,
		ID:     hash,
	})
	cf()
}

func (fsm *Store) OnLeave(hash uint64) {
	for _, v := range fsm.state.Voters {
		if v == hash {
			ctx, cf := context.WithTimeout(context.TODO(), ProposeConfTimeout)
			_ = fsm.raft.ProposeConfChange(ctx, raftpb.ConfChange{
				Type:   raftpb.ConfChangeRemoveNode,
				NodeID: hash,
				ID:     hash,
			})
			cf()
			return
		}
	}
}

func New(cfg Config) (*Store, error) {
	st := &Store{
		c:       cfg.Cluster,
		openTxn: map[int64]*badger.Txn{},
	}
	db, err := newDB(cfg)
	if err != nil {
		return nil, err
	}
	st.db = db

	store.Init(store.Config{
		DB:                  db,
		ConflictRetries:     cfg.ConflictRetries,
		ConflictMaxInterval: cfg.ConflictMaxInterval,
		BatchWorkers:        cfg.BatchWorkers,
		BatchSize:           cfg.BatchSize,
	})

	raftCfg := &raft.Config{
		ID:              crc64.Checksum(st.c.ServerID(), crcTable),
		ElectionTick:    ElectionTick,
		HeartbeatTick:   HeartbeatTick,
		Storage:         st.wal,
		MaxInflightMsgs: 100,
		MaxSizePerMsg:   4096,
	}

	_, cs, err := st.wal.InitialState()
	if err != nil {
		panic(err)
	}

	st.state = &cs

SnapshotLoop:
	snap, err := st.wal.Snapshot()
	switch err {
	case nil:
	case raft.ErrSnapshotTemporarilyUnavailable:
		time.Sleep(time.Second)
		goto SnapshotLoop
	default:
		panic(err)
	}

	if raft.IsEmptySnap(snap) {
		st.raft = raft.StartNode(raftCfg, nil)
	} else {
		st.raft = raft.RestartNode(raftCfg)
	}

	// Start Raft loop
	go st.raftLoop()

	ctx, cf := context.WithTimeout(context.TODO(), time.Second*10)
	defer cf()
	_ = st.raft.Campaign(ctx)

	st.c.Subscribe(st)

	return st, nil
}

func newDB(config Config) (*badger.DB, error) {
	opt := badger.DefaultOptions(filepath.Join(config.DirPath, "badger"))
	opt.Logger = nil
	return badger.Open(opt)
}

func (fsm *Store) ViewLocal(fn func(txn *store.LTxn) error) error {
	retry := defaultConflictRetries
Retry:
	err := fsm.db.View(fn)
	if err == badger.ErrConflict {
		if retry--; retry > 0 {
			metrics.IncCounter(metrics.CntStoreConflicts)
			time.Sleep(time.Duration(tools.RandomInt64(int64(defaultMaxInterval))))
			goto Retry
		}
	}
	return err
}

func (fsm *Store) UpdateLocal(fn func(txn *store.LTxn) error) error {
	retry := defaultConflictRetries
Retry:
	err := fsm.db.Update(fn)
	if err == badger.ErrConflict {
		if retry--; retry > 0 {
			metrics.IncCounter(metrics.CntStoreConflicts)
			time.Sleep(time.Duration(tools.RandomInt64(int64(defaultMaxInterval))))
			goto Retry
		}
	}
	return err
}

func (fsm *Store) DB() *store.LocalDB {
	return fsm.db
}

func (fsm *Store) Update(fn func(store.Txn) error) error {
	txn := fsm.newTxn(true)
	err := fsm.startTxn(txn)
	if err != nil {
		return err
	}

	err = fn(txn)
	return fsm.stopTxn(txn, err == nil)
}

func (fsm *Store) View(fn func(store.Txn) error) error {
	txn := fsm.newTxn(false)
	err := fsm.startTxn(txn)
	if err != nil {
		return err
	}

	err = fn(txn)
	return fsm.stopTxn(txn, err == nil)
}

func (fsm *Store) Shutdown() {
	if fsm.db != nil {
		_ = fsm.db.Close()
	}
}

func (fsm *Store) SetCluster(c cluster.Cluster) {
	fsm.c = c
}
