package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/internal/testEnv/pb/service"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"sort"
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
		DemoStartCmd, StartCmd, BenchCmd, ListCmd, Members,
		EchoCmd, EchoTunnelCmd, SetCmd, GetCmd,
		Trace, MemProf, LogLevel, ClusterCmd,
	)

	RootCmd.PersistentFlags().Int(FlagLogLevel, 0, "")
	RootCmd.PersistentFlags().String(FlagServerID, "", "")
	RootCmd.PersistentFlags().String(FlagTargetID, "", "")
	RootCmd.PersistentFlags().Uint64(FlagReplicaSet, 0, "")
	RootCmd.PersistentFlags().Int(FlagGossipPort, 800, "")
	RootCmd.PersistentFlags().Bool(FlagBootstrap, false, "")
	RootCmd.PersistentFlags().Int("n", 5, "")
	RootCmd.PersistentFlags().Int(FlagConcurrency, 1000, "")
	RootCmd.PersistentFlags().Int(FlagDemoReplicaFactor, 1, "")
	RootCmd.PersistentFlags().Int(FlagDemoReplica, 3, "")
	RootCmd.PersistentFlags().String(FlagKey, "k", "")
	RootCmd.PersistentFlags().String(FlagValue, "v", "")

}

var DemoStartCmd = &cobra.Command{
	Use: "demo_cluster_start",
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := cmd.Flags().GetInt(FlagDemoReplicaFactor)
		s, _ := cmd.Flags().GetInt(FlagDemoReplica)
		if s < 1 {
			s = 1
		}
		if s > 28 {
			s = 28
		}

		a := 'A'
		gossipPort := 700
		for i := 0; i < s; i++ {
			replicaSet := uint64(i + 1)
			for j := 0; j < n; j++ {
				startFunc(cmd, fmt.Sprintf("%c.%d", a, j), replicaSet, gossipPort, j == 0)
				gossipPort += 10
			}
			a++
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

		startFunc(cmd, serverID, replicaSet, gossipPort, raftBootstrap)
	},
}

func startFunc(cmd *cobra.Command, serverID string, replicaSet uint64, port int, bootstrap bool) {
	if _, ok := Edges[serverID]; !ok {
		opts := make([]edge.Option, 0)
		opts = append(opts,
			edge.WithTcpGateway(edge.TcpGatewayConfig{
				Concurrency:   10000,
				MaxIdleTime:   0,
				ListenAddress: "0.0.0.0:0",
			}),
		)

		if replicaSet != 0 {
			opts = append(opts,
				edge.WithDataDir(fmt.Sprintf("./_hdd/%s", serverID)),
				edge.WithGossipCluster(edge.GossipClusterConfig{
					ServerID:   []byte(serverID),
					Bootstrap:  bootstrap,
					ReplicaSet: replicaSet,
					GossipPort: port,
				}),
				edge.WithUdpTunnel(edge.UdpTunnelConfig{
					ListenAddress: "0.0.0.0:0",
					MaxBodySize:   4096,
				}),
			)
		}

		opts = append(opts, edge.WithDispatcher(&dispatcher{}))
		Edges[serverID] = edge.NewServer(serverID, opts...)
		service.RegisterSample(&SampleServer{es: Edges[serverID]}, Edges[serverID])

		Edges[serverID].Start()
		for _, e := range Edges {
			if e.GetServerID() != serverID {
				_, err := Edges[serverID].JoinCluster(e.Stats().Address)
				if err != nil {
					cmd.Println("Error On Join", err)
				}
				break
			}
		}
	}
}

var ListCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		listFunc(cmd)
	},
}

func listFunc(cmd *cobra.Command) {
	var rows []string
	rows = append(rows,
		"ServerID | ReplicaSet | Raft Members | Raft State |  Members | Tunnel | Gateway ",
		"------- | ----------- | ------ | ------ | -------- | ------- | ------",
	)

	ea := make([]*edge.Server, 0, len(Edges))
	for _, s := range Edges {
		ea = append(ea, s)
	}
	sort.Slice(ea, func(i, j int) bool {
		return strings.Compare(ea[i].GetServerID(), ea[j].GetServerID()) < 0
	})
	for _, s := range ea {
		edgeStats := s.Stats()
		rows = append(rows,
			fmt.Sprintf("%s | %d | %d | %s | %s(%s)", s.GetServerID(),
				edgeStats.ReplicaSet,
				edgeStats.Members,
				edgeStats.TunnelAddr,
				edgeStats.GatewayProtocol.String(), edgeStats.GatewayAddr,
			),
		)
	}

	cmd.Println(columnize.SimpleFormat(rows))
}

