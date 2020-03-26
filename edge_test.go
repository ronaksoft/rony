package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/cmd/cli-playground/msg"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"sync/atomic"
	"testing"
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
)

type testDispatcher struct {
}

func (t testDispatcher) DispatchUpdate(ctx *DispatchCtx, authID int64, envelope *UpdateEnvelope) {
	atomic.AddInt32(&receivedUpdates, 1)
}

func (t testDispatcher) DispatchMessage(ctx *DispatchCtx, authID int64, envelope *MessageEnvelope) {
	atomic.AddInt32(&receivedMessages, 1)
}

func (t testDispatcher) DispatchRequest(ctx *DispatchCtx, data []byte) (err error) {
	proto := &msg.ProtoMessage{}
	err = proto.Unmarshal(data)
	if err != nil {
		return
	}
	err = ctx.UnmarshalEnvelope(proto.Payload)
	if err != nil {
		return
	}
	ctx.SetAuthID(proto.AuthID)
	return
}

func initHandlers(edge *EdgeServer) {
	edge.AddHandler(100, func(ctx *RequestCtx, in *MessageEnvelope) {
		req := &pb.ReqSimple1{}
		res := &pb.ResSimple1{}
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(in.RequestID, "Invalid", "Proto")
			return
		}
		res.P1 = tools.StrToByte(req.P1)
		ctx.PushMessage(ctx.AuthID(), in.RequestID, 201, res)
	})
	edge.AddHandler(101, func(ctx *RequestCtx, in *MessageEnvelope) {
		req := &pb.ReqSimple1{}
		res := &pb.ResSimple1{}
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(in.RequestID, "Invalid", "Proto")
			return
		}
		res.P1 = tools.StrToByte(req.P1)

		ctx.PushMessage(ctx.AuthID(), in.RequestID, 201, res)
		for i := int64(10); i < 20; i++ {
			ctx.PushUpdate(ctx.AuthID(), i, 301, &pb.UpdateSimple1{
				P1: "EQ",
			})
		}
	})
}

func initEdgeServer(serverID string, clientPort int, opts ...Option) *EdgeServer {
	opts = append(opts,
		WithWebsocketGateway(websocketGateway.Config{
			NewConnectionWorkers: 1,
			MaxConcurrency:       10,
			MaxIdleTime:          0,
			ListenAddress:        fmt.Sprintf(":%d", clientPort),
		}),
	)
	edge := NewEdgeServer(serverID, &testDispatcher{}, opts...)
	initHandlers(edge)

	return edge
}

func BenchmarkEdgeServer(b *testing.B) {
	edgeServer := initEdgeServer("Adam", 8080, WithDataPath("./_hdd"))
	req := &pb.ReqSimple1{P1: fmt.Sprintf("%d", 100)}
	envelope := &MessageEnvelope{}
	envelope.RequestID = tools.RandomUint64()
	envelope.Constructor = 101
	envelope.Message, _ = req.Marshal()
	proto := &msg.ProtoMessage{}
	proto.AuthID = 100
	proto.Payload, _ = envelope.Marshal()
	bytes, _ := proto.Marshal()

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(1000)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			edgeServer.onGatewayMessage(mockGatewayConn{}, 0, bytes)
		}
	})
}
