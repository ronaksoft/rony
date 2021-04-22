package tcpGateway_test

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/internal/gateway"
	tcpGateway "github.com/ronaksoft/rony/internal/gateway/tcp"
	wsutil "github.com/ronaksoft/rony/internal/gateway/tcp/util"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
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
	gw.SetProxy(edge.MethodGet, "/x/:name",
		tcpGateway.CreateHandle(
			func(conn rony.Conn, ctx *gateway.RequestCtx) []byte {
				return tools.S2B(fmt.Sprintf("Received Get with Param: %s", conn.Get("name")))
			},
			func(data []byte) (*pools.ByteBuffer, map[string]string) {
				buf := pools.Buffer.GetCap(len(data))
				buf.Append(data)
				return buf, nil
			},
		),
	)
	gw.MessageHandler = func(c rony.Conn, streamID int64, data []byte, bypass bool) {
		if len(data) > 0 && data[0] == 'S' {
			time.Sleep(time.Duration(len(data)) * time.Second)
		}
		err := c.SendBinary(streamID, data)
		if err != nil {
			fmt.Println("MessageHandler:", err.Error())
		}
	}
	gw.ConnectHandler = func(c rony.Conn, kvs ...*rony.KeyValue) {
		// fmt.Println(c.ClientIP())
	}

	gw.Start()
	defer gw.Shutdown()
	Convey("Gateway Test", t, func(c C) {
		Convey("Ping", func(c C) {
			wsc, _, _, err := ws.Dial(context.Background(), fmt.Sprintf("ws://%s", hostPort))
			if err != nil {
				c.Println(err)
			}
			c.So(err, ShouldBeNil)
			for j := 0; j < 20; j++ {
				err := wsutil.WriteMessage(wsc, ws.StateClientSide, ws.OpPing, []byte{1, 2, 3, 4})
				if err != nil {
					c.Println(err)
				}
				c.So(err, ShouldBeNil)
				m, err := wsutil.ReadMessage(wsc, ws.StateClientSide, nil)
				c.So(err, ShouldBeNil)
				c.So(m, ShouldHaveLength, 1)
				c.So(m[0].OpCode, ShouldEqual, ws.OpPong)
			}
		})
		Convey("Websocket Parallel", func(c C) {
			wsc, _, _, err := ws.Dial(context.Background(), fmt.Sprintf("ws://%s", hostPort))
			if err != nil {
				c.Println(err)
			}
			c.So(err, ShouldBeNil)

			err = wsutil.WriteMessage(wsc, ws.StateClientSide, ws.OpBinary, []byte{'S', 1, 2, 3, 4})
			c.So(err, ShouldBeNil)

			wg := pools.AcquireWaitGroup()

			wg.Add(1)
			go func() {
				for i := 0; i < 5; i++ {
					err := wsutil.WriteMessage(wsc, ws.StateClientSide, ws.OpBinary, []byte{1, 2, 3, 4})
					if err != nil {
						c.Println(err)
					}
					c.So(err, ShouldBeNil)
					_, err = wsutil.ReadMessage(wsc, ws.StateClientSide, nil)
					if err != nil {
						c.Println(err)
					}
					c.So(err, ShouldBeNil)
				}
				wg.Done()
			}()

			wg.Wait()
			err = wsc.Close()
		})
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
		Convey("Http With Path", func(c C) {
			wg := pools.AcquireWaitGroup()
			wg.Add(1)
			go func() {
				httpc := &fasthttp.Client{}
				c.So(err, ShouldBeNil)
				for i := 0; i < 400; i++ {
					wg.Add(1)
					go func() {
						x := tools.RandomID(10)
						req := fasthttp.AcquireRequest()
						defer fasthttp.ReleaseRequest(req)
						res := fasthttp.AcquireResponse()
						defer fasthttp.ReleaseResponse(res)
						req.SetRequestURI(fmt.Sprintf("http://%s/x/%s", hostPort, x))
						req.Header.SetMethod(edge.MethodGet)
						err := httpc.Do(req, res)
						if err != nil {
							c.Println(err)
						}
						c.So(err, ShouldBeNil)
						c.So(res.Body(), ShouldResemble, tools.S2B(fmt.Sprintf("Received Get with Param: %s", x)))
						wg.Done()
					}()
				}
				wg.Done()
			}()

			wg.Wait()
			time.Sleep(time.Second)
		})

	})
}
