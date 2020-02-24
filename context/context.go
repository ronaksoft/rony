package context

import (
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"git.ronaksoftware.com/ronak/rony/msg"
	"github.com/gobwas/pool/pbytes"
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

// Context Values
const (
	CtxAuthKey   = "AUTH_KEY"
	CtxServerSeq = "S_SEQ"
	CtxClientSeq = "C_SEQ"
	CtxUser      = "USER"
	CtxStreamID  = "SID"
	CtxTemp      = "TEMP"
)

type Context struct {
	ConnID      uint64
	AuthID      int64
	QuickReturn bool
	NextChan    chan struct{}
	Stop        bool
	Blocking    bool

	// internals
	mtx         sync.RWMutex
	kv          map[uint32]interface{}
	MessageChan chan *messageDispatch
	UpdateChan  chan *updateDispatch
}

func New() *Context {
	return &Context{
		NextChan:    make(chan struct{}, 1),
		kv:          make(map[uint32]interface{}, 3),
		MessageChan: make(chan *messageDispatch, 3),
		UpdateChan:  make(chan *updateDispatch, 100),
	}
}

func (ctx *Context) Return() {
	if ctx.QuickReturn {
		ctx.NextChan <- struct{}{}
	}
}

func (ctx *Context) StopExecution() {
	ctx.Stop = true
	close(ctx.MessageChan)
	close(ctx.UpdateChan)
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
	ctx.MessageChan = make(chan *messageDispatch, 3)
	ctx.UpdateChan = make(chan *updateDispatch, 100)
	for k := range ctx.kv {
		delete(ctx.kv, k)
	}
}

type messageDispatch struct {
	AuthID   int64
	Envelope *msg.MessageEnvelope
}

func (ctx *Context) PushMessage(authID int64, requestID uint64, constructor int64, proto ProtoBufferMessage) {
	envelope := pools.AcquireMessageEnvelope()
	envelope.RequestID = requestID
	envelope.Constructor = constructor
	pbytes.Put(envelope.Message)
	envelope.Message = pbytes.GetLen(proto.Size())
	_, _ = proto.MarshalTo(envelope.Message)
	ctx.MessageChan <- &messageDispatch{
		AuthID:   authID,
		Envelope: envelope,
	}
}

func (ctx *Context) PushError(requestID uint64, code, item string) {
	ctx.PushMessage(ctx.AuthID, requestID, msg.C_Error, &msg.Error{
		Code:  code,
		Items: item,
	})
}

type updateDispatch struct {
	AuthID   int64
	Envelope *msg.UpdateEnvelope
}

func (ctx *Context) PushUpdate(authID int64, updateID int64, constructor int64, proto ProtoBufferMessage) {
	envelope := pools.AcquireUpdateEnvelope()
	envelope.Timestamp = time.Now().Unix()
	envelope.UpdateID = updateID
	if updateID != 0 {
		envelope.UCount = 1
	}
	envelope.Constructor = constructor
	pbytes.Put(envelope.Update)
	envelope.Update = pbytes.GetLen(proto.Size())
	_, _ = proto.MarshalTo(envelope.Update)
	ctx.UpdateChan <- &updateDispatch{
		AuthID:   authID,
		Envelope: envelope,
	}
}

// ProtoBufferMessage
type ProtoBufferMessage interface {
	proto.Marshaler
	proto.Sizer
	proto.Unmarshaler
	MarshalTo([]byte) (int, error)
}
