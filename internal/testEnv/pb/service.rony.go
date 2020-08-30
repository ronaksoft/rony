package pb

import (
	fmt "fmt"
	rony "git.ronaksoft.com/ronak/rony"
	edge "git.ronaksoft.com/ronak/rony/edge"
	edgeClient "git.ronaksoft.com/ronak/rony/edgeClient"
	registry "git.ronaksoft.com/ronak/rony/registry"
	proto "google.golang.org/protobuf/proto"
	sync "sync"
)

const C_Req1 int64 = 1772509555

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
	x.Reset()
	p.pool.Put(x)
}

var PoolReq1 = poolReq1{}

const C_Req2 int64 = 4038002889

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
	x.Reset()
	p.pool.Put(x)
}

var PoolReq2 = poolReq2{}

const C_Res1 int64 = 1536179185

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
	x.Reset()
	p.pool.Put(x)
}

var PoolRes1 = poolRes1{}

const C_Res2 int64 = 3264834123

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
	x.Reset()
	p.pool.Put(x)
}

var PoolRes2 = poolRes2{}

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
	x.Reset()
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
	x.Reset()
	p.pool.Put(x)
}

var PoolEchoResponse = poolEchoResponse{}

const C_AskRequest int64 = 3206229608

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
	x.Reset()
	p.pool.Put(x)
}

var PoolAskRequest = poolAskRequest{}

const C_AskResponse int64 = 489087205

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
	x.Reset()
	p.pool.Put(x)
}

var PoolAskResponse = poolAskResponse{}

func init() {
	registry.RegisterConstructor(1772509555, "Req1")
	registry.RegisterConstructor(4038002889, "Req2")
	registry.RegisterConstructor(1536179185, "Res1")
	registry.RegisterConstructor(3264834123, "Res2")
	registry.RegisterConstructor(1904100324, "EchoRequest")
	registry.RegisterConstructor(4192619139, "EchoResponse")
	registry.RegisterConstructor(3206229608, "AskRequest")
	registry.RegisterConstructor(489087205, "AskResponse")
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

func NewSampleServer(h ISample) *SampleWrapper {
	return &SampleWrapper{
		h: h,
	}
}

func (sw *SampleWrapper) Register(e *edge.Server) {
	e.AddHandler(C_Func1, sw.Func1Wrapper)
	e.AddHandler(C_Func2, sw.Func2Wrapper)
	e.AddHandler(C_Echo, sw.EchoWrapper)
	e.AddHandler(C_Ask, sw.AskWrapper)
}

func (sw *SampleWrapper) Func1Wrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolReq1.Get()
	defer PoolReq1.Put(req)
	res := PoolRes1.Get()
	defer PoolRes1.Put(res)
	err := proto.Unmarshal(in.Message, req)
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
	err := proto.Unmarshal(in.Message, req)
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
	err := proto.Unmarshal(in.Message, req)
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
	err := proto.Unmarshal(in.Message, req)
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
	c edgeClient.Client
}

func NewSampleClient(ec edgeClient.Client) *SampleClient {
	return &SampleClient{
		c: ec,
	}
}
func (c *SampleClient) Func1(req *Req1) (*Res1, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Func1, req)
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

func (c *SampleClient) Func2(req *Req2) (*Res2, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Func2, req)
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

func (c *SampleClient) Echo(req *EchoRequest) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Echo, req)
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

func (c *SampleClient) Ask(req *AskRequest) (*AskResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_Ask, req)
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
