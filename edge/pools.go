package edge

import (
	"github.com/ronaksoft/rony"
	"sync"
)

/*
   Creation Time: 2019 - Oct - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var requestCtxPool = sync.Pool{}

func acquireRequestCtx(dispatchCtx *DispatchCtx, quickReturn bool) *RequestCtx {
	var ctx *RequestCtx
	if v := requestCtxPool.Get(); v == nil {
		ctx = newRequestCtx(dispatchCtx.edge)
	} else {
		ctx = v.(*RequestCtx)
	}
	ctx.stop = false
	ctx.quickReturn = quickReturn
	ctx.dispatchCtx = dispatchCtx
	return ctx
}

func releaseRequestCtx(ctx *RequestCtx) {
	// Just to make sure channel is empty, or empty it if not
	select {
	case <-ctx.nextChan:
	default:
	}

	ctx.reqID = 0

	// Put back into the pool
	requestCtxPool.Put(ctx)
}

var dispatchCtxPool = sync.Pool{}

func acquireDispatchCtx(edge *Server, conn rony.Conn, streamID int64, serverID []byte, kind MessageKind) *DispatchCtx {
	var ctx *DispatchCtx
	if v := dispatchCtxPool.Get(); v == nil {
		ctx = newDispatchCtx(edge)
	} else {
		ctx = v.(*DispatchCtx)
	}
	ctx.conn = conn
	ctx.kind = kind
	ctx.streamID = streamID
	ctx.serverID = append(ctx.serverID[:0], serverID...)
	return ctx
}

func releaseDispatchCtx(ctx *DispatchCtx) {
	// Reset the Key-Value store
	ctx.reset()

	// Put back the context into the pool
	dispatchCtxPool.Put(ctx)
}
