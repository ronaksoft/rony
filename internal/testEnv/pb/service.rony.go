package pb

import (
	fmt "fmt"
	rony "github.com/ronaksoft/rony"
	edge "github.com/ronaksoft/rony/edge"
	edgec "github.com/ronaksoft/rony/edgec"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	cql "github.com/ronaksoft/rony/repo/cql"
	gocqlx "github.com/scylladb/gocqlx/v2"
	qb "github.com/scylladb/gocqlx/v2/qb"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_EchoRequest int64 = 1904100324

type poolEchoRequest struct {
	pool sync.Pool
}

func (p *poolEchoRequest) Get() *EchoRequest {
	x, ok := p.pool.Get().(*EchoRequest)
	if !ok {
		return &EchoRequest{}
	}
	return x
}

func (p *poolEchoRequest) Put(x *EchoRequest) {
	x.Int = 0
	x.Timestamp = 0
	x.ReplicaSet = 0
	p.pool.Put(x)
}

var PoolEchoRequest = poolEchoRequest{}

const C_EchoResponse int64 = 4192619139

type poolEchoResponse struct {
	pool sync.Pool
}

func (p *poolEchoResponse) Get() *EchoResponse {
	x, ok := p.pool.Get().(*EchoResponse)
	if !ok {
		return &EchoResponse{}
	}
	return x
}

func (p *poolEchoResponse) Put(x *EchoResponse) {
	x.Int = 0
	x.Responder = ""
	x.Timestamp = 0
	x.Delay = 0
	x.ServerID = ""
	p.pool.Put(x)
}

var PoolEchoResponse = poolEchoResponse{}

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
	registry.RegisterConstructor(1904100324, "EchoRequest")
	registry.RegisterConstructor(4192619139, "EchoResponse")
	registry.RegisterConstructor(3131464828, "Message1")
	registry.RegisterConstructor(598674886, "Message2")
	registry.RegisterConstructor(2074613123, "Model1")
	registry.RegisterConstructor(3802219577, "Model2")
	registry.RegisterConstructor(3073810188, "Echo")
	registry.RegisterConstructor(27569121, "EchoLeaderOnly")
	registry.RegisterConstructor(3809767204, "EchoTunnel")
	registry.RegisterConstructor(3639218737, "EchoDelay")
}

func (x *EchoRequest) DeepCopy(z *EchoRequest) {
	z.Int = x.Int
	z.Timestamp = x.Timestamp
	z.ReplicaSet = x.ReplicaSet
}

func (x *EchoResponse) DeepCopy(z *EchoResponse) {
	z.Int = x.Int
	z.Responder = x.Responder
	z.Timestamp = x.Timestamp
	z.Delay = x.Delay
	z.ServerID = x.ServerID
}

func (x *Message1) DeepCopy(z *Message1) {
	z.Param1 = x.Param1
	z.Param2 = x.Param2
	if x.M2 != nil {
		z.M2 = PoolMessage2.Get()
		x.M2.DeepCopy(z.M2)
	}
	for idx := range x.M2S {
		if x.M2S[idx] != nil {
			xx := PoolMessage2.Get()
			x.M2S[idx].DeepCopy(xx)
			z.M2S = append(z.M2S, xx)
		}
	}
}

func (x *Message2) DeepCopy(z *Message2) {
	z.Param1 = x.Param1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P3 = append(z.P3[:0], x.P3...)
	if x.M1 != nil {
		z.M1 = PoolMessage1.Get()
		x.M1.DeepCopy(z.M1)
	}
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

func (x *EchoRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_EchoRequest, x)
}

func (x *EchoResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_EchoResponse, x)
}

func (x *Message1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Message1, x)
}

func (x *Message2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Message2, x)
}

func (x *Model1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model1, x)
}

func (x *Model2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Model2, x)
}

