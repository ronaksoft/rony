// Code generated by Rony's protoc plugin; DO NOT EDIT.
// ProtoC ver. v3.17.3
// Rony ver. v0.16.4
// Source: service.proto

package service

import (
	bytes "bytes"
	context "context"
	fmt "fmt"
	http "net/http"
	sync "sync"

	rony "github.com/ronaksoft/rony"
	config "github.com/ronaksoft/rony/config"
	edge "github.com/ronaksoft/rony/edge"
	edgec "github.com/ronaksoft/rony/edgec"
	errors "github.com/ronaksoft/rony/errors"
	pools "github.com/ronaksoft/rony/pools"
	registry "github.com/ronaksoft/rony/registry"
	tools "github.com/ronaksoft/rony/tools"
	cobra "github.com/spf13/cobra"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _ = pools.Imported

const C_GetRequest uint64 = 8186060648624618456

type poolGetRequest struct {
	pool sync.Pool
}

func (p *poolGetRequest) Get() *GetRequest {
	x, ok := p.pool.Get().(*GetRequest)
	if !ok {
		x = &GetRequest{}
	}

	return x
}

func (p *poolGetRequest) Put(x *GetRequest) {
	if x == nil {
		return
	}

	x.Key = x.Key[:0]

	p.pool.Put(x)
}

var PoolGetRequest = poolGetRequest{}

func (x *GetRequest) DeepCopy(z *GetRequest) {
	z.Key = append(z.Key[:0], x.Key...)
}

func (x *GetRequest) Clone() *GetRequest {
	z := &GetRequest{}
	x.DeepCopy(z)
	return z
}

func (x *GetRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *GetRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *GetRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *GetRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryGetRequest() registry.Message {
	return &GetRequest{}
}

func (x *GetRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_GetRequest, x)
}

const C_GetResponse uint64 = 10382375233116730107

type poolGetResponse struct {
	pool sync.Pool
}

func (p *poolGetResponse) Get() *GetResponse {
	x, ok := p.pool.Get().(*GetResponse)
	if !ok {
		x = &GetResponse{}
	}

	return x
}

func (p *poolGetResponse) Put(x *GetResponse) {
	if x == nil {
		return
	}

	x.Key = x.Key[:0]
	x.Value = x.Value[:0]

	p.pool.Put(x)
}

var PoolGetResponse = poolGetResponse{}

func (x *GetResponse) DeepCopy(z *GetResponse) {
	z.Key = append(z.Key[:0], x.Key...)
	z.Value = append(z.Value[:0], x.Value...)
}

func (x *GetResponse) Clone() *GetResponse {
	z := &GetResponse{}
	x.DeepCopy(z)
	return z
}

func (x *GetResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *GetResponse) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *GetResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *GetResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryGetResponse() registry.Message {
	return &GetResponse{}
}

func (x *GetResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_GetResponse, x)
}

const C_SetRequest uint64 = 8181913290764647384

type poolSetRequest struct {
	pool sync.Pool
}

func (p *poolSetRequest) Get() *SetRequest {
	x, ok := p.pool.Get().(*SetRequest)
	if !ok {
		x = &SetRequest{}
	}

	return x
}

func (p *poolSetRequest) Put(x *SetRequest) {
	if x == nil {
		return
	}

	x.Key = x.Key[:0]
	x.Value = x.Value[:0]

	p.pool.Put(x)
}

var PoolSetRequest = poolSetRequest{}

func (x *SetRequest) DeepCopy(z *SetRequest) {
	z.Key = append(z.Key[:0], x.Key...)
	z.Value = append(z.Value[:0], x.Value...)
}

func (x *SetRequest) Clone() *SetRequest {
	z := &SetRequest{}
	x.DeepCopy(z)
	return z
}

func (x *SetRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *SetRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *SetRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *SetRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factorySetRequest() registry.Message {
	return &SetRequest{}
}

func (x *SetRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_SetRequest, x)
}

const C_SetResponse uint64 = 10382356249361281787

type poolSetResponse struct {
	pool sync.Pool
}

func (p *poolSetResponse) Get() *SetResponse {
	x, ok := p.pool.Get().(*SetResponse)
	if !ok {
		x = &SetResponse{}
	}

	return x
}

