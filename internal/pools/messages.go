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
	x := v.(*msg.MessageEnvelope)
	x.Message = x.Message[:0]
	x.Constructor = 0
	x.RequestID = 0
	return x
}

func ReleaseMessageEnvelope(x *msg.MessageEnvelope) {
	messageEnvelopePool.Put(x)
}

// UpdateEnvelope Pool
var updateEnvelopePool sync.Pool

func AcquireUpdateEnvelope() *msg.UpdateEnvelope {
	v := updateEnvelopePool.Get()
	if v == nil {
		return &msg.UpdateEnvelope{}
	}
	x := v.(*msg.UpdateEnvelope)
	x.Update = x.Update[:0]
	x.UpdateID = 0
	x.UCount = 0
	x.Timestamp = 0
	x.Constructor = 0
	return x
}

func ReleaseUpdateEnvelope(x *msg.UpdateEnvelope) {
	updateEnvelopePool.Put(x)
}

// ProtoMessage Pool
var protoMessagePool sync.Pool

func AcquireProtoMessage() *msg.ProtoMessage {
	v := protoMessagePool.Get()
	if v == nil {
		return &msg.ProtoMessage{}
	}
	x := v.(*msg.ProtoMessage)
	x.AuthID = 0
	x.Payload = x.Payload[:0]
	x.MessageKey = x.MessageKey[:0]
	return x
}

func ReleaseProtoMessage(x *msg.ProtoMessage) {
	protoMessagePool.Put(x)
}

var raftCommandPool sync.Pool

func AcquireRaftCommand() *msg.RaftCommand {
	v := raftCommandPool.Get()
	if v == nil {
		return &msg.RaftCommand{}
	}
	x := v.(*msg.RaftCommand)
	x.AuthID = 0
	x.UserID = 0
	x.Envelope = nil
	return x
}

func ReleaseRaftCommand(x *msg.RaftCommand) {
	raftCommandPool.Put(x)
}
