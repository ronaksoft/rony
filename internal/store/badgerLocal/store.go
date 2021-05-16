package badgerLocal

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/internal/store"
	staticStore "github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
	"path/filepath"
	"time"
)

/*
   Creation Time: 2021 - May - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Store struct {
	db *badger.DB
}

func New(cfg Config) (*Store, error) {
	st := &Store{}
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

func (s *Store) ViewLocal(fn func(txn *store.LTxn) error) error {
	retry := defaultConflictRetries
Retry:
	err := s.db.View(fn)
	if err == badger.ErrConflict {
		if retry--; retry > 0 {
			metrics.IncCounter(metrics.CntStoreConflicts)
			time.Sleep(time.Duration(tools.RandomInt64(int64(defaultMaxInterval))))
			goto Retry
		}
	}
	return err
}

func (s *Store) UpdateLocal(fn func(txn *store.LTxn) error) error {
	retry := defaultConflictRetries
Retry:
	err := s.db.Update(fn)
	if err == badger.ErrConflict {
		if retry--; retry > 0 {
			metrics.IncCounter(metrics.CntStoreConflicts)
			time.Sleep(time.Duration(tools.RandomInt64(int64(defaultMaxInterval))))
			goto Retry
		}
	}
	return err
}

func (s *Store) View(fn func(store.Txn) error) error {
	panic("BUG! not supported")
}

func (s *Store) Update(fn func(store.Txn) error) error {
	panic("BUG! not supported")
}

func (s *Store) Shutdown() {
	panic("implement me")
}

func (s *Store) DB() *store.DB {
	return s.db
}
