package cql

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"sync"
	"time"
)

/*
   Creation Time: 2020 - Aug - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type queryFactory struct {
	pool sync.Pool
	gen  queryBuilderFunc
}

type queryBuilderFunc func() *gocqlx.Queryx

var (
	speculativeAttempts = 10
	speculativeDelay    = time.Millisecond
)

func SetSpeculativeAttempts(attempts int) {
	speculativeAttempts = attempts
}

func SetSpeculativeDelay(delay time.Duration) {
	speculativeDelay = delay
}

// NewQueryFactory creates a query pool for cql queries
func NewQueryFactory(genFunc queryBuilderFunc) *queryFactory {
	return &queryFactory{
		pool: sync.Pool{},
		gen:  genFunc,
	}
}

func (qp *queryFactory) GetQuery() *gocqlx.Queryx {
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

func (qp *queryFactory) Put(q *gocqlx.Queryx) {
	qp.pool.Put(q)
}
