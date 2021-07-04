package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edgec"
	service "github.com/ronaksoft/rony/example/echo/rpc"
	"github.com/ronaksoft/rony/registry"
	"github.com/spf13/cobra"
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
		service.RegisterSampleCli(&service.SampleCli{}, wsc, cmd)

		return nil
	},
}
