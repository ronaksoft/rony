// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.16.22
// Source: sample.proto

package rpc

import (
	bytes "bytes"
	context "context"
	fmt "fmt"
	http "net/http"
	sync "sync"

	rony "github.com/ronaksoft/rony"
	edge "github.com/ronaksoft/rony/edge"
	edgec "github.com/ronaksoft/rony/edgec"
	errors "github.com/ronaksoft/rony/errors"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	tools "github.com/ronaksoft/rony/tools"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _ = pools.Imported

const C_InfoRequest uint64 = 393641345359739852

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

func factoryInfoRequest() registry.Message {
	return &InfoRequest{}
}

func (x *InfoRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_InfoRequest, x)
}

const C_InfoResponse uint64 = 10205569623592301307

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

func factoryInfoResponse() registry.Message {
	return &InfoResponse{}
}

func (x *InfoResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_InfoResponse, x)
}

const C_SampleInfoWithClientRedirect uint64 = 11847639677132253675
const C_SampleInfoWithServerRedirect uint64 = 11816518540128942571

// register constructors of the messages to the registry package
func init() {
	registry.Register(393641345359739852, "InfoRequest", factoryInfoRequest)
	registry.Register(10205569623592301307, "InfoResponse", factoryInfoResponse)
	registry.Register(11847639677132253675, "SampleInfoWithClientRedirect", factoryInfoRequest)
	registry.Register(11816518540128942571, "SampleInfoWithServerRedirect", factoryInfoRequest)

}

var _ = tools.TimeUnix()

type ISample interface {
	InfoWithClientRedirect(ctx *edge.RequestCtx, req *InfoRequest, res *InfoResponse) *rony.Error
	InfoWithServerRedirect(ctx *edge.RequestCtx, req *InfoRequest, res *InfoResponse) *rony.Error
}

func RegisterSample(h ISample, e *edge.Server, preHandlers ...edge.Handler) {
	w := sampleWrapper{
		h: h,
	}
	w.Register(e, func(c uint64) []edge.Handler {
		return preHandlers
	})
}

func RegisterSampleWithFunc(h ISample, e *edge.Server, handlerFunc func(c uint64) []edge.Handler) {
	w := sampleWrapper{
		h: h,
	}
	w.Register(e, handlerFunc)
}

type sampleWrapper struct {
	h ISample
}

func (sw *sampleWrapper) infoWithClientRedirectWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &InfoRequest{}
	res := &InfoResponse{}

	var err error
	if in.JsonEncoded {
		err = protojson.Unmarshal(in.Message, req)
	} else {
		err = proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	}
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.InfoWithClientRedirect(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_InfoResponse, res)
	}
}
func (sw *sampleWrapper) infoWithServerRedirectWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &InfoRequest{}
	res := &InfoResponse{}

	var err error
	if in.JsonEncoded {
		err = protojson.Unmarshal(in.Message, req)
	} else {
		err = proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	}
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.InfoWithServerRedirect(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_InfoResponse, res)
	}
}

func (sw *sampleWrapper) Register(e *edge.Server, handlerFunc func(c uint64) []edge.Handler) {
	if handlerFunc == nil {
		handlerFunc = func(c uint64) []edge.Handler {
			return nil
		}
	}
	e.SetHandler(
		edge.NewHandlerOptions().
			SetConstructor(C_SampleInfoWithClientRedirect).
			SetServiceName("Sample").
			SetMethodName("InfoWithClientRedirect").
			SetHandler(handlerFunc(C_SampleInfoWithClientRedirect)...).
			Append(sw.infoWithClientRedirectWrapper),
	)
	e.SetRestProxy(
		"get", "/info/client-redirect/:ReplicaSet/:RandomText",
		edge.NewRestProxy(sw.infoWithClientRedirectRestClient, sw.infoWithClientRedirectRestServer),
	)
	e.SetHandler(
		edge.NewHandlerOptions().
			SetConstructor(C_SampleInfoWithServerRedirect).
			SetServiceName("Sample").
			SetMethodName("InfoWithServerRedirect").
			SetHandler(handlerFunc(C_SampleInfoWithServerRedirect)...).
			Append(sw.infoWithServerRedirectWrapper),
	)
	e.SetRestProxy(
		"get", "/info/server-redirect/:ReplicaSet/:RandomText",
		edge.NewRestProxy(sw.infoWithServerRedirectRestClient, sw.infoWithServerRedirectRestServer),
	)
}

func TunnelRequestSampleInfoWithClientRedirect(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *InfoRequest, res *InfoResponse,
	kvs ...*rony.KeyValue,
) error {
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
func TunnelRequestSampleInfoWithServerRedirect(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *InfoRequest, res *InfoResponse,
	kvs ...*rony.KeyValue,
) error {
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
	req := &InfoRequest{}
	req.ReplicaSet = tools.StrToUInt64(tools.GetString(conn.Get("ReplicaSet"), "0"))
	req.RandomText = tools.GetString(conn.Get("RandomText"), "")

	ctx.Fill(conn.ConnID(), C_SampleInfoWithClientRedirect, req)
	return nil
}
func (sw *sampleWrapper) infoWithClientRedirectRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	conn.WriteHeader("Content-Type", "application/json")
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
	req := &InfoRequest{}
	req.ReplicaSet = tools.StrToUInt64(tools.GetString(conn.Get("ReplicaSet"), "0"))
	req.RandomText = tools.GetString(conn.Get("RandomText"), "")

	ctx.Fill(conn.ConnID(), C_SampleInfoWithServerRedirect, req)
	return nil
}
func (sw *sampleWrapper) infoWithServerRedirectRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	conn.WriteHeader("Content-Type", "application/json")
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

type ISampleClient interface {
	InfoWithClientRedirect(ctx context.Context, req *InfoRequest, kvs ...*rony.KeyValue) (*InfoResponse, error)
	InfoWithServerRedirect(ctx context.Context, req *InfoRequest, kvs ...*rony.KeyValue) (*InfoResponse, error)
}

type SampleClient struct {
	name string
	c    edgec.Client
}

func NewSampleClient(name string, ec edgec.Client) *SampleClient {
	return &SampleClient{
		name: name,
		c:    ec,
	}
}
func (c *SampleClient) InfoWithClientRedirect(
	ctx context.Context, req *InfoRequest, kvs ...*rony.KeyValue,
) (*InfoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleInfoWithClientRedirect, req, kvs...)
	err := c.c.Send(ctx, out, in)
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
func (c *SampleClient) InfoWithServerRedirect(
	ctx context.Context, req *InfoRequest, kvs ...*rony.KeyValue,
) (*InfoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleInfoWithServerRedirect, req, kvs...)
	err := c.c.Send(ctx, out, in)
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