var EchoCmd = &cobra.Command{
	Use: "echo",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		n, _ := cmd.Flags().GetInt("n")
		if len(serverID) == 0 {
			cmd.Println("Needs ServerID, e.g. echo --serverID First.01")
			return
		}
		e1 := Edges[serverID]
		if e1 == nil {
			cmd.Println("Invalid Args")
			return
		}
		gatewayAddrs := e1.Stats().GatewayAddr
		if len(gatewayAddrs) == 0 {
			cmd.Println("No Gateway Addr", gatewayAddrs)
			return
		}
		ec := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort: gatewayAddrs[0],
			Handler: func(m *rony.MessageEnvelope) {
				cmd.Println("Uncaught Response", registry.ConstructorName(m.Constructor), m.RequestID)
			},
			Secure:         false,
			ContextTimeout: time.Second * 3,
		})
		err := ec.Start()
		if err != nil {
			cmd.Println(err)
			return
		}
		defer ec.Close()
		c := service.NewSampleClient(ec)
		req := service.PoolEchoRequest.Get()
		defer service.PoolEchoRequest.Put(req)
		req.Int = tools.RandomInt64(0)
		req.Timestamp = tools.NanoTime()
		var cnt int64
		wg := sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				res, err := c.Echo(req)
				switch err {
				case nil:
					atomic.AddInt64(&cnt, 1)
					cmd.Println(res)
				default:
					cmd.Println("Error:", err)
				}

				wg.Done()
			}()
		}
		wg.Wait()
		cmd.Println("N:", cnt)
	},
}

var EchoTunnelCmd = &cobra.Command{
	Use: "echo-tunnel",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		replicaSet, _ := cmd.Flags().GetUint64(FlagReplicaSet)
		n, _ := cmd.Flags().GetInt("n")
		if len(serverID) == 0 {
			cmd.Println("Needs ServerID, e.g. echo --serverID First.01")
			return
		}
		e1 := Edges[serverID]
		if e1 == nil {
			cmd.Println("Invalid Args")
			return
		}
		gatewayAddrs := e1.Stats().GatewayAddr
		if len(gatewayAddrs) == 0 {
			cmd.Println("No Gateway Addr", gatewayAddrs)
			return
		}
		ec := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort: gatewayAddrs[0],
			Handler: func(m *rony.MessageEnvelope) {
				cmd.Println("Uncaught Response", registry.ConstructorName(m.Constructor), m.RequestID)
			},
			Secure:          false,
			RequestTimeout:  time.Second * 10,
			RequestMaxRetry: 1,
		})
		err := ec.Start()
		if err != nil {
			cmd.Println(err)
			return
		}
		defer ec.Close()
		c := service.NewSampleClient(ec)

		var cnt int64
		wg := sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				req := service.PoolEchoRequest.Get()
				defer service.PoolEchoRequest.Put(req)
				req.Int = tools.RandomInt64(0)
				req.Timestamp = tools.NanoTime()
				req.ReplicaSet = replicaSet
				res, err := c.EchoTunnel(req)
				switch err {
				case nil:
					atomic.AddInt64(&cnt, 1)
					cmd.Println("Res:", res)
				default:
					cmd.Println("Error:", err)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		cmd.Println("N:", cnt)
	},
}

var SetCmd = &cobra.Command{
	Use: "set",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		n, _ := cmd.Flags().GetInt("n")
		key, _ := cmd.Flags().GetString("key")
		val, _ := cmd.Flags().GetString("val")
		if len(serverID) == 0 {
			cmd.Println("Needs ServerID, e.g. echo --serverID First.01")
			return
		}
		e1 := Edges[serverID]
		if e1 == nil {
			cmd.Println("Invalid Args")
			return
		}
		gatewayAddrs := e1.Stats().GatewayAddr
		if len(gatewayAddrs) == 0 {
			cmd.Println("No Gateway Addr", gatewayAddrs)
			return
		}
		ec := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort: gatewayAddrs[0],
			Handler: func(m *rony.MessageEnvelope) {
				cmd.Println("Uncaught Response", registry.ConstructorName(m.Constructor), m.RequestID)
			},
			Secure:         false,
			ContextTimeout: time.Second * 3,
		})
		err := ec.Start()
		if err != nil {
			cmd.Println(err)
			return
		}
		defer ec.Close()
		c := service.NewSampleClient(ec)
		req := service.PoolSetRequest.Get()
		defer service.PoolSetRequest.Put(req)
		req.Key = tools.S2B(key)
		req.Value = tools.S2B(val)

		var cnt int64
		wg := sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				res, err := c.Set(req)
				switch err {
				case nil:
					atomic.AddInt64(&cnt, 1)
					cmd.Println(res)
				default:
					cmd.Println("Error:", err)
				}

				wg.Done()
			}()
		}
		wg.Wait()
		cmd.Println("N:", cnt)
	},
}

var GetCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		n, _ := cmd.Flags().GetInt("n")
		key, _ := cmd.Flags().GetString("key")
		if len(serverID) == 0 {
			cmd.Println("Needs ServerID, e.g. echo --serverID First.01")
			return
		}
		e1 := Edges[serverID]
		if e1 == nil {
			cmd.Println("Invalid Args")
			return
		}
		gatewayAddrs := e1.Stats().GatewayAddr
		if len(gatewayAddrs) == 0 {
			cmd.Println("No Gateway Addr", gatewayAddrs)
			return
		}
		ec := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort: gatewayAddrs[0],
			Handler: func(m *rony.MessageEnvelope) {
				cmd.Println("Uncaught Response", registry.ConstructorName(m.Constructor), m.RequestID)
			},
			Secure:         false,
			ContextTimeout: time.Second * 3,
		})
		err := ec.Start()
		if err != nil {
			cmd.Println(err)
			return
		}
		defer ec.Close()
		c := service.NewSampleClient(ec)
		req := service.PoolGetRequest.Get()
		defer service.PoolGetRequest.Put(req)
		req.Key = tools.S2B(key)

		var cnt int64
		wg := sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				res, err := c.Get(req)
				switch err {
				case nil:
					atomic.AddInt64(&cnt, 1)
					cmd.Println(res)
				default:
					cmd.Println("Error:", err)
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
			cmd.Println("Needs ServerID, e.g. echo --serverID First.01")
			return
		}
		e1 := Edges[serverID]
		if e1 == nil {
			cmd.Println("Invalid Args")
			return
		}
		gatewayAddr := e1.Stats().GatewayAddr
		if len(gatewayAddr) == 0 {
			cmd.Println("No Gateway Addr", gatewayAddr)
			return
		}
		parts := strings.Split(gatewayAddr[0], ":")
		if len(parts) != 2 {
			cmd.Println("Invalid Gateway Addr", gatewayAddr)
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
					SeedHostPort: fmt.Sprintf("127.0.0.1:%s", parts[1]),
					Handler: func(m *rony.MessageEnvelope) {
						atomic.AddInt64(&delayedCnt, 1)
					},
				})
				err := ec.Start()
				if err != nil {
					cmd.Println(err)
					return
				}
				c := service.NewSampleClient(ec)
				for i := 0; i < 10; i++ {
					startTime := tools.CPUTicks()
					req := service.PoolEchoRequest.Get()
					req.Int = tools.RandomInt64(0)
					req.Timestamp = tools.CPUTicks()
					atomic.AddInt64(&reqCnt, 1)
					_, err := c.Echo(req)
					service.PoolEchoRequest.Put(req)
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
			cmd.Println("Needs ServerID, e.g. echo First.01")
			return
		}
		e1 := Edges[args[0]]
		if e1 == nil {
			cmd.Println("Invalid Args")
			return
		}

		var rows []string
		rows = append(rows,
			"ServerID | ReplicaSet   | Raft State |  Raft Port | Gateway ",
			"------- | ----------- | ------ |  -------- | ------- ",
		)

		for _, m := range e1.Cluster().Members() {
			rows = append(rows,
				fmt.Sprintf("%s | %d |  %v",
					m.ServerID(), m.ReplicaSet(), m.GatewayAddr(),
				),
			)
		}
		cmd.Println(columnize.SimpleFormat(rows))
	},
}

var ClusterCmd = &cobra.Command{
	Use: "cluster-info",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		replicaSet, _ := cmd.Flags().GetUint64(FlagReplicaSet)
		if len(serverID) == 0 {
			cmd.Println("Needs ServerID, e.g. echo --serverID First.01")
			return
		}
		e1 := Edges[serverID]
		if e1 == nil {
			cmd.Println("Invalid Args")
			return
		}
		gatewayAddrs := e1.Stats().GatewayAddr
		if len(gatewayAddrs) == 0 {
			cmd.Println("No Gateway Addr", gatewayAddrs)
			return
		}
		ec := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort: gatewayAddrs[0],
			Handler: func(m *rony.MessageEnvelope) {
				cmd.Println("Uncaught Response", registry.ConstructorName(m.Constructor), m.RequestID)
			},
			Secure:          false,
			RequestTimeout:  time.Second * 10,
			RequestMaxRetry: 1,
		})
		err := ec.Start()
		if err != nil {
			cmd.Println(err)
			return
		}
		defer ec.Close()

		var (
			edges *rony.Edges
		)
		if replicaSet != 0 {
			edges, err = ec.ClusterInfo(replicaSet)
		} else {
			edges, err = ec.ClusterInfo()
		}

		if err != nil {
			cmd.Println(err)
			return
		}
		for idx, n := range edges.Nodes {
			cmd.Println(fmt.Sprintf("%d: %s[%d] -- %v", idx, n.ServerID, n.ReplicaSet, n.HostPorts))
		}

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

var LogLevel = &cobra.Command{
	Use: "log-level",
	Run: func(cmd *cobra.Command, args []string) {
		level, _ := cmd.Flags().GetInt(FlagLogLevel)
		rony.SetLogLevel(level)
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
