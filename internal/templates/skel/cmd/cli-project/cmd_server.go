package main

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cluster"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/gateway"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var ServerCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.BindCmdFlags(cmd)
		if err != nil {
			return err
		}

		s := NewServer(
			config.GetString("server.id"),
			edge.WithTcpGateway(edge.TcpGatewayConfig{
				Concurrency:   100,
				ListenAddress: ":80",
				MaxBodySize:   0,
				MaxIdleTime:   time.Minute,
			}),
			edge.WithGossipCluster(edge.GossipClusterConfig{
				Bootstrap:  config.GetBool("replica.bootstrap"),
				GossipPort: config.GetInt("gossip.port"),
				RaftPort:   config.GetInt("replica.port"),
				ReplicaSet: config.GetUint64("replica.set"),
				Mode:       cluster.Mode(config.GetString("replica.mode")),
			}),
		)
		err = s.Start()
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

func (s *Server) OnMessage(ctx *edge.DispatchCtx, envelope *rony.MessageEnvelope) {
	panic("implement me")
}

func (s *Server) Interceptor(ctx *edge.DispatchCtx, data []byte) (err error) {
	panic("implement me")
}

func (s *Server) Done(ctx *edge.DispatchCtx) {
	panic("implement me")
}

func (s *Server) OnOpen(conn gateway.Conn, kvs ...*rony.KeyValue) {
	panic("implement me")
}

func (s *Server) OnClose(conn gateway.Conn) {
	panic("implement me")
}
