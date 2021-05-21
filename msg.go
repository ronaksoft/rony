package rony

import (
	"fmt"
	"github.com/ronaksoft/rony/pools"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2020 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. --go_out=paths=source_relative:. imsg.proto msg.proto options.proto
//go:generate protoc -I=. --gorony_out=paths=source_relative,option=no_edge_dep:. imsg.proto msg.proto
func init() {}

/*
	Extra methods for MessageEnvelope
*/

func (x *MessageEnvelope) Clone() *MessageEnvelope {
	c := PoolMessageEnvelope.Get()
	c.Constructor = x.Constructor
	c.RequestID = x.RequestID
	c.Message = append(c.Message[:0], x.Message...)
	c.Auth = append(c.Auth[:0], x.Auth...)
	if cap(c.Header) >= len(x.Header) {
		c.Header = c.Header[:len(x.Header)]
	} else {
		c.Header = make([]*KeyValue, len(x.Header))
	}
	for idx, kv := range x.Header {
		if c.Header[idx] == nil {
			c.Header[idx] = &KeyValue{}
		}
		kv.DeepCopy(c.Header[idx])
	}
	return c
}

func (x *MessageEnvelope) Fill(reqID uint64, constructor int64, p proto.Message, kvs ...*KeyValue) {
	x.RequestID = reqID
	x.Constructor = constructor
	x.Header = append(x.Header[:0], kvs...)

	buf := pools.Buffer.FromProto(p)
	x.Message = append(x.Message[:0], *buf.Bytes()...)
	pools.Buffer.Put(buf)
}

func (x *MessageEnvelope) Get(key, defaultVal string) string {
	for _, kv := range x.Header {
		if kv.Key == key {
			return kv.Value
		}
	}
	return defaultVal
}

func (x *MessageEnvelope) Set(KVs ...*KeyValue) {
	x.Header = append(x.Header[:0], KVs...)
}

/*
	Extra methods for MessageContainer
*/
func (x *MessageContainer) Add(reqID uint64, constructor int64, p proto.Message, kvs ...*KeyValue) {
	me := PoolMessageEnvelope.Get()
	me.Fill(reqID, constructor, p, kvs...)
	x.Envelopes = append(x.Envelopes, me)
	x.Length += 1
}

/*
	Extra methods for TunnelMessage
*/

func (x *TunnelMessage) Fill(senderID []byte, senderReplicaSet uint64, e *MessageEnvelope, kvs ...*KeyValue) {
	x.SenderID = append(x.SenderID[:0], senderID...)
	x.SenderReplicaSet = senderReplicaSet
	x.Store = append(x.Store[:0], kvs...)
	if x.Envelope == nil {
		x.Envelope = PoolMessageEnvelope.Get()
	}
	e.DeepCopy(x.Envelope)
}

/*
	Extra methods for Error
*/
func (x *Error) Error() string {
	if len(x.Description) > 0 {
		return fmt.Sprintf("%s:%s (%s)", x.Code, x.Items, x.Description)
	} else {
		return fmt.Sprintf("%s:%s", x.Code, x.Items)
	}
}

func (x *Error) Expand() (string, string) {
	return x.Code, x.Items
}

func (x *Error) ToEnvelope(me *MessageEnvelope) {
	me.Fill(me.RequestID, C_Error, x)
}
