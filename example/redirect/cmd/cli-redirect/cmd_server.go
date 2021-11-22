package main

import (
	"os"
	"runtime"
	"time"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/example/redirect/rpc"
	"github.com/spf13/cobra"
)

var edgeServer *edge.Server

var ServerCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.BindCmdFlags(cmd)
		if err != nil {
			return errors.Wrap("bind flag:")(err)
		}

		// Instantiate the edge server
		edgeServer = edge.NewServer(
			config.GetString("server.id"),
			edge.WithDataDir(config.GetString("data.path")),
			edge.WithTcpGateway(edge.TcpGatewayConfig{
				Concurrency:   runtime.NumCPU() * 100,
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
		rpc.RegisterSample(&rpc.Sample{}, edgeServer)

		// Start the edge server components
		edgeServer.Start()

		time.Sleep(time.Second * 3)

		if config.GetString("seed") != "" {
			n, err := edgeServer.Cluster().Join(config.GetString("seed"))
			if err != nil {
				cmd.Println("Error On Joining Cluster:", err)
			} else {
				cmd.Println("Joined", config.GetString("seed"), n)
			}
		}

		// Wait until a shutdown signal received.
		edgeServer.WaitForSignal(os.Kill, os.Interrupt)
		_ = edgeServer.Cluster().Leave()
		edgeServer.Shutdown()

		return nil
	},
}
