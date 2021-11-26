package edge

import (
	"context"
	"sync"
	"time"

	"github.com/ronaksoft/rony/registry"
	"go.opentelemetry.io/otel/attribute"

	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/log"
	"github.com/ronaksoft/rony/tools"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

// RequestCtx holds the context of an RPC handler
type RequestCtx struct {
	dispatchCtx *DispatchCtx
	edge        *Server
	reqID       uint64
	nextChan    chan struct{}
	quickReturn bool
	stop        bool
	err         *rony.Error
	// cfg holds all cancel functions which will be called at the
	// end of the lifecycle of the RequestCtx
	ctx context.Context
	cf  func()
}

func newRequestCtx() *RequestCtx {
	return &RequestCtx{
		nextChan: make(chan struct{}, 1),
	}
}

// NextCanGo is useful only if the main request is rony.MessageContainer with SyncRun flag set to TRUE.
// Then by calling this function it lets the next request (if any) executed concurrently.
func (ctx *RequestCtx) NextCanGo() {
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
	ctx.StopExecution()
}

// PushRedirectSession redirects the client to another server with replicaSet for the rest
// of the request for this session. The interpretation of session is defined by the application
// and does not have any side effect on internal states of the Edge server.
// This function call StopExecution internally and further handlers would not be called.
func (ctx *RequestCtx) PushRedirectSession(replicaSet uint64) {
	ctx.pushRedirect(rony.RedirectReason_ReplicaSetSession, replicaSet)
}

// PushRedirectRequest redirects the client to another server with replicaSet for only
// this request. This function call StopExecution internally and further handlers would
// not be called.
func (ctx *RequestCtx) PushRedirectRequest(replicaSet uint64) {
	ctx.pushRedirect(rony.RedirectReason_ReplicaSetRequest, replicaSet)
}

// PushMessage is a wrapper func for PushCustomMessage.
func (ctx *RequestCtx) PushMessage(constructor uint64, proto proto.Message) {
	ctx.PushCustomMessage(ctx.ReqID(), constructor, proto)
}

// PushCustomMessage would do different actions based on the Conn. If connection is persistent (i.e. Websocket)
// then Rony sends the message down to the wire real time. If the connection is not persistent
// (i.e. HTTP) then Rony push the message into a temporary buffer and at the end of the life-time
// of the RequestCtx pushes the message as response.
// If there was multiple message in the buffer then Rony creates a MessageContainer and wrap
// all those messages in that container. Client needs to unwrap MessageContainer when ever they
// see it in the response.
func (ctx *RequestCtx) PushCustomMessage(
	requestID uint64, constructor uint64, message proto.Message, kvs ...*rony.KeyValue,
) {
	envelope := rony.PoolMessageEnvelope.Get()
	envelope.Fill(requestID, constructor, message, kvs...)

	switch ctx.Kind() {
	case TunnelMessage:
		ctx.dispatchCtx.BufferPush(envelope)
	case GatewayMessage:
		if ctx.edge.tracer != nil {
			trace.SpanFromContext(ctx.ctx).
				AddEvent("message",
					trace.WithAttributes(
						semconv.MessageTypeSent,
						semconv.MessageIDKey.Int64(int64(requestID)),
						attribute.String("rony.constructor", registry.C(constructor)),
					),
				)
		}

		if ctx.Conn().Persistent() {
			_ = ctx.edge.dispatcher.Encode(ctx.dispatchCtx.conn, ctx.dispatchCtx.streamID, envelope)
			rony.PoolMessageEnvelope.Put(envelope)
		} else {
			ctx.dispatchCtx.BufferPush(envelope)
		}
	}
}

func (ctx *RequestCtx) PushError(err *rony.Error) {
	ctx.err = err.Clone()
	ctx.PushMessage(rony.C_Error, err)
	ctx.stop = true
}

func (ctx *RequestCtx) PushCustomError(code, item string, desc string) {
	ctx.PushMessage(
		rony.C_Error,
		&rony.Error{
			Code:        code,
			Items:       item,
			Description: desc,
		},
	)
	ctx.stop = true
}

func (ctx *RequestCtx) Error() *rony.Error {
	return ctx.err
}

func (ctx *RequestCtx) Cluster() rony.Cluster {
	return ctx.edge.cluster
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

func (ctx *RequestCtx) TryTunnelRequest(
	attempts int, retryWait time.Duration, replicaSet uint64,
	req, res *rony.MessageEnvelope,
) error {
	return ctx.edge.TryTunnelRequest(attempts, retryWait, replicaSet, req, res)
}

func (ctx *RequestCtx) TunnelRequest(replicaSet uint64, req, res *rony.MessageEnvelope) error {
	return ctx.edge.TryTunnelRequest(1, 0, replicaSet, req, res)
}

// Log returns a logger
func (ctx *RequestCtx) Log() log.Logger {
	return ctx.edge.logger
}

// Span returns a tracer span. Don't End the span, since it will be
// closed automatically at the end of RequestCtx lifecycle.
func (ctx *RequestCtx) Span() trace.Span {
	return trace.SpanFromContext(ctx.ctx)
}

func (ctx *RequestCtx) ReplicaSet() uint64 {
	if ctx.edge.cluster == nil {
		return 0
	}

	return ctx.edge.cluster.ReplicaSet()
}

func (ctx *RequestCtx) Router() rony.Router {
	return ctx.edge.router
}

func (ctx *RequestCtx) Context() context.Context {
	return ctx.ctx
}

var requestCtxPool = sync.Pool{}

func acquireRequestCtx(dispatchCtx *DispatchCtx, quickReturn bool) *RequestCtx {
	var ctx *RequestCtx
	if v := requestCtxPool.Get(); v == nil {
		ctx = newRequestCtx()
	} else {
		ctx = v.(*RequestCtx)
	}
	ctx.stop = false
	ctx.quickReturn = quickReturn
	ctx.dispatchCtx = dispatchCtx
	ctx.edge = dispatchCtx.edge
	ctx.err = nil
	ctx.ctx, ctx.cf = context.WithCancel(dispatchCtx.ctx)

	return ctx
}

func releaseRequestCtx(ctx *RequestCtx) {
	// Just to make sure channel is empty, or empty it if not
	select {
	case <-ctx.nextChan:
	default:
	}

	// call cancel func
	ctx.cf()

	// reset request id
	ctx.reqID = 0

	// Put back into the pool
	requestCtxPool.Put(ctx)
}
