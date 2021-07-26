// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.15.8
// Rony ver. v0.12.18
// Source: sample.proto

package rpc

import (
	bytes "bytes"
	fmt "fmt"
	rony "github.com/ronaksoft/rony"
	edge "github.com/ronaksoft/rony/edge"
	edgec "github.com/ronaksoft/rony/edgec"
	errors "github.com/ronaksoft/rony/errors"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	tools "github.com/ronaksoft/rony/tools"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	http "net/http"
	sync "sync"
)

var _ = pools.Imported

const C_InfoRequest int64 = 1180205313

type poolInfoRequest struct {
	pool sync.Pool
}

func (p *poolInfoRequest) Get() *InfoRequest {
	x, ok := p.pool.Get().(*InfoRequest)
	if !ok {
		x = &InfoRequest{}
	}

	return x
}

func (p *poolInfoRequest) Put(x *InfoRequest) {
	if x == nil {
		return
	}

	x.ReplicaSet = 0
	x.RandomText = ""

	p.pool.Put(x)
}

var PoolInfoRequest = poolInfoRequest{}

func (x *InfoRequest) DeepCopy(z *InfoRequest) {
	z.ReplicaSet = x.ReplicaSet
	z.RandomText = x.RandomText
}

func (x *InfoRequest) Clone() *InfoRequest {
	z := &InfoRequest{}
	x.DeepCopy(z)
	return z
}

func (x *InfoRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *InfoRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *InfoRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *InfoRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *InfoRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_InfoRequest, x)
}

const C_InfoResponse int64 = 699497142

type poolInfoResponse struct {
	pool sync.Pool
}

func (p *poolInfoResponse) Get() *InfoResponse {
	x, ok := p.pool.Get().(*InfoResponse)
	if !ok {
		x = &InfoResponse{}
	}

	return x
}

func (p *poolInfoResponse) Put(x *InfoResponse) {
	if x == nil {
		return
	}

	x.ServerID = ""
	x.RandomText = ""

	p.pool.Put(x)
}

var PoolInfoResponse = poolInfoResponse{}

func (x *InfoResponse) DeepCopy(z *InfoResponse) {
	z.ServerID = x.ServerID
	z.RandomText = x.RandomText
}

func (x *InfoResponse) Clone() *InfoResponse {
	z := &InfoResponse{}
	x.DeepCopy(z)
	return z
}

func (x *InfoResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *InfoResponse) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *InfoResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *InfoResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func (x *InfoResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_InfoResponse, x)
}

const C_SampleInfoWithClientRedirect int64 = 1298563006
const C_SampleInfoWithServerRedirect int64 = 1831149695

func init() {
	registry.RegisterConstructor(1180205313, "InfoRequest")
	registry.RegisterConstructor(699497142, "InfoResponse")
	registry.RegisterConstructor(1298563006, "SampleInfoWithClientRedirect")
	registry.RegisterConstructor(1831149695, "SampleInfoWithServerRedirect")
}

var _ = tools.TimeUnix()

type ISample interface {
	InfoWithClientRedirect(ctx *edge.RequestCtx, req *InfoRequest, res *InfoResponse)
	InfoWithServerRedirect(ctx *edge.RequestCtx, req *InfoRequest, res *InfoResponse)
}

func RegisterSample(h ISample, e *edge.Server, preHandlers ...edge.Handler) {
	w := sampleWrapper{
		h: h,
	}
	w.Register(e, func(c int64) []edge.Handler {
		return preHandlers
	})
}

func RegisterSampleWithFunc(h ISample, e *edge.Server, handlerFunc func(c int64) []edge.Handler) {
	w := sampleWrapper{
		h: h,
	}
	w.Register(e, handlerFunc)
}

type sampleWrapper struct {
	h ISample
}

func (sw *sampleWrapper) infoWithClientRedirectWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolInfoRequest.Get()
	defer PoolInfoRequest.Put(req)
	res := PoolInfoResponse.Get()
	defer PoolInfoResponse.Put(res)

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	sw.h.InfoWithClientRedirect(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_InfoResponse, res)
	}
}
func (sw *sampleWrapper) infoWithServerRedirectWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := PoolInfoRequest.Get()
	defer PoolInfoRequest.Put(req)
	res := PoolInfoResponse.Get()
	defer PoolInfoResponse.Put(res)

	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	sw.h.InfoWithServerRedirect(ctx, req, res)
	if !ctx.Stopped() {
		ctx.PushMessage(C_InfoResponse, res)
	}
}

func (sw *sampleWrapper) Register(e *edge.Server, handlerFunc func(c int64) []edge.Handler) {
	if handlerFunc == nil {
		handlerFunc = func(c int64) []edge.Handler {
			return nil
		}
	}
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_SampleInfoWithClientRedirect).
			SetHandler(handlerFunc(C_SampleInfoWithClientRedirect)...).
			Append(sw.infoWithClientRedirectWrapper),
	)
	e.SetRestProxy(
		"get", "/info/client-redirect/:ReplicaSet/:RandomText",
		edge.NewRestProxy(sw.infoWithClientRedirectRestClient, sw.infoWithClientRedirectRestServer),
	)
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_SampleInfoWithServerRedirect).
			SetHandler(handlerFunc(C_SampleInfoWithServerRedirect)...).
			Append(sw.infoWithServerRedirectWrapper),
	)
	e.SetRestProxy(
		"get", "/info/server-redirect/:ReplicaSet/:RandomText",
		edge.NewRestProxy(sw.infoWithServerRedirectRestClient, sw.infoWithServerRedirectRestServer),
	)
}

