package edgetest

import (
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

type Server struct {
	edge *edge.Server
	gw   *dummyGateway.Gateway
}

func NewServer(serverID string, d edge.Dispatcher) *Server {
	s := &Server{}
	s.edge = edge.NewServer(serverID, d,
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

func (s *Server) Context() *conn {
	return newConn(s.gw)
}
