package tcpGateway_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gobwas/ws"
	"github.com/ronaksoft/rony"
	tcpGateway "github.com/ronaksoft/rony/internal/gateway/tcp"
	wsutil "github.com/ronaksoft/rony/internal/gateway/tcp/util"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/valyala/fasthttp"
)

/*
   Creation Time: 2020 - Aug - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	hostPort = "127.0.0.1:8080"
	gw       *tcpGateway.Gateway
)

func TestMain(m *testing.M) {
	_ = os.RemoveAll("./_hdd")
	_ = os.MkdirAll("./_hdd", os.ModePerm)

	gw = setupGateway()
	gw.Start()

	// Wait for server to come up
	time.Sleep(time.Second * 3)

	code := m.Run()
	gw.Shutdown()
	_ = os.RemoveAll("./_hdd")
	os.Exit(code)
}

func setupGateway() *tcpGateway.Gateway {
	//rony.SetLogLevel(log.DebugLevel)
	gw, err := tcpGateway.New(tcpGateway.Config{
		Concurrency:   1000,
		ListenAddress: hostPort,
		MaxBodySize:   0,
		MaxIdleTime:   0,
		ExternalAddrs: []string{hostPort},
	})
	if err != nil {
		panic(err)
	}

	gw.Subscribe(
		&testGatewayDelegate{
			onConnectFunc: func(c rony.Conn, kvs ...*rony.KeyValue) {
				//fmt.Println("Open", c.ClientIP())
			},
			onMessageFunc: func(c rony.Conn, streamID int64, data []byte) {
				if len(data) > 0 && data[0] == 'S' {
					time.Sleep(time.Duration(len(data)) * time.Second)
				}
				err := c.WriteBinary(streamID, data)
				if err != nil {
					fmt.Println("MessageHandler:", err.Error())
				}
			},
			onClose: func(c rony.Conn) {
				//fmt.Println("Closed", c.ClientIP())
			},
		},
	)

	return gw
}

func TestGateway(t *testing.T) {
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
			c.So(err, ShouldBeNil)
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
					c.So(err, ShouldBeNil)
					wg.Done()
				}()
			}
			wg.Wait()
			time.Sleep(time.Second)
		})
		Convey("Http With Normal Handler", func(c C) {
			wg := pools.AcquireWaitGroup()
			wg.Add(1)
			go func() {
				httpc := &fasthttp.Client{}
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
		})
		Convey("Http With Path", func(c C) {
			wg := pools.AcquireWaitGroup()
			wg.Add(1)
			go func() {
				httpc := &fasthttp.Client{}
				for i := 0; i < 400; i++ {
					wg.Add(1)
					go func() {
						x := tools.RandomID(10)
						req := fasthttp.AcquireRequest()
						defer fasthttp.ReleaseRequest(req)
						res := fasthttp.AcquireResponse()
						defer fasthttp.ReleaseResponse(res)
						req.SetRequestURI(fmt.Sprintf("http://%s/x/%s", hostPort, x))
						req.Header.SetMethod(rony.MethodGet)
						err := httpc.Do(req, res)
						if err != nil {
							c.Println(err)
						}
						c.So(err, ShouldBeNil)
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

type testGatewayDelegate struct {
	onConnectFunc func(c rony.Conn, kvs ...*rony.KeyValue)
	onMessageFunc func(c rony.Conn, streamID int64, data []byte)
	onClose       func(c rony.Conn)
}

func (t *testGatewayDelegate) OnConnect(c rony.Conn, kvs ...*rony.KeyValue) {
	t.onConnectFunc(c, kvs...)
}

func (t *testGatewayDelegate) OnMessage(c rony.Conn, streamID int64, data []byte) {
	t.onMessageFunc(c, streamID, data)
}

func (t *testGatewayDelegate) OnClose(c rony.Conn) {
	t.onClose(c)
}

func TestGateway2(t *testing.T) {
	msg := []byte(tools.RandomID(1024))
	Convey("Gateway Test", t, func(c C) {
		for j := 0; j < 100; j++ {
			wsc, _, _, err := ws.Dial(context.Background(), fmt.Sprintf("ws://%s", hostPort))
			c.So(err, ShouldBeNil)

			var m []wsutil.Message
			for i := 0; i < 10; i++ {
				err = wsutil.WriteMessage(wsc, ws.StateClientSide, ws.OpBinary, msg)
				c.So(err, ShouldBeNil)
				m, err = wsutil.ReadMessage(wsc, ws.StateClientSide, m[:0])
				c.So(err, ShouldBeNil)
				c.So(m, ShouldHaveLength, 1)
				c.So(m[0].Payload, ShouldResemble, msg)
			}

			err = wsc.Close()
			c.So(err, ShouldBeNil)
		}
	})
}

func BenchmarkGateway(b *testing.B) {
	msg := []byte(tools.RandomID(1024))

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1)
	b.RunParallel(func(pb *testing.PB) {
		wsc, _, _, err := ws.Dial(context.Background(), fmt.Sprintf("ws://%s", hostPort))
		if err != nil {
			panic(err)
		}

		var m []wsutil.Message
		for pb.Next() {
			err = wsutil.WriteMessage(wsc, ws.StateClientSide, ws.OpBinary, msg)
			if err != nil {
				b.Log(err)
			}
			m, err = wsutil.ReadMessage(wsc, ws.StateClientSide, m[:0])
			if err != nil {
				b.Log(err)
			}
		}

		_ = wsc.Close()
	})
}
