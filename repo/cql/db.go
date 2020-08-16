package cql

import (
	"github.com/scylladb/gocqlx/v2"
)

/*
   Creation Time: 2020 - Aug - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func Session() gocqlx.Session {
	panic("implement it")
}

func Exec(q *gocqlx.Queryx) error {
	return q.Exec()
}

func Scan(q *gocqlx.Queryx, dest ...interface{}) error {
	return q.Scan(dest...)
}
