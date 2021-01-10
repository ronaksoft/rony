package edgetest

import (
	"github.com/ronaksoft/rony/edge"
	dummyGateway "github.com/ronaksoft/rony/gateway/dummy"
	"github.com/ronaksoft/rony/pools"
	"time"
)

/*
   Creation Time: 2020 - Dec - 09
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Server struct {
	edge *edge.Server
	gw   *dummyGateway.Gateway
	ctxs []*context
}

func NewServer(serverID string, d edge.Dispatcher) *Server {
	s := &Server{}
	s.edge = edge.NewServer(serverID,
		edge.WithDispatcher(d),
		edge.WithTestGateway(edge.DummyGatewayConfig{
			Exposer: func(gw *dummyGateway.Gateway) {
				s.gw = gw
			},
		}),
	)
	return s
}

func (s *Server) Start() {
	_ = s.edge.StartGateway()
}

func (s *Server) Shutdown() {
	s.edge.Shutdown()
}

func (s *Server) RealEdge() *edge.Server {
	return s.edge
}

func (s *Server) SetPreHandlers(h ...edge.Handler) {
	s.edge.SetPreHandlers(h...)
}

func (s *Server) SetHandlers(constructor int64, h ...edge.Handler) {
	s.edge.SetHandlers(constructor, true, h...)
}

func (s *Server) SetPostHandlers(h ...edge.Handler) {
	s.edge.SetPostHandlers(h...)
}

func (s *Server) Context() *context {
	return newContext(s.gw)
}

func (s *Server) AppendContext() *context {
	ctx := newContext(s.gw)
	s.ctxs = append(s.ctxs, ctx)
	return ctx
}

func (s *Server) RunAll(timeout time.Duration) error {
	mErr := &multiError{}
	wg := pools.AcquireWaitGroup()
	for _, ctx := range s.ctxs {
		wg.Add(1)
		go func(ctx *context) {
			err := ctx.Run(timeout)
			mErr.AddError(err)
			wg.Done()
		}(ctx)
	}
	wg.Wait()
	if mErr.HasError() {
		return mErr
	}
	return nil
}
