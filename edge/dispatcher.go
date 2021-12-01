package edge

import (
	"encoding/json"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/pools"
)

/*
   Creation Time: 2020 - Nov - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Dispatcher interface {
	// Encode will be called on the outgoing messages to encode them into the connection.
	// it is responsible for write data to conn
	Encode(conn rony.Conn, streamID int64, me *rony.MessageEnvelope) error
	// Decode decodes the incoming wire messages and converts it to a rony.MessageEnvelope
	Decode(data []byte, me *rony.MessageEnvelope) error
	// Done will be called when the context has been finished, this lets cleaning up, or in case you need to flush the
	// messages and updates in one go.
	Done(ctx *DispatchCtx)
	// OnOpen will be called when a new connection has been opened
	OnOpen(conn rony.Conn, kvs ...*rony.KeyValue)
	// OnClose will be called when a connection is closed
	OnClose(conn rony.Conn)
}

// DefaultDispatcher is a default implementation of Dispatcher. You only need to set OnMessageFunc with
type DefaultDispatcher struct{}

func (s *DefaultDispatcher) Encode(conn rony.Conn, streamID int64, me *rony.MessageEnvelope) error {
	buf := pools.Buffer.FromProto(me)
	_ = conn.WriteBinary(streamID, *buf.Bytes())
	pools.Buffer.Put(buf)

	return nil
}

func (s *DefaultDispatcher) Decode(data []byte, me *rony.MessageEnvelope) error {
	return me.Unmarshal(data)
}

func (s *DefaultDispatcher) Done(ctx *DispatchCtx) {
	ctx.BufferPopAll(
		func(envelope *rony.MessageEnvelope) {
			buf := pools.Buffer.FromProto(envelope)
			_ = ctx.Conn().WriteBinary(ctx.StreamID(), *buf.Bytes())
			pools.Buffer.Put(buf)
		},
	)
}

func (s *DefaultDispatcher) OnOpen(conn rony.Conn, kvs ...*rony.KeyValue) {
	for _, kv := range kvs {
		conn.Set(kv.Key, kv.Value)
	}
}

func (s *DefaultDispatcher) OnClose(conn rony.Conn) {
	// Do nothing
}

type JSONDispatcher struct{}

func (j *JSONDispatcher) Encode(conn rony.Conn, streamID int64, me *rony.MessageEnvelope) error {
	b, err := json.Marshal(me)
	if err != nil {
		return err
	}

	_ = conn.WriteBinary(streamID, b)

	return nil
}

func (j *JSONDispatcher) Decode(data []byte, me *rony.MessageEnvelope) error {
	return json.Unmarshal(data, me)
}

func (j *JSONDispatcher) Done(ctx *DispatchCtx) {
	ctx.BufferPopAll(
		func(envelope *rony.MessageEnvelope) {
			b, err := json.Marshal(envelope)
			if err != nil {
				return
			}

			_ = ctx.conn.WriteBinary(ctx.StreamID(), b)
		},
	)
}

func (j *JSONDispatcher) OnOpen(conn rony.Conn, kvs ...*rony.KeyValue) {
	for _, kv := range kvs {
		conn.Set(kv.Key, kv.Value)
	}
}

func (j *JSONDispatcher) OnClose(conn rony.Conn) {
	// DO NOTHING
}
