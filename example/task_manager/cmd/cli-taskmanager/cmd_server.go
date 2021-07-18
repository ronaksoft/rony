package main

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/example/task_manager/rpc"
	"github.com/spf13/cobra"
	"os"
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
				ListenAddress: config.GetString("gateway.listen"),
				MaxIdleTime:   config.GetDuration("idle-time"),
				Protocol:      rony.TCP,
				ExternalAddrs: config.GetStringSlice("gateway.advertise.url"),
			}),
			edge.WithGossipCluster(edge.GossipClusterConfig{
				Bootstrap:  config.GetBool("bootstrap"),
				ReplicaSet: config.GetUint64("replica-set"),
				GossipPort: config.GetInt("gossip.port"),
			}),
			edge.WithUdpTunnel(edge.UdpTunnelConfig{
				ListenAddress: config.GetString("tunnel.listen"),
				ExternalAddrs: config.GetStringSlice("tunnel.advertise.url"),
			}),
		)

		// Register the implemented service into the edge server
		authService := rpc.NewAuth(edgeServer.Store())
		rpc.RegisterAuth(authService, edgeServer)
		taskService := rpc.NewTaskManager(edgeServer.Store())
		rpc.RegisterTaskManager(taskService, edgeServer, authService.CheckSession, authService.MustAuthorized)

		// Start the edge server components
		edgeServer.Start()

		// Wait until a shutdown signal received.
		edgeServer.ShutdownWithSignal(os.Kill, os.Interrupt)
		return nil
	},
}
