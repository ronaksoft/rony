package rony

import (
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gogo/protobuf/proto"
	"hash/crc32"
	"sync"
	"time"
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
	CtxServerID = "SID"
)

type Context struct {
	ConnID      uint64
	AuthID      int64
	QuickReturn bool
	NextChan    chan struct{}
	Stop        bool

	// internals
	mtx         sync.RWMutex
	kv          map[uint32]interface{}
	CarrierChan chan *carrier
}

func New() *Context {
	return &Context{
		NextChan:    make(chan struct{}, 1),
		kv:          make(map[uint32]interface{}, 3),
		CarrierChan: make(chan *carrier, 10),
	}
}

func (ctx *Context) Return() {
	if ctx.QuickReturn {
		ctx.NextChan <- struct{}{}
	}
}

func (ctx *Context) StopExecution() {
	ctx.Stop = true
	close(ctx.CarrierChan)
}

func (ctx *Context) Set(key string, v interface{}) {
	ctx.mtx.Lock()
	ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))] = v
	ctx.mtx.Unlock()
}

func (ctx *Context) Get(key string) interface{} {
	ctx.mtx.RLock()
	v := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))]
	ctx.mtx.RUnlock()
	return v
}

func (ctx *Context) GetBytes(key string, defaultValue []byte) []byte {
	v, ok := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))].([]byte)
	if ok {
		return v
	}
	return defaultValue
}

func (ctx *Context) GetInt64(key string, defaultValue int64) int64 {
	v, ok := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))].(int64)
	if ok {
		return v
	}
	return defaultValue
}

func (ctx *Context) GetBool(key string) bool {
	v, ok := ctx.kv[crc32.ChecksumIEEE(tools.StrToByte(key))].(bool)
	if ok {
		return v
	}
	return false
}

func (ctx *Context) Clear() {
	ctx.CarrierChan = make(chan *carrier, 10)
	for k := range ctx.kv {
		delete(ctx.kv, k)
	}
}

func (ctx *Context) PushMessage(authID int64, requestID uint64, constructor int64, proto ProtoBufferMessage) {
	envelope := acquireMessageEnvelope()
	envelope.RequestID = requestID
	envelope.Constructor = constructor
	protoSize := proto.Size()
	if protoSize > cap(envelope.Message) {
		pools.Bytes.Put(envelope.Message)
		envelope.Message = pools.Bytes.GetLen(proto.Size())
	} else {
		envelope.Message = envelope.Message[:protoSize]
	}
	_, _ = proto.MarshalTo(envelope.Message)

	ctx.CarrierChan <- acquireMessageCarrier(authID, envelope)
}

func (ctx *Context) PushError(requestID uint64, code, item string) {
	ctx.PushMessage(ctx.AuthID, requestID, C_Error, &Error{
		Code:  code,
		Items: item,
	})
}

func (ctx *Context) PushClusterMessage(serverID string, authID int64, requestID uint64, constructor int64, proto ProtoBufferMessage) {
	envelope := acquireMessageEnvelope()
	envelope.RequestID = requestID
	envelope.Constructor = constructor
	protoSize := proto.Size()
	if protoSize > cap(envelope.Message) {
		pools.Bytes.Put(envelope.Message)
		envelope.Message = pools.Bytes.GetLen(proto.Size())
	} else {
		envelope.Message = envelope.Message[:protoSize]
	}
	_, _ = proto.MarshalTo(envelope.Message)
	ctx.CarrierChan <- acquireClusterMessageCarrier(authID, serverID, envelope)
}

func (ctx *Context) PushUpdate(authID int64, updateID int64, constructor int64, proto ProtoBufferMessage) {
	envelope := acquireUpdateEnvelope()
	envelope.Timestamp = time.Now().Unix()
	envelope.UpdateID = updateID
	if updateID != 0 {
		envelope.UCount = 1
	}
	envelope.Constructor = constructor
	protoSize := proto.Size()
	if protoSize > cap(envelope.Update) {
		pools.Bytes.Put(envelope.Update)
		envelope.Update = pools.Bytes.GetLen(proto.Size())
	} else {
		envelope.Update = envelope.Update[:protoSize]
	}
	_, _ = proto.MarshalTo(envelope.Update)
	ctx.CarrierChan <- acquireUpdateCarrier(authID, envelope)
}

const (
	_ byte = iota
	carrierMessage
	carrierCluster
	carrierUpdate
)

type carrier struct {
	kind            byte
	AuthID          int64
	ServerID        []byte
	MessageEnvelope *MessageEnvelope
	UpdateEnvelope  *UpdateEnvelope
}

// ProtoBufferMessage
type ProtoBufferMessage interface {
	proto.Marshaler
	proto.Sizer
	proto.Unmarshaler
	MarshalTo([]byte) (int, error)
}
