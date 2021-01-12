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
	err := config.Init("<%= projectName() %>")
	if err != nil {
		fmt.Println("config initialization had error:", err)
	}

	// Define the configs if this executable is running as a server instance
	// Set the flags as config parameters
	config.SetCmdFlags(ServerCmd,
		config.StringFlag("serverID", tools.RandomID(12), ""),
		config.StringFlag("gatewayListen", "0.0.0.0:80", ""),
		config.StringSliceFlag("gatewayAdvertiseUrl", nil, ""),
		config.StringFlag("tunnelListen", "0.0.0.0:81", ""),
		config.StringSliceFlag("tunnelAdvertiseUrl", nil, ""),
		config.DurationFlag("idleTime", time.Minute, ""),
		config.IntFlag("raftPort", 7080, ""),
		config.UInt64Flag("replicaSet", 1, ""),
		config.IntFlag("gossipPort", 7081, ""),
		config.StringFlag("dataPath", "./_hdd", ""),
		config.BoolFlag("bootstrap", false, ""),
	)

	// Define the configs if this executable is running as a server instance
	config.SetCmdFlags(ClientCmd,
		config.StringFlag("host", "127.0.0.1", "the host of the seed server"),
		config.IntFlag("port", 80, "the port of the seed server"),
	)

	RootCmd.AddCommand(ServerCmd, ClientCmd)
	err = RootCmd.Execute()
	if err != nil {
		fmt.Println("we got error:", err)
	}
}

var RootCmd = &cobra.Command{
	Use: "<%= projectName() %>",
}
