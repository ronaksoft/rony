package badgerStore

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/tools"
	"io"
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
//go:generate protoc -I=. --gorony_out=paths=source_relative:. commands.proto

type Config struct {
}

// Store is the finite state machine which will be used when Raft is enabled.
type Store struct {
	c  cluster.Cluster
	db *badger.DB

	// configs
	conflictRetry         int
	conflictRetryInterval time.Duration
}

func New(cfg *Config) *Store {
	return &Store{}
}

// Update executes a function, creating and managing a read-write transaction
// for the user. Error returned by the function is relayed by the Update method.
// It retries in case of badger.ErrConflict returned.
func (fsm *Store) Update(fn func(txn *Txn) error) (err error) {
	retry := fsm.conflictRetry
Retry:
	err = fsm.db.Update(fn)
	if err == badger.ErrConflict {
		if retry--; retry > 0 {
			metrics.IncCounter(metrics.CntStoreConflicts)
			time.Sleep(time.Duration(tools.RandomInt64(int64(fsm.conflictRetryInterval))))
			goto Retry
		}
	}
	return
}

// View executes a function creating and managing a read-only transaction for the user. Error
// returned by the function is relayed by the View method. It retries in case of badger.ErrConflict returned.
func (fsm *Store) View(fn func(txn *Txn) error) (err error) {
	retry := fsm.conflictRetry
Retry:
	err = fsm.db.View(fn)
	if err == badger.ErrConflict {
		if retry--; retry > 0 {
			metrics.IncCounter(metrics.CntStoreConflicts)
			time.Sleep(time.Duration(tools.RandomInt64(int64(fsm.conflictRetryInterval))))
			goto Retry
		}
	}
	return
}

func (fsm *Store) Apply(raftLog *raft.Log) interface{} {
	storeCmd := PoolStoreCommand.Get()
	err := storeCmd.Unmarshal(raftLog.Data)
	if err != nil {
		return err
	}

	switch storeCmd.Type {
	case StoreCommandType_Set:
	case StoreCommandType_Get:
	case StoreCommandType_Delete:

	}

	PoolStoreCommand.Put(storeCmd)
	return nil
}

func (fsm *Store) Snapshot() (raft.FSMSnapshot, error) {
	return &SnapshotBuilder{}, nil
}

func (fsm *Store) Restore(rd io.ReadCloser) error {

	return nil
}

// SnapshotBuilder is used for snapshot of Raft logs
type SnapshotBuilder struct{}

func (s SnapshotBuilder) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write([]byte{})
	if err != nil {
		return err
	}
	return sink.Close()
}

func (s SnapshotBuilder) Release() {
}
