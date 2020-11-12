package edge

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/gateway"
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

// Dispatcher
type Dispatcher interface {
	// All the input arguments are valid in the function context, if you need to pass 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it.
	OnMessage(ctx *DispatchCtx, envelope *rony.MessageEnvelope, kvs ...*rony.KeyValue)
	// All the input arguments are valid in the function context, if you need to pass 'data' or 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it. If 'err' is not nil then envelope will be
	// discarded, it is the user's responsibility to send back appropriate message using 'conn'
	// Note that conn IS NOT nil in any circumstances.
	Interceptor(ctx *DispatchCtx, data []byte, kvs ...gateway.KeyValue) (err error)
	// This will be called when the context has been finished, this lets cleaning up, or in case you need to flush the
	// messages and updates in one go.
	Done(ctx *DispatchCtx)
	// This will be called when a new connection has been opened
	OnOpen(conn gateway.Conn)
	// This will be called when a connection is closed
	OnClose(conn gateway.Conn)
}

// SimpleDispatcher is a naive implementation of Dispatcher. You only need to set OnMessageFunc with
type SimpleDispatcher struct {
	OnMessageFunc func(ctx *DispatchCtx, envelope *rony.MessageEnvelope, kvs ...*rony.KeyValue)
}

func (s *SimpleDispatcher) OnMessage(ctx *DispatchCtx, envelope *rony.MessageEnvelope, kvs ...*rony.KeyValue) {
	if s.OnMessageFunc != nil {
		s.OnMessageFunc(ctx, envelope, kvs...)
		return
	}

	mo := proto.MarshalOptions{UseCachedSize: true}
	eb := pools.Bytes.GetCap(mo.Size(envelope))
	eb, _ = mo.MarshalAppend(eb, envelope)
	_ = ctx.Conn().SendBinary(ctx.StreamID(), eb)
	pools.Bytes.Put(eb)
}

func (s *SimpleDispatcher) Interceptor(ctx *DispatchCtx, data []byte, kvs ...gateway.KeyValue) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (s *SimpleDispatcher) Done(ctx *DispatchCtx) {
	// Do nothing
}

func (s *SimpleDispatcher) OnOpen(conn gateway.Conn) {
	// Do nothing
}

func (s *SimpleDispatcher) OnClose(conn gateway.Conn) {
	panic("implement me")
}