func (p *poolSetResponse) Put(x *SetResponse) {
	if x == nil {
		return
	}

	x.OK = false

	p.pool.Put(x)
}

var PoolSetResponse = poolSetResponse{}

func (x *SetResponse) DeepCopy(z *SetResponse) {
	z.OK = x.OK
}

func (x *SetResponse) Clone() *SetResponse {
	z := &SetResponse{}
	x.DeepCopy(z)
	return z
}

func (x *SetResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *SetResponse) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *SetResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *SetResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factorySetResponse() registry.Message {
	return &SetResponse{}
}

func (x *SetResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_SetResponse, x)
}

const C_EchoRequest uint64 = 634453994073422796

type poolEchoRequest struct {
	pool sync.Pool
}

func (p *poolEchoRequest) Get() *EchoRequest {
	x, ok := p.pool.Get().(*EchoRequest)
	if !ok {
		x = &EchoRequest{}
	}

	return x
}

func (p *poolEchoRequest) Put(x *EchoRequest) {
	if x == nil {
		return
	}

	x.Int = 0
	x.Timestamp = 0
	x.ReplicaSet = 0

	p.pool.Put(x)
}

var PoolEchoRequest = poolEchoRequest{}

func (x *EchoRequest) DeepCopy(z *EchoRequest) {
	z.Int = x.Int
	z.Timestamp = x.Timestamp
	z.ReplicaSet = x.ReplicaSet
}

func (x *EchoRequest) Clone() *EchoRequest {
	z := &EchoRequest{}
	x.DeepCopy(z)
	return z
}

func (x *EchoRequest) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *EchoRequest) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EchoRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *EchoRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryEchoRequest() registry.Message {
	return &EchoRequest{}
}

func (x *EchoRequest) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_EchoRequest, x)
}

const C_EchoResponse uint64 = 10208763112635265787

type poolEchoResponse struct {
	pool sync.Pool
}

func (p *poolEchoResponse) Get() *EchoResponse {
	x, ok := p.pool.Get().(*EchoResponse)
	if !ok {
		x = &EchoResponse{}
	}

	return x
}

func (p *poolEchoResponse) Put(x *EchoResponse) {
	if x == nil {
		return
	}

	x.Int = 0
	x.Responder = ""
	x.Timestamp = 0
	x.Delay = 0
	x.ServerID = ""

	p.pool.Put(x)
}

var PoolEchoResponse = poolEchoResponse{}

func (x *EchoResponse) DeepCopy(z *EchoResponse) {
	z.Int = x.Int
	z.Responder = x.Responder
	z.Timestamp = x.Timestamp
	z.Delay = x.Delay
	z.ServerID = x.ServerID
}

func (x *EchoResponse) Clone() *EchoResponse {
	z := &EchoResponse{}
	x.DeepCopy(z)
	return z
}

func (x *EchoResponse) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *EchoResponse) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *EchoResponse) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *EchoResponse) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryEchoResponse() registry.Message {
	return &EchoResponse{}
}

func (x *EchoResponse) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_EchoResponse, x)
}

const C_Message1 uint64 = 1806736971761742569

type poolMessage1 struct {
	pool sync.Pool
}

func (p *poolMessage1) Get() *Message1 {
	x, ok := p.pool.Get().(*Message1)
	if !ok {
		x = &Message1{}
	}

	x.M2 = PoolMessage2.Get()

	return x
}

func (p *poolMessage1) Put(x *Message1) {
	if x == nil {
		return
	}

	x.Param1 = 0
	x.Param2 = ""
	PoolMessage2.Put(x.M2)
	for _, z := range x.M2S {
		PoolMessage2.Put(z)
	}
	x.M2S = x.M2S[:0]

	p.pool.Put(x)
}

var PoolMessage1 = poolMessage1{}

