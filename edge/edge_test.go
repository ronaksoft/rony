package edge_test

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	dummyGateway "github.com/ronaksoft/rony/gateway/dummy"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
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

func BenchmarkStandaloneSerial(b *testing.B) {
	edgeServer := testEnv.InitEdgeServerWithWebsocket("Adam", 8080, edge.WithDataPath("./_hdd/adam"))
	pb.RegisterSample(testEnv.Handlers{}, edgeServer)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(1000)
	conn := dummyGateway.Conn{}
	req := pb.EchoRequest{
		Int:       100,
		Bool:      false,
		Timestamp: 32809238402,
	}
	e := &rony.MessageEnvelope{}
	e.Fill(100, pb.C_Echo, &req)

	reqBytes, _ := proto.Marshal(e)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			edgeServer.OnGatewayMessage(&conn, 0, reqBytes)
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
