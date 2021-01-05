package edgec_test

import (
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	. "github.com/smartystreets/goconvey/convey"
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

func TestClient_Connect(t *testing.T) {
	Convey("Websocket Client Tests", t, func(c C) {
		e := testEnv.InitEdgeServerWithWebsocket("Test.01", 8081, 10)
		pb.RegisterSample(&testEnv.Handlers{
			ServerID: e.GetServerID(),
		}, e)

		err := e.StartCluster()
		if err != nil && err != edge.ErrClusterNotSet {
			t.Fatal(err)
		}
		err = e.StartGateway()
		c.So(err, ShouldBeNil)
		defer e.Shutdown()

		Convey("One Connection With Concurrent Request", func(c C) {
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
					c.So(res.Int, ShouldEqual, 123)
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
			clnt := pb.NewSampleClient(wsc)
			err = wsc.Start()
			c.So(err, ShouldBeNil)

			for i := 0; i < 5; i++ {
				res, err := clnt.Echo(&pb.EchoRequest{Int: 123})
				c.So(err, ShouldBeNil)
				c.So(res.Int, ShouldEqual, 123)
				time.Sleep(time.Second * 3)
			}
		})

	})

}
