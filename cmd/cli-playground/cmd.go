package main

import (
	context2 "context"
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/cmd/cli-playground/msg"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/pkg/profile"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
	"sync"
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

var Edges map[string]*rony.EdgeServer

func init() {
	Edges = make(map[string]*rony.EdgeServer)
	RootCmd.AddCommand(
		StartCmd, StopCmd, ListCmd, JoinCmd, EchoCmd, BenchCmd, ClusterMessageCmd,
		MembersCmd,
	)

	RootCmd.PersistentFlags().String(FlagServerID, "Node", "")
	RootCmd.PersistentFlags().String(FlagReplicaSet, "", "")
	RootCmd.PersistentFlags().Int(FlagGossipPort, 801, "")
	RootCmd.PersistentFlags().Bool(FlagBootstrap, false, "")

}

var StartCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		replicaSet, _ := cmd.Flags().GetUint32(FlagReplicaSet)
		gossipPort, _ := cmd.Flags().GetInt(FlagGossipPort)
		raftBootstrap, _ := cmd.Flags().GetBool(FlagBootstrap)
		startFunc(serverID, replicaSet, gossipPort, raftBootstrap)
	},
}

func startFunc(serverID string, replicaSet uint32, port int, bootstrap bool) {
	if _, ok := Edges[serverID]; !ok {
		opts := make([]rony.Option, 0)
		opts = append(opts,
			rony.WithWebsocketGateway(websocketGateway.Config{
				NewConnectionWorkers: 1,
				MaxConcurrency:       1000,
				MaxIdleTime:          0,
				ListenAddress:        "0.0.0.0:0",
			}),
			rony.WithDataPath(filepath.Join("./_hdd", serverID)),
			rony.WithGossipPort(port),
		)

		if replicaSet != 0 {
			opts = append(opts, rony.WithReplicaSet(replicaSet, port*10, bootstrap))
		}

		Edges[serverID] = rony.NewEdgeServer(serverID, &dispatcher{}, opts...)
		Edges[serverID].AddHandler(msg.C_EchoRequest, GenEchoHandler(serverID))
		Edges[serverID].AddHandler(msg.C_AskRequest, GenAskHandler(serverID))
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
}

var StopCmd = &cobra.Command{
	Use: "stop",
	Run: func(cmd *cobra.Command, args []string) {
		for id, s := range Edges {
			s.Shutdown()
			fmt.Println(id, "Shutdown!")
		}
	},
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
		"ServerID | Raft Members | Raft State |  Members | Membership State | Gateway ",
		"------- | ------ | ------ | -------- | ------- | ------",
	)

	for id, s := range Edges {
		edgeStats := s.Stats()
		rows = append(rows,
			fmt.Sprintf("%s | %d | %s | %d | %d | %s(%s)", id,
				edgeStats.RaftMembers, edgeStats.RaftState,
				edgeStats.Members, edgeStats.MembershipScore,
				edgeStats.GatewayProtocol, edgeStats.GatewayAddr,
			),
		)
	}

	fmt.Println(columnize.SimpleFormat(rows))
}

var JoinCmd = &cobra.Command{
	Use: "join",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Needs Two ServerID, e.g. join First.01 First.02")
			return
		}
		joinFunc(args[0], args[1])
	},
}

