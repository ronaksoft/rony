package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/example/task_manager/modules/auth"
	"github.com/ronaksoft/rony/example/task_manager/modules/task"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
	"github.com/spf13/cobra"
	"os"
)

var (
	_sessionID string
)

func getSessionID() *rony.KeyValue {
	return &rony.KeyValue{
		Key:   "SessionID",
		Value: _sessionID,
	}
}

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
		task.RegisterTaskManagerCli(&TaskManagerCLI{}, wsc, shellCmd)
		tools.RunShell(shellCmd)

		return nil
	},
}

type TaskManagerCLI struct{}

func (t *TaskManagerCLI) Register(cli *auth.AuthClient, cmd *cobra.Command, args []string) error {
	req := &auth.RegisterRequest{
		Username:  config.GetString("username"),
		FirstName: config.GetString("firstName"),
		LastName:  config.GetString("lastName"),
		Password:  config.GetString("password"),
	}
	res, err := cli.Register(req)
	if err != nil {
		return err
	}
	_sessionID = res.GetSessionID()
	cmd.Println("Response: ", res.String())

	return nil
}

func (t *TaskManagerCLI) Login(cli *auth.AuthClient, cmd *cobra.Command, args []string) error {
	req := &auth.LoginRequest{
		Username: config.GetString("username"),
		Password: config.GetString("password"),
	}
	res, err := cli.Login(req)
	if err != nil {
		return err
	}
	_sessionID = res.GetSessionID()
	cmd.Println("Response: ", res.String())

	return nil
}

func (t *TaskManagerCLI) Create(cli *task.TaskManagerClient, cmd *cobra.Command, args []string) error {
	req := &task.CreateRequest{
		Title:   config.GetString("title"),
		TODOs:   config.GetStringSlice("tODOs"),
		DueDate: config.GetInt64("dueDate"),
	}
	res, err := cli.Create(req, getSessionID())
	if err != nil {
		return err
	}
	cmd.Println("Response: ", res.String())

	return nil
}

func (t *TaskManagerCLI) Get(cli *task.TaskManagerClient, cmd *cobra.Command, args []string) error {
	req := &task.GetRequest{
		TaskID: config.GetInt64("TaskID"),
	}
	res, err := cli.Get(req, getSessionID())
	if err != nil {
		return err
	}
	cmd.Println("Response: ", res.String())

	return nil
}

func (t *TaskManagerCLI) Delete(cli *task.TaskManagerClient, cmd *cobra.Command, args []string) error {
	req := &task.DeleteRequest{
		TaskID: config.GetInt64("TaskID"),
	}
	res, err := cli.Delete(req, getSessionID())
	if err != nil {
		return err
	}
	cmd.Println("Response: ", res.String())

	return nil
}

func (t *TaskManagerCLI) List(cli *task.TaskManagerClient, cmd *cobra.Command, args []string) error {
	req := &task.ListRequest{
		Offset: config.GetInt32("offset"),
		Limit:  config.GetInt32("limit"),
	}
	res, err := cli.List(req, getSessionID())
	if err != nil {
		return err
	}
	cmd.Println("Response: ", res.String())

	return nil
}
