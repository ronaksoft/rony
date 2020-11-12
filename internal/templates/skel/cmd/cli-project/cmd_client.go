package main

import (
	"github.com/ronaksoft/rony/config"
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
		/**
			c := service.NewSampleServiceClient(
				edgec.NewWebsocket(
					edgec.Config{
						HostPort: config.GetString("server.hostport"),
					},
				),
			)
		 **/

		return nil
	},
}
