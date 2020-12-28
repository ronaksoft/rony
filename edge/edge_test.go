package edge_test

import (
	"flag"
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	. "github.com/smartystreets/goconvey/convey"
	"os"
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
	edgeServer *edge.Server
)

func TestMain(m *testing.M) {
	edgeServer = testEnv.InitEdgeServerWithWebsocket("Adam", 8080, 1000, edge.WithDataPath("./_hdd/adam"))

	pb.RegisterSample(testEnv.Handlers{}, edgeServer)
	_ = edgeServer.StartCluster()
	edgeServer.StartGateway()
	flag.Parse()
	code := m.Run()

	time.Sleep(time.Second * 10)

	edgeServer.Shutdown()
	os.Exit(code)
}

func BenchmarkServer(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	// b.SetParallelism(10)

	echoRequest := pb.EchoRequest{
		Int:       100,
		Bool:      false,
		Timestamp: 32809238402,
	}

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			edgeClient := edgec.NewWebsocket(edgec.WebsocketConfig{
				HostPort:     "127.0.0.1:8080",
				IdleTimeout:  time.Second,
				DialTimeout:  time.Second,
				ForceConnect: true,
				Handler:      func(m *rony.MessageEnvelope) {},
				// RequestMaxRetry: 1,
				// RequestTimeout:  time.Second,
				// ContextTimeout:  time.Second,
			})

			req := rony.PoolMessageEnvelope.Get()
			res := rony.PoolMessageEnvelope.Get()
			req.Fill(edgeClient.GetRequestID(), pb.C_Echo, &echoRequest)
			_ = edgeClient.Send(req, res)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
	})
}

func TestWithTestGateway(t *testing.T) {
	Convey("EdgeTest Gateway", t, func(c C) {
		s := testEnv.InitTestServer("TestServer")
		s.SetHandlers(pb.C_EchoRequest, testEnv.EchoSimple)
		s.Start()
		defer s.Shutdown()

		err := s.Context().
			Request(pb.C_EchoRequest, &pb.EchoRequest{
				Int:       100,
				Bool:      true,
				Timestamp: 123,
			}).
			Expect(pb.C_EchoResponse, func(b []byte, auth []byte, kv ...*rony.KeyValue) error {
				x := &pb.EchoResponse{}
				err := x.Unmarshal(b)
				if err != nil {
					return err
				}
				if x.Int != 100 {
					return fmt.Errorf("int not equal")
				}
				if x.Timestamp != 123 {
					return fmt.Errorf("timestamp not equal")
				}
				if !x.Bool {
					return fmt.Errorf("bool not equal")
				}
				return nil
			}).
			Run(time.Second)
		c.So(err, ShouldBeNil)
	})

}
