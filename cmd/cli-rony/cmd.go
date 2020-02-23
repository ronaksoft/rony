package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"path/filepath"
)

/*
   Creation Time: 2019 - Oct - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var Edges map[string]*rony.EdgeServer

func init() {
	RootCmd.AddCommand(StartCmd, StopCmd, ListCmd, JoinCmd)

	RootCmd.PersistentFlags().String(FlagBundleID, "Test", "")
	RootCmd.PersistentFlags().String(FlagInstanceID, "01", "")
	RootCmd.PersistentFlags().Int(FlagGossipPort, 801, "")
	RootCmd.PersistentFlags().Bool(FlagBootstrap, false, "")

}

var StartCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		bundleID, _ := cmd.Flags().GetString(FlagBundleID)
		instanceID, _ := cmd.Flags().GetString(FlagInstanceID)
		gossipPort, _ := cmd.Flags().GetInt(FlagGossipPort)
		raftBootstrap, _ := cmd.Flags().GetBool(FlagBootstrap)

		serverID := fmt.Sprintf("%s.%s", bundleID, instanceID)
		if _, ok := Edges[serverID]; !ok {
			Edges[serverID] = rony.NewEdgeServer(bundleID, instanceID,
				rony.WithWebsocketGateway(websocketGateway.Config{
					NewConnectionWorkers: 1,
					MaxConcurrency:       1,
					MaxIdleTime:          0,
					ListenAddress:        "0.0.0.0:0",
				}),
				rony.WithDataPath(filepath.Join("./_hdd", serverID)),
				rony.WithRaft(gossipPort*10, raftBootstrap),
				rony.WithGossipPort(gossipPort),
			)
			err := Edges[serverID].Run()
			if err != nil {
				fmt.Println(err)
				delete(Edges, serverID)
				return
			}
			StopCmd.AddCommand(&cobra.Command{
				Use: serverID,
				Run: func(cmd *cobra.Command, args []string) {
					s, ok := Edges[serverID]
					if ok {
						s.Shutdown()
						delete(Edges, serverID)
					}
					StopCmd.RemoveCommand(cmd)
				},
			})
		}
	},
}

var StopCmd = &cobra.Command{
	Use: "stop",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var ListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		var rows []string
		rows = append(rows,
			"ServerID | Raft Members | Raft State | Serf Members | Serf State | Gateway ",
			"------- | ------ | ------ | -------- | ------- | ------",
		)

		for id, s := range Edges {
			edgeStats := s.Stats()
			rows = append(rows,
				fmt.Sprintf("%s | %d | %s | %d | %s | %s(%s)", id,
					edgeStats.RaftMembers, edgeStats.RaftState,
					edgeStats.SerfMembers, edgeStats.SerfState,
					edgeStats.GatewayProtocol, edgeStats.GatewayAddr,
				),
			)
		}

		fmt.Println(columnize.SimpleFormat(rows))
	},
}

var JoinCmd = &cobra.Command{
	Use: "join",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Needs Two ServerID")
			return
		}
		e1 := Edges[args[0]]
		e2 := Edges[args[1]]
		if e1 == nil || e2 == nil {
			fmt.Println("Invalid Args")
			return
		}
		e1Stats := e1.Stats()
		e2Stats := e2.Stats()
		fmt.Println("Joining ", e1Stats.Address, "--->", e2Stats.Address)
		err := e1.Join(e2Stats.Address)

		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	Edges = make(map[string]*rony.EdgeServer)
}
