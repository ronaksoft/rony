package rony

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/registry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
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
//go:generate protoc -I=. --gorony_out=paths=source_relative,rony_opt=no_edge_dep:. msg.proto
func init() {}

/*
	Extra methods for MessageEnvelope
*/

func (x *MessageEnvelope) Fill(reqID uint64, constructor uint64, p proto.Message, kvs ...*KeyValue) {
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

func (x *MessageEnvelope) Append(key, val string) {
	x.Header = append(x.Header, &KeyValue{
		Key:   key,
		Value: val,
	})
}

func (x *MessageEnvelope) Unwrap() (protoreflect.Message, error) {
	m, err := registry.Unwrap(x)
	if err != nil {
		return nil, err
	}

	return m.ProtoReflect(), nil
}

// Carrier wraps the MessageEnvelope into an envelope carrier which is useful for cross-cutting
// tracing. With this method, instrumentation packages can pass the trace info into the wire.
// Rony has builtin support, and you do not need to use this function explicitly.
// Check edge.WithTracer option to enable tracing
func (x *MessageEnvelope) Carrier() *envelopeCarrier {
	return &envelopeCarrier{
		e: x,
	}
}

func (x *MessageEnvelope) MarshalJSON() ([]byte, error) {
	m, err := registry.Unwrap(x)
	if err != nil {
		return nil, err
	}

	je := MessageEnvelopeJSON{
		RequestID:   x.RequestID,
		Constructor: registry.C(x.Constructor),
	}

	if len(x.Header) > 0 {
		je.Header = map[string]string{}
		for _, kv := range x.Header {
			je.Header[kv.Key] = kv.Value
		}
	}

	je.Message, err = m.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return json.Marshal(je)
}

func (x *MessageEnvelope) UnmarshalJSON(b []byte) error {
	je := MessageEnvelopeJSON{}
	err := json.Unmarshal(b, &je)
	if err != nil {
		return err
	}

	x.RequestID = je.RequestID
	x.Constructor = registry.N(je.Constructor)

	m, err := registry.Get(x.Constructor)
	if err != nil {
		return err
	}

	for k, v := range je.Header {
		x.Append(k, v)
	}

	err = m.UnmarshalJSON(je.Message)
	if err != nil {
		return err
	}

	x.Message, err = proto.Marshal(m)
	if err != nil {
		return err
	}

	return nil
}

// envelopeCarrier is an adapted for MessageEnvelope to implement propagation.TextMapCarrier interface
type envelopeCarrier struct {
	e *MessageEnvelope
}

func (e envelopeCarrier) Get(key string) string {
	return e.e.Get(key, "")
}

func (e envelopeCarrier) Set(key string, value string) {
	e.e.Append(key, value)
}

func (e envelopeCarrier) Keys() []string {
	var keys = make([]string, len(e.e.Header))
	for idx := range e.e.Header {
		keys[idx] = e.e.Header[idx].Key
	}

	return keys
}

/*
	Extra methods for MessageContainer
*/
func (x *MessageContainer) Add(reqID uint64, constructor uint64, p proto.Message, kvs ...*KeyValue) {
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
	// me.Fill(me.RequestID, C_Error, x)
}

// MessageEnvelopeJSON is the JSON representation of MessageEnvelope.
type MessageEnvelopeJSON struct {
	RequestID   uint64            `json:"requestId,omitempty"`
	Header      map[string]string `json:"header,omitempty"`
	Constructor string            `json:"constructor"`
	Message     json.RawMessage   `json:"message,omitempty"`
}
