package rony_test

import (
	"context"
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/gateway/ws/util"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/ws"
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

var (
	clientPort = 8080
)

func init() {
	_ = os.MkdirAll("./_hdd", os.ModePerm)
	testEnv.Init()
	testEnv.EdgeServer = testEnv.InitEdgeServer("Adam", clientPort,
		edge.WithDataPath("./_hdd"),
	)
	err := testEnv.EdgeServer.Run()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
}

func TestEdgeServerSimple(t *testing.T) {
	Convey("Simple Edge", t, func(c C) {
		conn, _, _, err := ws.Dial(context.Background(), fmt.Sprintf("ws://127.0.0.1:%d", clientPort))
		c.So(err, ShouldBeNil)
		for i := int64(1); i <= 10; i++ {
			req := &pb.ReqSimple1{P1: tools.StrToByte(tools.Int64ToStr(i))}
			envelope := &rony.MessageEnvelope{
				RequestID:   tools.RandomUint64(),
				Constructor: 101,
			}
			envelope.Message, _ = req.Marshal()
			proto := &pb.ProtoMessage{}
			proto.AuthID = i
			proto.Payload, _ = envelope.Marshal()
			bytes, _ := proto.Marshal()
			err = wsutil.WriteMessage(conn, ws.StateClientSide, ws.OpBinary, bytes)
			c.So(err, ShouldBeNil)
		}
		testEnv.EdgeServer.Shutdown()
	})
}

func TestEdgeServerRaft(t *testing.T) {
	Convey("Replicated Edge", t, func(c C) {
		clientPort1 := 8081
		edge1 := testEnv.InitEdgeServer("Raft.01", clientPort1,
			edge.WithDataPath("./_hdd/edge01"),
			edge.WithReplicaSet(1, 9091, true),
			edge.WithGossipPort(9081),
		)
		clientPort2 := 8082
		edge2 := testEnv.InitEdgeServer("Raft.02", clientPort2,
			edge.WithDataPath("./_hdd/edge02"),
			edge.WithReplicaSet(1, 9092, false),
			edge.WithGossipPort(9082),
		)
		clientPort3 := 8083
		edge3 := testEnv.InitEdgeServer("Raft.03", clientPort3,
			edge.WithDataPath("./_hdd/edge03"),
			edge.WithReplicaSet(1, 9093, false),
			edge.WithGossipPort(9083),
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

func BenchmarkServer(b *testing.B) {
	req := &pb.ReqSimple1{P1: tools.StrToByte(tools.Int64ToStr(3232343434))}
	envelope := &rony.MessageEnvelope{}
	envelope.RequestID = tools.RandomUint64()
	envelope.Constructor = 101
	envelope.Message, _ = req.Marshal()
	proto := &pb.ProtoMessage{}
	proto.AuthID = tools.RandomInt64(0)
	proto.Payload, _ = envelope.Marshal()
	bytes, _ := proto.Marshal()
	clientPort := 8080

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		var (
			ms []wsutil.Message
		)
		conn, _, _, err := ws.Dial(context.Background(), fmt.Sprintf("ws://127.0.0.1:%d", clientPort))
		if err != nil {
			b.Fatal(err)
		}
		for p.Next() {
			err = wsutil.WriteMessage(conn, ws.StateClientSide, ws.OpBinary, bytes)
			if err != nil {
				b.Error("Write", err)
			}
			ms, err = wsutil.ReadMessage(conn, ws.StateClientSide, ms)
			if err != nil {
				b.Error("Read", err)
			}
			ms = ms[:0]
		}
	})
	b.StopTimer()
}
