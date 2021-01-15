package tcpGateway_test

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/ronaksoft/rony"
	tcpGateway "github.com/ronaksoft/rony/gateway/tcp"
	wsutil "github.com/ronaksoft/rony/gateway/tcp/util"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/pools"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/valyala/fasthttp"
	"net/http"
	"testing"
	"time"
)

/*
   Creation Time: 2020 - Aug - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func init() {
	testEnv.Init()
}

func TestGateway(t *testing.T) {
	// rony.SetLogLevel(-1)
	hostPort := "127.0.0.1:8080"
	gw, err := tcpGateway.New(tcpGateway.Config{
		Concurrency:   1000,
		ListenAddress: hostPort,
		MaxBodySize:   0,
		MaxIdleTime:   0,
		ExternalAddrs: []string{hostPort},
	})
	if err != nil {
		t.Fatal(err)
	}
	gw.MessageHandler = func(c rony.Conn, streamID int64, data []byte) {
		err := c.SendBinary(streamID, data)
		if err != nil {
			fmt.Println("MessageHandler:", err.Error())
		}
	}

	gw.Start()
	defer gw.Shutdown()
	Convey("Gateway Test", t, func(c C) {
		Convey("Websocket / With Normal Handler", func(c C) {
			wg := pools.AcquireWaitGroup()
			for i := 0; i < 50; i++ {
				wg.Add(1)
				go func() {
					wsc, _, _, err := ws.Dial(context.Background(), fmt.Sprintf("ws://%s", hostPort))
					if err != nil {
						c.Println(err)
					}
					c.So(err, ShouldBeNil)
					for j := 0; j < 20; j++ {
						err := wsutil.WriteMessage(wsc, ws.StateClientSide, ws.OpBinary, []byte{1, 2, 3, 4})
						if err != nil {
							c.Println(err)
						}
						c.So(err, ShouldBeNil)
					}
					err = wsc.Close()
					wg.Done()
				}()
			}
			wg.Wait()
			time.Sleep(time.Second)
			c.Println("Total Connections", gw.TotalConnections())
		})
		Convey("Http With Normal Handler", func(c C) {
			wg := pools.AcquireWaitGroup()
			wg.Add(1)
			go func() {
				httpc := &fasthttp.Client{}
				c.So(err, ShouldBeNil)
				for i := 0; i < 400; i++ {
					wg.Add(1)
					go func() {
						req := fasthttp.AcquireRequest()
						defer fasthttp.ReleaseRequest(req)
						res := fasthttp.AcquireResponse()
						defer fasthttp.ReleaseResponse(res)
						req.SetHost(hostPort)
						req.Header.SetMethod(http.MethodPost)
						req.SetBody([]byte{1, 2, 3, 4})
						err := httpc.Do(req, res)
						if err != nil {
							c.Println(err)
						}
						c.So(err, ShouldBeNil)
						c.So(res.Body(), ShouldResemble, []byte{1, 2, 3, 4})
						wg.Done()
					}()
				}
				wg.Done()
			}()

			wg.Wait()
			time.Sleep(time.Second)
			c.Println("Total Connections", gw.TotalConnections())
		})
	})

}
