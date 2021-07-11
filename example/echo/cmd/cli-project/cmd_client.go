package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/errors"
	service "github.com/ronaksoft/rony/example/echo/rpc"
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
			return errors.Wrap("bind flag:")(err)
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
			return errors.Wrap("websocket client:")(err)
		}

		// Instantiate the client stub code and set its underlying client connection
		service.RegisterSampleCli(&service.SampleCli{}, wsc, ShellCmd)

		ShellCmd.AddCommand(ExitCmd)
		p := prompt.New(tools.PromptExecutor(ShellCmd), tools.PromptCompleter(ShellCmd))
		p.Run()
		return nil
	},
}

var ShellCmd = &cobra.Command{}

var ExitCmd = &cobra.Command{
	Use: "exit",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}
