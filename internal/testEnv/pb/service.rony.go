package pb

import (
	fmt "fmt"
	rony "github.com/ronaksoft/rony"
	edge "github.com/ronaksoft/rony/edge"
	edgec "github.com/ronaksoft/rony/edgec"
	registry "github.com/ronaksoft/rony/registry"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_Req1 int64 = 36419086

type poolReq1 struct {
	pool sync.Pool
}

func (p *poolReq1) Get() *Req1 {
	x, ok := p.pool.Get().(*Req1)
	if !ok {
		return &Req1{}
	}
	return x
}

func (p *poolReq1) Put(x *Req1) {
	x.Item1 = 0
	p.pool.Put(x)
}

var PoolReq1 = poolReq1{}

const C_Req2 int64 = 2602756020

type poolReq2 struct {
	pool sync.Pool
}

func (p *poolReq2) Get() *Req2 {
	x, ok := p.pool.Get().(*Req2)
	if !ok {
		return &Req2{}
	}
	return x
}

func (p *poolReq2) Put(x *Req2) {
	x.Item1 = ""
	p.pool.Put(x)
}

var PoolReq2 = poolReq2{}

const C_Res1 int64 = 807261324

type poolRes1 struct {
	pool sync.Pool
}

func (p *poolRes1) Get() *Res1 {
	x, ok := p.pool.Get().(*Res1)
	if !ok {
		return &Res1{}
	}
	return x
}

func (p *poolRes1) Put(x *Res1) {
	x.Item1 = 0
	p.pool.Put(x)
}

var PoolRes1 = poolRes1{}

const C_Res2 int64 = 2836694326

type poolRes2 struct {
	pool sync.Pool
}

func (p *poolRes2) Get() *Res2 {
	x, ok := p.pool.Get().(*Res2)
	if !ok {
		return &Res2{}
	}
	return x
}

func (p *poolRes2) Put(x *Res2) {
	x.Item1 = ""
	p.pool.Put(x)
}

var PoolRes2 = poolRes2{}

const C_EchoRequest int64 = 3196215161

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
	x.Bool = false
	x.Timestamp = 0
	p.pool.Put(x)
}

var PoolEchoRequest = poolEchoRequest{}

const C_EchoResponse int64 = 2006407526

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
	x.Bool = false
	x.Timestamp = 0
	x.Delay = 0
	x.ServerID = ""
	p.pool.Put(x)
}

var PoolEchoResponse = poolEchoResponse{}

const C_AskRequest int64 = 4187105886

type poolAskRequest struct {
	pool sync.Pool
}

func (p *poolAskRequest) Get() *AskRequest {
	x, ok := p.pool.Get().(*AskRequest)
	if !ok {
		return &AskRequest{}
	}
	return x
}

func (p *poolAskRequest) Put(x *AskRequest) {
	x.ServerID = ""
	p.pool.Put(x)
}

var PoolAskRequest = poolAskRequest{}

const C_AskResponse int64 = 3537566840

type poolAskResponse struct {
	pool sync.Pool
}

func (p *poolAskResponse) Get() *AskResponse {
	x, ok := p.pool.Get().(*AskResponse)
	if !ok {
		return &AskResponse{}
	}
	return x
}

func (p *poolAskResponse) Put(x *AskResponse) {
	x.Coordinator = ""
	x.Responder = ""
	p.pool.Put(x)
}

var PoolAskResponse = poolAskResponse{}

func init() {
	registry.RegisterConstructor(36419086, "pb.Req1")
	registry.RegisterConstructor(2602756020, "pb.Req2")
	registry.RegisterConstructor(807261324, "pb.Res1")
	registry.RegisterConstructor(2836694326, "pb.Res2")
	registry.RegisterConstructor(3196215161, "pb.EchoRequest")
	registry.RegisterConstructor(2006407526, "pb.EchoResponse")
	registry.RegisterConstructor(4187105886, "pb.AskRequest")
	registry.RegisterConstructor(3537566840, "pb.AskResponse")
}

func (x *Req1) DeepCopy(z *Req1) {
	z.Item1 = x.Item1
}

func (x *Req2) DeepCopy(z *Req2) {
	z.Item1 = x.Item1
}