func (x *Message1) DeepCopy(z *Message1) {
	z.Param1 = x.Param1
	z.Param2 = x.Param2
	if x.M2 != nil {
		if z.M2 == nil {
			z.M2 = PoolMessage2.Get()
		}
		x.M2.DeepCopy(z.M2)
	} else {
		PoolMessage2.Put(z.M2)
		z.M2 = nil
	}
	for idx := range x.M2S {
		if x.M2S[idx] == nil {
			continue
		}
		xx := PoolMessage2.Get()
		x.M2S[idx].DeepCopy(xx)
		z.M2S = append(z.M2S, xx)
	}
}

func (x *Message1) Clone() *Message1 {
	z := &Message1{}
	x.DeepCopy(z)
	return z
}

func (x *Message1) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Message1) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Message1) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Message1) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryMessage1() registry.Message {
	return &Message1{}
}

func (x *Message1) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Message1, x)
}

const C_Message2 uint64 = 2000391755738673897

type poolMessage2 struct {
	pool sync.Pool
}

func (p *poolMessage2) Get() *Message2 {
	x, ok := p.pool.Get().(*Message2)
	if !ok {
		x = &Message2{}
	}

	x.M1 = PoolMessage1.Get()

	return x
}

func (p *poolMessage2) Put(x *Message2) {
	if x == nil {
		return
	}

	x.Param1 = 0
	x.P2 = x.P2[:0]
	x.P3 = x.P3[:0]
	PoolMessage1.Put(x.M1)

	p.pool.Put(x)
}

var PoolMessage2 = poolMessage2{}

func (x *Message2) DeepCopy(z *Message2) {
	z.Param1 = x.Param1
	z.P2 = append(z.P2[:0], x.P2...)
	z.P3 = append(z.P3[:0], x.P3...)
	if x.M1 != nil {
		if z.M1 == nil {
			z.M1 = PoolMessage1.Get()
		}
		x.M1.DeepCopy(z.M1)
	} else {
		PoolMessage1.Put(z.M1)
		z.M1 = nil
	}
}

func (x *Message2) Clone() *Message2 {
	z := &Message2{}
	x.DeepCopy(z)
	return z
}

func (x *Message2) Unmarshal(b []byte) error {
	return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
}

func (x *Message2) Marshal() ([]byte, error) {
	return proto.Marshal(x)
}

func (x *Message2) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, x)
}

func (x *Message2) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(x)
}

func factoryMessage2() registry.Message {
	return &Message2{}
}

func (x *Message2) PushToContext(ctx *edge.RequestCtx) {
	ctx.PushMessage(C_Message2, x)
}

const C_SampleEcho uint64 = 5610266072904040111
const C_SampleSet uint64 = 859899950148792148
const C_SampleGet uint64 = 859879334305771348
const C_SampleEchoTunnel uint64 = 16129836997487988187
const C_SampleEchoInternal uint64 = 8481593834425277560
const C_SampleEchoDelay uint64 = 5258191516040289195

// register constructors of the messages to the registry package
func init() {
	registry.Register(8186060648624618456, "GetRequest", factoryGetRequest)
	registry.Register(10382375233116730107, "GetResponse", factoryGetResponse)
	registry.Register(8181913290764647384, "SetRequest", factorySetRequest)
	registry.Register(10382356249361281787, "SetResponse", factorySetResponse)
	registry.Register(634453994073422796, "EchoRequest", factoryEchoRequest)
	registry.Register(10208763112635265787, "EchoResponse", factoryEchoResponse)
	registry.Register(1806736971761742569, "Message1", factoryMessage1)
	registry.Register(2000391755738673897, "Message2", factoryMessage2)
	registry.Register(5610266072904040111, "SampleEcho", factoryEchoRequest)
	registry.Register(859899950148792148, "SampleSet", factorySetRequest)
	registry.Register(859879334305771348, "SampleGet", factoryGetRequest)
	registry.Register(16129836997487988187, "SampleEchoTunnel", factoryEchoRequest)
	registry.Register(8481593834425277560, "SampleEchoInternal", factoryEchoRequest)
	registry.Register(5258191516040289195, "SampleEchoDelay", factoryEchoRequest)

}

var _ = tools.TimeUnix()

