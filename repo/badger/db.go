package badger

import (
	"github.com/dgraph-io/badger/v2"
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
	_DB                    *badger.DB
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
	db, err := New(config)
	if err != nil {
		return err
	}
	_DB = db
	_ConflictRetry = config.ConflictRetries
	_ConflictRetryInterval = config.ConflictMaxInterval
	return nil
}

func New(config Config) (*badger.DB, error) {
	opt := badger.DefaultOptions(config.DirPath)
	return badger.Open(opt)
}

func DB() *badger.DB {
	return _DB
}

func Update(fn func(txn *badger.Txn) error) (err error) {
	for retry := _ConflictRetry; retry > 0; retry-- {
		err = _DB.Update(fn)
		switch err {
		case nil:
			return nil
		case badger.ErrConflict:
		default:
			return
		}
		time.Sleep(time.Duration(tools.RandomInt64(int64(_ConflictRetryInterval))))
	}
	return
}

func View(fn func(txn *badger.Txn) error) (err error) {
	for retry := _ConflictRetry; retry > 0; retry-- {
		err = _DB.View(fn)
		switch err {
		case nil:
			return nil
		case badger.ErrConflict:
		default:
			return
		}
		time.Sleep(time.Duration(tools.RandomInt64(int64(_ConflictRetryInterval))))
	}
	return
}
