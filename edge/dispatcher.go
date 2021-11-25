package edge

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/pools"
	"google.golang.org/protobuf/proto"
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
	Encode(conn rony.Conn, steamID int64, me *rony.MessageEnvelope) error
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

// defaultDispatcher is a default implementation of Dispatcher. You only need to set OnMessageFunc with
type defaultDispatcher struct{}

func (s *defaultDispatcher) Encode(conn rony.Conn, streamID int64, me *rony.MessageEnvelope) error {
	buf := pools.Buffer.FromProto(me)
	mo := proto.MarshalOptions{UseCachedSize: true}
	bb, _ := mo.MarshalAppend(*buf.Bytes(), me)
	buf.SetBytes(&bb)
	_ = conn.WriteBinary(streamID, *buf.Bytes())
	pools.Buffer.Put(buf)

	return nil
}

func (s *defaultDispatcher) Decode(data []byte, me *rony.MessageEnvelope) error {
	return me.Unmarshal(data)
}

func (s *defaultDispatcher) Done(ctx *DispatchCtx) {
	ctx.BufferPopAll(func(envelope *rony.MessageEnvelope) {
		buf := pools.Buffer.FromProto(envelope)
		_ = ctx.Conn().WriteBinary(ctx.StreamID(), *buf.Bytes())
		pools.Buffer.Put(buf)
	})
}

func (s *defaultDispatcher) OnOpen(conn rony.Conn, kvs ...*rony.KeyValue) {
	for _, kv := range kvs {
		conn.Set(kv.Key, kv.Value)
	}
}

func (s *defaultDispatcher) OnClose(conn rony.Conn) {
	// Do nothing
}
