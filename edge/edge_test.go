package edge_test

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	"github.com/ronaksoft/rony/registry"
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
			Request(pb.C_Echo, &pb.EchoRequest{
				Int:       100,
				Timestamp: 123,
			}).
			ErrorHandler(func(constructor int64, e *rony.Error) {
				c.Println(registry.ConstructorName(constructor), "-->", e.Code, e.Items, e.Template)
			}).
			Expect(pb.C_EchoResponse, func(b []byte, auth []byte, kv ...*rony.KeyValue) error {
				x := &pb.EchoResponse{}
				err := x.Unmarshal(b)
				c.So(err, ShouldBeNil)
				c.So(x.Int, ShouldEqual, 100)
				c.So(x.Timestamp, ShouldEqual, 123)
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
