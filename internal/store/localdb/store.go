package localdb

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/store"
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

var (
	vlogTicker          *time.Ticker // runs every 1m, check size of vlog and run GC conditionally.
	mandatoryVlogTicker *time.Ticker // runs every 10m, we always run vlog GC.
)

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

	vlogTicker = time.NewTicker(time.Minute)
	mandatoryVlogTicker = time.NewTicker(time.Minute * 10)
	go runVlogGC(db, 1<<30)

	store.Init(store.Config{
		DB:                  db,
		ConflictRetries:     cfg.ConflictRetries,
		ConflictMaxInterval: cfg.ConflictMaxInterval,
		BatchWorkers:        cfg.BatchWorkers,
		BatchSize:           cfg.BatchSize,
	})

	return st, nil
}

func newDB(config Config) (*badger.DB, error) {
	opt := badger.DefaultOptions(filepath.Join(config.DirPath, "badger"))
	opt.Logger = nil
	return badger.Open(opt)
}

func runVlogGC(db *badger.DB, threshold int64) {
	// Get initial size on start.
	_, lastVlogSize := db.Size()

	runGC := func() {
		var err error
		for err == nil {
			// If a GC is successful, immediately run it again.
			err = db.RunValueLogGC(0.7)
		}
		_, lastVlogSize = db.Size()
	}

	for {
		select {
		case <-vlogTicker.C:
			_, currentVlogSize := db.Size()
			if currentVlogSize < lastVlogSize+threshold {
				continue
			}
			runGC()
		case <-mandatoryVlogTicker.C:
			runGC()
		}
	}
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
	if vlogTicker != nil {
		vlogTicker.Stop()
	}
	if mandatoryVlogTicker != nil {
		mandatoryVlogTicker.Stop()
	}
	if s.db != nil {
		_ = s.db.Close()
	}
}

func (s *Store) DB() *store.LocalDB {
	return s.db
}
