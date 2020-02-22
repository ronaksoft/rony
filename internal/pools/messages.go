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

// func AcquireRaftCommand() *msg.
