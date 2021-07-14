package pools

import (
	"github.com/scylladb/gocqlx/v2"
	"sync"
)

/*
   Creation Time: 2021 - Jul - 14
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type QueryPool struct {
	pool sync.Pool
	gen  queryBuilderFunc
}

type queryBuilderFunc func() *gocqlx.Queryx

// NewQueryPool creates a new query pool
func NewQueryPool(genFunc queryBuilderFunc) *QueryPool {
	return &QueryPool{
		gen: genFunc,
	}
}

func (qp *QueryPool) GetQuery() *gocqlx.Queryx {
	q, ok := qp.pool.Get().(*gocqlx.Queryx)
	if !ok {
		q = qp.gen()
		return q
	}
	return q
}

func (qp *QueryPool) Put(q *gocqlx.Queryx) {
	qp.pool.Put(q)
}
