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

//go:generate protoc -I=. --go_out=paths=source_relative:. msg.proto options.proto
//go:generate protoc -I=. --gorony_out=paths=source_relative,option=no_edge_dep:. msg.proto
func init() {}

/*
	Extra methods for MessageEnvelope
*/

func (x *MessageEnvelope) Fill(reqID uint64, constructor int64, p proto.Message, kvs ...*KeyValue) {
	x.RequestID = reqID
	x.Constructor = constructor

	// Fill Header
	if cap(x.Header) >= len(kvs) {
		x.Header = x.Header[:len(kvs)]
	} else {
		x.Header = make([]*KeyValue, len(kvs))
	}
	for idx, kv := range kvs {
		if x.Header[idx] == nil {
			x.Header[idx] = &KeyValue{}
		}
		kv.DeepCopy(x.Header[idx])
	}

	// Fill Message
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