func TunnelRequestSampleInfoWithClientRedirect(ctx *edge.RequestCtx, replicaSet uint64, req *InfoRequest, res *InfoResponse, kvs ...*rony.KeyValue) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleInfoWithClientRedirect, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_InfoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return errors.ErrUnexpectedTunnelResponse
	}
}
func TunnelRequestSampleInfoWithServerRedirect(ctx *edge.RequestCtx, replicaSet uint64, req *InfoRequest, res *InfoResponse, kvs ...*rony.KeyValue) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleInfoWithServerRedirect, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_InfoResponse:
		_ = res.Unmarshal(in.GetMessage())
		return nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return errors.ErrUnexpectedTunnelResponse
	}
}

func (sw *sampleWrapper) infoWithClientRedirectRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := PoolInfoRequest.Get()
	defer PoolInfoRequest.Put(req)
	req.ReplicaSet = tools.StrToUInt64(tools.GetString(conn.Get("ReplicaSet"), "0"))
	req.RandomText = tools.GetString(conn.Get("RandomText"), "")

	ctx.FillEnvelope(conn.ConnID(), C_SampleInfoWithClientRedirect, req)
	return nil
}
func (sw *sampleWrapper) infoWithClientRedirectRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_InfoResponse:
			x := &InfoResponse{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			b, err = x.MarshalJSON()
			if err != nil {
				return
			}
			err = conn.WriteBinary(ctx.StreamID(), b)
			return
		case rony.C_Error:
			x := &rony.Error{}
			_ = x.Unmarshal(envelope.Message)
			err = x
			return
		case rony.C_Redirect:
			x := &rony.Redirect{}
			_ = x.Unmarshal(envelope.Message)
			if len(x.Edges) == 0 || len(x.Edges[0].HostPorts) == 0 {
				break
			}
			switch x.Reason {
			case rony.RedirectReason_ReplicaSetSession:
				conn.Redirect(http.StatusPermanentRedirect, x.Edges[0].HostPorts[0])
			case rony.RedirectReason_ReplicaSetRequest:
				conn.Redirect(http.StatusTemporaryRedirect, x.Edges[0].HostPorts[0])
			}
			return
		}
		err = errors.ErrUnexpectedResponse
	}) {
		err = errors.ErrInternalServer
	}

	return
}

func (sw *sampleWrapper) infoWithServerRedirectRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := PoolInfoRequest.Get()
	defer PoolInfoRequest.Put(req)
	req.ReplicaSet = tools.StrToUInt64(tools.GetString(conn.Get("ReplicaSet"), "0"))
	req.RandomText = tools.GetString(conn.Get("RandomText"), "")

	ctx.FillEnvelope(conn.ConnID(), C_SampleInfoWithServerRedirect, req)
	return nil
}
func (sw *sampleWrapper) infoWithServerRedirectRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_InfoResponse:
			x := &InfoResponse{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			b, err = x.MarshalJSON()
			if err != nil {
				return
			}
			err = conn.WriteBinary(ctx.StreamID(), b)
			return
		case rony.C_Error:
			x := &rony.Error{}
			_ = x.Unmarshal(envelope.Message)
			err = x
			return
		case rony.C_Redirect:
			x := &rony.Redirect{}
			_ = x.Unmarshal(envelope.Message)
			if len(x.Edges) == 0 || len(x.Edges[0].HostPorts) == 0 {
				break
			}
			switch x.Reason {
			case rony.RedirectReason_ReplicaSetSession:
				conn.Redirect(http.StatusPermanentRedirect, x.Edges[0].HostPorts[0])
			case rony.RedirectReason_ReplicaSetRequest:
				conn.Redirect(http.StatusTemporaryRedirect, x.Edges[0].HostPorts[0])
			}
			return
		}
		err = errors.ErrUnexpectedResponse
	}) {
		err = errors.ErrInternalServer
	}

	return
}

type SampleClient struct {
	c edgec.Client
}

func NewSampleClient(ec edgec.Client) *SampleClient {
	return &SampleClient{
		c: ec,
	}
}
func (c *SampleClient) InfoWithClientRedirect(req *InfoRequest, kvs ...*rony.KeyValue) (*InfoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleInfoWithClientRedirect, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_InfoResponse:
		x := &InfoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.Message)
		return nil, x
	default:
		return nil, fmt.Errorf("unknown message :%d", in.GetConstructor())
	}
}
func (c *SampleClient) InfoWithServerRedirect(req *InfoRequest, kvs ...*rony.KeyValue) (*InfoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleInfoWithServerRedirect, req, kvs...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_InfoResponse:
		x := &InfoResponse{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.Message)
		return nil, x
	default:
		return nil, fmt.Errorf("unknown message :%d", in.GetConstructor())
	}
}

var _ = bytes.MinRead