func (x *EchoRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EchoResponse) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Message1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Message2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Model2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EchoRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *EchoResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Message1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Message2) Unmarshal(b []byte) error {
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
p_1 	 blob,
id 	 bigint,
shard_key 	 int,
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
		Columns("p_1", "id", "shard_key", "data").
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

	q.Bind(x.P1, x.ID, x.ShardKey, b)
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

const C_Echo int64 = 3073810188
const C_EchoLeaderOnly int64 = 27569121
const C_EchoTunnel int64 = 3809767204
const C_EchoDelay int64 = 3639218737

type ISample interface {
	Echo(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
	EchoLeaderOnly(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
	EchoTunnel(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
	EchoDelay(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
}

type SampleWrapper struct {
	h ISample
}

func RegisterSample(h ISample, e *edge.Server) {
	w := SampleWrapper{
		h: h,
	}
	w.Register(e)
}

func (sw *SampleWrapper) Register(e *edge.Server) {
	e.SetHandlers(C_Echo, false, sw.EchoWrapper)
	e.SetHandlers(C_EchoLeaderOnly, true, sw.EchoLeaderOnlyWrapper)
	e.SetHandlers(C_EchoTunnel, true, sw.EchoTunnelWrapper)
	e.SetHandlers(C_EchoDelay, true, sw.EchoDelayWrapper)
}

func (sw *SampleWrapper) EchoWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.Echo(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func (sw *SampleWrapper) EchoLeaderOnlyWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.EchoLeaderOnly(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func (sw *SampleWrapper) EchoTunnelWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.EchoTunnel(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func (sw *SampleWrapper) EchoDelayWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolEchoRequest.Get()
	defer PoolEchoRequest.Put(req)
	res := PoolEchoResponse.Get()
	defer PoolEchoResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.EchoDelay(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}

func ExecuteRemoteEcho(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_EchoRequest, req)
	err := ctx.ExecuteRemote(replicaSet, false, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_EchoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return edge.ErrUnexpectedTunnelResponse
	}
}

func ExecuteRemoteEchoLeaderOnly(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_EchoRequest, req)
	err := ctx.ExecuteRemote(replicaSet, true, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_EchoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return edge.ErrUnexpectedTunnelResponse
	}
}

func ExecuteRemoteEchoTunnel(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_EchoRequest, req)
	err := ctx.ExecuteRemote(replicaSet, true, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_EchoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return edge.ErrUnexpectedTunnelResponse
	}
}

func ExecuteRemoteEchoDelay(ctx *edge.RequestCtx, replicaSet uint64, req *EchoRequest, res *EchoResponse) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_EchoRequest, req)
	err := ctx.ExecuteRemote(replicaSet, true, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_EchoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return edge.ErrUnexpectedTunnelResponse
	}
}

type SampleClient struct {
	c edgec.Client
}

func NewSampleClient(ec edgec.Client) *SampleClient {
	return &SampleClient{
		c: ec,
	}
}

func (c *SampleClient) Echo(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Echo, req, kvs...)
	err := c.c.Send(out, in, false)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = proto.Unmarshal(in.Message, x)
		return nil, fmt.Errorf("%s:%s", x.GetCode(), x.GetItems())
	default:
		return nil, fmt.Errorf("unknown message: %d", in.GetConstructor())
	}
}

func (c *SampleClient) EchoLeaderOnly(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_EchoLeaderOnly, req, kvs...)
	err := c.c.Send(out, in, true)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = proto.Unmarshal(in.Message, x)
		return nil, fmt.Errorf("%s:%s", x.GetCode(), x.GetItems())
	default:
		return nil, fmt.Errorf("unknown message: %d", in.GetConstructor())
	}
}

func (c *SampleClient) EchoTunnel(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_EchoTunnel, req, kvs...)
	err := c.c.Send(out, in, true)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = proto.Unmarshal(in.Message, x)
		return nil, fmt.Errorf("%s:%s", x.GetCode(), x.GetItems())
	default:
		return nil, fmt.Errorf("unknown message: %d", in.GetConstructor())
	}
}

func (c *SampleClient) EchoDelay(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_EchoDelay, req, kvs...)
	err := c.c.Send(out, in, true)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = proto.Unmarshal(in.Message, x)
		return nil, fmt.Errorf("%s:%s", x.GetCode(), x.GetItems())
	default:
		return nil, fmt.Errorf("unknown message: %d", in.GetConstructor())
	}
}
