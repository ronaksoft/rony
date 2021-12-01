package edge

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/proto"
)

// DispatchCtx holds the context of the dispatcher's request. Each DispatchCtx could
// hold one or many RequestCtx. DispatchCtx lives until the last of its RequestCtx children.
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
	ctx context.Context
	cf  func()
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

func (ctx *DispatchCtx) FillEnvelope(e *rony.MessageEnvelope) {
	if ctx.reqFilled {
		panic("BUG!!! request has been already filled")
	}
	ctx.reqFilled = true
	e.DeepCopy(ctx.req)
}

func (ctx *DispatchCtx) Fill(
	requestID uint64, constructor uint64, p proto.Message, kv ...*rony.KeyValue,
) {
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

// Kind identifies that this dispatch context is generated from Tunnel or Gateway. This helps
// developer to apply different strategies based on the source of the incoming message
func (ctx *DispatchCtx) Kind() MessageKind {
	return ctx.kind
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
	rony.PoolMessageEnvelope.Put(me)

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

func acquireDispatchCtx(
	edge *Server, conn rony.Conn,
	streamID int64, serverID []byte, kind MessageKind,
) *DispatchCtx {
	var ctx *DispatchCtx
	if v := dispatchCtxPool.Get(); v == nil {
		ctx = newDispatchCtx(edge)
	} else {
		ctx = v.(*DispatchCtx)
		ctx.req.Reset()
	}
	ctx.conn = conn
	ctx.kind = kind
	ctx.streamID = streamID
	ctx.serverID = append(ctx.serverID[:0], serverID...)
	ctx.ctx, ctx.cf = context.WithCancel(context.TODO())

	return ctx
}

func releaseDispatchCtx(ctx *DispatchCtx) {
	// call cancel func
	ctx.cf()

	// Reset the Key-Value store
	ctx.reset()

	// Put back the context into the pool
	dispatchCtxPool.Put(ctx)
}
