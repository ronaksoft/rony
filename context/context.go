package context

import (
	"git.ronaksoftware.com/ronak/rony/tools"
	"hash"
	"hash/crc32"
	"hash/crc64"
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

var crcGen hash.Hash64

func init() {
	crcGen = crc64.New(crc64.MakeTable(crc64.ISO))
}

type Context struct {
	AuthID      int64
	UserID      int64
	ConnID      uint64
	QuickReturn bool
	NextChan    chan struct{}
	Stop        bool
	Blocking    bool

	// internals
	mtx sync.RWMutex
	kv  map[uint32]interface{}
}

func New() *Context {
	return &Context{
		NextChan: make(chan struct{}, 1),
		kv:       make(map[uint32]interface{}, 3),
	}
}

func (ctx *Context) Return() {
	if ctx.QuickReturn {
		ctx.NextChan <- struct{}{}
	}
}

func (ctx *Context) StopExecution() {
	ctx.Stop = true
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
	for k := range ctx.kv {
		delete(ctx.kv, k)
	}
}
