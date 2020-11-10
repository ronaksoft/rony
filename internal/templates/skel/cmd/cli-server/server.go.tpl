package main

import (
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
	return s
}

func (s *Server) Start() error {
	err := s.e.StartCluster()
	if err != nil {
		return err
	}
	s.e.StartGateway()
	return nil
}

func (s *Server) Shutdown(signals ...os.Signal) {
	s.e.ShutdownWithSignal(signals...)
}

func (s *Server) OnMessage(ctx *edge.DispatchCtx, envelope *rony.MessageEnvelope, kvs *rony.KeyValue) {
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
