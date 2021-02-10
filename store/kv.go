package store

import (
	"github.com/dgraph-io/badger/v3"
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
	_DB                    *badger.DB
	_Flusher               *tools.FlusherPool
	_ConflictRetry         int
	_ConflictRetryInterval time.Duration
)

func MustInit(config Config) {
	err := Init(config)
	if err != nil {
		panic(err)
	}
}

func Init(config Config) error {
	db, err := newDB(config)
	if err != nil {
		return err
	}
	_DB = db
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
	_ConflictRetry = config.ConflictRetries
	_ConflictRetryInterval = config.ConflictMaxInterval

	_Flusher = tools.NewFlusherPool(int32(config.BatchWorkers), int32(config.BatchSize), writeFlushFunc)
	return nil
}

func newDB(config Config) (*badger.DB, error) {
	opt := badger.DefaultOptions(filepath.Join(config.DirPath, "badger"))
	return badger.Open(opt)
}

func writeFlushFunc(targetID string, entries []tools.FlushEntry) {
	wb := _DB.NewWriteBatch()
	for idx := range entries {
		_ = wb.SetEntry(entries[idx].Value().(*badger.Entry))
	}
	_ = wb.Flush()
}

// DB returns the underlying object of the database
func DB() *badger.DB {
	return _DB
}

// Update executes a function, creating and managing a read-write transaction
// for the user. Error returned by the function is relayed by the Update method.
// It retries in case of badger.ErrConflict returned.
func Update(fn func(txn *badger.Txn) error) (err error) {
	for retry := _ConflictRetry; retry > 0; retry-- {
		err = _DB.Update(fn)
		switch err {
		case badger.ErrConflict:
		default:
			return
		}
		time.Sleep(time.Duration(tools.RandomInt64(int64(_ConflictRetryInterval))))
	}
	return
}

// View executes a function creating and managing a read-only transaction for the user. Error
// returned by the function is relayed by the View method. It retries in case of badger.ErrConflict returned.
func View(fn func(txn *badger.Txn) error) (err error) {
	for retry := _ConflictRetry; retry > 0; retry-- {
		err = _DB.View(fn)
		switch err {
		case badger.ErrConflict:
		default:
			return
		}
		time.Sleep(time.Duration(tools.RandomInt64(int64(_ConflictRetryInterval))))
	}
	return
}

// BatchWrite is the helper function to set the entry backed by an internal flusher. Which makes writes
// faster but there is no guarantee that write has been done successfully, since we bypass the errors.
// TODO:: Maybe improve the flusher structure to return error in the case
func BatchWrite(e *badger.Entry) {
	_Flusher.EnterAndWait("", tools.NewEntry(e))
}
