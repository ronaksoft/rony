package edge

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/gateway"
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"hash/crc32"
	"reflect"
	"sync"
)

/*
   Creation Time: 2019 - Jun - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	_ byte = iota
	gatewayMessage
	clusterMessage
)

// DispatchCtx
type DispatchCtx struct {
	streamID int64
	serverID []byte
	conn     gateway.Conn
	req      *rony.MessageEnvelope
	edge     *Server
	kind     byte
}

func newDispatchCtx(edge *Server) *DispatchCtx {
	return &DispatchCtx{
		edge: edge,
		req:  &rony.MessageEnvelope{},
	}
}

func (ctx *DispatchCtx) Debug() {
	fmt.Println("###")
	t := reflect.Indirect(reflect.ValueOf(ctx))
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Type().Field(i).Name, t.Type().Field(i).Offset, t.Type().Field(i).Type.Size())
	}
}

func (ctx *DispatchCtx) Conn() gateway.Conn {
	return ctx.conn
}

func (ctx *DispatchCtx) StreamID() int64 {
	return ctx.streamID
}

func (ctx *DispatchCtx) FillEnvelope(requestID uint64, constructor int64, payload []byte) {
	ctx.req.RequestID = requestID
	ctx.req.Constructor = constructor
	ctx.req.Message = append(ctx.req.Message[:0], payload...)
}

func (ctx *DispatchCtx) Get(key string) interface{} {
	return ctx.Conn().Get(key)
}

func (ctx *DispatchCtx) Set(key string, val interface{}) {
	ctx.Conn().Set(key, val)
}

func (ctx *DispatchCtx) UnmarshalEnvelope(data []byte) error {
	uo := proto.UnmarshalOptions{
		Merge: true,
	}
	return uo.Unmarshal(data, ctx.req)
}

// RequestCtx
type RequestCtx struct {
	dispatchCtx *DispatchCtx
	quickReturn bool
	nextChan    chan struct{}
	stop        bool
	mtx         sync.RWMutex
	kv          map[uint32]interface{}
}

func newRequestCtx() *RequestCtx {
	return &RequestCtx{
		nextChan: make(chan struct{}, 1),
		kv:       make(map[uint32]interface{}, 3),
	}
}

func (ctx *RequestCtx) reset() {
	for k := range ctx.kv {
		delete(ctx.kv, k)
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

func (ctx *RequestCtx) Conn() gateway.Conn {
	return ctx.dispatchCtx.Conn()
}

func (ctx *RequestCtx) ReqID() uint64 {
	return ctx.dispatchCtx.req.GetRequestID()
}

func (ctx *RequestCtx) StopExecution() {
	ctx.stop = true
}

func (ctx *RequestCtx) Stopped() bool {
	return ctx.stop
}

func (ctx *RequestCtx) Set(key string, v interface{}) {
	ctx.mtx.Lock()
	ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))] = v
	ctx.mtx.Unlock()
}

func (ctx *RequestCtx) Get(key string) interface{} {
	ctx.mtx.RLock()
	v := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))]
	ctx.mtx.RUnlock()
	return v
}

func (ctx *RequestCtx) GetBytes(key string, defaultValue []byte) []byte {
	v, ok := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))].([]byte)
	if ok {
		return v
	}
	return defaultValue
}

func (ctx *RequestCtx) GetString(key string, defaultValue string) string {
	v := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))]
	switch x := v.(type) {
	case []byte:
		return tools.ByteToStr(x)
	case string:
		return x
	default:
		return defaultValue
	}
}

func (ctx *RequestCtx) GetInt64(key string, defaultValue int64) int64 {
	v, ok := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))].(int64)
	if ok {
		return v
	}
	return defaultValue
}

func (ctx *RequestCtx) GetBool(key string) bool {
	v, ok := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))].(bool)
	if ok {
		return v
	}
	return false
}

func (ctx *RequestCtx) PushMessage(constructor int64, proto proto.Message) {
	ctx.PushCustomMessage(ctx.ReqID(), constructor, proto)
}

func (ctx *RequestCtx) PushCustomMessage(requestID uint64, constructor int64, proto proto.Message, kvs ...*rony.KeyValue) {
	envelope := acquireMessageEnvelope()
	envelope.Fill(requestID, constructor, proto)
	ctx.dispatchCtx.edge.dispatcher.OnMessage(ctx.dispatchCtx, envelope, kvs...)
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

func (ctx *RequestCtx) PushClusterMessage(serverID string, requestID uint64, constructor int64, proto proto.Message, kvs ...*rony.KeyValue) {
	envelope := acquireMessageEnvelope()
	envelope.Fill(requestID, constructor, proto)
	err := ctx.dispatchCtx.edge.ClusterSend(tools.StrToByte(serverID), envelope, kvs...)
	if err != nil {
		log.Error("ClusterMessage Error",
			zap.Bool("GatewayRequest", ctx.dispatchCtx.conn != nil),
			zap.Uint64("RequestID", envelope.GetRequestID()),
			zap.String("C", ctx.dispatchCtx.edge.getConstructorName(envelope.GetConstructor())),
			zap.Error(err),
		)
	}
	releaseMessageEnvelope(envelope)
}
