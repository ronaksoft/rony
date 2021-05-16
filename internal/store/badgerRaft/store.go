package badgerRaft

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/internal/store"
	"github.com/ronaksoft/rony/pools"
	staticStore "github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
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

// Store is the finite state machine which will be used when Raft is enabled.
type Store struct {
	c  cluster.Cluster
	db *badger.DB

	// configs
	conflictRetry         int
	conflictRetryInterval time.Duration
	openTxnMtx            sync.RWMutex
	openTxn               map[int64]*badger.Txn
}

func New(cfg Config) (*Store, error) {
	st := &Store{
		openTxn: map[int64]*badger.Txn{},
	}
	db, err := newDB(cfg)
	if err != nil {
		return nil, err
	}
	st.db = db

	err = staticStore.Init(staticStore.Config{
		DirPath:             cfg.DirPath,
		ConflictRetries:     cfg.ConflictRetries,
		ConflictMaxInterval: cfg.ConflictMaxInterval,
		BatchWorkers:        cfg.BatchWorkers,
		BatchSize:           cfg.BatchSize,
	})
	if err != nil {
		return nil, err
	}
	return st, nil
}

func newDB(config Config) (*badger.DB, error) {
	opt := badger.DefaultOptions(filepath.Join(config.DirPath, "badger"))
	opt.Logger = nil
	return badger.Open(opt)
}

func (fsm *Store) newTxn(update bool) *Txn {
	return &Txn{
		ID:     tools.RandomInt64(0),
		store:  fsm,
		update: update,
	}
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

func (fsm *Store) DB() *store.DB {
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

func (fsm *Store) startTxn(txn *Txn) error {
	req := PoolStartTxn.Get()
	defer PoolStartTxn.Put(req)
	req.ID = txn.ID
	req.Update = txn.update
	b := pools.Buffer.FromProto(req)
	f := txn.store.c.RaftApply(*b.Bytes())
	pools.Buffer.Put(b)

	err := f.Error()
	if err != nil {
		return err
	}

	res := f.Response()
	if res == nil {
		return nil
	}

	switch x := res.(type) {
	case error:
		return x
	default:
		return ErrUnknown
	}
}

func (fsm *Store) stopTxn(txn *Txn, commit bool) error {
	req := PoolStopTxn.Get()
	defer PoolStopTxn.Put(req)
	req.ID = txn.ID
	req.Commit = commit
	b := pools.Buffer.FromProto(req)
	f := txn.store.c.RaftApply(*b.Bytes())
	pools.Buffer.Put(b)

	err := f.Error()
	if err != nil {
		return err
	}

	res := f.Response()
	if res == nil {
		return nil
	}

	switch x := res.(type) {
	case error:
		return x
	default:
		return ErrUnknown
	}
}

func (fsm *Store) Shutdown() {
	_ = fsm.db.Close()
}
