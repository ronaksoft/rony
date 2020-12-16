package rony

import (
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

//go:generate protoc -I=. --go_out=paths=source_relative:. msg.proto imsg.proto
//go:generate protoc -I=. --gorony_out=paths=source_relative:. msg.proto imsg.proto
var (
	ConstructorNames = map[int64]string{}
)

func ErrorMessage(out *MessageEnvelope, reqID uint64, errCode, errItem string) {
	errMessage := PoolError.Get()
	errMessage.Code = errCode
	errMessage.Items = errItem
	out.Fill(reqID, C_Error, errMessage)
	PoolError.Put(errMessage)
	return
}

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

	mo := proto.MarshalOptions{
		UseCachedSize: true,
	}
	b := pools.Bytes.GetCap(mo.Size(p))
	b, _ = mo.MarshalAppend(b, p)
	x.Message = append(x.Message[:0], b...)
	pools.Bytes.Put(b)
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

func (x *RaftCommand) Fill(senderID []byte, e *MessageEnvelope, kvs ...*KeyValue) {
	x.Sender = append(x.Sender[:0], senderID...)
	x.Store = append(x.Store[:0], kvs...)
	e.DeepCopy(x.Envelope)
}

func (x *ClusterMessage) Fill(senderID []byte, e *MessageEnvelope, kvs ...*KeyValue) {
	x.Sender = append(x.Sender[:0], senderID...)
	x.Store = append(x.Store[:0], kvs...)
	e.DeepCopy(x.Envelope)
}
