package edgec_test

import (
	"flag"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb/service"
	"github.com/ronaksoft/rony/log"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"sync"
	"testing"
	"time"
)

/*
   Creation Time: 2020 - Jul - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func init() {
	testEnv.Init()
}

func TestMain(m *testing.M) {
	edgeServer := testEnv.EdgeServer("Adam", 8081, 1000)
	rony.SetLogLevel(log.WarnLevel)
	service.RegisterSample(
		&service.Sample{
			ServerID: edgeServer.GetServerID(),
		},
		edgeServer,
	)

	edgeServer.Start()

	flag.Parse()
	code := m.Run()
	edgeServer.Shutdown()
	os.Exit(code)
}

func TestClient_Connect(t *testing.T) {
	Convey("Websocket Client Tests", t, func(c C) {
		Convey("One Connection With Concurrent Request", func(c C) {
			wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
				SeedHostPort: "127.0.0.1:8081",
				Handler: func(m *rony.MessageEnvelope) {
					c.Println("Received Uncaught Message", registry.C(m.Constructor))
				},
			})
			clnt := service.NewSampleClient(wsc)
			err := wsc.Start()
			_, _ = c.Println(wsc.ConnInfo())
			c.So(err, ShouldBeNil)

			wg := sync.WaitGroup{}
			for i := 0; i < 200; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					x := tools.RandomInt64(0)
					res, err := clnt.Echo(&service.EchoRequest{Int: x})
					if err != nil {
						c.Println(err.Error())
					}
					// _ = res
					c.So(err, ShouldBeNil)
					c.So(res.Int, ShouldEqual, x)
				}()
			}
			wg.Wait()
			err = wsc.Close()
			c.So(err, ShouldBeNil)
		})
		Convey("One Connection With Slow Data-Rate", func(c C) {
			wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
				SeedHostPort: "127.0.0.1:8081",
			})
			clnt := service.NewSampleClient(wsc)
			err := wsc.Start()
			c.So(err, ShouldBeNil)

			for i := 0; i < 5; i++ {
				res, err := clnt.Echo(&service.EchoRequest{Int: 123})
				c.So(err, ShouldBeNil)
				c.So(res.Int, ShouldEqual, 123)
				time.Sleep(time.Second * 2)
			}
		})
		Convey("One Connection With Slow Request", func(c C) {
			wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
				SeedHostPort: "127.0.0.1:8081",
			})
			clnt := service.NewSampleClient(wsc)
			err := wsc.Start()
			c.So(err, ShouldBeNil)

			for i := 0; i < 5; i++ {
				res, err := clnt.EchoDelay(&service.EchoRequest{Int: 123})
				c.So(err, ShouldBeNil)
				c.So(res.Int, ShouldEqual, 123)
			}
		})
	})
}
