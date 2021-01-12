package kv

import (
	"github.com/dgraph-io/badger/v2"
	"github.com/ronaksoft/rony/tools"
	"github.com/tidwall/buntdb"
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
	_Index                 *buntdb.DB
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
	idx, err := newIndex(config)
	if err != nil {
		return err
	}
	_Index = idx
	_DB = db
	if config.ConflictRetries == 0 {
		config.ConflictRetries = defaultConflictRetries
	}
	if config.ConflictMaxInterval == 0 {
		config.ConflictMaxInterval = defaultMaxInterval
	}
	_ConflictRetry = config.ConflictRetries
	_ConflictRetryInterval = config.ConflictMaxInterval
	return nil
}

func newDB(config Config) (*badger.DB, error) {
	opt := badger.DefaultOptions(filepath.Join(config.DirPath, "badger"))
	return badger.Open(opt)
}

func newIndex(config Config) (*buntdb.DB, error) {
	idx, err := buntdb.Open(filepath.Join(config.DirPath, "bunt"))
	if err != nil {
		return nil, err
	}
	_ = idx.Shrink()
	return idx, nil
}

func DB() *badger.DB {
	return _DB
}

func Index() *buntdb.DB {
	return _Index
}

func UpdateIndex(fn func(tx *buntdb.Tx) error) (err error) {
	return _Index.Update(fn)
}

func ViewIndex(fn func(tx *buntdb.Tx) error) error {
	return _Index.View(fn)
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
