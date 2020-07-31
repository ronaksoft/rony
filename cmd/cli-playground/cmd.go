package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/edgeClient"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
	"sync/atomic"
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
		BatchStartCmd, StartCmd, EchoCmd, BenchCmd, ListCmd, Members,
	)

	RootCmd.PersistentFlags().String(FlagServerID, "Node", "")
	RootCmd.PersistentFlags().Uint64(FlagReplicaSet, 0, "")
	RootCmd.PersistentFlags().Int(FlagGossipPort, 800, "")
	RootCmd.PersistentFlags().Bool(FlagBootstrap, false, "")
	RootCmd.PersistentFlags().Int("n", 5, "")
}

var BatchStartCmd = &cobra.Command{
	Use: "batch_start",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		replicaSet, _ := cmd.Flags().GetUint64(FlagReplicaSet)
		gossipPort, _ := cmd.Flags().GetInt(FlagGossipPort)
		raftBootstrap, _ := cmd.Flags().GetBool(FlagBootstrap)
		n, _ := cmd.Flags().GetInt("n")
		for i := 0; i < n; i++ {
			startFunc(fmt.Sprintf("%s.%d", serverID, i), replicaSet, gossipPort, raftBootstrap)
			gossipPort++
			raftBootstrap = false
		}
	},
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
				NewConnectionWorkers: 10,
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
		h := pb.NewSampleServer(&SampleServer{})
		h.Register(Edges[serverID])
		err := Edges[serverID].RunCluster()
		if err != nil {
			fmt.Println(err)
			delete(Edges, serverID)
			return
		}
		Edges[serverID].RunGateway()
		for _, e := range Edges {
			if e.GetServerID() != serverID {
				err = Edges[serverID].JoinCluster(e.Stats().Address)
				if err != nil {
					fmt.Println("Error On Join", err)
				}
				break
			}

		}
	}
}

var ListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		listFunc()
	},
}

func listFunc() {
	var rows []string
	rows = append(rows,
		"ServerID | ReplicaSet | Raft Members | Raft State |  Members | Membership Score | Gateway ",
		"------- | ----------- | ------ | ------ | -------- | ------- | ------",
	)

	for id, s := range Edges {
		edgeStats := s.Stats()
		rows = append(rows,
			fmt.Sprintf("%s | %d | %d | %s | %d | %d | %s(%s)", id,
				edgeStats.ReplicaSet, edgeStats.RaftMembers, edgeStats.RaftState,
				edgeStats.Members, edgeStats.MembershipScore,
				edgeStats.GatewayProtocol, edgeStats.GatewayAddr,
			),
		)
	}

	fmt.Println(columnize.SimpleFormat(rows))
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
		gatewayAddrs := e1.Stats().GatewayAddr
		if len(gatewayAddrs) == 0 {
			fmt.Println("No Gateway Addr", gatewayAddrs)
			return
		}
		ec := edgeClient.NewWebsocket(edgeClient.Config{
			HostPort: fmt.Sprintf("%s", gatewayAddrs[0]),
			Handler: func(m *rony.MessageEnvelope) {
				cmd.Print(m)
			},
			Secure: false,
		})
		defer ec.Close()
		c := pb.NewSampleClient(ec)
		req := pb.PoolEchoRequest.Get()
		defer pb.PoolEchoRequest.Put(req)
		req.Int = tools.RandomInt64(0)
		req.Bool = true
		req.Timestamp = time.Now().UnixNano()

		res, err := c.Echo(req)
		switch err {
		case nil:
		default:
			cmd.Println("Error:", err)
			return
		}
		cmd.Println(res)
	},
}

var BenchCmd = &cobra.Command{
	Use: "bench",
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
		var (
			dd   int64
			cnt  int64
			gcnt int64
		)
		waitGroup := pools.AcquireWaitGroup()
		for i := 0; i < 100; i++ {
			waitGroup.Add(1)
			go func(idx int) {
				defer waitGroup.Done()
				ec := edgeClient.NewWebsocket(edgeClient.Config{
					HostPort: fmt.Sprintf("127.0.0.1:%s", parts[1]),
					Handler: func(m *rony.MessageEnvelope) {
						cmd.Println(m.Constructor, m.RequestID, idx)
						atomic.AddInt64(&gcnt, 1)
					},
				})
				c := pb.NewSampleClient(ec)
				for i := 0; i < 100; i++ {
					startTime := time.Now()
					req := pb.PoolEchoRequest.Get()
					req.Int = tools.RandomInt64(0)
					req.Bool = true
					req.Timestamp = time.Now().UnixNano()
					_, err := c.Echo(req)
					pb.PoolEchoRequest.Put(req)
					if err != nil {
						cmd.Println("Error:", idx, i, err)
						continue
					}
					d := time.Now().Sub(startTime)
					atomic.AddInt64(&cnt, 1)
					atomic.AddInt64(&dd, int64(d))
				}
			}(i)
		}
		waitGroup.Wait()
		cmd.Println("Avg Response:", time.Duration(dd/cnt), cnt, gcnt)
	},
}

var Members = &cobra.Command{
	Use: "members",
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

		var rows []string
		rows = append(rows,
			"ServerID | ReplicaSet   | Raft State |  Raft Port | Gateway ",
			"------- | ----------- | ------ |  -------- | ------- ",
		)

		for _, m := range e1.ClusterMembers() {
			rows = append(rows,
				fmt.Sprintf("%s | %d | %s | %d | %v",
					m.ServerID, m.ReplicaSet, m.RaftState.String(), m.RaftPort, m.GatewayAddr,
				),
			)
		}
		cmd.Println(columnize.SimpleFormat(rows))
	},
}
