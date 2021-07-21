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
		)

		// Register the implemented service into the edge server

		taskService := rpc.NewTaskManager(edgeServer.Store())
		rpc.RegisterTaskManagerWithFunc(taskService, edgeServer, func(c int64) []edge.Handler {
			switch c {
			case rpc.C_TaskManagerLogin, rpc.C_TaskManagerRegister:
			default:
				return []edge.Handler{taskService.CheckSession, taskService.MustAuthorized}
			}
			return nil
		})

		// Start the edge server components
		edgeServer.Start()

		// Wait until a shutdown signal received.
		edgeServer.WaitForSignal(os.Kill, os.Interrupt)
		edgeServer.Shutdown()
		return nil
	},
}
