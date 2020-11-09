package edge_test

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	dummyGateway "github.com/ronaksoft/rony/gateway/dummy"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
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
   Copyright Ronak Software Group 2020
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
			err := edgeServer.StartCluster()
			if err != nil {
				panic(err)
			}
			if bootstrap {
				time.Sleep(time.Second)
			}
			edgeServer.StartGateway()
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
	pb.RegisterSample(testEnv.Handlers{}, edgeServer)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(1000)
	conn := dummyGateway.Conn{}
	req := pb.EchoRequest{
		Int:       100,
		Bool:      false,
		Timestamp: 32809238402,
	}
	e := &rony.MessageEnvelope{}
	e.Fill(100, pb.C_Echo, &req)

	reqBytes, _ := proto.Marshal(e)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			edgeServer.OnGatewayMessage(&conn, 0, reqBytes)
		}
	})
}
