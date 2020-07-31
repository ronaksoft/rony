// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: service.proto

package pb

import (
	fmt "fmt"
	rony "git.ronaksoftware.com/ronak/rony"
	edge "git.ronaksoftware.com/ronak/rony/edge"
	edgeClient "git.ronaksoftware.com/ronak/rony/edgeClient"
	pools "git.ronaksoftware.com/ronak/rony/pools"
	proto "github.com/gogo/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

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

func NewSampleServer(h ISample) SampleWrapper {
	return SampleWrapper{
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
	err := req.Unmarshal(in.Message)
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
	err := req.Unmarshal(in.Message)
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
	err := req.Unmarshal(in.Message)
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
	err := req.Unmarshal(in.Message)
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
	b := pools.Bytes.GetLen(req.Size())
	req.MarshalToSizedBuffer(b)
	out.RequestID = c.c.GetRequestID()
	out.Constructor = C_Func1
	out.Message = append(out.Message[:0], b...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.Constructor {
	case C_Res1:
		x := &Res1{}
		_ = x.Unmarshal(in.Message)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.Message)
		return nil, fmt.Errorf("%s:%s", x.Code, x.Items)
	default:
		return nil, fmt.Errorf("unknown message: %d", in.Constructor)
	}
}
func (c *SampleClient) Func2(req *Req2) (*Res2, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	b := pools.Bytes.GetLen(req.Size())
	req.MarshalToSizedBuffer(b)
	out.RequestID = c.c.GetRequestID()
	out.Constructor = C_Func2
	out.Message = append(out.Message[:0], b...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.Constructor {
	case C_Res2:
		x := &Res2{}
		_ = x.Unmarshal(in.Message)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.Message)
		return nil, fmt.Errorf("%s:%s", x.Code, x.Items)
	default:
		return nil, fmt.Errorf("unknown message: %d", in.Constructor)
	}
}
func (c *SampleClient) Echo(req *EchoRequest) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	b := pools.Bytes.GetLen(req.Size())
	req.MarshalToSizedBuffer(b)
	out.RequestID = c.c.GetRequestID()
	out.Constructor = C_Echo
	out.Message = append(out.Message[:0], b...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.Constructor {
	case C_EchoResponse:
		x := &EchoResponse{}
		_ = x.Unmarshal(in.Message)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.Message)
		return nil, fmt.Errorf("%s:%s", x.Code, x.Items)
	default:
		return nil, fmt.Errorf("unknown message: %d", in.Constructor)
	}
}
func (c *SampleClient) Ask(req *AskRequest) (*AskResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	b := pools.Bytes.GetLen(req.Size())
	req.MarshalToSizedBuffer(b)
	out.RequestID = c.c.GetRequestID()
	out.Constructor = C_Ask
	out.Message = append(out.Message[:0], b...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.Constructor {
	case C_AskResponse:
		x := &AskResponse{}
		_ = x.Unmarshal(in.Message)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.Message)
		return nil, fmt.Errorf("%s:%s", x.Code, x.Items)
	default:
		return nil, fmt.Errorf("unknown message: %d", in.Constructor)
	}
}