package edge

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/proto"
	"reflect"
	"sync"
	"time"
)

/*
   Creation Time: 2019 - Jun - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type MessageKind byte

const (
	_ MessageKind = iota
	GatewayMessage
	TunnelMessage
)

var (
	messageKindNames = map[MessageKind]string{
		GatewayMessage: "GatewayMessage",
		TunnelMessage:  "TunnelMessage",
	}
)

func (c MessageKind) String() string {
	return messageKindNames[c]
}

// DispatchCtx holds the context of the dispatcher's request. Each DispatchCtx could holds one or many RequestCtx.
// DispatchCtx lives until the last of its RequestCtx children.
type DispatchCtx struct {
	edge      *Server
	streamID  int64
	serverID  []byte
	conn      rony.Conn
	req       *rony.MessageEnvelope
	reqFilled bool
	kind      MessageKind
	buf       *tools.LinkedList
	// KeyValue Store Parameters
	mtx sync.RWMutex
	kv  map[string]interface{}
}

func newDispatchCtx(edge *Server) *DispatchCtx {
	return &DispatchCtx{
		edge: edge,
		req:  &rony.MessageEnvelope{},
		kv:   make(map[string]interface{}, 3),
		buf:  tools.NewLinkedList(),
	}
}

func (ctx *DispatchCtx) reset() {
	ctx.reqFilled = false
	for k := range ctx.kv {
		delete(ctx.kv, k)
	}
	ctx.buf.Reset()
}

func (ctx *DispatchCtx) ServerID() string {
	return string(ctx.serverID)
}

func (ctx *DispatchCtx) Debug() {
	fmt.Println("###")
	t := reflect.Indirect(reflect.ValueOf(ctx))
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Type().Field(i).Name, t.Type().Field(i).Offset, t.Type().Field(i).Type.Size())
	}
}

func (ctx *DispatchCtx) Conn() rony.Conn {
	return ctx.conn
}

func (ctx *DispatchCtx) StreamID() int64 {
	return ctx.streamID
}

func (ctx *DispatchCtx) CopyEnvelope(e *rony.MessageEnvelope) {
	if ctx.reqFilled {
		panic("BUG!!! request has been already filled")
	}
	ctx.reqFilled = true
	e.DeepCopy(ctx.req)
}

func (ctx *DispatchCtx) FillEnvelope(requestID uint64, constructor int64, p proto.Message, kv ...*rony.KeyValue) {
	if ctx.reqFilled {
		panic("BUG!!! request has been already filled")
	}
	ctx.reqFilled = true
	ctx.req.Fill(requestID, constructor, p, kv...)
}

func (ctx *DispatchCtx) Set(key string, v interface{}) {
	ctx.mtx.Lock()
	ctx.kv[key] = v
	ctx.mtx.Unlock()
}

func (ctx *DispatchCtx) Get(key string) interface{} {
	ctx.mtx.RLock()
	v := ctx.kv[key]
	ctx.mtx.RUnlock()
	return v
}

func (ctx *DispatchCtx) GetBytes(key string, defaultValue []byte) []byte {
	v, ok := ctx.Get(key).([]byte)
	if ok {
		return v
	}
	return defaultValue
}

func (ctx *DispatchCtx) GetString(key string, defaultValue string) string {
	v := ctx.Get(key)
	switch x := v.(type) {
	case []byte:
		return tools.ByteToStr(x)
	case string:
		return x
	default:
		return defaultValue
	}
}

// Kind identifies that this dispatch context is generated from Tunnel or Gateway. This helps
// developer to apply different strategies based on the source of the incoming message
func (ctx *DispatchCtx) Kind() MessageKind {
	return ctx.kind
}

func (ctx *DispatchCtx) GetInt64(key string, defaultValue int64) int64 {
	v, ok := ctx.Get(key).(int64)
	if ok {
		return v
	}
	return defaultValue
}

func (ctx *DispatchCtx) GetBool(key string) bool {
	v, ok := ctx.Get(key).(bool)
	if ok {
		return v
	}
	return false
}

func (ctx *DispatchCtx) UnmarshalEnvelope(data []byte) error {
	return proto.Unmarshal(data, ctx.req)
}

func (ctx *DispatchCtx) BufferPush(m *rony.MessageEnvelope) {
	ctx.buf.Append(m)
}

func (ctx *DispatchCtx) BufferPop(f func(envelope *rony.MessageEnvelope)) bool {
	me, _ := ctx.buf.PickHeadData().(*rony.MessageEnvelope)
	if me == nil {
		return false
	}
	f(me)
	return true
}

func (ctx *DispatchCtx) BufferPopAll(f func(envelope *rony.MessageEnvelope)) {
	for ctx.BufferPop(f) {
	}
}

func (ctx *DispatchCtx) BufferSize() int32 {
	return ctx.buf.Size()
}

var dispatchCtxPool = sync.Pool{}

func acquireDispatchCtx(edge *Server, conn rony.Conn, streamID int64, serverID []byte, kind MessageKind) *DispatchCtx {
	var ctx *DispatchCtx
	if v := dispatchCtxPool.Get(); v == nil {
		ctx = newDispatchCtx(edge)
	} else {
		ctx = v.(*DispatchCtx)
	}
	ctx.conn = conn
	ctx.kind = kind
	ctx.streamID = streamID
	ctx.serverID = append(ctx.serverID[:0], serverID...)
	return ctx
}

func releaseDispatchCtx(ctx *DispatchCtx) {
	// Reset the Key-Value store
	ctx.reset()

	// Put back the context into the pool
	dispatchCtxPool.Put(ctx)
}

// RequestCtx holds the context of an RPC handler
type RequestCtx struct {
	dispatchCtx *DispatchCtx
	edge        *Server
	reqID       uint64
	nextChan    chan struct{}
	quickReturn bool
	stop        bool
}

func newRequestCtx(edge *Server) *RequestCtx {
	return &RequestCtx{
		nextChan: make(chan struct{}, 1),
		edge:     edge,
	}
}

func (ctx *RequestCtx) Return() {
	if ctx.quickReturn {
		ctx.nextChan <- struct{}{}
	}
}

func (ctx *RequestCtx) ConnID() uint64 {
	if ctx.dispatchCtx.Conn() != nil {
		return ctx.dispatchCtx.Conn().ConnID()
	}
	return 0
}

func (ctx *RequestCtx) Conn() rony.Conn {
	return ctx.dispatchCtx.Conn()
}

func (ctx *RequestCtx) ReqID() uint64 {
	return ctx.reqID
}

func (ctx *RequestCtx) ServerID() string {
	return string(ctx.dispatchCtx.serverID)
}

func (ctx *RequestCtx) Kind() MessageKind {
	return ctx.dispatchCtx.kind
}

func (ctx *RequestCtx) StopExecution() {
	ctx.stop = true
}

func (ctx *RequestCtx) Stopped() bool {
	return ctx.stop
}

func (ctx *RequestCtx) Store() rony.Store {
	return ctx.edge.store
}

func (ctx *RequestCtx) Set(key string, v interface{}) {
	ctx.dispatchCtx.mtx.Lock()
	ctx.dispatchCtx.kv[key] = v
	ctx.dispatchCtx.mtx.Unlock()
}

func (ctx *RequestCtx) Get(key string) interface{} {
	return ctx.dispatchCtx.Get(key)
}

func (ctx *RequestCtx) GetBytes(key string, defaultValue []byte) []byte {
	return ctx.dispatchCtx.GetBytes(key, defaultValue)
}

func (ctx *RequestCtx) GetString(key string, defaultValue string) string {
	return ctx.dispatchCtx.GetString(key, defaultValue)
}

func (ctx *RequestCtx) GetInt64(key string, defaultValue int64) int64 {
	return ctx.dispatchCtx.GetInt64(key, defaultValue)
}

func (ctx *RequestCtx) GetBool(key string) bool {
	return ctx.dispatchCtx.GetBool(key)
}

func (ctx *RequestCtx) pushRedirect(reason rony.RedirectReason, replicaSet uint64) {
	r := rony.PoolRedirect.Get()
	r.Reason = reason
	r.WaitInSec = 0

	members := ctx.Cluster().MembersByReplicaSet(replicaSet)
	for _, m := range members {
		ni := m.Proto(rony.PoolEdge.Get())
		r.Edges = append(r.Edges, ni)
	}

	ctx.PushMessage(rony.C_Redirect, r)
	rony.PoolRedirect.Put(r)
}

func (ctx *RequestCtx) PushRedirectSession(replicaSet uint64) {
	ctx.pushRedirect(rony.RedirectReason_ReplicaSetSession, replicaSet)
}

func (ctx *RequestCtx) PushRedirectRequest(replicaSet uint64) {
	ctx.pushRedirect(rony.RedirectReason_ReplicaSetRequest, replicaSet)
}

func (ctx *RequestCtx) PushMessage(constructor int64, proto proto.Message) {
	ctx.PushCustomMessage(ctx.ReqID(), constructor, proto)
}

func (ctx *RequestCtx) PushCustomMessage(requestID uint64, constructor int64, message proto.Message, kvs ...*rony.KeyValue) {
	envelope := rony.PoolMessageEnvelope.Get()
	envelope.Fill(requestID, constructor, message, kvs...)

	switch ctx.Kind() {
	case TunnelMessage:
		ctx.dispatchCtx.BufferPush(envelope.Clone())
	case GatewayMessage:
		if ctx.Conn().Persistent() {
			ctx.edge.dispatcher.OnMessage(ctx.dispatchCtx, envelope)
		} else {
			ctx.dispatchCtx.BufferPush(envelope.Clone())
		}
	}

	rony.PoolMessageEnvelope.Put(envelope)
}

func (ctx *RequestCtx) PushError(err *rony.Error) {
	ctx.PushMessage(rony.C_Error, err)
	ctx.stop = true
}

func (ctx *RequestCtx) PushCustomError(code, item string, desc string) {
	ctx.PushMessage(rony.C_Error, &rony.Error{
		Code:        code,
		Items:       item,
		Description: desc,
	})
	ctx.stop = true
}

func (ctx *RequestCtx) Cluster() rony.Cluster {
	return ctx.dispatchCtx.edge.cluster
}

func (ctx *RequestCtx) ClusterEdges(replicaSet uint64, edges *rony.Edges) (*rony.Edges, error) {
	req := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(req)
	res := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(res)
	req.Fill(tools.RandomUint64(0), rony.C_GetAllNodes, &rony.GetAllNodes{})
	err := ctx.TunnelRequest(replicaSet, req, res)
	if err != nil {
		return nil, err
	}
	switch res.Constructor {
	case rony.C_Edges:
		if edges == nil {
			edges = &rony.Edges{}
		}
		_ = edges.Unmarshal(res.GetMessage())
		return edges, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(res.GetMessage())
		return nil, x
	default:
		return nil, errors.ErrUnexpectedTunnelResponse
	}
}

func (ctx *RequestCtx) TryTunnelRequest(attempts int, retryWait time.Duration, replicaSet uint64, req, res *rony.MessageEnvelope) error {
	return ctx.edge.TryTunnelRequest(attempts, retryWait, replicaSet, req, res)
}

func (ctx *RequestCtx) TunnelRequest(replicaSet uint64, req, res *rony.MessageEnvelope) error {
	return ctx.edge.TryTunnelRequest(1, 0, replicaSet, req, res)
}

func (ctx *RequestCtx) Log() log.Logger {
	return log.DefaultLogger
}

func (ctx *RequestCtx) ReplicaSet() uint64 {
	if ctx.edge.cluster == nil {
		return 0
	}
	return ctx.edge.cluster.ReplicaSet()
}

var requestCtxPool = sync.Pool{}

func acquireRequestCtx(dispatchCtx *DispatchCtx, quickReturn bool) *RequestCtx {
	var ctx *RequestCtx
	if v := requestCtxPool.Get(); v == nil {
		ctx = newRequestCtx(dispatchCtx.edge)
	} else {
		ctx = v.(*RequestCtx)
	}
	ctx.stop = false
	ctx.quickReturn = quickReturn
	ctx.dispatchCtx = dispatchCtx
	return ctx
}

func releaseRequestCtx(ctx *RequestCtx) {
	// Just to make sure channel is empty, or empty it if not
	select {
	case <-ctx.nextChan:
	default:
	}

	ctx.reqID = 0

	// Put back into the pool
	requestCtxPool.Put(ctx)
}
