package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/edgeClient"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
	"time"
)

/*
   Creation Time: 2019 - Oct - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var Edges map[string]*edge.Server

func init() {
	Edges = make(map[string]*edge.Server)
	RootCmd.AddCommand(
		StartCmd, EchoCmd,
	)

	RootCmd.PersistentFlags().String(FlagServerID, "Node", "")
	RootCmd.PersistentFlags().Uint64(FlagReplicaSet, 0, "")
	RootCmd.PersistentFlags().Int(FlagGossipPort, 801, "")
	RootCmd.PersistentFlags().Bool(FlagBootstrap, false, "")

}

var StartCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		replicaSet, _ := cmd.Flags().GetUint64(FlagReplicaSet)
		gossipPort, _ := cmd.Flags().GetInt(FlagGossipPort)
		raftBootstrap, _ := cmd.Flags().GetBool(FlagBootstrap)
		startFunc(serverID, replicaSet, gossipPort, raftBootstrap)
	},
}

func startFunc(serverID string, replicaSet uint64, port int, bootstrap bool) {
	if _, ok := Edges[serverID]; !ok {
		opts := make([]edge.Option, 0)
		opts = append(opts,
			edge.WithWebsocketGateway(websocketGateway.Config{
				NewConnectionWorkers: 1,
				MaxConcurrency:       1000,
				MaxIdleTime:          0,
				ListenAddress:        "0.0.0.0:0",
			}),
			edge.WithDataPath(filepath.Join("./_hdd", serverID)),
			edge.WithGossipPort(port),
		)

		if replicaSet != 0 {
			opts = append(opts, edge.WithReplicaSet(replicaSet, port*10, bootstrap))
		}

		Edges[serverID] = edge.NewServer(serverID, &dispatcher{}, opts...)
		h := pb.NewSampleServer(&Handlers{})
		h.Register(Edges[serverID])
		err := Edges[serverID].RunCluster()
		if err != nil {
			fmt.Println(err)
			delete(Edges, serverID)
			return
		}
		Edges[serverID].RunGateway()
	}
}

var EchoCmd = &cobra.Command{
	Use: "echo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Needs ServerID, e.g. echo First.01")
			return
		}
		e1 := Edges[args[0]]
		if e1 == nil {
			fmt.Println("Invalid Args")
			return
		}
		gatewayAddr := e1.Stats().GatewayAddr
		if len(gatewayAddr) == 0 {
			fmt.Println("No Gateway Addr", gatewayAddr)
			return
		}
		parts := strings.Split(gatewayAddr[0], ":")
		if len(parts) != 2 {
			fmt.Println("Invalid Gateway Addr", gatewayAddr)
			return
		}
		ec := edgeClient.NewWebsocket(edgeClient.Config{
			HostPort: fmt.Sprintf("ws://127.0.0.1:%s", parts[1]),
			Handler: func(m *rony.MessageEnvelope) {
				cmd.Print(m)
			},
		})
		c := pb.NewSampleClient(ec)

		req := pb.PoolEchoRequest.Get()
		defer pb.PoolEchoRequest.Put(req)
		req.Int = tools.RandomInt64(0)
		req.Bool = true
		req.Timestamp = time.Now().UnixNano()
		res, err := c.Echo(req)
		if err != nil {
			cmd.Println("Error:", err)
			return
		}
		cmd.Println(res)
	},
}
