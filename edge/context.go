package edge

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
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
   Copyright Ronak Software Group 2018
*/

const (
	_ byte = iota
	carrierMessage
	carrierCluster
	carrierUpdate
)

type carrier struct {
	ServerID        []byte
	AuthID          int64
	MessageEnvelope *rony.MessageEnvelope
	UpdateEnvelope  *rony.UpdateEnvelope
	kind            byte
}

const (
	_ byte = iota
	gatewayMessage
	clusterMessage
)

// DispatchCtx
type DispatchCtx struct {
	streamID int64
	authID   int64
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

func (ctx *DispatchCtx) GetAuthID() int64 {
	return ctx.authID
}

func (ctx *DispatchCtx) SetAuthID(authID int64) {
	ctx.authID = authID
}

func (ctx *DispatchCtx) FillEnvelope(requestID uint64, constructor int64, payload []byte) {
	ctx.req.RequestID = requestID
	ctx.req.Constructor = constructor
	if len(payload) > cap(ctx.req.Message) {
		pools.Bytes.Put(ctx.req.Message)
		ctx.req.Message = pools.Bytes.GetCap(len(payload))
	}
	ctx.req.Message = append(ctx.req.Message, payload...)
}

func (ctx *DispatchCtx) UnmarshalEnvelope(data []byte) error {
	return ctx.req.Unmarshal(data)
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
	if ctx.dispatchCtx.conn != nil {
		return ctx.dispatchCtx.conn.GetConnID()
	}
	return 0
}

func (ctx *RequestCtx) AuthID() int64 {
	return ctx.dispatchCtx.authID
}

func (ctx *RequestCtx) StopExecution() {
	ctx.stop = true
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

func (ctx *RequestCtx) PushMessage(authID int64, requestID uint64, constructor int64, proto ProtoBufferMessage) {
	envelope := acquireMessageEnvelope()
	envelope.RequestID = requestID
	envelope.Constructor = constructor
	protoSize := proto.Size()
	if protoSize > cap(envelope.Message) {
		pools.Bytes.Put(envelope.Message)
		envelope.Message = pools.Bytes.GetLen(protoSize)
	} else {
		envelope.Message = envelope.Message[:protoSize]
	}
	_, err := proto.MarshalToSizedBuffer(envelope.Message)
	if err != nil {
		log.Error("Error On Marshaling Message", zap.Error(err))
		return
	}

	ctx.dispatchCtx.edge.dispatcher.DispatchMessage(ctx.dispatchCtx, authID, envelope)
	releaseMessageEnvelope(envelope)
}

func (ctx *RequestCtx) PushError(requestID uint64, code, item string) {
	ctx.PushMessage(ctx.dispatchCtx.authID, requestID, rony.C_Error, &rony.Error{
		Code:  code,
		Items: item,
	})
}

func (ctx *RequestCtx) PushClusterMessage(serverID string, authID int64, requestID uint64, constructor int64, proto ProtoBufferMessage) {
	envelope := acquireMessageEnvelope()
	envelope.RequestID = requestID
	envelope.Constructor = constructor
	protoSize := proto.Size()
	if protoSize > cap(envelope.Message) {
		pools.Bytes.Put(envelope.Message)
		envelope.Message = pools.Bytes.GetLen(protoSize)
	} else {
		envelope.Message = envelope.Message[:protoSize]
	}
	_, _ = proto.MarshalToSizedBuffer(envelope.Message)

	err := ctx.dispatchCtx.edge.ClusterSend(tools.StrToByte(serverID), authID, envelope)
	if err != nil {
		log.Error("ClusterMessage Error",
			zap.Bool("GatewayRequest", ctx.dispatchCtx.conn != nil),
			zap.Int64("AuthID", authID),
			zap.Uint64("RequestID", envelope.RequestID),
			zap.String("C", ctx.dispatchCtx.edge.getConstructorName(envelope.Constructor)),
			zap.Error(err),
		)
	}
	releaseMessageEnvelope(envelope)
}

func (ctx *RequestCtx) PushUpdate(authID int64, updateID int64, constructor, ts int64, proto ProtoBufferMessage) {
	envelope := acquireUpdateEnvelope()
	envelope.Timestamp = ts
	envelope.UpdateID = updateID
	if updateID != 0 {
		envelope.UCount = 1
	}
	envelope.Constructor = constructor
	protoSize := proto.Size()
	if protoSize > cap(envelope.Update) {
		pools.Bytes.Put(envelope.Update)
		envelope.Update = pools.Bytes.GetLen(protoSize)
	} else {
		envelope.Update = envelope.Update[:protoSize]
	}
	_, err := proto.MarshalToSizedBuffer(envelope.Update)
	if err != nil {
		log.Error("Error On Marshaling Update", zap.Error(err))
		return
	}

	ctx.dispatchCtx.edge.dispatcher.DispatchUpdate(ctx.dispatchCtx, authID, envelope)
	releaseUpdateEnvelope(envelope)
}

// ProtoBufferMessage
type ProtoBufferMessage interface {
	proto.Marshaler
	proto.Sizer
	proto.Unmarshaler
	MarshalTo([]byte) (int, error)
	MarshalToSizedBuffer([]byte) (int, error)
}
