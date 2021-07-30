package main

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/example/task_manager/modules/auth"
	"github.com/ronaksoft/rony/example/task_manager/modules/task"
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
		taskService := &rpc.TaskManager{
			ModuleBase: task.New(edgeServer.Store()),
		}
		authService := &rpc.Auth{
			ModuleBase: auth.New(edgeServer.Store()),
		}

		auth.RegisterAuth(authService, edgeServer)
		task.RegisterTaskManager(taskService, edgeServer, authService.CheckSession, authService.MustAuthorized)

		// Start the edge server components
		edgeServer.Start()

		// Wait until a shutdown signal received.
		edgeServer.WaitForSignal(os.Kill, os.Interrupt)
		edgeServer.Shutdown()
		return nil
	},
}
