package edgeClient_test

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgeClient"
	"github.com/ronaksoft/rony/gateway"
	tcpGateway "github.com/ronaksoft/rony/gateway/tcp"
	log "github.com/ronaksoft/rony/internal/logger"
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

func (s server) OnMessage(ctx *edge.DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	b, _ := proto.Marshal(envelope)
	_ = ctx.Conn().SendBinary(ctx.StreamID(), b)
}

func (s server) Prepare(ctx *edge.DispatchCtx, data []byte, kvs ...gateway.KeyValue) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (s server) Done(ctx *edge.DispatchCtx) {
}

func (s server) OnOpen(conn gateway.Conn) {
	fmt.Println("Connected")
}

func (s server) OnClose(conn gateway.Conn) {
	fmt.Println("Disconnected")
}

func (s server) Func1(ctx *edge.RequestCtx, req *pb.Req1, res *pb.Res1) {
	res.Item1 = req.Item1
}

func (s server) Func2(ctx *edge.RequestCtx, req *pb.Req2, res *pb.Res2) {

}

func (s server) Echo(ctx *edge.RequestCtx, req *pb.EchoRequest, res *pb.EchoResponse) {

}

func (s server) Ask(ctx *edge.RequestCtx, req *pb.AskRequest, res *pb.AskResponse) {

}

func newTestServer() *server {
	log.SetLevel(log.WarnLevel)
	s := &server{}
	s.e = edge.NewServer("Test.01", s,
		edge.WithTcpGateway(
			tcpGateway.Config{
				Concurrency:   10,
				ListenAddress: ":8081",
				MaxBodySize:   0,
				MaxIdleTime:   0,
				Protocol:      tcpGateway.Auto,
			},
		),
	)
	w := pb.NewSampleServer(s)
	w.Register(s.e)

	return s
}
func TestClient_Connect(t *testing.T) {
	testEnv.Init()
	s := newTestServer()
	go func() {
		s.e.RunGateway()
		_ = s.e.RunCluster()
	}()

	pb.NewSampleServer(&server{})
	c := pb.NewSampleClient(edgeClient.NewWebsocket(edgeClient.Config{
		HostPort: "127.0.0.1:8081",
	}))
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := c.Func1(&pb.Req1{Item1: 123})
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(res)
		}()
	}
	wg.Wait()

}
