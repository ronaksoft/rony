package main

import (
	"github.com/spf13/cobra"
	"os"
	"time"
)

/*
   Creation Time: 2020 - Mar - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func init() {
	RootCmd.AddCommand(DemoCmd)
	DemoCmd.AddCommand(DemoRaftCmd, DemoClusterBroadcastCmd, DemoClusterMessageCmd)
}

var DemoCmd = &cobra.Command{
	Use: "demo",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var DemoRaftCmd = &cobra.Command{
	Use: "raft",
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.RemoveAll("./_hdd")
		startFunc("Raft.01", "RS01", 801, true)
		startFunc("Raft.02", "RS01", 802, false)
		startFunc("Raft.03", "RS01", 803, false)
		listFunc()
		joinFunc("Raft.02", "Raft.01")
		joinFunc("Raft.03", "Raft.02")
		time.Sleep(time.Second)
		listFunc()
	},
}

var DemoClusterMessageCmd = &cobra.Command{
	Use: "clusterMessage",
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.RemoveAll("./_hdd")
		startFunc("Cluster.01", "", 801, true)
		startFunc("Cluster.02", "", 802, true)
		startFunc("Cluster.03", "", 803, true)
		joinFunc("Cluster.01", "Cluster.02")
		joinFunc("Cluster.01", "Cluster.03")
		time.Sleep(time.Second * 3)
		listFunc()
		clusterMessage("First.01", "Second.01")
		clusterMessage("Third.01", "First.01")
		for _, e := range Edges {
			e.Shutdown()
		}
	},
}

var DemoClusterBroadcastCmd = &cobra.Command{
	Use: "clusterBroadcast",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
