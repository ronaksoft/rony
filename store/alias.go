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

var (
	ErrKeyNotFound         = badger.ErrKeyNotFound
	DefaultIteratorOptions = badger.DefaultIteratorOptions
)

type (
	Item    = badger.Item
	Entry   = badger.Entry
	LocalDB = badger.DB
)

func NewEntry(key, value []byte) *Entry {
	return badger.NewEntry(key, value)
}
