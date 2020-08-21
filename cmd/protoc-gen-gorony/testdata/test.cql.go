package packageName

import (
	"git.ronaksoftware.com/ronak/rony/pools"
	"git.ronaksoftware.com/ronak/rony/repo/cql"
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
id 	 bigint,
shard_key 	 int,
p1 	 blob,
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

var _Model2InsertFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Insert(TableModel2).
		Columns(ColID, ColShardKey, ColP1, colData).
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

	q.Bind(x.ID, x.ShardKey, x.P1, data)
	err = cql.Exec(q)
	return err
}
