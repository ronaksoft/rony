package main

import (
	"github.com/ronaksoft/rony/cluster"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/gateway"
	"github.com/ronaksoft/rony/tools"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var edgeServer *edge.Server

var ServerCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.BindCmdFlags(cmd)
		if err != nil {
			return err
		}

		// Instantiate the edge server
		edgeServer = edge.NewServer(
			config.GetString("server.id"),
			edge.WithDataDir(config.GetString("data.path")),
			edge.WithTcpGateway(edge.TcpGatewayConfig{
				Concurrency:   runtime.NumCPU() * 100,
				ListenAddress: config.GetString("gateway.listen"),
				MaxIdleTime:   config.GetDuration("idle-time"),
				Protocol:      gateway.Http,
				ExternalAddrs: config.GetStringSlice("gateway.advertise.url"),
			}),
			edge.WithGossipCluster(edge.GossipClusterConfig{
				ServerID:   tools.StrToByte(config.GetString("server.id")),
				Bootstrap:  config.GetBool("bootstrap"),
				RaftPort:   config.GetInt("raft.port"),
				ReplicaSet: config.GetUint64("replica-set"),
				Mode:       cluster.MultiReplica,
				GossipPort: config.GetInt("gossip.port"),
			}),
			edge.WithUdpTunnel(edge.UdpTunnelConfig{
				ServerID:      config.GetString("server.id"),
				Concurrency:   runtime.NumCPU() * 100,
				ListenAddress: config.GetString("tunnel.listen"),
				ExternalAddrs: config.GetStringSlice("tunnel.advertise.url"),
			}),
		)

		// Register the implemented service into the edge server
		// service.RegisterSampleService(&service.SampleService{}, edgeServer)

		// Start the edge server components
		edgeServer.Start()

		// Wait until a shutdown signal received.
		edgeServer.ShutdownWithSignal(os.Kill, os.Interrupt)
		return nil
	},
}