func (x *Res1) DeepCopy(z *Res1) {
	z.Item1 = x.Item1
}

func (x *Res2) DeepCopy(z *Res2) {
	z.Item1 = x.Item1
}

func (x *EchoRequest) DeepCopy(z *EchoRequest) {
	z.Int = x.Int
	z.Bool = x.Bool
	z.Timestamp = x.Timestamp
}

func (x *EchoResponse) DeepCopy(z *EchoResponse) {
	z.Int = x.Int
	z.Bool = x.Bool
	z.Timestamp = x.Timestamp
	z.Delay = x.Delay
	z.ServerID = x.ServerID
}

func (x *AskRequest) DeepCopy(z *AskRequest) {
	z.ServerID = x.ServerID
}

func (x *AskResponse) DeepCopy(z *AskResponse) {
	z.Coordinator = x.Coordinator
	z.Responder = x.Responder
}

func (x *Req1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Req1, x)
}

func (x *Req2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Req2, x)
}

func (x *Res1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Res1, x)
}

func (x *Res2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Res2, x)
}

func (x *EchoRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_EchoRequest, x)
}

func (x *EchoResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_EchoResponse, x)
}

func (x *AskRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_AskRequest, x)
}

func (x *AskResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_AskResponse, x)
}

func (x *Req1) MarshalTo(b []byte) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(b, x)
}

func (x *Req2) MarshalTo(b []byte) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(b, x)
}

func (x *Res1) MarshalTo(b []byte) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(b, x)
}

func (x *Res2) MarshalTo(b []byte) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(b, x)
}

func (x *EchoRequest) MarshalTo(b []byte) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(b, x)
}

func (x *EchoResponse) MarshalTo(b []byte) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(b, x)
}

func (x *AskRequest) MarshalTo(b []byte) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(b, x)
}

func (x *AskResponse) MarshalTo(b []byte) ([]byte, error) {
	return proto.MarshalOptions{}.MarshalAppend(b, x)
}

func (x *Req1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Req2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Res1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *Res2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *EchoRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *EchoResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *AskRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

func (x *AskResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{}.Unmarshal(b, x)
}

const C_Func1 int64 = 272094254
const C_Func2 int64 = 2302576020
const C_Echo int64 = 3073810188
const C_Ask int64 = 1349233664

type ISample interface {
	Func1(ctx *edge.RequestCtx, req *Req1, res *Res1)
	Func2(ctx *edge.RequestCtx, req *Req2, res *Res2)
	Echo(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse)
	Ask(ctx *edge.RequestCtx, req *AskRequest, res *AskResponse)
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
	e.SetHandlers(C_Func1, true, sw.Func1Wrapper)
	e.SetHandlers(C_Func2, true, sw.Func2Wrapper)
	e.SetHandlers(C_Echo, true, sw.EchoWrapper)
	e.SetHandlers(C_Ask, true, sw.AskWrapper)
}

func (sw *SampleWrapper) Func1Wrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolReq1.Get()
	defer PoolReq1.Put(req)
	res := PoolRes1.Get()
	defer PoolRes1.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.Func1(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_Res1, res)
	}
}

func (sw *SampleWrapper) Func2Wrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolReq2.Get()
	defer PoolReq2.Put(req)
	res := PoolRes2.Get()
	defer PoolRes2.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.Func2(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_Res2, res)
	}
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

func (sw *SampleWrapper) AskWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolAskRequest.Get()
	defer PoolAskRequest.Put(req)
	res := PoolAskResponse.Get()
	defer PoolAskResponse.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	sw.h.Ask(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_AskResponse, res)
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

func (c *SampleClient) Func1(req *Req1, kvs ...*rony.KeyValue) (*Res1, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Func1, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_Res1:
		x := &Res1{}
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

func (c *SampleClient) Func2(req *Req2, kvs ...*rony.KeyValue) (*Res2, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Func2, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_Res2:
		x := &Res2{}
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

func (c *SampleClient) Echo(req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Echo, req, kvs...)
	err := c.c.Send(out, in)
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

func (c *SampleClient) Ask(req *AskRequest, kvs ...*rony.KeyValue) (*AskResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Ask, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_AskResponse:
		x := &AskResponse{}
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