func joinFunc(serverID1, serverID2 string) {
	e1 := Edges[serverID1]
	e2 := Edges[serverID2]
	if e1 == nil || e2 == nil {
		fmt.Println("Invalid Args")
		return
	}
	e1Stats := e1.Stats()
	e2Stats := e2.Stats()
	fmt.Println("Joining ", e1Stats.Address, "--->", e2Stats.Address)
	err := e1.JoinCluster(e2Stats.Address)

	if err != nil {
		fmt.Println(err)
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
		parts := strings.Split(gatewayAddr, ":")
		if len(parts) != 2 {
			fmt.Println("Invalid Gateway Addr", gatewayAddr)
			return
		}
		conn, _, _, err := ws.Dial(context2.Background(), fmt.Sprintf("ws://127.0.0.1:%s", parts[1]))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		req := msg.PoolEchoRequest.Get()
		defer msg.PoolEchoRequest.Put(req)
		req.Int = tools.RandomInt64(0)
		req.Bool = true
		req.Timestamp = time.Now().UnixNano()

		envelope := &rony.MessageEnvelope{
			Constructor: msg.C_EchoRequest,
			RequestID:   tools.RandomUint64(),
			Message:     nil,
		}
		envelope.Message, _ = req.Marshal()
		proto := &msg.ProtoMessage{}
		proto.AuthID = 1000
		proto.Payload, _ = envelope.Marshal()
		bytes, _ := proto.Marshal()
		err = wsutil.WriteClientBinary(conn, bytes)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Sent:", proto.AuthID, envelope.RequestID, msg.ConstructorNames[envelope.Constructor])

		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		resBytes, err := wsutil.ReadServerBinary(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = proto.Unmarshal(resBytes)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = envelope.Unmarshal(proto.Payload)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Received:", proto.AuthID, envelope.RequestID, msg.ConstructorNames[envelope.Constructor])
		switch envelope.Constructor {
		case rony.C_Error:
			res := rony.Error{}
			err = res.Unmarshal(envelope.Message)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Error:", res.Code, res.Items)
		case msg.C_EchoResponse:
			res := msg.EchoResponse{}
			err = res.Unmarshal(envelope.Message)
			if err != nil {
				fmt.Println(err)
				return
			}
			if res.Bool != req.Bool || res.Int != req.Int {
				fmt.Println("ERROR!!! In Response")
			}
			fmt.Println("Delay:", time.Duration(res.Delay))
		}
	},
}

var BenchCmd = &cobra.Command{
	Use: "bench",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Needs ServerID, e.g. echo First.01 <Count>")
			return
		}
		e1 := Edges[args[0]]
		if e1 == nil {
			fmt.Println("Invalid Args")
			return
		}
		var count int32 = 100
		if len(args) == 2 {
			count = tools.StrToInt32(args[1])
		}
		gatewayAddr := e1.Stats().GatewayAddr
		parts := strings.Split(gatewayAddr, ":")
		if len(parts) != 2 {
			fmt.Println("Invalid Gateway Addr", gatewayAddr)
			return
		}

		profiler := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
		startTime := time.Now()
		waitGroup := &sync.WaitGroup{}
		for i := int64(1); i <= int64(count); i++ {
			waitGroup.Add(1)
			go func(i int64) {
				benchRoutine(i*1000, int(count), parts[1])
				waitGroup.Done()
			}(i)
			time.Sleep(time.Millisecond)
		}
		waitGroup.Wait()
		d := time.Now().Sub(startTime)
		t := count * count
		if count > 50 {
			t = count * 50
		}
		fmt.Println("Total Time:", d, ", ", t)
		fmt.Println("Avg:", int(float64(t)/d.Seconds()))
		profiler.Stop()
	},
}

func benchRoutine(authID int64, count int, port string) {
	if count > 50 {
		count = 50
	}
	conn, _, _, err := ws.Dial(context2.Background(), fmt.Sprintf("ws://127.0.0.1:%s", port))
	if err != nil {
		fmt.Println(authID, "Connect", err)
		return
	}
	defer conn.Close()
	for i := 0; i < count; i++ {
		req := &msg.EchoRequest{
			Int:       tools.RandomInt64(0),
			Bool:      true,
			Timestamp: time.Now().UnixNano(),
		}
		reqBytes, _ := req.Marshal()
		envelope := &rony.MessageEnvelope{
			Constructor: msg.C_EchoRequest,
			RequestID:   tools.RandomUint64(),
			Message:     reqBytes,
		}
		proto := &msg.ProtoMessage{
			AuthID: authID,
		}
		proto.Payload, _ = envelope.Marshal()
		bytes, _ := proto.Marshal()
		err = wsutil.WriteClientBinary(conn, bytes)
		if err != nil {
			fmt.Println(authID, "Write:", err)
			return
		}

		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		resBytes, err := wsutil.ReadServerBinary(conn)
		if err != nil {
			fmt.Println(authID, "Read", err)
			return
		}
		err = proto.Unmarshal(resBytes)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = envelope.Unmarshal(proto.Payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch envelope.Constructor {
		case rony.C_Error:
			res := rony.Error{}
			err = res.Unmarshal(envelope.Message)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("ERROR!!! In Response", authID)
		case msg.C_EchoResponse:
			res := msg.EchoResponse{}
			err = res.Unmarshal(envelope.Message)
			if err != nil {
				fmt.Println(err)
				return
			}
			if res.Bool != req.Bool || res.Int != req.Int {
				fmt.Println("MisMatch!!! In Response", authID)
			}
			// fmt.Println("Delay:", time.Duration(res.Delay))
		default:
			fmt.Println("UNKNOWN!!!!")
		}

	}

}

var ClusterMessageCmd = &cobra.Command{
	Use: "clusterMessage",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Needs ServerID, e.g. echo First Second <Count>")
			return
		}
		e1 := Edges[args[0]]
		e2 := Edges[args[1]]
		if e1 == nil || e2 == nil {
			fmt.Println("Invalid Args")
			return
		}
		var count int32 = 100
		if len(args) == 3 {
			count = tools.StrToInt32(args[2])
		}
		gatewayAddr := e1.Stats().GatewayAddr
		parts := strings.Split(gatewayAddr, ":")
		if len(parts) != 2 {
			fmt.Println("Invalid Gateway Addr", gatewayAddr)
			return
		}


		startTime := time.Now()
		waitGroup := &sync.WaitGroup{}
		for i := int64(1); i <= int64(count); i++ {
			waitGroup.Add(1)
			go func(i int64) {
				benchClusterMessage(i*1000, int(count), parts[1])
				waitGroup.Done()
			}(i)
			time.Sleep(time.Millisecond)
		}
		waitGroup.Wait()
		d := time.Now().Sub(startTime)
		t := count * count
		if count > 50 {
			t = count * 50
		}
		fmt.Println("Total Time:", d, ", ", t)
		fmt.Println("Avg:", int(float64(t)/d.Seconds()))
	},
}