type ISample interface {
	Echo(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) *rony.Error
	Set(ctx *edge.RequestCtx, req *SetRequest, res *SetResponse) *rony.Error
	Get(ctx *edge.RequestCtx, req *GetRequest, res *GetResponse) *rony.Error
	EchoTunnel(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) *rony.Error
	EchoInternal(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) *rony.Error
	EchoDelay(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) *rony.Error
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

func (sw *sampleWrapper) echoWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &EchoRequest{}
	res := &EchoResponse{}

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

	rErr := sw.h.Echo(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}
func (sw *sampleWrapper) setWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &SetRequest{}
	res := &SetResponse{}

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

	rErr := sw.h.Set(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_SetResponse, res)
	}
}
func (sw *sampleWrapper) getWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &GetRequest{}
	res := &GetResponse{}

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

	rErr := sw.h.Get(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_GetResponse, res)
	}
}
func (sw *sampleWrapper) echoTunnelWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &EchoRequest{}
	res := &EchoResponse{}

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

	rErr := sw.h.EchoTunnel(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}
func (sw *sampleWrapper) echoInternalWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &EchoRequest{}
	res := &EchoResponse{}

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

	rErr := sw.h.EchoInternal(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
	}
}
func (sw *sampleWrapper) echoDelayWrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := &EchoRequest{}
	res := &EchoResponse{}

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

	rErr := sw.h.EchoDelay(ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
		ctx.PushMessage(C_EchoResponse, res)
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
			SetConstructor(C_SampleEcho).
			SetServiceName("Sample").
			SetMethodName("Echo").
			SetHandler(handlerFunc(C_SampleEcho)...).
			Append(sw.echoWrapper),
	)
	e.SetRestProxy(
		"get", "/echo/:replica_set",
		edge.NewRestProxy(sw.echoRestClient, sw.echoRestServer),
	)
	e.SetHandler(
		edge.NewHandlerOptions().
			SetConstructor(C_SampleSet).
			SetServiceName("Sample").
			SetMethodName("Set").
			SetHandler(handlerFunc(C_SampleSet)...).
			Append(sw.setWrapper),
	)
	e.SetRestProxy(
		"post", "/set",
		edge.NewRestProxy(sw.setRestClient, sw.setRestServer),
	)
	e.SetHandler(
		edge.NewHandlerOptions().
			SetConstructor(C_SampleGet).
			SetServiceName("Sample").
			SetMethodName("Get").
			SetHandler(handlerFunc(C_SampleGet)...).
			Append(sw.getWrapper),
	)
	e.SetRestProxy(
		"get", "/req/:Key/something",
		edge.NewRestProxy(sw.getRestClient, sw.getRestServer),
	)
	e.SetHandler(
		edge.NewHandlerOptions().
			SetConstructor(C_SampleEchoTunnel).
			SetServiceName("Sample").
			SetMethodName("EchoTunnel").
			SetHandler(handlerFunc(C_SampleEchoTunnel)...).
			Append(sw.echoTunnelWrapper),
	)
	e.SetRestProxy(
		"get", "/echo_tunnel/:X/:YY",
		edge.NewRestProxy(sw.echoTunnelRestClient, sw.echoTunnelRestServer),
	)
	e.SetHandler(
		edge.NewHandlerOptions().
			SetConstructor(C_SampleEchoInternal).
			SetServiceName("Sample").
			SetMethodName("EchoInternal").
			SetHandler(handlerFunc(C_SampleEchoInternal)...).
			Append(sw.echoInternalWrapper).TunnelOnly(),
	)
	e.SetHandler(
		edge.NewHandlerOptions().
			SetConstructor(C_SampleEchoDelay).
			SetServiceName("Sample").
			SetMethodName("EchoDelay").
			SetHandler(handlerFunc(C_SampleEchoDelay)...).
			Append(sw.echoDelayWrapper),
	)
}

