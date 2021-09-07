package main

import (
	"fmt"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/tools"
	"github.com/spf13/cobra"
	"time"
)

func main() {
	// Initialize the config package
	config.Init()

	// Define the configs if this executable is running as a server instance
	// Set the flags as config parameters
	config.SetCmdFlags(ServerCmd,
		config.StringFlag("server.id", tools.RandomID(12), ""),
		config.StringFlag("gateway.listen", "0.0.0.0:80", ""),
		config.StringSliceFlag("gateway.advertise.url", nil, ""),
		config.StringFlag("tunnel.listen", "0.0.0.0:81", ""),
		config.StringSliceFlag("tunnel.advertise.url", nil, ""),
		config.DurationFlag("idle-time", time.Minute, ""),
		config.Uint64Flag("replica-set", 1, ""),
		config.IntFlag("gossip.port", 7081, ""),
		config.StringFlag("data.path", "./_hdd", ""),
		config.BoolFlag("bootstrap", false, ""),
	)

	// Define the configs if this executable is running as a server instance
	config.SetCmdFlags(ClientCmd,
		config.StringFlag("host", "127.0.0.1", "the host of the seed server"),
		config.IntFlag("port", 80, "the port of the seed server"),
	)

	RootCmd.AddCommand(ServerCmd, ClientCmd)
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println("we got error:", err)
	}
}

var RootCmd = &cobra.Command{
	Use: "<%= projectName() %>",
}
