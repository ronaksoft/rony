package edgetest

import (
	"github.com/ronaksoft/rony/edge"
	dummyGateway "github.com/ronaksoft/rony/gateway/dummy"
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
		edge.WithTestGateway(dummyGateway.Config{
			Exposer: func(gw *dummyGateway.Gateway) {
				s.gw = gw
			},
		}),
	)
	return s
}

func (s *Server) Start() {
	s.edge.StartGateway()
}

func (s *Server) Shutdown() {
	s.edge.Shutdown()
}

func (s *Server) SetHandlers(constructor int64, h ...edge.Handler) {
	s.edge.SetHandlers(constructor, true, h...)
}

func (s *Server) Context() *conn {
	return newConn(s.gw)
}