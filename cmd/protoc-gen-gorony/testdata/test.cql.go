package packageName

import (
	"git.ronaksoft.com/ronak/rony/pools"
	"git.ronaksoft.com/ronak/rony/repo/cql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"google.golang.org/protobuf/proto"
)

const (
	colData    = "data"
	colVersion = "ver"
)

// Tables
const (
	TableModel1 = "t_model1"
	ViewModel1  = "mv_model1_by_shard_key"
	TableModel2 = "t_model2"
	ViewModel2  = "mv_model2_by_p1"
)

// Columns
const (
	ColID       = "id"
	ColShardKey = "shard_key"
	ColP1       = "p1"
)

// CQLs Create
const (
	_Model1CqlCreateTable = `
CREATE TABLE IF NOT EXISTS t_model1 (
id 	 int,
shard_key 	 int,
data 	 blob,
PRIMARY KEY (id, shard_key)
)
`
	_Model2CqlCreateTable = `
CREATE TABLE IF NOT EXISTS t_model2 (
shard_key 	 int,
p1 	 blob,
id 	 bigint,
data 	 blob,
PRIMARY KEY ((id, shard_key), p1)
)
`
)

var _Model1InsertFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Insert(TableModel1).
		Columns(ColID, ColShardKey, colData).
		Query(cql.Session())
})

func Model1Insert(x *Model1) (err error) {
	q := _Model1InsertFactory.GetQuery()
	defer _Model1InsertFactory.Put(q)

	mo := proto.MarshalOptions{UseCachedSize: true}
	data := pools.Bytes.GetCap(mo.Size(x))
	defer pools.Bytes.Put(data)

	data, err = mo.MarshalAppend(data, x)
	if err != nil {
		return err
	}

	q.Bind(x.ID, x.ShardKey, data)
	err = cql.Exec(q)
	return err
}

var _Model1GetFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Select(TableModel1).
		Columns(colData).
		Where(qb.Eq(ColID), qb.Eq(ColShardKey)).
		Query(cql.Session())
})

func Model1Get(id int32, shard_key int32) (m *Model1, err error) {
	q := _Model1GetFactory.GetQuery()
	defer _Model1GetFactory.Put(q)

	data := pools.Bytes.GetCap(512)
	defer pools.Bytes.Put(data)

	q.Bind(id, shard_key)
	err = cql.Scan(q, data)
	if err != nil {
		return
	}

	err = proto.UnmarshalOptions{Merge: true}.Unmarshal(data, m)
	return m, err
}

var _Model2InsertFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Insert(TableModel2).
		Columns(ColShardKey, ColP1, ColID, colData).
		Query(cql.Session())
})

func Model2Insert(x *Model2) (err error) {
	q := _Model2InsertFactory.GetQuery()
	defer _Model2InsertFactory.Put(q)

	mo := proto.MarshalOptions{UseCachedSize: true}
	data := pools.Bytes.GetCap(mo.Size(x))
	defer pools.Bytes.Put(data)

	data, err = mo.MarshalAppend(data, x)
	if err != nil {
		return err
	}

	q.Bind(x.ShardKey, x.P1, x.ID, data)
	err = cql.Exec(q)
	return err
}

var _Model2GetFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Select(TableModel2).
		Columns(colData).
		Where(qb.Eq(ColShardKey), qb.Eq(ColP1), qb.Eq(ColID)).
		Query(cql.Session())
})

func Model2Get(shard_key int32, p1 string, id int64) (m *Model2, err error) {
	q := _Model2GetFactory.GetQuery()
	defer _Model2GetFactory.Put(q)

	data := pools.Bytes.GetCap(512)
	defer pools.Bytes.Put(data)

	q.Bind(shard_key, p1, id)
	err = cql.Scan(q, data)
	if err != nil {
		return
	}

	err = proto.UnmarshalOptions{Merge: true}.Unmarshal(data, m)
	return m, err
}
