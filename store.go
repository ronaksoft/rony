package rony

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - Jul - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type (
	StoreLocalTxn = badger.Txn
)

type StoreTxn interface {
	Delete(alloc *tools.Allocator, keyParts ...interface{}) error
	Set(alloc *tools.Allocator, val []byte, keyParts ...interface{}) error
	Get(alloc *tools.Allocator, keyParts ...interface{}) ([]byte, error)
	Exists(alloc *tools.Allocator, keyParts ...interface{}) bool
}

type Store interface {
	View(fn func(StoreTxn) error) error
	Update(fn func(StoreTxn) error) error
	ViewLocal(fn func(txn *StoreLocalTxn) error) error
	UpdateLocal(fn func(txn *StoreLocalTxn) error) error
	Shutdown()
}
