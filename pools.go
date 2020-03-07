package rony

import (
	"git.ronaksoftware.com/ronak/rony/gateway"
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
	x.Sender = ""
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
	x.ConnID = 0
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

var ctxPool = sync.Pool{}

func acquireContext(conn gateway.Conn, authID int64, quickReturn bool) *Context {
	var ctx *Context
	if v := ctxPool.Get(); v == nil {
		ctx = New()
	} else {
		ctx = v.(*Context)
		ctx.Clear()
		// Just to make sure channel is empty, or empty it if not
		select {
		case <-ctx.NextChan:
		default:
		}
	}
	ctx.Stop = false
	ctx.QuickReturn = quickReturn
	if conn != nil {
		ctx.ConnID = conn.GetConnID()
	} else {
		ctx.ConnID = 0
	}
	ctx.AuthID = authID
	return ctx
}

func releaseContext(ctx *Context) {
	ctxPool.Put(ctx)
}
