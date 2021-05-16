package badgerStore

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/tools"
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
	openTxnMtx            sync.RWMutex
	openTxn               map[int64]*badger.Txn
}

func New(cfg *Config) *Store {
	return &Store{
		openTxn: map[int64]*badger.Txn{},
	}
}

func (fsm *Store) newTxn() *Txn {
	return &Txn{
		ID:    tools.RandomInt64(0),
		store: fsm,
	}
}
