package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	tcpGateway "github.com/ronaksoft/rony/gateway/tcp"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2019 - Oct - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var Edges map[string]*edge.Server

func init() {
	Edges = make(map[string]*edge.Server)
	RootCmd.AddCommand(
		BatchStartCmd, StartCmd, EchoCmd, BenchCmd, ListCmd, Members,
		Trace, MemProf,
	)

	RootCmd.PersistentFlags().String(FlagServerID, "", "")
	RootCmd.PersistentFlags().Uint64(FlagReplicaSet, 0, "")
	RootCmd.PersistentFlags().Int(FlagGossipPort, 800, "")
	RootCmd.PersistentFlags().Bool(FlagBootstrap, false, "")
	RootCmd.PersistentFlags().Int("n", 5, "")
	RootCmd.PersistentFlags().Int(FlagConcurrency, 1000, "")
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
			edge.WithTcpGateway(tcpGateway.Config{
				Concurrency:   10000,
				MaxIdleTime:   0,
				ListenAddress: "0.0.0.0:0",
			}),
			edge.WithDataPath(filepath.Join("./_hdd", serverID)),
			edge.WithGossipPort(port),
		)

		if replicaSet != 0 {
			opts = append(opts, edge.WithReplicaSet(replicaSet, port*10, bootstrap))
		}

		Edges[serverID] = edge.NewServer(serverID, &dispatcher{}, opts...)
		pb.RegisterSample(&SampleServer{}, Edges[serverID])

		err := Edges[serverID].StartCluster()
		if err != nil {
			fmt.Println(err)
			delete(Edges, serverID)
			return
		}
		Edges[serverID].StartGateway()
		for _, e := range Edges {
			if e.GetServerID() != serverID {
				_, err = Edges[serverID].JoinCluster(e.Stats().Address)
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
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		if len(serverID) == 0 {
			fmt.Println("Needs ServerID, e.g. echo --serverID First.01")
			return
		}
		e1 := Edges[serverID]
		if e1 == nil {
			fmt.Println("Invalid Args")
			return
		}
		gatewayAddrs := e1.Stats().GatewayAddr
		if len(gatewayAddrs) == 0 {
			fmt.Println("No Gateway Addr", gatewayAddrs)
			return
		}
		ec := edgec.NewWebsocket(edgec.WebsocketConfig{
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
		var cnt int64
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				res, err := c.Echo(req)
				switch err {
				case nil:
					atomic.AddInt64(&cnt, 1)
					cmd.Println(res)
				default:
					cmd.Println("Error:", err)
					return
				}

				wg.Done()
			}()
		}
		wg.Wait()
		cmd.Println("N:", cnt)
	},
}

var BenchCmd = &cobra.Command{
	Use: "bench",
	Run: func(cmd *cobra.Command, args []string) {
		concurrency, _ := cmd.Flags().GetInt(FlagConcurrency)
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		if len(serverID) == 0 {
			fmt.Println("Needs ServerID, e.g. echo --serverID First.01")
			return
		}
		e1 := Edges[serverID]
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
			dd         int64
			reqCnt     int64
			resCnt     int64
			delayedCnt int64
		)

		f, err := os.Create("rony-bench.trace")
		if err != nil {
			cmd.Println(err)
			return
		}
		err = trace.Start(f)
		if err != nil {
			cmd.Println("Error On TraceStart", err)
		}
		defer trace.Stop()

		benchStart := tools.CPUTicks()
		waitGroup := pools.AcquireWaitGroup()
		for i := 0; i < concurrency; i++ {
			waitGroup.Add(1)
			go func(idx int) {
				defer waitGroup.Done()
				ec := edgec.NewWebsocket(edgec.WebsocketConfig{
					HostPort: fmt.Sprintf("127.0.0.1:%s", parts[1]),
					Handler: func(m *rony.MessageEnvelope) {
						atomic.AddInt64(&delayedCnt, 1)
					},
				})
				c := pb.NewSampleClient(ec)
				for i := 0; i < 10; i++ {
					startTime := tools.CPUTicks()
					req := pb.PoolEchoRequest.Get()
					req.Int = tools.RandomInt64(0)
					req.Bool = true
					req.Timestamp = tools.CPUTicks()
					atomic.AddInt64(&reqCnt, 1)
					_, err := c.Echo(req)
					pb.PoolEchoRequest.Put(req)
					if err != nil {
						cmd.Println("Error:", idx, i, err)
						continue
					}
					d := time.Duration(tools.CPUTicks() - startTime)
					atomic.AddInt64(&resCnt, 1)
					atomic.AddInt64(&dd, int64(d))
				}
			}(i)
		}
		waitGroup.Wait()
		cmd.Println("Bench Test Time:", time.Duration(tools.CPUTicks()-benchStart))
		cmd.Println("Avg Response:", time.Duration(dd/resCnt))
		cmd.Println("Reqs / Res / Delayed", reqCnt, resCnt, delayedCnt)
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

var Trace = &cobra.Command{
	Use: "trace",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Create("rony-bench.trace")
		if err != nil {
			cmd.Println(err)
			return
		}
		err = trace.Start(f)
		if err != nil {
			cmd.Println("Error On TraceStart", err)
		}
		time.Sleep(time.Second * 10)
		trace.Stop()
		_ = f.Close()
	},
}

var MemProf = &cobra.Command{
	Use: "mem-prof",
	Run: func(cmd *cobra.Command, args []string) {
		f1, err := os.Create("rony-mem.out")
		if err != nil {
			cmd.Println(err)
			return
		}
		pprof.Lookup("heap").WriteTo(f1, 0)
		_ = f1.Close()

		f2, err := os.Create("rony-alloc.out")
		if err != nil {
			cmd.Println(err)
			return
		}
		pprof.Lookup("allocs").WriteTo(f2, 0)
		_ = f2.Close()
	},
}