func TunnelRequestSampleEcho(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *EchoRequest, res *EchoResponse,
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleEcho, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
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
		return errors.ErrUnexpectedTunnelResponse
	}
}
func TunnelRequestSampleSet(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *SetRequest, res *SetResponse,
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleSet, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_SetResponse:
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
func TunnelRequestSampleGet(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *GetRequest, res *GetResponse,
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleGet, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case C_GetResponse:
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
func TunnelRequestSampleEchoTunnel(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *EchoRequest, res *EchoResponse,
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleEchoTunnel, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
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
		return errors.ErrUnexpectedTunnelResponse
	}
}
func TunnelRequestSampleEchoInternal(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *EchoRequest, res *EchoResponse,
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleEchoInternal, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
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
		return errors.ErrUnexpectedTunnelResponse
	}
}
func TunnelRequestSampleEchoDelay(
	ctx *edge.RequestCtx, replicaSet uint64,
	req *EchoRequest, res *EchoResponse,
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_SampleEchoDelay, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
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
		return errors.ErrUnexpectedTunnelResponse
	}
}

func (sw *sampleWrapper) echoRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := &EchoRequest{}
	if len(conn.Body()) > 0 {
		err := req.UnmarshalJSON(conn.Body())
		if err != nil {
			return err
		}
	}
	req.ReplicaSet = tools.StrToUInt64(tools.GetString(conn.Get("replica_set"), "0"))

	ctx.Fill(conn.ConnID(), C_SampleEcho, req)
	return nil
}
func (sw *sampleWrapper) echoRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_EchoResponse:
			x := &EchoResponse{}
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

func (sw *sampleWrapper) setRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := &SetRequest{}
	err := req.Unmarshal(conn.Body())
	if err != nil {
		return err
	}

	ctx.Fill(conn.ConnID(), C_SampleSet, req)
	return nil
}
func (sw *sampleWrapper) setRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_SetResponse:
			x := &SetResponse{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			b, err = x.Marshal()
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

func (sw *sampleWrapper) getRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := &GetRequest{}
	req.Key = tools.S2B(tools.GetString(conn.Get("Key"), ""))

	ctx.Fill(conn.ConnID(), C_SampleGet, req)
	return nil
}
func (sw *sampleWrapper) getRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_GetResponse:
			x := &GetResponse{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			b, err = x.Marshal()
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

func (sw *sampleWrapper) echoTunnelRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := &EchoRequest{}
	err := req.Unmarshal(conn.Body())
	if err != nil {
		return err
	}
	req.Int = tools.StrToInt64(tools.GetString(conn.Get("X"), "0"))
	req.Timestamp = tools.StrToInt64(tools.GetString(conn.Get("YY"), "0"))

	ctx.Fill(conn.ConnID(), C_SampleEchoTunnel, req)
	return nil
}
func (sw *sampleWrapper) echoTunnelRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_EchoResponse:
			x := &EchoResponse{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			b, err = x.Marshal()
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

func (sw *sampleWrapper) echoInternalRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := &EchoRequest{}

	ctx.Fill(conn.ConnID(), C_SampleEchoInternal, req)
	return nil
}
func (sw *sampleWrapper) echoInternalRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_EchoResponse:
			x := &EchoResponse{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			b, err = x.Marshal()
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

func (sw *sampleWrapper) echoDelayRestClient(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	req := &EchoRequest{}

	ctx.Fill(conn.ConnID(), C_SampleEchoDelay, req)
	return nil
}
func (sw *sampleWrapper) echoDelayRestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case C_EchoResponse:
			x := &EchoResponse{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			b, err = x.Marshal()
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
	Echo(ctx context.Context, req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error)
	Set(ctx context.Context, req *SetRequest, kvs ...*rony.KeyValue) (*SetResponse, error)
	Get(ctx context.Context, req *GetRequest, kvs ...*rony.KeyValue) (*GetResponse, error)
	EchoTunnel(ctx context.Context, req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error)
	EchoDelay(ctx context.Context, req *EchoRequest, kvs ...*rony.KeyValue) (*EchoResponse, error)
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
func (c *SampleClient) Echo(
	ctx context.Context, req *EchoRequest, kvs ...*rony.KeyValue,
) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleEcho, req, kvs...)
	err := c.c.Send(ctx, out, in)
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
		_ = x.Unmarshal(in.Message)
		return nil, x
	default:
		return nil, fmt.Errorf("unknown message :%d", in.GetConstructor())
	}
}
func (c *SampleClient) Set(
	ctx context.Context, req *SetRequest, kvs ...*rony.KeyValue,
) (*SetResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleSet, req, kvs...)
	err := c.c.Send(ctx, out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_SetResponse:
		x := &SetResponse{}
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
func (c *SampleClient) Get(
	ctx context.Context, req *GetRequest, kvs ...*rony.KeyValue,
) (*GetResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleGet, req, kvs...)
	err := c.c.Send(ctx, out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	case C_GetResponse:
		x := &GetResponse{}
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
func (c *SampleClient) EchoTunnel(
	ctx context.Context, req *EchoRequest, kvs ...*rony.KeyValue,
) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleEchoTunnel, req, kvs...)
	err := c.c.Send(ctx, out, in)
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
		_ = x.Unmarshal(in.Message)
		return nil, x
	default:
		return nil, fmt.Errorf("unknown message :%d", in.GetConstructor())
	}
}
func (c *SampleClient) EchoDelay(
	ctx context.Context, req *EchoRequest, kvs ...*rony.KeyValue,
) (*EchoResponse, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_SampleEchoDelay, req, kvs...)
	err := c.c.Send(ctx, out, in)
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
		_ = x.Unmarshal(in.Message)
		return nil, x
	default:
		return nil, fmt.Errorf("unknown message :%d", in.GetConstructor())
	}
}

