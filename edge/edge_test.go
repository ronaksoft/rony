package edge_test

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgetest"
	dummyGateway "github.com/ronaksoft/rony/internal/gateway/dummy"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb/service"
	"github.com/ronaksoft/rony/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/encoding/protojson"
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
	s *edgetest.Server
)

func TestMain(m *testing.M) {
	s = testEnv.TestServer("TestServer")
	service.RegisterSample(&service.Sample{ServerID: "TestServer"}, s.RealEdge())
	s.Start()
	defer s.Shutdown()

	m.Run()
}

func TestWithTestGateway(t *testing.T) {
	Convey("EdgeTest Gateway", t, func(c C) {
		err := s.RPC().
			Request(service.C_SampleEcho, &service.EchoRequest{
				Int:       100,
				Timestamp: 123,
			}).
			ErrorHandler(func(constructor uint64, e *rony.Error) {
				c.Println(registry.ConstructorName(constructor), "-->", e.Code, e.Items, e.Description)
			}).
			Expect(service.C_EchoResponse, func(b []byte, kv ...*rony.KeyValue) error {
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
	Convey("Edge With RestProxy", t, func(c C) {
		Convey("Manual", func(c C) {
			s.RealEdge().SetRestProxy(
				rony.MethodGet, "/x/:value",
				edge.NewRestProxy(
					func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
						req := &service.EchoRequest{
							Int: tools.StrToInt64(tools.GetString(conn.Get("value"), "0")),
						}
						ctx.Fill(conn.ConnID(), service.C_SampleEcho, req)
						return nil
					},
					func(conn rony.RestConn, ctx *edge.DispatchCtx) error {
						ctx.BufferPopAll(func(envelope *rony.MessageEnvelope) {
							c.So(envelope.Constructor, ShouldEqual, service.C_EchoResponse)
							x := &service.EchoResponse{}
							err := x.Unmarshal(envelope.Message)
							c.So(err, ShouldBeNil)
							err = conn.WriteBinary(ctx.StreamID(), tools.S2B(tools.Int64ToStr(x.Int)))
							c.So(err, ShouldBeNil)
						})

						return nil
					},
				),
			)
			service.RegisterSample(&service.Sample{ServerID: "TestServer"}, s.RealEdge())
			s.Start()
			defer s.Shutdown()

			err := s.REST().
				Request(rony.MethodGet, "/x/123", nil).
				Expect(func(b []byte, kv ...*rony.KeyValue) error {
					c.So(string(b), ShouldResemble, "123")
					return nil
				}).
				RunShort()
			c.So(err, ShouldBeNil)
		})
		Convey("JSON", func(c C) {
			s.RealEdge().SetRestProxy(rony.MethodPost, "/echo", service.EchoRest)
			service.RegisterSample(&service.Sample{ServerID: "TestServer"}, s.RealEdge())

			req := &service.EchoRequest{
				Int:        tools.RandomInt64(0),
				Timestamp:  tools.NanoTime(),
				ReplicaSet: tools.RandomUint64(0),
			}
			reqJSON, err := protojson.Marshal(req)
			c.So(err, ShouldBeNil)
			err = s.REST().
				Request(rony.MethodPost, "/echo", reqJSON).
				Expect(func(b []byte, kv ...*rony.KeyValue) error {
					res := &service.EchoResponse{}
					err = protojson.Unmarshal(b, res)
					c.Println(string(b))
					c.Println(res)
					c.So(err, ShouldBeNil)
					c.So(res.Int, ShouldEqual, req.Int)
					c.So(res.Timestamp, ShouldEqual, req.Timestamp)
					return nil
				}).
				RunShort()
			c.So(err, ShouldBeNil)
		})
		Convey("JSON and Binding", func(c C) {
			s.RealEdge().SetRestProxy(rony.MethodGet, "/echo/:value/:ts", service.EchoRestBinding)
			service.RegisterSample(&service.Sample{ServerID: "TestServer"}, s.RealEdge())

			value := tools.RandomInt64(0)
			ts := tools.NanoTime()
			err := s.REST().
				Request(rony.MethodGet, fmt.Sprintf("/echo/%d/%d", value, ts), nil).
				Expect(func(b []byte, kv ...*rony.KeyValue) error {
					res := &service.EchoResponse{}
					err := protojson.Unmarshal(b, res)
					c.Println(string(b))
					c.Println(res)
					c.So(err, ShouldBeNil)
					c.So(res.Int, ShouldEqual, value)
					c.So(res.Timestamp, ShouldEqual, ts)
					return nil
				}).
				RunShort()
			c.So(err, ShouldBeNil)
		})
	})
}

func BenchmarkEdge(b *testing.B) {
	rony.SetLogLevel(log.WarnLevel)
	e := testEnv.EdgeServer(tools.RandomID(0), 8080, 1000, edge.WithInMemoryStore(true))
	service.RegisterSample(&service.Sample{}, e)
	e.Start()
	defer e.Shutdown()

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := service.PoolEchoRequest.Get()
			req.Int = tools.RandomInt64(0)
			me := rony.PoolMessageEnvelope.Get()
			me.Fill(tools.RandomUint64(0), service.C_SampleEcho, req)
			buf := pools.Buffer.FromProto(me)
			e.OnGatewayMessage(dummyGateway.NewConn(func(connID uint64, streamID int64, data []byte, hdr map[string]string) {
				me := rony.PoolMessageEnvelope.Get()
				err := me.Unmarshal(data)
				if err != nil {
					b.Fatal(err)
				}
				if me.Constructor != service.C_EchoResponse {
					b.Fatal("invalid constructor")
				}
				rony.PoolMessageEnvelope.Put(me)
			}), 0, *buf.Bytes())
			pools.Buffer.Put(buf)
			service.PoolEchoRequest.Put(req)
			rony.PoolMessageEnvelope.Put(me)
		}
	})
}
