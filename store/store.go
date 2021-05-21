package store

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/tools"
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
	db                    *badger.DB
	flusher               *tools.FlusherPool
	conflictRetry         int
	conflictRetryInterval time.Duration
)

func Init(config Config) {
	if db != nil {
		return
	}
	db = config.DB

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
}

func writeFlushFunc(_ string, entries []tools.FlushEntry) {
	wb := db.NewWriteBatch()
	for idx := range entries {
		_ = wb.SetEntry(entries[idx].Value().(*Entry))
	}
	_ = wb.Flush()
}

// DB returns the underlying object of the database
func DB() *LocalDB {
	return db
}

// Update executes a function, creating and managing a read-write transaction
// for the user. Error returned by the function is relayed by the Update method.
// It retries in case of badger.ErrConflict returned.
func Update(fn func(txn *LTxn) error) (err error) {
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
func View(fn func(txn *LTxn) error) (err error) {
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
