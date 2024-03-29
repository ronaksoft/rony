package main

import (
	"fmt"
	"os"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools/cliutil"
	"github.com/spf13/cobra"
)

var ClientCmd = &cobra.Command{
	Use: "client",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.BindCmdFlags(cmd)
		if err != nil {
			return errors.WrapText("bind flag:")(err)
		}

		// Sample code for creating a client
		// Instantiate a websocket connection, to use http connection we could use edgec.NewHttp
		wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort: fmt.Sprintf("%s:%d", config.GetString("host"), config.GetInt("port")),
			Handler: func(m *rony.MessageEnvelope) {
				fmt.Println(m.RequestID, registry.C(m.Constructor))
			},
		})

		// Start the websocket connection manager
		err = wsc.Start()
		if err != nil {
			return errors.WrapText("websocket client:")(err)
		}

		ShellCmd.AddCommand(ExitCmd)
		cliutil.RunShell(ShellCmd)

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
