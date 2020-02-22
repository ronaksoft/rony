package websocketGateway_test

import (
	"context"
	"git.ronaksoftware.com/ronak/rony/gateway"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"sync"
	"testing"
	"time"
)

/*
   Creation Time: 2019 - May - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	gw *websocketGateway.Gateway
)

func init() {
	testEnv.Init()
}

func TestGateway(t *testing.T) {
	Convey("Test Websocket Gateway", t, func() {
		Convey("Run the Server", func(c C) {
			var err error
			gw, err = websocketGateway.New(websocketGateway.Config{
				CloseHandler: func(c gateway.Conn) {},
				MessageHandler: func(conn gateway.Conn, streamID int64, data []byte) {
					c.So(data, ShouldHaveLength, 4)
					err := conn.(*websocketGateway.Conn).SendBinary(streamID, []byte{1, 2, 3, 4})
					c.So(err, ShouldBeNil)
				},
				ConnectHandler: func(connID uint64) {},
				FlushFunc: func(c gateway.Conn) [][]byte {
					return nil
				},
				MaxIdleTime:          time.Second * 3,
				NewConnectionWorkers: 1,
				MaxConcurrency:       1,
				ListenAddress:        ":81",
			})
			c.So(err, ShouldBeNil)
			gw.Run()
			time.Sleep(time.Second)
		})

		Convey("Run the Clients", func(c C) {
			wg := sync.WaitGroup{}
			for j := 0; j < 1; j++ {
				wg.Add(1)
				go func(j int) {
					var conn net.Conn
					var err error
					for i := 0; i < 5; i++ {
						conn, _, _, err = ws.Dial(context.Background(), "ws://localhost:81")
						c.So(err, ShouldBeNil)
					}
					for i := 0; i < 5; i++ {
						time.Sleep(time.Second * 3)
						err := wsutil.WriteClientBinary(conn, tools.StrToByte("ABCD"))
						c.So(err, ShouldBeNil)
					}
					keepGoing := true
					var m []wsutil.Message
					for keepGoing {
						_ = conn.SetDeadline(time.Now().Add(time.Second))
						m, err = wsutil.ReadServerMessage(conn, m)
						if err != nil {
							keepGoing = false
						}
					}
					c.So(m, ShouldHaveLength, 5)
					_ = conn.Close()
					wg.Done()
				}(j)
			}
			wg.Wait()
		})

		Convey("Run Idle Client", func(c C) {
			conn, _, _, err := ws.Dial(context.Background(), "ws://localhost:81")
			c.So(err, ShouldBeNil)
			_, _ = wsutil.ReadServerMessage(conn, nil)
			_ = conn.Close()
		})
		Convey("Wait To Finish", func(c C) {
			time.Sleep(time.Second * 3)
		})
	})
}
