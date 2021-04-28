package edge

import (
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
	// OnMessage is called on every message pushed to RequestCtx.
	// All the input arguments are valid in the function context, if you need to pass 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it.
	// **NOTE**: This is not called on non-persistent connections, instead it is user's responsibility to check the
	// ctx.BufferSize() and ctx.BufferPop() functions
	OnMessage(ctx *DispatchCtx, envelope *rony.MessageEnvelope)
	// Interceptor is called before any handler called.
	// All the input arguments are valid in the function context, if you need to pass 'data' or 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it. If 'err' is not nil then envelope will be
	// discarded, it is the user's responsibility to send back appropriate message using 'conn'
	// Note that conn IS NOT nil in any circumstances.
	Interceptor(ctx *DispatchCtx, data []byte) (err error)
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

func (s *defaultDispatcher) OnMessage(ctx *DispatchCtx, envelope *rony.MessageEnvelope) {
	buf := pools.Buffer.FromProto(envelope)
	_ = ctx.Conn().SendBinary(ctx.StreamID(), *buf.Bytes())
	pools.Buffer.Put(buf)
}

func (s *defaultDispatcher) Interceptor(ctx *DispatchCtx, data []byte) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (s *defaultDispatcher) Done(ctx *DispatchCtx) {
	// Do nothing
}

func (s *defaultDispatcher) OnOpen(conn rony.Conn, kvs ...*rony.KeyValue) {
	// Do nothing
}

func (s *defaultDispatcher) OnClose(conn rony.Conn) {
	// Do nothing
}
