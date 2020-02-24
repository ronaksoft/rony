package pools

import (
	"git.ronaksoftware.com/ronak/rony/msg"
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

func AcquireMessageEnvelope() *msg.MessageEnvelope {
	v := messageEnvelopePool.Get()
	if v == nil {
		return &msg.MessageEnvelope{}
	}
	return v.(*msg.MessageEnvelope)
}

func ReleaseMessageEnvelope(x *msg.MessageEnvelope) {
	x.Message = x.Message[:0]
	x.Constructor = 0
	x.RequestID = 0
	messageEnvelopePool.Put(x)
}

// UpdateEnvelope Pool
var updateEnvelopePool sync.Pool

func AcquireUpdateEnvelope() *msg.UpdateEnvelope {
	v := updateEnvelopePool.Get()
	if v == nil {
		return &msg.UpdateEnvelope{}
	}
	return v.(*msg.UpdateEnvelope)
}

func ReleaseUpdateEnvelope(x *msg.UpdateEnvelope) {
	x.Update = x.Update[:0]
	x.UpdateID = 0
	x.UCount = 0
	x.Timestamp = 0
	x.Constructor = 0
	updateEnvelopePool.Put(x)
}

// ProtoMessage Pool
var protoMessagePool sync.Pool

func AcquireProtoMessage() *msg.ProtoMessage {
	v := protoMessagePool.Get()
	if v == nil {
		return &msg.ProtoMessage{}
	}
	return v.(*msg.ProtoMessage)
}

func ReleaseProtoMessage(x *msg.ProtoMessage) {
	x.AuthID = 0
	x.Payload = x.Payload[:0]
	x.MessageKey = x.MessageKey[:0]
	protoMessagePool.Put(x)
}

var raftCommandPool sync.Pool

func AcquireRaftCommand() *msg.RaftCommand {
	v := raftCommandPool.Get()
	if v == nil {
		return &msg.RaftCommand{}
	}
	return v.(*msg.RaftCommand)
}

func ReleaseRaftCommand(x *msg.RaftCommand) {
	x.ConnID = 0
	x.AuthID = 0
	x.Envelope = nil
	raftCommandPool.Put(x)
}
