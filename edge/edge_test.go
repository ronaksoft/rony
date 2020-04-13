package edge_test

import (
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
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
			err := edgeServer.Run()
			if err != nil {
				panic(err)
			}
			if bootstrap {
				time.Sleep(time.Second)
			}
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
	panic("should not be here")
}

func BenchmarkEdgeServerMessageSerial(b *testing.B) {
	edgeServer := testEnv.InitEdgeServerWithWebsocket("Adam", 8080, edge.WithDataPath("./_hdd/adam"))

	req := &pb.ReqSimple1{P1: tools.StrToByte(tools.Int64ToStr(100))}
	envelope := &rony.MessageEnvelope{}
	envelope.RequestID = tools.RandomUint64()
	envelope.Constructor = 101
	envelope.Message, _ = req.Marshal()
	proto := &pb.ProtoMessage{}
	proto.AuthID = 100
	proto.Payload, _ = envelope.Marshal()
	bytes, _ := proto.Marshal()

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(1000)
	conn := testEnv.MockGatewayConn{}

	for i := 0; i < b.N; i++ {
		edgeServer.HandleGatewayMessage(&conn, 0, bytes)
	}
}

func BenchmarkEdgeServerMessageParallel(b *testing.B) {
	edgeServer := testEnv.InitEdgeServerWithWebsocket("Adam", 8080, edge.WithDataPath("./_hdd/adam"))
	req := &pb.ReqSimple1{P1: tools.StrToByte(tools.Int64ToStr(100))}
	envelope := &rony.MessageEnvelope{}
	envelope.RequestID = tools.RandomUint64()
	envelope.Constructor = 101
	envelope.Message, _ = req.Marshal()
	proto := &pb.ProtoMessage{}
	proto.AuthID = 100
	proto.Payload, _ = envelope.Marshal()
	bytes, _ := proto.Marshal()

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(1000)
	conn := testEnv.MockGatewayConn{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			edgeServer.HandleGatewayMessage(&conn, 0, bytes)
		}
	})
}

func BenchmarkEdgeServerWithRaftMessageSerial(b *testing.B) {
	log.SetLevel(log.ErrorLevel)
	initRaftWithWebsocket()

	req := &pb.ReqSimple1{P1: tools.StrToByte(tools.Int64ToStr(100))}
	envelope := &rony.MessageEnvelope{}
	envelope.RequestID = tools.RandomUint64()
	envelope.Constructor = 101
	envelope.Message, _ = req.Marshal()
	proto := &pb.ProtoMessage{}
	proto.AuthID = 100
	proto.Payload, _ = envelope.Marshal()
	bytes, _ := proto.Marshal()

	// b.SetParallelism(1000)
	conn := testEnv.MockGatewayConn{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raftLeader.HandleGatewayMessage(&conn, 0, bytes)
	}
	b.StopTimer()
}

func BenchmarkEdgeServerWithRaftMessageParallel(b *testing.B) {
	log.SetLevel(log.ErrorLevel)
	initRaftWithWebsocket()

	req := &pb.ReqSimple1{P1: tools.StrToByte(tools.Int64ToStr(100))}
	envelope := &rony.MessageEnvelope{}
	envelope.RequestID = tools.RandomUint64()
	envelope.Constructor = 101
	envelope.Message, _ = req.Marshal()
	proto := &pb.ProtoMessage{}
	proto.AuthID = 100
	proto.Payload, _ = envelope.Marshal()
	bytes, _ := proto.Marshal()

	b.SetParallelism(1000)
	conn := testEnv.MockGatewayConn{}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			raftLeader.HandleGatewayMessage(&conn, 0, bytes)
		}
	})
	b.StopTimer()
}
