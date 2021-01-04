package edgec_test

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/gateway"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
	"sync"
	"testing"
)

/*
   Creation Time: 2020 - Jul - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type server struct {
	e *edge.Server
}

func (s server) OnMessage(ctx *edge.DispatchCtx, envelope *rony.MessageEnvelope) {
	b, _ := proto.Marshal(envelope)
	_ = ctx.Conn().SendBinary(ctx.StreamID(), b)
}

func (s server) Interceptor(ctx *edge.DispatchCtx, data []byte) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (s server) Done(ctx *edge.DispatchCtx) {
}

func (s server) OnOpen(conn gateway.Conn, kvs ...gateway.KeyValue) {
	fmt.Println("Connected")
}

func (s server) OnClose(conn gateway.Conn) {
	fmt.Println("Disconnected")
}

func newTestServer() *server {
	// log.SetLevel(log.WarnLevel)

	s := &server{
		e: testEnv.InitEdgeServerWithWebsocket("Test.01", 8081, 10),
	}
	pb.RegisterSample(&testEnv.Handlers{
		ServerID: s.e.GetServerID(),
	}, s.e)

	return s
}

func TestClient_Connect(t *testing.T) {
	Convey("Client Connect", t, func(c C) {
		testEnv.Init()
		s := newTestServer()
		err := s.e.StartCluster()
		if err != nil && err != edge.ErrClusterNotSet {
			t.Fatal(err)
		}
		err = s.e.StartGateway()
		c.So(err, ShouldBeNil)

		wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort: "127.0.0.1:8081",
		})
		clnt := pb.NewSampleClient(wsc)
		err = wsc.Start()
		c.So(err, ShouldBeNil)

		wg := sync.WaitGroup{}
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				res, err := clnt.Echo(&pb.EchoRequest{Int: 123})
				c.So(err, ShouldBeNil)
				c.Println(res)
			}()
		}
		wg.Wait()
	})

}
