package rest

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - Apr - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type restHandler struct {
	ctx             *Context
	requestHandler  RequestHandler
	responseHandler ResponseHandler
}

func (h *restHandler) Release() {
	// TODO:: implement it
}

func (h *restHandler) OnRequest(conn rony.Conn, ctx *gateway.RequestCtx, writer gateway.BodyWriter) {
	h.ctx = &Context{
		reqCtx: ctx,
		conn:   conn,
	}

	err := h.requestHandler(h.ctx, writer)
	if err != nil {
		_ = conn.SendBinary(0, tools.StrToByte(err.Error()))
		return
	}

	return
}

func (h *restHandler) OnResponse(data []byte, bodyWriter gateway.BodyWriter, hdrWriter *gateway.HeaderWriter) {
	me := rony.PoolMessageEnvelope.Get()
	_ = me.Unmarshal(data)
	h.responseHandler(me, bodyWriter, hdrWriter)
	rony.PoolMessageEnvelope.Put(me)
	return
}
