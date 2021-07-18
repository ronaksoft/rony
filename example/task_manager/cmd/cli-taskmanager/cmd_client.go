package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/example/task_manager/rpc"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
	"github.com/spf13/cobra"
	"os"
)

var ClientCmd = &cobra.Command{
	Use: "client",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.BindCmdFlags(cmd)
		if err != nil {
			return err
		}

		// Sample code for creating a client
		// Instantiate a websocket connection, to use http connection we could use edgec.NewHttp
		wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort: fmt.Sprintf("%s:%d", config.GetString("host"), config.GetInt("port")),
			Handler: func(m *rony.MessageEnvelope) {
				fmt.Println(m.RequestID, registry.ConstructorName(m.Constructor))
			},
		})

		// Start the websocket connection manager
		err = wsc.Start()
		if err != nil {
			return err
		}

		// Instantiate the client stub code and set its underlying client connection
		exitCmd := &cobra.Command{
			Use: "exit",
			Run: func(cmd *cobra.Command, args []string) {
				os.Exit(1)
			},
		}

		var shellCmd = &cobra.Command{}
		shellCmd.AddCommand(exitCmd)
		rpc.RegisterAuthCli(&AuthCLI{}, wsc, shellCmd)
		rpc.RegisterTaskManagerCli(&TaskManagerCLI{}, wsc, shellCmd)
		tools.RunShell(shellCmd)
		return nil
	},
}

type AuthCLI struct {
}

func (a *AuthCLI) Register(cli *rpc.AuthClient, cmd *cobra.Command, args []string) error {
	req := &rpc.RegisterRequest{
		Username:  config.GetString("username"),
		FirstName: config.GetString("firstName"),
		LastName:  config.GetString("lastName"),
		Password:  config.GetString("password"),
	}
	res, err := cli.Register(req)
	if err != nil {
		cmd.Println("Error: ", err)
		return err
	}
	cmd.Println("Response: ", res.String())
	return nil
}

func (a *AuthCLI) Login(cli *rpc.AuthClient, cmd *cobra.Command, args []string) error {
	req := &rpc.LoginRequest{
		Username:  config.GetString("username"),
		Password:  config.GetString("password"),
	}
	res, err := cli.Login(req)
	if err != nil {
		cmd.Println("Error: ", err)
		return err
	}
	cmd.Println("Response: ", res.String())
	return nil
}

type TaskManagerCLI struct {}

func (t *TaskManagerCLI) Create(cli *rpc.TaskManagerClient, cmd *cobra.Command, args []string) error {
	panic("implement me")
}

func (t *TaskManagerCLI) Get(cli *rpc.TaskManagerClient, cmd *cobra.Command, args []string) error {
	panic("implement me")
}

func (t *TaskManagerCLI) Delete(cli *rpc.TaskManagerClient, cmd *cobra.Command, args []string) error {
	panic("implement me")
}

func (t *TaskManagerCLI) List(cli *rpc.TaskManagerClient, cmd *cobra.Command, args []string) error {
	panic("implement me")
}

