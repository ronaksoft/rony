package edge

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/gateway"
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

var messageEnvelopePool sync.Pool

func acquireMessageEnvelope() *rony.MessageEnvelope {
	v := messageEnvelopePool.Get()
	if v == nil {
		return &rony.MessageEnvelope{}
	}
	return v.(*rony.MessageEnvelope)
}

func releaseMessageEnvelope(x *rony.MessageEnvelope) {
	x.Message = x.Message[:0]
	x.Constructor = 0
	x.RequestID = 0
	x.Header = x.Header[:0]
	x.Auth = x.Auth[:0]
	messageEnvelopePool.Put(x)
}

var tunnelMessagePool sync.Pool

func acquireTunnelMessage() *rony.TunnelMessage {
	v := tunnelMessagePool.Get()
	if v == nil {
		return &rony.TunnelMessage{
			Envelope: &rony.MessageEnvelope{},
		}
	}
	return v.(*rony.TunnelMessage)
}

func releaseTunnelMessage(x *rony.TunnelMessage) {
	x.Store = x.Store[:0]
	x.SenderID = x.SenderID[:0]
	x.SenderReplicaSet = 0
	x.Envelope.Constructor = 0
	x.Envelope.Message = x.Envelope.Message[:0]
	x.Envelope.RequestID = 0
	tunnelMessagePool.Put(x)
}

var raftCommandPool sync.Pool

func acquireRaftCommand() *rony.RaftCommand {
	v := raftCommandPool.Get()
	if v == nil {
		return &rony.RaftCommand{
			Envelope: &rony.MessageEnvelope{},
		}
	}
	return v.(*rony.RaftCommand)
}

func releaseRaftCommand(x *rony.RaftCommand) {
	x.Sender = x.Sender[:0]
	x.Store = x.Store[:0]
	x.Envelope.Message = x.Envelope.Message[:0]
	x.Envelope.RequestID = 0
	x.Envelope.Constructor = 0
	raftCommandPool.Put(x)
}

var waitGroupPool sync.Pool

func acquireWaitGroup() *sync.WaitGroup {
	wgv := waitGroupPool.Get()
	if wgv == nil {
		return &sync.WaitGroup{}
	}

	return wgv.(*sync.WaitGroup)
}

func releaseWaitGroup(wg *sync.WaitGroup) {
	waitGroupPool.Put(wg)
}

var requestCtxPool = sync.Pool{}

func acquireRequestCtx(dispatchCtx *DispatchCtx, quickReturn bool) *RequestCtx {
	var ctx *RequestCtx
	if v := requestCtxPool.Get(); v == nil {
		ctx = newRequestCtx()
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

func acquireDispatchCtx(edge *Server, conn gateway.Conn, streamID int64, serverID []byte, kind ContextKind) *DispatchCtx {
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

	dispatchCtxPool.Put(ctx)
}
