package cql

import (
	"github.com/gocql/gocql"
	"time"
)

/*
   Creation Time: 2020 - Aug - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var cqls []string

func AddCqlQuery(q string) {
	cqls = append(cqls, q)
}

func RunCqlQueries() error {
	for _, q := range cqls {
	L:
		err := Session().ExecStmt(q)
		switch err {
		case nil:
		case gocql.ErrTimeoutNoResponse:
			time.Sleep(time.Millisecond * 100)
			goto L
		default:
			return err

		}
	}
	return nil
}
