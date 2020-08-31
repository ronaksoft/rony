package packageName

import (
	pools "git.ronaksoft.com/ronak/rony/pools"
	registry "git.ronaksoft.com/ronak/rony/registry"
	cql "git.ronaksoft.com/ronak/rony/repo/cql"
	gocqlx "github.com/scylladb/gocqlx/v2"
	qb "github.com/scylladb/gocqlx/v2/qb"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_Message1 int64 = 3131464828

type poolMessage1 struct {
	pool sync.Pool
}

func (p *poolMessage1) Get() *Message1 {
	x, ok := p.pool.Get().(*Message1)
	if !ok {
		return &Message1{}
	}
	return x
}

func (p *poolMessage1) Put(x *Message1) {
	x.Param1 = 0
	x.Param2 = ""
	if x.M2 != nil {
		PoolMessage2.Put(x.M2)
		x.M2 = nil
	}
	for idx := range x.M2S {
		if x.M2S[idx] != nil {
			PoolMessage2.Put(x.M2S[idx])
			x.M2S = nil
		}
	}
	x.M2S = x.M2S[:0]
	p.pool.Put(x)
}

var PoolMessage1 = poolMessage1{}

const C_Message2 int64 = 598674886

type poolMessage2 struct {
	pool sync.Pool
}

func (p *poolMessage2) Get() *Message2 {
	x, ok := p.pool.Get().(*Message2)
	if !ok {
		return &Message2{}
	}
	return x
}

func (p *poolMessage2) Put(x *Message2) {
	x.Param1 = 0
	x.P2 = x.P2[:0]
	x.P3 = x.P3[:0]
	if x.M1 != nil {
		PoolMessage1.Put(x.M1)
		x.M1 = nil
	}
	p.pool.Put(x)
}

var PoolMessage2 = poolMessage2{}

const C_Model1 int64 = 2074613123

type poolModel1 struct {
	pool sync.Pool
}

func (p *poolModel1) Get() *Model1 {
	x, ok := p.pool.Get().(*Model1)
	if !ok {
		return &Model1{}
	}
	return x
}

func (p *poolModel1) Put(x *Model1) {
	x.ID = 0
	x.ShardKey = 0
	x.P1 = ""
	x.P2 = x.P2[:0]
	x.P5 = 0
	p.pool.Put(x)
}

var PoolModel1 = poolModel1{}

const C_Model2 int64 = 3802219577

type poolModel2 struct {
	pool sync.Pool
}

func (p *poolModel2) Get() *Model2 {
	x, ok := p.pool.Get().(*Model2)
	if !ok {
		return &Model2{}
	}
	return x
}

func (p *poolModel2) Put(x *Model2) {
	x.ID = 0
	x.ShardKey = 0
	x.P1 = ""
	x.P2 = x.P2[:0]
	x.P5 = 0
	p.pool.Put(x)
}

var PoolModel2 = poolModel2{}

func init() {
	registry.RegisterConstructor(3131464828, "Message1")
	registry.RegisterConstructor(598674886, "Message2")
	registry.RegisterConstructor(2074613123, "Model1")
	registry.RegisterConstructor(3802219577, "Model2")
}

// Tables
const (
	TableModel1          = "model_1"
	ViewModel1ByShardKey = "model_1_by_shard_key"
	TableModel2          = "model_2"
	ViewModel2ByP1       = "model_2_by_p_1"
)

func init() {
	cql.AddCqlQuery(`
CREATE TABLE IF NOT EXISTS model_1 (
shard_key 	 int,
id 	 int,
data 	 blob,
PRIMARY KEY (id, shard_key)
) WITH CLUSTERING ORDER BY (shard_key ASC);
`)
	cql.AddCqlQuery(`
CREATE MATERIALIZED VIEW IF NOT EXISTS model_1_by_shard_key AS 
SELECT *
FROM model_1
WHERE shard_key IS NOT null
AND id IS NOT null
PRIMARY KEY (shard_key, id)
`)
	cql.AddCqlQuery(`
CREATE TABLE IF NOT EXISTS model_2 (
id 	 bigint,
shard_key 	 int,
p_1 	 blob,
data 	 blob,
PRIMARY KEY ((id, shard_key), p_1)
) WITH CLUSTERING ORDER BY (p_1 DESC);
`)
	cql.AddCqlQuery(`
CREATE MATERIALIZED VIEW IF NOT EXISTS model_2_by_p_1 AS 
SELECT *
FROM model_2
WHERE p_1 IS NOT null
AND shard_key IS NOT null
AND id IS NOT null
PRIMARY KEY (p_1, shard_key, id)
`)
}

var _Model1InsertFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Insert(TableModel1).
		Columns("shard_key", "id", "data").
		Query(cql.Session())
})

