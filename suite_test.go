package rony_test

import (
	context2 "context"
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/context"
	"git.ronaksoftware.com/ronak/rony/gateway"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"git.ronaksoftware.com/ronak/rony/msg"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

/*
   Creation Time: 2020 - Feb - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	receivedMessages int32
	receivedUpdates  int32
)

type testDispatcher struct {
}

func (t testDispatcher) DispatchUpdate(conn gateway.Conn, streamID, authID int64, envelope *msg.UpdateEnvelope) {
	atomic.AddInt32(&receivedUpdates, 1)
}

func (t testDispatcher) DispatchMessage(conn gateway.Conn, streamID, authID int64, envelope *msg.MessageEnvelope) {
	log.Warn("Message Received", zap.Uint64("ReqID", envelope.RequestID))
	atomic.AddInt32(&receivedMessages, 1)
}

func (t testDispatcher) DispatchRequest(conn gateway.Conn, streamID int64, data []byte, envelope *msg.MessageEnvelope) (authID int64, err error) {
	proto := pools.AcquireProtoMessage()
	err = proto.Unmarshal(data)
	if err != nil {
		return
	}
	err = envelope.Unmarshal(proto.Payload)
	if err != nil {
		return
	}
	authID = proto.AuthID
	pools.ReleaseProtoMessage(proto)
	return
}

func initHandlers(edge *rony.EdgeServer) {
	edge.AddHandler(100, func(ctx *context.Context, in *msg.MessageEnvelope) {
		req := &pb.ReqSimple1{}
		res := &pb.ResSimple1{}
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(in.RequestID, "Invalid", "Proto")
			return
		}
		res.P1 = tools.StrToByte(req.P1)
		ctx.PushMessage(ctx.AuthID, in.RequestID, 201, res)
	})
	edge.AddHandler(101, func(ctx *context.Context, in *msg.MessageEnvelope) {
		req := &pb.ReqSimple1{}
		res := &pb.ResSimple1{}
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(in.RequestID, "Invalid", "Proto")
			return
		}
		res.P1 = tools.StrToByte(req.P1)

		ctx.PushMessage(ctx.AuthID, in.RequestID, 201, res)
		for i := int64(10); i < 20; i++ {
			ctx.PushUpdate(ctx.AuthID, i, 301, &pb.UpdateSimple1{
				P1: fmt.Sprintf("%d", i),
			})
		}
	})
}

func initEdgeServer(bundleID, instanceID string, clientPort int, opts ...rony.Option) *rony.EdgeServer {
	opts = append(opts,
		rony.WithWebsocketGateway(websocketGateway.Config{
			NewConnectionWorkers: 1,
			MaxConcurrency:       10,
			MaxIdleTime:          0,
			ListenAddress:        fmt.Sprintf(":%d", clientPort),
		}),
	)
	edge := rony.NewEdgeServer(bundleID, instanceID, &testDispatcher{}, opts...)
	initHandlers(edge)

	return edge
}

func init() {
	_ = os.MkdirAll("./_hdd", os.ModePerm)
	testEnv.Init()
}

func TestEdgeServerSimple(t *testing.T) {
	Convey("Simple Edge", t, func(c C) {
		clientPort := 8080
		edge := initEdgeServer("Test", "01", clientPort,
			rony.WithDataPath("./_hdd"),
		)
		err := edge.Run()
		c.So(err, ShouldBeNil)
		time.Sleep(time.Second)
		conn, _, _, err := ws.Dial(context2.Background(), fmt.Sprintf("ws://127.0.0.1:%d", clientPort))
		c.So(err, ShouldBeNil)
		for i := int64(1); i <= 10; i++ {
			req := &pb.ReqSimple1{P1: fmt.Sprintf("%d", i)}
			envelope := pools.AcquireMessageEnvelope()
			envelope.RequestID = tools.RandomUint64()
			envelope.Constructor = 101
			envelope.Message, _ = req.Marshal()
			proto := pools.AcquireProtoMessage()
			proto.AuthID = i
			proto.Payload, _ = envelope.Marshal()
			bytes, _ := proto.Marshal()
			err = wsutil.WriteClientBinary(conn, bytes)
			c.So(err, ShouldBeNil)
		}
		edge.Shutdown()
	})
}

func TestEdgeServerRaft(t *testing.T) {
	Convey("Replicated Edge", t, func(c C) {
		clientPort1 := 8081
		edge1 := initEdgeServer("Test", "01", clientPort1,
			rony.WithDataPath("./_hdd/edge01"),
			rony.WithRaft(9091, true),
			rony.WithGossipPort(9081),
		)
		clientPort2 := 8082
		edge2 := initEdgeServer("Test", "02", clientPort2,
			rony.WithDataPath("./_hdd/edge02"),
			rony.WithRaft(9092, false),
			rony.WithGossipPort(9082),
		)
		clientPort3 := 8083
		edge3 := initEdgeServer("Test", "03", clientPort3,
			rony.WithDataPath("./_hdd/edge03"),
			rony.WithRaft(9093, false),
			rony.WithGossipPort(9083),
		)

		// Run Edge 01
		err := edge1.Run()
		c.So(err, ShouldBeNil)
		// Run Edge 02
		err = edge2.Run()
		c.So(err, ShouldBeNil)
		// Run Edge 03
		err = edge3.Run()
		c.So(err, ShouldBeNil)

		// Join Nodes
		err = edge1.JoinCluster("127.0.0.1:9082")
		c.So(err, ShouldBeNil)
		err = edge1.JoinCluster("127.0.0.1:9083")
		c.So(err, ShouldBeNil)

		// conn, _, _, err := ws.Dial(context2.Background(), fmt.Sprintf("ws://127.0.0.1:%d", clientPort1))
		// c.So(err, ShouldBeNil)
		// c.So(conn, ShouldNotBeNil)
		// for i := 1; i < 11; i++ {
		// 	req := &pb.ReqSimple1{P1: fmt.Sprintf("%d", i)}
		// 	reqBytes, _ := req.Marshal()
		// 	msgIn := pools.AcquireMessageEnvelope()
		// 	msgIn.RequestID = uint64(i)
		// 	msgIn.Constructor = 101
		// 	msgIn.Message = reqBytes
		// 	msgInBytes, _ := msgIn.Marshal()
		// 	err = wsutil.WriteClientBinary(conn, msgInBytes)
		// 	c.So(err, ShouldBeNil)
		// }
		// time.Sleep(time.Second * 3)
		// edge1.Shutdown()
		// time.Sleep(time.Second)
		// edge2.Shutdown()
		// time.Sleep(time.Second)
		// edge3.Shutdown()
		// c.So(receivedMessages, ShouldEqual, 30)
		// c.So(receivedUpdates, ShouldEqual, 300)

		time.Sleep(10 * time.Second)
	})
}
