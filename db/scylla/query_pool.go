package scylla

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"sync"
)

/*
   Creation Time: 2019 - Sep - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type queryPool struct {
	pool sync.Pool
	gen  queryBuilderFunc
}

type queryBuilderFunc func() *gocqlx.Queryx

// NewQueryPool exposed only for internal tests
func NewQueryPool(genFunc queryBuilderFunc) *queryPool {
	return newQueryPool(genFunc)
}

func newQueryPool(genFunc queryBuilderFunc) *queryPool {
	qp := new(queryPool)
	qp.gen = genFunc
	return qp
}

func (qp *queryPool) GetQuery() *gocqlx.Queryx {
	q, ok := qp.pool.Get().(*gocqlx.Queryx)
	if !ok {
		q = qp.gen()
		q.SetSpeculativeExecutionPolicy(&gocql.SimpleSpeculativeExecution{
			NumAttempts:  speculativeAttempts,
			TimeoutDelay: speculativeDelay,
		})
		return q
	}
	return q
}

func (qp *queryPool) Put(q *gocqlx.Queryx) {
	qp.pool.Put(q)
}
