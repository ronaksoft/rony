package edgec_test

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/gateway"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
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
	testEnv.Init()
	s := newTestServer()
	err := s.e.StartCluster()
	if err != nil {
		t.Fatal(err)
	}
	err = s.e.StartGateway()
	if err != nil {
		t.Fatal(err)
	}

	wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
		SeedHostPort: "127.0.0.1:8081",
	})
	c := pb.NewSampleClient(wsc)
	err = wsc.Start()
	if err != nil {
		t.Fatal(err)
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := c.Echo(&pb.EchoRequest{Int: 123})
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(res)
		}()
	}
	wg.Wait()

}
