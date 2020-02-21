package scylla

import (
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

/*
   Creation Time: 2019 - Sep - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	AuthCountAll     *queryPool
	AuthInsert       *queryPool
	AuthSelect       *queryPool
	AuthUpdateBind   *queryPool
	AuthUpdateUnbind *queryPool
	AuthKeySelect    *queryPool
	AuthKeyCount     *queryPool
)

func initAuth() {
	// AuthSelect
	AuthSelect = newQueryPool(func() *gocqlx.Queryx {
		stmt9, names9 := qb.Select(DbTableAuth).
			Columns(CnAuthKey, CnUserID).
			Where(qb.Eq(CnAuthID)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt9), names9)
	})

	// AuthUpdateBind
	AuthUpdateBind = newQueryPool(func() *gocqlx.Queryx {
		stmt10, names10 := qb.Update(DbTableAuth).
			Where(qb.Eq(CnAuthID)).
			Set(CnUserID, CnAuthKey).ToCql()
		return gocqlx.Query(Session().Query(stmt10), names10)
	})

	// AuthUpdateUnbind
	AuthUpdateUnbind = newQueryPool(func() *gocqlx.Queryx {
		stmt12, names12 := qb.Update(DbTableAuth).
			Set(CnUserID).
			Where(qb.Eq(CnAuthID)).ToCql()
		return gocqlx.Query(Session().Query(stmt12), names12)
	})

	// AuthCountAll
	AuthCountAll = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbTableAuth).CountAll().Where(qb.Eq(CnAuthID)).ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	// AuthInsert
	AuthInsert = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Insert(DbTableAuth).
			Columns(CnAuthID, CnAuthKey).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

}
func initAuthKey() {
	AuthKeyCount = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbViewAuthByKey).
			CountAll().
			Where(qb.Eq(CnAuthKey)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})

	AuthKeySelect = newQueryPool(func() *gocqlx.Queryx {
		stmt, names := qb.Select(DbViewAuthByKey).
			Columns(CnAuthID, CnUserID).
			Where(qb.Eq(CnAuthKey)).
			ToCql()
		return gocqlx.Query(Session().Query(stmt), names)
	})
}
