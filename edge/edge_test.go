package edge_test

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb/service"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
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

func TestWithTestGateway(t *testing.T) {
	Convey("EdgeTest Gateway", t, func(c C) {
		s := testEnv.InitTestServer("TestServer")
		service.RegisterSample(&testEnv.Handlers{ServerID: "TestServer"}, s.RealEdge())
		s.Start()
		defer s.Shutdown()

		err := s.Context().
			Request(service.C_SampleEcho, &service.EchoRequest{
				Int:       100,
				Timestamp: 123,
			}).
			ErrorHandler(func(constructor int64, e *rony.Error) {
				c.Println(registry.ConstructorName(constructor), "-->", e.Code, e.Items, e.Description)
			}).
			Expect(service.C_EchoResponse, func(b []byte, auth []byte, kv ...*rony.KeyValue) error {
				x := &service.EchoResponse{}
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

func TestRestProxy(t *testing.T) {
	Convey("Eddge With RestProxy", t, func(c C) {
		s := testEnv.InitTestServer("TestServer")
		s.RealEdge().SetRestProxy(
			rony.MethodGet, "/x/:value",
			edge.NewRestProxy(
				func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
					req := &service.EchoRequest{
						Int: tools.StrToInt64(conn.Get("value").(string)),
					}
					reqB, _ := proto.Marshal(req)
					ctx.FillEnvelope(conn.ConnID(), service.C_EchoRequest, reqB, nil)
					return nil
				},
				func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
					for me := ctx.BufferPop(); me != nil; me = ctx.BufferPop() {
						c.So(me.Constructor, ShouldEqual, service.C_EchoResponse)
						x := &service.EchoResponse{}
						err := x.Unmarshal(me.Message)
						c.So(err, ShouldBeNil)
						err = conn.WriteBinary(ctx.StreamID(), tools.S2B(tools.Int64ToStr(x.Int)))
						c.So(err, ShouldBeNil)
					}
					return nil
				},
			),
		)
		service.RegisterSample(&testEnv.Handlers{ServerID: "TestServer"}, s.RealEdge())
		s.Start()
		defer s.Shutdown()

	})
}

func TestConcurrent(t *testing.T) {
	Convey("Concurrent", t, func(c C) {

	})
}

func BenchmarkEdge(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {

	})
}
