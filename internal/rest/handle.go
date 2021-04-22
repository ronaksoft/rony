package rest

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/pools"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2021 - Apr - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Handle struct {
	ctx        *Context
	onRequest  func(ctx *Context)
	onResponse func(envelope *rony.MessageEnvelope) (*pools.ByteBuffer, map[string]string)
	inputBuf   *pools.ByteBuffer
}

func (h *Handle) OnRequest(conn rony.Conn, ctx *gateway.RequestCtx) []byte {
	h.ctx = &Context{
		reqCtx: ctx,
		conn:   conn,
	}

	h.ctx.me = rony.PoolMessageEnvelope.Get()
	h.onRequest(h.ctx)
	mo := proto.MarshalOptions{UseCachedSize: true}
	h.inputBuf = pools.Buffer.GetCap(mo.Size(h.ctx.me))
	rony.PoolMessageEnvelope.Put(h.ctx.me)
	return *h.inputBuf.Bytes()
}

func (h *Handle) OnResponse(data []byte) (*pools.ByteBuffer, map[string]string) {
	me := rony.PoolMessageEnvelope.Get()
	_ = me.Unmarshal(data)
	out, m := h.onResponse(me)
	rony.PoolMessageEnvelope.Put(me)
	pools.Buffer.Put(h.inputBuf)
	return out, m
}
