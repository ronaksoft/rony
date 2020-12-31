package tcp_test

import (
	"context"
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/gateway"
	tcpGateway "github.com/ronaksoft/rony/gateway/tcp"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
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
	rony.SetLogLevel(0)
	gw, err := tcpGateway.New(tcpGateway.Config{
		Concurrency:   1000,
		ListenAddress: "0.0.0.0:88",
		MaxBodySize:   0,
		MaxIdleTime:   0,
		Protocol:      tcpGateway.Auto,
		ExternalAddrs: []string{"127.0.0.1:88"},
	})

	if err != nil {
		t.Fatal(err)
	}

	gw.Start()
	defer gw.Shutdown()
	Convey("Gateway Test", t, func(c C) {
		Convey("Websocket / With Normal Handler", func(c C) {
			gw.MessageHandler = func(c gateway.Conn, streamID int64, data []byte) {
				e := &rony.MessageEnvelope{}
				_ = e.Unmarshal(data)
				out, _ := e.Marshal()
				err := c.SendBinary(streamID, out)
				if err != nil {
					fmt.Println("MessageHandler:", err.Error())
				}
			}
			wg := pools.AcquireWaitGroup()
			for i := 0; i < 500; i++ {
				wg.Add(1)
				go func() {
					wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
						HostPort:        "127.0.0.1:88",
						IdleTimeout:     time.Second,
						DialTimeout:     time.Second,
						Handler:         nil,
						Header:          nil,
						Secure:          false,
						ForceConnect:    true,
						RequestMaxRetry: 10,
						RequestTimeout:  time.Second,
						ContextTimeout:  time.Second,
					})
					req := &rony.MessageEnvelope{
						Constructor: 100,
						RequestID:   100,
						Message:     []byte{1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1},
						Auth:        nil,
						Header:      nil,
					}
					res := &rony.MessageEnvelope{}
					for j := 0; j < 20; j++ {
						ctx, cf := context.WithTimeout(context.TODO(), time.Second*5)
						err = wsc.SendWithContext(ctx, req, res)
						if err != nil {
							c.Println(err)
						}
						c.So(err, ShouldBeNil)
						cf()
					}
					err = wsc.Close()
					wg.Done()
				}()
			}
			wg.Wait()
			time.Sleep(time.Second)
			c.Println("Total Connections", gw.TotalConnections())
		})
		Convey("Websocket / With Panic Handler", func(c C) {
			gw.MessageHandler = func(c gateway.Conn, streamID int64, data []byte) {
				err := c.SendBinary(streamID, nil)
				if err != nil {
					fmt.Println("MessageHandler:", err.Error())
				}
			}

			for i := 0; i < 10; i++ {
				wsc := edgec.NewWebsocket(edgec.WebsocketConfig{
					HostPort:        "127.0.0.1:88",
					IdleTimeout:     time.Second,
					DialTimeout:     time.Second,
					Handler:         nil,
					Header:          nil,
					Secure:          false,
					ForceConnect:    true,
					RequestMaxRetry: 10,
					RequestTimeout:  time.Second,
					ContextTimeout:  time.Second,
				})
				err = wsc.Close()
				c.So(err, ShouldBeNil)
			}
			time.Sleep(time.Second * 2)
			c.Println("Total Connections", gw.TotalConnections())
		})
		Convey("Http With Normal Handler", func(c C) {
			gw.MessageHandler = func(c gateway.Conn, streamID int64, data []byte) {
				e := &rony.MessageEnvelope{}
				_ = e.Unmarshal(data)
				out, _ := e.Marshal()
				err := c.SendBinary(streamID, out)
				if err != nil {
					fmt.Println("MessageHandler:", err.Error())
				}
			}

			wg := pools.AcquireWaitGroup()
			wg.Add(1)
			go func() {
				httpc := edgec.NewHttp(edgec.HttpConfig{
					HostPort: "127.0.0.1:88",
					Header:   nil,
				})
				for i := 0; i < 1000; i++ {
					wg.Add(1)
					go func() {
						req := &rony.MessageEnvelope{
							Constructor: tools.RandomInt64(0),
							RequestID:   httpc.GetRequestID(),
							Message:     tools.StrToByte(tools.RandomID(32)),
							Auth:        nil,
							Header:      nil,
						}
						res := &rony.MessageEnvelope{}
						err = httpc.Send(req, res)
						if err != nil {
							c.Println(err)
						}
						c.So(res.Constructor, ShouldEqual, req.Constructor)
						c.So(res.RequestID, ShouldEqual, req.RequestID)
						c.So(res.Message, ShouldResemble, req.Message)
						c.So(err, ShouldBeNil)
						wg.Done()
					}()
				}
				err = httpc.Close()
				c.So(err, ShouldBeNil)
				wg.Done()
			}()

			wg.Wait()
			time.Sleep(time.Second)
			c.Println("Total Connections", gw.TotalConnections())
		})
	})

}