func Model1Insert(x *Model1) (err error) {
	q := _Model1InsertFactory.GetQuery()
	defer _Model1InsertFactory.Put(q)

	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Bytes.GetCap(mo.Size(x))
	defer pools.Bytes.Put(b)

	b, err = mo.MarshalAppend(b, x)
	if err != nil {
		return err
	}

	q.Bind(x.ShardKey, x.ID, b)
	err = cql.Exec(q)
	return err
}

var _Model1GetFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Select(TableModel1).
		Columns("data").
		Where(qb.Eq("id")).
		Query(cql.Session())
})

func Model1Get(id int32, x *Model1) (*Model1, error) {
	if x == nil {
		x = &Model1{}
	}
	q := _Model1GetFactory.GetQuery()
	defer _Model1GetFactory.Put(q)

	b := pools.Bytes.GetCap(512)
	defer pools.Bytes.Put(b)

	q.Bind(id)
	err := cql.Scan(q, &b)
	if err != nil {
		return x, err
	}

	err = proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
	return x, err
}

var _Model1ListByShardKeyFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Select(ViewModel1ByShardKey).
		Columns("data").
		Where(qb.Eq("shard_key")).
		Query(cql.Session())
})

func Model1ListByShardKey(shardKey int32, limit int32, f func(x *Model1) bool) error {
	q := _Model1ListByShardKeyFactory.GetQuery()
	defer _Model1ListByShardKeyFactory.Put(q)

	b := pools.Bytes.GetCap(512)
	defer pools.Bytes.Put(b)

	q.Bind(shardKey)
	iter := q.Iter()
	for iter.Scan(&b) {
		x := PoolModel1.Get()
		err := proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
		if err != nil {
			PoolModel1.Put(x)
			return err
		}
		if limit--; limit <= 0 || !f(x) {
			PoolModel1.Put(x)
			break
		}
		PoolModel1.Put(x)
	}
	return nil
}

var _Model2InsertFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Insert(TableModel2).
		Columns("id", "shard_key", "p_1", "data").
		Query(cql.Session())
})

func Model2Insert(x *Model2) (err error) {
	q := _Model2InsertFactory.GetQuery()
	defer _Model2InsertFactory.Put(q)

	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Bytes.GetCap(mo.Size(x))
	defer pools.Bytes.Put(b)

	b, err = mo.MarshalAppend(b, x)
	if err != nil {
		return err
	}

	q.Bind(x.ID, x.ShardKey, x.P1, b)
	err = cql.Exec(q)
	return err
}

var _Model2GetFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Select(TableModel2).
		Columns("data").
		Where(qb.Eq("id"), qb.Eq("shard_key")).
		Query(cql.Session())
})

func Model2Get(id int64, shardKey int32, x *Model2) (*Model2, error) {
	if x == nil {
		x = &Model2{}
	}
	q := _Model2GetFactory.GetQuery()
	defer _Model2GetFactory.Put(q)

	b := pools.Bytes.GetCap(512)
	defer pools.Bytes.Put(b)

	q.Bind(id, shardKey)
	err := cql.Scan(q, &b)
	if err != nil {
		return x, err
	}

	err = proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
	return x, err
}

var _Model2ListByP1Factory = cql.NewQueryFactory(func() *gocqlx.Queryx {
	return qb.Select(ViewModel2ByP1).
		Columns("data").
		Where(qb.Eq("p_1")).
		Query(cql.Session())
})

func Model2ListByP1(p1 string, limit int32, f func(x *Model2) bool) error {
	q := _Model2ListByP1Factory.GetQuery()
	defer _Model2ListByP1Factory.Put(q)

	b := pools.Bytes.GetCap(512)
	defer pools.Bytes.Put(b)

	q.Bind(p1)
	iter := q.Iter()
	for iter.Scan(&b) {
		x := PoolModel2.Get()
		err := proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
		if err != nil {
			PoolModel2.Put(x)
			return err
		}
		if limit--; limit <= 0 || !f(x) {
			PoolModel2.Put(x)
			break
		}
		PoolModel2.Put(x)
	}
	return nil
}
