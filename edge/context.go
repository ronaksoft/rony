package edge

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cluster"
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

type ContextKind byte

const (
	_ ContextKind = iota
	GatewayMessage
	ReplicaMessage
	TunnelMessage
)

// Conn defines the Connection interface
type Conn interface {
	ConnID() uint64
	ClientIP() string
	SendBinary(streamID int64, data []byte) error
	Persistent() bool
	Get(key string) interface{}
	Set(key string, val interface{})
}

// DispatchCtx
type DispatchCtx struct {
	streamID          int64
	serverID          []byte
	conn              Conn
	req               *rony.MessageEnvelope
	cluster           cluster.Cluster
	gatewayDispatcher GatewayDispatcher
	tunnelDispatcher  TunnelDispatcher
	kind              ContextKind
	buf               *tools.LinkedList
	// KeyValue Store Parameters
	mtx sync.RWMutex
	kv  map[string]interface{}
}

func newDispatchCtx(edge *Server) *DispatchCtx {
	return &DispatchCtx{
		cluster:           edge.cluster,
		gatewayDispatcher: edge.gatewayDispatcher,
		tunnelDispatcher:  edge.tunnelDispatcher,
		req:               &rony.MessageEnvelope{},
		kv:                make(map[string]interface{}, 3),
		buf:               tools.NewLinkedList(),
	}
}

func (ctx *DispatchCtx) reset() {
	for k := range ctx.kv {
		delete(ctx.kv, k)
	}
	ctx.buf.Reset()
}

func (ctx *DispatchCtx) Debug() {
	fmt.Println("###")
	t := reflect.Indirect(reflect.ValueOf(ctx))
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Type().Field(i).Name, t.Type().Field(i).Offset, t.Type().Field(i).Type.Size())
	}
}

func (ctx *DispatchCtx) Conn() Conn {
	return ctx.conn
}

func (ctx *DispatchCtx) StreamID() int64 {
	return ctx.streamID
}

func (ctx *DispatchCtx) FillEnvelope(requestID uint64, constructor int64, payload []byte, auth []byte, kv ...*rony.KeyValue) {
	ctx.req.RequestID = requestID
	ctx.req.Constructor = constructor
	ctx.req.Message = append(ctx.req.Message[:0], payload...)
	ctx.req.Auth = append(ctx.req.Auth[:0], auth...)
	if cap(ctx.req.Header) >= len(kv) {
		ctx.req.Header = ctx.req.Header[:len(kv)]
	} else {
		ctx.req.Header = make([]*rony.KeyValue, len(kv))
	}
	for idx, kv := range kv {
		if ctx.req.Header[idx] == nil {
			ctx.req.Header[idx] = &rony.KeyValue{}
		}
		kv.DeepCopy(ctx.req.Header[idx])
	}
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

func (ctx *DispatchCtx) UnmarshalEnvelope(data []byte) error {
	return proto.Unmarshal(data, ctx.req)
}

func (ctx *DispatchCtx) Kind() ContextKind {
	return ctx.kind
}

func (ctx *DispatchCtx) Push(m *rony.MessageEnvelope) {
	ctx.buf.Append(m)
}

func (ctx *DispatchCtx) Pop() *rony.MessageEnvelope {
	v := ctx.buf.PickHeadData()
	if v != nil {
		return v.(*rony.MessageEnvelope)
	}
	return nil
}

// RequestCtx
type RequestCtx struct {
	dispatchCtx *DispatchCtx
	reqID       uint64
	nextChan    chan struct{}
	quickReturn bool
	stop        bool
}

func newRequestCtx() *RequestCtx {
	return &RequestCtx{
		nextChan: make(chan struct{}, 1),
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

func (ctx *RequestCtx) Conn() Conn {
	return ctx.dispatchCtx.Conn()
}

func (ctx *RequestCtx) ReqID() uint64 {
	return ctx.reqID
}

func (ctx *RequestCtx) ServerID() string {
	return string(ctx.dispatchCtx.serverID)
}

func (ctx *RequestCtx) Kind() ContextKind {
	return ctx.dispatchCtx.Kind()
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

func (ctx *RequestCtx) PushRedirectLeader() {
	ctxCluster := ctx.dispatchCtx.cluster
	if leaderID := ctxCluster.RaftLeaderID(); leaderID == "" {
		ctx.PushError(rony.ErrCodeUnavailable, rony.ErrItemRaftLeader)
	} else {
		r := rony.PoolRedirect.Get()
		r.Reason = rony.RedirectReason_ReplicaMaster
		r.WaitInSec = 0
		members := ctxCluster.RaftMembers(ctxCluster.ReplicaSet())
		nodeInfos := make([]*rony.NodeInfo, 0, len(members))
		for _, m := range members {
			ni := m.Proto(rony.PoolNodeInfo.Get())
			if ni.Leader {
				r.Leader = ni
			} else {
				r.Followers = append(r.Followers, ni)
			}
			nodeInfos = append(nodeInfos, ni)
		}
		ctx.PushMessage(rony.C_Redirect, r)
		for _, p := range nodeInfos {
			rony.PoolNodeInfo.Put(p)
		}
		rony.PoolRedirect.Put(r)
	}
}

func (ctx *RequestCtx) PushRedirectShard(shard uint32, wait time.Duration) {
	// TODO:: implement it
}

func (ctx *RequestCtx) PushMessage(constructor int64, proto proto.Message) {
	ctx.PushCustomMessage(ctx.ReqID(), constructor, proto)
}

func (ctx *RequestCtx) PushCustomMessage(requestID uint64, constructor int64, proto proto.Message, kvs ...*rony.KeyValue) {
	if ctx.dispatchCtx.kind == ReplicaMessage {
		return
	}

	envelope := acquireMessageEnvelope()
	envelope.Fill(requestID, constructor, proto)
	envelope.Header = append(envelope.Header[:0], kvs...)

	switch ctx.dispatchCtx.kind {
	case GatewayMessage:
		ctx.dispatchCtx.gatewayDispatcher.OnMessage(ctx.dispatchCtx, envelope)
	case TunnelMessage:
		ctx.dispatchCtx.tunnelDispatcher.OnMessage(ctx.dispatchCtx, envelope)
	}

	releaseMessageEnvelope(envelope)
}

func (ctx *RequestCtx) PushError(code, item string) {
	ctx.PushCustomError(code, item, "", nil, "", nil)
}

func (ctx *RequestCtx) PushCustomError(code, item string, enTxt string, enItems []string, localTxt string, localItems []string) {
	ctx.PushMessage(rony.C_Error, &rony.Error{
		Code:               code,
		Items:              item,
		TemplateItems:      enItems,
		Template:           enTxt,
		LocalTemplate:      localTxt,
		LocalTemplateItems: localItems,
	})
	ctx.stop = true
}
