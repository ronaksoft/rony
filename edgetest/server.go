package edgetest

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	dummyGateway "github.com/ronaksoft/rony/internal/gateway/dummy"
)

/*
   Creation Time: 2020 - Dec - 09
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type CheckFunc func(b []byte, kv ...*rony.KeyValue) error

var (
	connID uint64 = 1
)

type Server struct {
	edge *edge.Server
	gw   *dummyGateway.Gateway
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

func (s *Server) SetGlobalPreHandlers(h ...edge.Handler) {
	s.edge.SetGlobalPreHandlers(h...)
}

func (s *Server) SetHandlers(constructor int64, h ...edge.Handler) {
	s.edge.SetHandler(edge.NewHandlerOptions().SetConstructor(constructor).SetHandler(h...))
}

func (s *Server) SetGlobalPostHandlers(h ...edge.Handler) {
	s.edge.SetGlobalPostHandlers(h...)
}

func (s *Server) RPC() *rpcCtx {
	return newRPCContext(s.gw)
}

func (s *Server) REST() *restCtx {
	return newRESTContext(s.gw)
}

func (s *Server) Scenario() *scenario {
	return newScenario()
}
