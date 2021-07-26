package rony

import (
	"github.com/dgraph-io/badger/v3"
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
	LocalDB  = badger.DB
	StoreTxn = badger.Txn
)

type Store interface {
	View(fn func(*StoreTxn) error) error
	Update(fn func(*StoreTxn) error) error
	LocalDB() *LocalDB
	Shutdown()
}
