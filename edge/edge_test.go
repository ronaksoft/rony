package edge_test

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	. "github.com/smartystreets/goconvey/convey"
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

func TestWithTestGateway(t *testing.T) {
	Convey("EdgeTest Gateway", t, func(c C) {
		s := testEnv.InitTestServer("TestServer")
		pb.RegisterSample(&testEnv.Handlers{ServerID: "TestServer"}, s.RealEdge())
		s.Start()
		defer s.Shutdown()

		err := s.Context().
			Request(pb.C_EchoRequest, &pb.EchoRequest{
				Int:       100,
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
				return nil
			}).
			Run(time.Second)
		c.So(err, ShouldBeNil)
	})

}

func TestConcurrent(t *testing.T) {
	Convey("Concurrent", t, func(c C) {

	})
}
