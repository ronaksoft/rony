package store

import (
	"github.com/dgraph-io/badger/v3"
)

/*
   Creation Time: 2021 - Feb - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type (
	Store = badger.DB
	Txn   = badger.Txn
	Entry = badger.Entry
)

var (
	ErrKeyNotFound         = badger.ErrKeyNotFound
	DefaultIteratorOptions = badger.DefaultIteratorOptions
)

func NewEntry(key, value []byte) *Entry {
	return badger.NewEntry(key, value)
}
