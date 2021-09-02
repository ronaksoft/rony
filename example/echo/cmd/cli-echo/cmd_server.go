package main

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	service "github.com/ronaksoft/rony/example/echo/rpc"
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
		)

		// Register the implemented service into the edge server
		service.RegisterSample(&service.Sample{}, edgeServer)

		// Start the edge server components
		edgeServer.Start()

		// Wait until a shutdown signal received.
		edgeServer.WaitForSignal(os.Kill, os.Interrupt)
		edgeServer.Shutdown()

		return nil
	},
}