func benchClusterMessage(authID int64, count int, port string, serverID string) {
	if count > 50 {
		count = 50
	}
	conn, _, _, err := ws.Dial(context2.Background(), fmt.Sprintf("ws://127.0.0.1:%s", port))
	if err != nil {
		fmt.Println(authID, "Connect", err)
		return
	}
	defer conn.Close()
	for i := 0; i < count; i++ {
		req := &msg.AskRequest{
			ServerID: serverID,
		}
		reqBytes, _ := req.Marshal()
		envelope := &rony.MessageEnvelope{
			Constructor: msg.C_AskRequest,
			RequestID:   tools.RandomUint64(),
			Message:     reqBytes,
		}
		proto := &msg.ProtoMessage{
			AuthID: authID,
		}
		proto.Payload, _ = envelope.Marshal()
		bytes, _ := proto.Marshal()
		err = wsutil.WriteClientBinary(conn, bytes)
		if err != nil {
			fmt.Println(authID, "Write:", err)
			return
		}

		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		resBytes, err := wsutil.ReadServerBinary(conn)
		if err != nil {
			fmt.Println(authID, "Read", err)
			return
		}
		err = proto.Unmarshal(resBytes)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = envelope.Unmarshal(proto.Payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch envelope.Constructor {
		case rony.C_Error:
			res := rony.Error{}
			err = res.Unmarshal(envelope.Message)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("ERROR!!! In Response", authID)
		case msg.C_AskResponse:
			res := msg.AskResponse{}
			err = res.Unmarshal(envelope.Message)
			if err != nil {
				fmt.Println(err)
				return
			}

			// fmt.Println("Delay:", time.Duration(res.Delay))
		default:
			fmt.Println("UNKNOWN!!!!")
		}

	}

}

var MembersCmd = &cobra.Command{
	Use: "members",
	Run: func(cmd *cobra.Command, args []string) {
		serverID, _ := cmd.Flags().GetString(FlagServerID)
		s, ok := Edges[serverID]
		if !ok {
			return
		}

		var rows []string
		rows = append(rows,
			"ServerID | RaftState | ReplicaSet | Raft Port | Gossip Addr | Gateway Addr",
			"--------- | ------------ | ------------ | ---------------- | ----------- | ---------",
		)

		for _, cm := range s.ClusterMembers() {
			rows = append(rows,
				fmt.Sprintf("%s | %s | %d | %d | %s(%d) | %s",
					cm.ServerID, cm.RaftState.String(), cm.ReplicaSet, cm.RaftPort, cm.Addr.String(), cm.Port, cm.GatewayAddr,
				),
			)
		}

		fmt.Println(columnize.SimpleFormat(rows))

	},
}
