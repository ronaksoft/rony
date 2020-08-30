package main

import (
    "<%= goPackage() %>/service"
	"git.ronaksoft.com/ronak/rony"
	"git.ronaksoft.com/ronak/rony/edge"
	"git.ronaksoft.com/ronak/rony/gateway"
	"os"
)

/*
   Copyright Ronak Software Group 2020
*/

type Server struct {
	e *edge.Server
}

func NewServer(serverID string, opts ...edge.Option) *Server {
	s := &Server{}
	s.e = edge.NewServer(serverID, s, opts...)

	// Register the task service to the server
	service.NewTaskServer(&service.Task{}).Register(s.e)

	return s
}

func (s *Server) Start() error {
	err := s.e.RunCluster()
	if err != nil {
		return err
	}
	s.e.RunGateway()
	return nil
}

func (s *Server) Shutdown(signals ...os.Signal) {
	s.e.ShutdownWithSignal(signals...)
}

func (s *Server) OnMessage(ctx *edge.DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	panic("implement me")
}

func (s *Server) Prepare(ctx *edge.DispatchCtx, data []byte, kvs ...gateway.KeyValue) (err error) {
	panic("implement me")
}

func (s *Server) Done(ctx *edge.DispatchCtx) {
	panic("implement me")
}

func (s *Server) OnOpen(conn gateway.Conn) {
	panic("implement me")
}

func (s *Server) OnClose(conn gateway.Conn) {
	panic("implement me")
}
