package rony

import (
	"git.ronaksoftware.com/ronak/rony/gateway"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"sync"
)

/*
   Creation Time: 2019 - Oct - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// MessageEnvelope Pool
var messageEnvelopePool sync.Pool

func acquireMessageEnvelope() *MessageEnvelope {
	v := messageEnvelopePool.Get()
	if v == nil {
		return &MessageEnvelope{}
	}
	return v.(*MessageEnvelope)
}

func releaseMessageEnvelope(x *MessageEnvelope) {
	x.Message = x.Message[:0]
	x.Constructor = 0
	x.RequestID = 0
	messageEnvelopePool.Put(x)
}

// UpdateEnvelope Pool
var updateEnvelopePool sync.Pool

func acquireUpdateEnvelope() *UpdateEnvelope {
	v := updateEnvelopePool.Get()
	if v == nil {
		return &UpdateEnvelope{}
	}
	return v.(*UpdateEnvelope)
}

func releaseUpdateEnvelope(x *UpdateEnvelope) {
	x.Update = x.Update[:0]
	x.UpdateID = 0
	x.UCount = 0
	x.Timestamp = 0
	x.Constructor = 0
	updateEnvelopePool.Put(x)
}

var clusterMessagePool sync.Pool

func acquireClusterMessage() *ClusterMessage {
	v := clusterMessagePool.Get()
	if v == nil {
		return &ClusterMessage{}
	}
	return v.(*ClusterMessage)
}

func releaseClusterMessage(x *ClusterMessage) {
	x.AuthID = 0
	x.Sender = x.Sender[:0]
	x.Envelope.Constructor = 0
	x.Envelope.Message = x.Envelope.Message[:0]
	x.Envelope.RequestID = 0
	clusterMessagePool.Put(x)
}

var raftCommandPool sync.Pool

func acquireRaftCommand() *RaftCommand {
	v := raftCommandPool.Get()
	if v == nil {
		return &RaftCommand{}
	}
	return v.(*RaftCommand)
}

func releaseRaftCommand(x *RaftCommand) {
	x.Sender = x.Sender[:0]
	x.AuthID = 0
	x.Envelope = nil
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
		ctx.reset()
		// Just to make sure channel is empty, or empty it if not
		select {
		case <-ctx.nextChan:
		default:
		}
	}
	ctx.stop = false
	ctx.quickReturn = quickReturn
	ctx.dispatchCtx = dispatchCtx
	return ctx
}

func releaseRequestCtx(ctx *RequestCtx) {
	requestCtxPool.Put(ctx)
}

var dispatchCtxPool = sync.Pool{}

func acquireDispatchCtx(edge *EdgeServer, conn gateway.Conn, streamID int64, authID int64, serverID []byte) *DispatchCtx {
	var ctx *DispatchCtx
	if v := dispatchCtxPool.Get(); v == nil {
		ctx = newDispatchCtx(edge)
	} else {
		ctx = v.(*DispatchCtx)
	}
	ctx.conn = conn
	if ctx.conn == nil {
		ctx.kind = gatewayMessage
	} else {
		ctx.kind = clusterMessage
	}
	ctx.streamID = streamID
	ctx.authID = authID
	if len(serverID) > cap(ctx.serverID) {
		pools.Bytes.Put(ctx.serverID)
		ctx.serverID = pools.Bytes.GetCap(len(serverID))
	}
	ctx.serverID = append(ctx.serverID, serverID...)
	return ctx
}

func releaseDispatchCtx(ctx *DispatchCtx) {
	dispatchCtxPool.Put(ctx)
}
