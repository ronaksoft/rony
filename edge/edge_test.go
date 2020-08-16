package edge_test

import (
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"google.golang.org/protobuf/proto"
	"path/filepath"
	"testing"
	"time"
)

/*
   Creation Time: 2020 - Mar - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	raftServers map[string]*edge.Server
	raftLeader  *edge.Server
)

func initRaftWithWebsocket() {
	if raftServers == nil {
		raftServers = make(map[string]*edge.Server)
	}
	ids := []string{"AdamRaft", "EveRaft", "AbelRaft"}
	for idx, id := range ids {
		if raftServers[id] == nil {
			bootstrap := false
			if idx == 0 {
				bootstrap = true
			}
			edgeServer := testEnv.InitEdgeServerWithWebsocket(id, 8080+idx,
				edge.WithDataPath(filepath.Join("./_hdd/", id)),
				edge.WithReplicaSet(1, 9080+idx, bootstrap),
				edge.WithGossipPort(7080+idx),
			)
			err := edgeServer.RunCluster()
			if err != nil {
				panic(err)
			}
			if bootstrap {
				time.Sleep(time.Second)
			}
			edgeServer.RunGateway()
			raftServers[id] = edgeServer
		}
	}
	time.Sleep(time.Second)
	for _, id := range ids {
		if raftServers[id].Stats().RaftState == "Leader" {
			raftLeader = raftServers[id]
			return
		}
	}
	panic("BUG!! should not be here")
}

func BenchmarkStandaloneSerial(b *testing.B) {
	edgeServer := testEnv.InitEdgeServerWithWebsocket("Adam", 8080, edge.WithDataPath("./_hdd/adam"))
	sampleServer := pb.NewSampleServer(testEnv.Handlers{})
	sampleServer.Register(edgeServer)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(1000)
	conn := testEnv.MockGatewayConn{}
	req := pb.EchoRequest{
		Int:       100,
		Bool:      false,
		Timestamp: 32809238402,
	}
	e := &rony.MessageEnvelope{}
	e.Fill(100, pb.C_Echo, &req)

	reqBytes, _ := proto.Marshal(e)
	for i := 0; i < b.N; i++ {
		edgeServer.HandleGatewayMessage(&conn, 0, reqBytes)
	}
}
