package store

import (
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - May - 15
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Txn interface {
	Delete(alloc *tools.Allocator, keyParts ...interface{}) error
	Set(alloc *tools.Allocator, val []byte, keyParts ...interface{}) error
	Get(alloc *tools.Allocator, keyParts ...interface{}) ([]byte, error)
	Exists(alloc *tools.Allocator, keyParts ...interface{}) bool
}

type Store interface {
	View(fn func(Txn) error) error
	Update(fn func(Txn) error) error
	Shutdown()
}
