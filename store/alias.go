package store

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - Feb - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	ErrKeyNotFound         = badger.ErrKeyNotFound
	DefaultIteratorOptions = badger.DefaultIteratorOptions
)

func NewEntry(key, value []byte) *Entry {
	return badger.NewEntry(key, value)
}

type (
	LocalDB = badger.DB
	LTxn    = badger.Txn
	Entry   = badger.Entry
)
type Txn interface {
	Delete(alloc *tools.Allocator, keyParts ...interface{}) error
	Set(alloc *tools.Allocator, val []byte, keyParts ...interface{}) error
	Get(alloc *tools.Allocator, keyParts ...interface{}) ([]byte, error)
	Exists(alloc *tools.Allocator, keyParts ...interface{}) bool
}

type Store interface {
	View(fn func(Txn) error) error
	Update(fn func(Txn) error) error
	ViewLocal(fn func(txn *LTxn) error) error
	UpdateLocal(fn func(txn *LTxn) error) error
	DB() *LocalDB
	Shutdown()
}
