package edge

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"path/filepath"
	"sync/atomic"
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

type mockGatewayConn struct{}

func (m mockGatewayConn) GetAuthID() int64 {
	return 0
}

func (m mockGatewayConn) GetAuthKey() []byte {
	return nil
}

func (m mockGatewayConn) GetConnID() uint64 {
	return 0
}

func (m mockGatewayConn) GetClientIP() string {
	return ""
}

func (m mockGatewayConn) GetUserID() int64 {
	return 0
}

func (m mockGatewayConn) SendBinary(streamID int64, data []byte) error {
	return nil
}

func (m mockGatewayConn) SetAuthID(int64) {
	return
}

func (m mockGatewayConn) SetAuthKey([]byte) {
	return
}

func (m mockGatewayConn) SetUserID(int64) {
	return
}

func (m mockGatewayConn) Flush() {
	return
}

func (m mockGatewayConn) Persistent() bool {
	return true
}

var (
	receivedMessages int32
	receivedUpdates  int32
	raftServers      map[string]*Server
	raftLeader       *Server
)

type testDispatcher struct {
}

func (t testDispatcher) DispatchUpdate(ctx *DispatchCtx, authID int64, envelope *rony.UpdateEnvelope) {
	atomic.AddInt32(&receivedUpdates, 1)
}

func (t testDispatcher) DispatchMessage(ctx *DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	atomic.AddInt32(&receivedMessages, 1)
}

func (t testDispatcher) DispatchRequest(ctx *DispatchCtx, data []byte) (err error) {
	proto := pb.PoolProtoMessage.Get()
	err = proto.Unmarshal(data)
	if err != nil {
		return
	}
	err = ctx.UnmarshalEnvelope(proto.Payload)
	if err != nil {
		return
	}
	ctx.SetAuthID(proto.AuthID)
	pb.PoolProtoMessage.Put(proto)
	return
}

func initHandlers(edge *Server) {
	edge.AddHandler(100, func(ctx *RequestCtx, in *rony.MessageEnvelope) {
		req := &pb.ReqSimple1{}
		res := &pb.ResSimple1{}
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(in.RequestID, "Invalid", "Proto")
			return
		}
		res.P1 = req.P1
		ctx.PushMessage(ctx.AuthID(), in.RequestID, 201, res)
	})
	edge.AddHandler(101, func(ctx *RequestCtx, in *rony.MessageEnvelope) {
		req := pb.PoolReqSimple1.Get()
		defer pb.PoolReqSimple1.Put(req)
		res := pb.PoolResSimple1.Get()
		defer pb.PoolResSimple1.Put(res)
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(in.RequestID, "Invalid", "Proto")
			return
		}
		res.P1 = req.P1

		u := pb.PoolUpdateSimple1.Get()
		defer pb.PoolUpdateSimple1.Put(u)
		ctx.PushMessage(ctx.AuthID(), in.RequestID, 201, res)
		for i := int64(10); i < 20; i++ {
			ctx.PushUpdate(ctx.AuthID(), i, 301, tools.TimeUnix(), u)
		}
	})
}

func initEdgeServer(serverID string, clientPort int, opts ...Option) *Server {
	opts = append(opts,
		WithWebsocketGateway(websocketGateway.Config{
			NewConnectionWorkers: 1,
			MaxConcurrency:       10,
			MaxIdleTime:          0,
			ListenAddress:        fmt.Sprintf(":%d", clientPort),
		}),
	)
	edge := NewServer(serverID, &testDispatcher{}, opts...)
	initHandlers(edge)

	return edge
}

func initRaft() {
	if raftServers == nil {
		raftServers = make(map[string]*Server)
	}
	ids := []string{"AdamRaft", "EveRaft", "AbelRaft"}
	for idx, id := range ids {
		if raftServers[id] == nil {
			bootstrap := false
			if idx == 0 {
				bootstrap = true
			}
			edgeServer := initEdgeServer(id, 8080+idx,
				WithDataPath(filepath.Join("./_hdd/", id)),
				WithReplicaSet(1, 9080+idx, bootstrap),
				WithGossipPort(7080+idx),
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
	edgeServer := initEdgeServer("Adam", 8080, WithDataPath("./_hdd/adam"))

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
	conn := mockGatewayConn{}

	for i := 0; i < b.N; i++ {
		edgeServer.onGatewayMessage(&conn, 0, bytes)
	}
}

func BenchmarkEdgeServerMessageParallel(b *testing.B) {
	edgeServer := initEdgeServer("Adam", 8080, WithDataPath("./_hdd/adam"))
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
	conn := mockGatewayConn{}
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			edgeServer.onGatewayMessage(&conn, 0, bytes)
		}
	})
}

func BenchmarkEdgeServerWithRaftMessageSerial(b *testing.B) {
	log.SetLevel(log.ErrorLevel)
	initRaft()

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
	conn := mockGatewayConn{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raftLeader.onGatewayMessage(&conn, 0, bytes)
	}
	b.StopTimer()
}

func BenchmarkEdgeServerWithRaftMessageParallel(b *testing.B) {
	log.SetLevel(log.ErrorLevel)
	initRaft()

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
	conn := mockGatewayConn{}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			raftLeader.onGatewayMessage(&conn, 0, bytes)
		}
	})
	b.StopTimer()
}
