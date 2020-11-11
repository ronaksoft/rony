package main

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/gateway"
	tcpGateway "github.com/ronaksoft/rony/gateway/tcp"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var ServerCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := NewServer(
			config.GetString("server.id"),
			edge.WithTcpGateway(tcpGateway.Config{
				Concurrency:   100,
				ListenAddress: "",
				MaxBodySize:   0,
				MaxIdleTime:   time.Minute,
				Protocol:      tcpGateway.Auto,
			}),
			edge.WithGossipPort(config.GetInt("gossip.port")),
			edge.WithReplicaSet(config.GetUint64("replica.set"), config.GetInt("replica.port"), config.GetBool("replica.bootstrap")),
		)
		err := s.Start()
		if err != nil {
			panic(err)
		}

		// Wait for Kill or Interrupt signal to shutdown the server gracefully
		s.Shutdown(os.Kill, os.Interrupt)
		return nil
	},
}

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

func (s *Server) OnMessage(ctx *edge.DispatchCtx, envelope *rony.MessageEnvelope, kvs ...*rony.KeyValue) {
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