func prepareSampleCommand(cmd *cobra.Command, c edgec.Client) (*SampleClient, error) {
	// Bind current flags to the registered flags in config package
	err := config.BindCmdFlags(cmd)
	if err != nil {
		return nil, err
	}

	return NewSampleClient("Sample", c), nil
}

var genSampleEchoCmd = func(h ISampleCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "echo",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareSampleCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Echo(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.Int64Flag("int", tools.StrToInt64(""), ""),
		config.Int64Flag("timestamp", tools.StrToInt64(""), ""),
		config.Uint64Flag("replicaSet", tools.StrToUInt64(""), ""),
	)
	return cmd
}
var genSampleSetCmd = func(h ISampleCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "set",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareSampleCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Set(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.StringFlag("key", "", ""),
		config.StringFlag("value", "", ""),
	)
	return cmd
}
var genSampleGetCmd = func(h ISampleCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "get",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareSampleCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.Get(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.StringFlag("key", "somekey", "enter some random key"),
	)
	return cmd
}
var genSampleEchoTunnelCmd = func(h ISampleCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "echo-tunnel",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareSampleCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.EchoTunnel(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.Int64Flag("int", tools.StrToInt64(""), ""),
		config.Int64Flag("timestamp", tools.StrToInt64(""), ""),
		config.Uint64Flag("replicaSet", tools.StrToUInt64(""), ""),
	)
	return cmd
}
var genSampleEchoDelayCmd = func(h ISampleCli, c edgec.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use: "echo-delay",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli, err := prepareSampleCommand(cmd, c)
			if err != nil {
				return err
			}
			return h.EchoDelay(cli, cmd, args)
		},
	}
	config.SetFlags(cmd,
		config.Int64Flag("int", tools.StrToInt64(""), ""),
		config.Int64Flag("timestamp", tools.StrToInt64(""), ""),
		config.Uint64Flag("replicaSet", tools.StrToUInt64(""), ""),
	)
	return cmd
}

type ISampleCli interface {
	Echo(cli *SampleClient, cmd *cobra.Command, args []string) error
	Set(cli *SampleClient, cmd *cobra.Command, args []string) error
	Get(cli *SampleClient, cmd *cobra.Command, args []string) error
	EchoTunnel(cli *SampleClient, cmd *cobra.Command, args []string) error
	EchoDelay(cli *SampleClient, cmd *cobra.Command, args []string) error
}

func RegisterSampleCli(h ISampleCli, c edgec.Client, rootCmd *cobra.Command) {
	subCommand := &cobra.Command{
		Use: "Sample",
	}
	subCommand.AddCommand(
		genSampleEchoCmd(h, c),
		genSampleSetCmd(h, c),
		genSampleGetCmd(h, c),
		genSampleEchoTunnelCmd(h, c),
		genSampleEchoDelayCmd(h, c),
	)

	rootCmd.AddCommand(subCommand)
}

var _ = bytes.MinRead
