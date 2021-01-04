package tcp_test

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/gateway"
	tcpGateway "github.com/ronaksoft/rony/gateway/tcp"
	wsutil "github.com/ronaksoft/rony/gateway/tcp/util"
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
				err := c.SendBinary(streamID, data)
				if err != nil {
					fmt.Println("MessageHandler:", err.Error())
				}
			}
			wg := pools.AcquireWaitGroup()
			for i := 0; i < 50; i++ {
				wg.Add(1)
				go func() {
					wsc, _, _, err := ws.Dial(context.Background(), "ws://127.0.0.1:88")
					if err != nil {
						c.Println(err)
					}
					c.So(err, ShouldBeNil)
					for j := 0; j < 20; j++ {
						err := wsutil.WriteMessage(wsc, ws.StateClientSide, ws.OpBinary, []byte{1,2,3,4})
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
						err := httpc.Send(req, res)
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
