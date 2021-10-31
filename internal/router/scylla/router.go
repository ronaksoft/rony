package scyllaRouter

import (
	"github.com/ronaksoft/rony/pools/querypool"
	"github.com/ronaksoft/rony/tools"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
)

type Router struct {
	s  gocqlx.Session
	qp map[string]*querypool.QueryPool
	t  *table.Table
	v  map[string]*table.Table
}

func New(s gocqlx.Session) *Router {
	r := &Router{
		s: s,
		t: table.New(table.Metadata{
			Name:    "entities",
			Columns: []string{"entity_id", "replica_set", "created_on", "edited_on"},
			PartKey: []string{"entity_id"},
			SortKey: []string{},
		}),
	}
	r.qp = map[string]*querypool.QueryPool{
		"insertIF": querypool.New(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Unique().Query(s)
		}),
		"insert": querypool.New(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Query(s)
		}),
		"update": querypool.New(func() *gocqlx.Queryx {
			return r.t.UpdateBuilder().Set("replica_set", "edited_on").Query(s)
		}),
		"delete": querypool.New(func() *gocqlx.Queryx {
			return r.t.DeleteBuilder().Query(s)
		}),
		"get": querypool.New(func() *gocqlx.Queryx {
			return r.t.GetQuery(s, "replica_set", "created_on", "edited_on")
		}),
	}

	return r
}

func (r *Router) Set(entityID string, replicaSet uint64, replace bool) error {
	var q *gocqlx.Queryx
	if replace {
		q = r.qp["insertIF"].GetQuery()
		defer r.qp["insertIF"].Put(q)
	} else {
		q = r.qp["insert"].GetQuery()
		defer r.qp["insert"].Put(q)
	}

	q.Bind(entityID, replicaSet, tools.TimeUnix(), tools.TimeUnix())

	return q.Exec()
}

func (r *Router) Get(entityID string) (replicaSet uint64, err error) {
	q := r.qp["get"].GetQuery()
	defer r.qp["get"].Put(q)

	q.Bind(entityID)

	err = q.Scan(&replicaSet, nil, nil)

	return
}
