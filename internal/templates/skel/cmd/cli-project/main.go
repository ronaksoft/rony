package main

import (
	"github.com/ronaksoft/rony/config"
	"github.com/spf13/cobra"
)

func main() {
	// Initialize the config package to hold configurable parameters
	_ = config.Init("")

	// Define the configs if this executable is running as a server instance
	config.SetCmdFlags(ServerCmd,
		config.StringFlag("server.id", "YourServerID", "a unique id for this server instance in the whole cluster"),
		config.IntFlag("gossip.port", 2374, "the port server listen for gossip protocol"),
		config.UInt64Flag("replica.set", 1, "the replica set"),
		config.IntFlag("replica.port", 2000, "the port server listen for raft protocol"),
		config.BoolFlag("replica.bootstrap", true, "if this server is the first node of the replica set"),
	)

	// Define the configs if this executable is running as a server instance
	config.SetCmdFlags(ClientCmd,
		config.StringFlag("server.hostport", "localhost:8080", "the host:port of the seed server"),
	)

	RootCmd.AddCommand(ServerCmd, ClientCmd)
	_ = RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use: "<%= projectName() %>",
}
