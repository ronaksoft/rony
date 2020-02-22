package rony_test

import (
	context2 "context"
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/context"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"git.ronaksoftware.com/ronak/rony/msg"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	. "github.com/smartystreets/goconvey/convey"
	"os"
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
	opts = append(opts, rony.WithWebsocketGateway(websocketGateway.Config{
		NewConnectionWorkers: 1,
		MaxConcurrency:       10,
		MaxIdleTime:          0,
		ListenAddress:        fmt.Sprintf(":%d", clientPort),
	}))
	edge := rony.NewEdgeServer(bundleID, instanceID, opts...)
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
		for i := 0; i < 10; i++ {
			req := &pb.ReqSimple1{P1: fmt.Sprintf("%d", i)}
			reqBytes, _ := req.Marshal()
			err = wsutil.WriteClientBinary(conn, reqBytes)
			c.So(err, ShouldBeNil)
		}
		edge.Shutdown()
	})
}

func TestEdgeServerRaft(t *testing.T) {
	Convey("Replicated Edge", t, func(c C) {
		clientPort1 := 8081
		edge1 := initEdgeServer("Test", "01", clientPort1,
			rony.WithDataPath("./_hdd"),
			rony.WithRaft(9091, true),
		)
		clientPort2 := 8082
		edge2 := initEdgeServer("Test", "02", clientPort2,
			rony.WithDataPath("./_hdd"),
			rony.WithRaft(9092, true),
		)
		clientPort3 := 8083
		edge3 := initEdgeServer("Test", "03", clientPort3,
			rony.WithDataPath("./_hdd"),
			rony.WithRaft(9093, true),
		)

		// Run Edge 01
		err := edge1.Run()
		c.So(err, ShouldBeNil)
		c.Println("Wait for node1 to initialize")

		// Run Edge 02
		err = edge2.Run()
		c.So(err, ShouldBeNil)
		err = edge2.Join(edge1.GetServerID(), "127.0.0.1:9091")
		c.So(err, ShouldBeNil)

		// Run Edge 03
		err = edge3.Run()
		c.So(err, ShouldBeNil)
		err = edge3.Join(edge1.GetServerID(), "127.0.0.1:9091")
		c.So(err, ShouldBeNil)

		conn, _, _, err := ws.Dial(context2.Background(), fmt.Sprintf("ws://127.0.0.1:%d", clientPort1))
		c.So(err, ShouldBeNil)
		for i := 0; i < 10; i++ {
			req := &pb.ReqSimple1{P1: fmt.Sprintf("%d", i)}
			reqBytes, _ := req.Marshal()
			err = wsutil.WriteClientBinary(conn, reqBytes)
			c.So(err, ShouldBeNil)
		}
		edge1.Shutdown()
		edge2.Shutdown()
		edge3.Shutdown()
	})
}
