package store

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/tools"
	"path/filepath"
	"time"
)

/*
   Creation Time: 2020 - Nov - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	db                    *Store
	flusher               *tools.FlusherPool
	conflictRetry         int
	conflictRetryInterval time.Duration
	vlogTicker            *time.Ticker // runs every 1m, check size of vlog and run GC conditionally.
	mandatoryVlogTicker   *time.Ticker // runs every 10m, we always run vlog GC.
)

func MustInit(config Config) {
	err := Init(config)
	if err != nil {
		panic(err)
	}

}

func Init(config Config) (err error) {
	if db != nil {
		return
	}
	db, err = newDB(config)
	if err != nil {
		return
	}
	if config.ConflictRetries == 0 {
		config.ConflictRetries = defaultConflictRetries
	}
	if config.ConflictMaxInterval == 0 {
		config.ConflictMaxInterval = defaultMaxInterval
	}
	if config.BatchWorkers == 0 {
		config.BatchWorkers = defaultBatchWorkers
	}
	if config.BatchSize == 0 {
		config.BatchSize = defaultBatchSize
	}
	conflictRetry = config.ConflictRetries
	conflictRetryInterval = config.ConflictMaxInterval
	flusher = tools.NewFlusherPool(int32(config.BatchWorkers), int32(config.BatchSize), writeFlushFunc)
	vlogTicker = time.NewTicker(time.Minute)
	mandatoryVlogTicker = time.NewTicker(time.Minute * 10)
	go runVlogGC(db, 1<<30)
	return nil
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

func newDB(config Config) (*badger.DB, error) {
	opt := badger.DefaultOptions(filepath.Join(config.DirPath, "badger"))
	return badger.Open(opt)
}

func writeFlushFunc(targetID string, entries []tools.FlushEntry) {
	wb := db.NewWriteBatch()
	for idx := range entries {
		_ = wb.SetEntry(entries[idx].Value().(*Entry))
	}
	_ = wb.Flush()
}

// DB returns the underlying object of the database
func DB() *badger.DB {
	return db
}

// Shutdown stops all the background go-routines and closed the underlying database
func Shutdown() {
	if vlogTicker != nil {
		vlogTicker.Stop()
	}
	if mandatoryVlogTicker != nil {
		mandatoryVlogTicker.Stop()
	}
	if db != nil {
		_ = db.Close()
	}
}

// Update executes a function, creating and managing a read-write transaction
// for the user. Error returned by the function is relayed by the Update method.
// It retries in case of badger.ErrConflict returned.
func Update(fn func(txn *Txn) error) (err error) {
	retry := conflictRetry
Retry:
	err = db.Update(fn)
	if err == badger.ErrConflict {
		if retry--; retry > 0 {
			metrics.IncCounter(metrics.CntStoreConflicts)
			time.Sleep(time.Duration(tools.RandomInt64(int64(conflictRetryInterval))))
			goto Retry
		}
	}
	return
}

// View executes a function creating and managing a read-only transaction for the user. Error
// returned by the function is relayed by the View method. It retries in case of badger.ErrConflict returned.
func View(fn func(txn *Txn) error) (err error) {
	retry := conflictRetry
Retry:
	err = db.View(fn)
	if err == badger.ErrConflict {
		if retry--; retry > 0 {
			metrics.IncCounter(metrics.CntStoreConflicts)
			time.Sleep(time.Duration(tools.RandomInt64(int64(conflictRetryInterval))))
			goto Retry
		}
	}
	return
}

// BatchWrite is the helper function to set the entry backed by an internal flusher. Which makes writes
// faster but there is no guarantee that write has been done successfully, since we bypass the errors.
// TODO:: Maybe improve the flusher structure to return error in the case
func BatchWrite(e *Entry) {
	flusher.EnterAndWait("", tools.NewEntry(e))
}
