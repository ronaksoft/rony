package model

import (
	edge "github.com/ronaksoft/rony/edge"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	cql "github.com/ronaksoft/rony/repo/cql"
	gocqlx "github.com/scylladb/gocqlx/v2"
	qb "github.com/scylladb/gocqlx/v2/qb"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_Hook int64 = 74116203

type poolHook struct {
	pool sync.Pool
}

func (p *poolHook) Get() *Hook {
	x, ok := p.pool.Get().(*Hook)
	if !ok {
		return &Hook{}
	}
	return x
}

func (p *poolHook) Put(x *Hook) {
	x.ClientID = ""
	x.ID = ""
	x.Timestamp = ""
	x.HookUrl = ""
	x.Fired = false
	x.Success = false
	p.pool.Put(x)
}

var PoolHook = poolHook{}

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
	registry.RegisterConstructor(74116203, "Hook")
	registry.RegisterConstructor(2074613123, "Model1")
	registry.RegisterConstructor(3802219577, "Model2")
}

func (x *Hook) DeepCopy(z *Hook) {
	z.ClientID = x.ClientID
	z.ID = x.ID
	z.Timestamp = x.Timestamp
	z.HookUrl = x.HookUrl
	z.Fired = x.Fired
	z.Success = x.Success
}

func (x *Model1) DeepCopy(z *Model1) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
}

func (x *Model2) DeepCopy(z *Model2) {
	z.ID = x.ID
	z.ShardKey = x.ShardKey
	z.P1 = x.P1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P5 = x.P5
}

func (x *Hook) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Hook, x)
}

func (x *Model1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model1, x)
}

func (x *Model2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model2, x)
}

func (x *Hook) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Hook) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Model2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
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
id 	 int,
shard_key 	 int,
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
		Columns("id", "shard_key", "data").
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

	q.Bind(x.ID, x.ShardKey, b)
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
