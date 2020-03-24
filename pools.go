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

var carrierPool = sync.Pool{}

func acquireMessageCarrier(authID int64, e *MessageEnvelope) *carrier {
	cv, _ := carrierPool.Get().(*carrier)
	if cv == nil {
		return &carrier{
			kind:            carrierMessage,
			AuthID:          authID,
			ServerID:        nil,
			MessageEnvelope: e,
			UpdateEnvelope:  nil,
		}
	}
	cv.kind = carrierMessage
	cv.AuthID = authID
	cv.MessageEnvelope = e
	return cv
}

func acquireClusterMessageCarrier(authID int64, serverID string, e *MessageEnvelope) *carrier {
	cv, _ := carrierPool.Get().(*carrier)
	if cv == nil {
		return &carrier{
			kind:            carrierCluster,
			AuthID:          authID,
			ServerID:        []byte(serverID),
			MessageEnvelope: e,
			UpdateEnvelope:  nil,
		}
	}
	if len(serverID) > cap(cv.ServerID) {
		pools.Bytes.Put(cv.ServerID)
		pools.Bytes.GetCap(len(serverID))
	}
	cv.kind = carrierCluster
	cv.ServerID = append(cv.ServerID, serverID...)
	cv.AuthID = authID
	cv.MessageEnvelope = e
	return cv
}

func acquireUpdateCarrier(authID int64, e *UpdateEnvelope) *carrier {
	cv, _ := carrierPool.Get().(*carrier)
	if cv == nil {
		return &carrier{
			kind:            carrierUpdate,
			AuthID:          authID,
			ServerID:        nil,
			MessageEnvelope: nil,
			UpdateEnvelope:  e,
		}
	}
	cv.kind = carrierUpdate
	cv.AuthID = authID
	cv.UpdateEnvelope = e
	return cv
}

func releaseCarrier(c *carrier) {
	c.MessageEnvelope = nil
	c.UpdateEnvelope = nil
	c.ServerID = c.ServerID[:0]
	carrierPool.Put(c)
}
