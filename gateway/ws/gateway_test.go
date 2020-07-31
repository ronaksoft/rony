package websocketGateway_test

import (
	"context"
	"git.ronaksoftware.com/ronak/rony/gateway"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	"git.ronaksoftware.com/ronak/rony/gateway/ws/util"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/ws"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
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
	runtime.GOMAXPROCS(1)
	Convey("Test Websocket Gateway", t, func() {
		Convey("Run the Server", func(c C) {
			var err error
			gw, err = websocketGateway.New(websocketGateway.Config{
				MaxIdleTime:          time.Second * 3,
				NewConnectionWorkers: 1,
				MaxConcurrency:       16,
				ListenAddress:        ":81",
			})
			c.So(err, ShouldBeNil)

			gw.MessageHandler = func(conn gateway.Conn, streamID int64, data []byte, kvs ...gateway.KeyValue) {
				// c.So(data, ShouldHaveLength, 4)
				_ = conn.(*websocketGateway.Conn).SendBinary(streamID, []byte{1, 2, 3, 4})
				// c.So(err, ShouldBeNil)
			}
			gw.Run()
			time.Sleep(time.Second)
		})

		Convey("Run the Clients", func(c C) {
			var td, tc int64
			wg := sync.WaitGroup{}
			for j := 0; j < 500; j++ {
				wg.Add(1)
				go func(j int) {
					var conn net.Conn
					var err error
					for try := 0; try < 5; try++ {
						conn, _, _, err = ws.Dial(context.Background(), "ws://localhost:81")
						if err == nil {
							break
						}
					}
					c.So(err, ShouldBeNil)
					m := make([]wsutil.Message, 0, 10)
					for i := 0; i < 50; i++ {
						err := wsutil.WriteMessage(conn, ws.StateClientSide, ws.OpBinary, tools.StrToByte("ABCD"))
						if err != nil {
							c.Println(err)
						}
						// c.So(err, ShouldBeNil)
						st := time.Now()
						_ = conn.SetDeadline(time.Now().Add(time.Second))
						m, err = wsutil.ReadMessage(conn, ws.StateClientSide, m)
						atomic.AddInt64(&tc, 1)
						atomic.AddInt64(&td, int64(time.Now().Sub(st)))
					}

					_ = conn.Close()
					wg.Done()
				}(j)
			}
			wg.Wait()
			c.Println("Total:", tc)
			c.Println("Average", time.Duration(td/tc))
		})

		// Convey("Run Idle Client", func(c C) {
		// 	conn, _, _, err := ws.Dial(context.Background(), "ws://localhost:81")
		// 	c.So(err, ShouldBeNil)
		// 	_, _ = wsutil.ReadServerMessage(conn, nil)
		// 	_ = conn.Close()
		// })
		Convey("Wait To Finish", func(c C) {
			time.Sleep(time.Second * 3)
		})
	})
}

func BenchmarkGatewaySerial(b *testing.B) {
	var err error
	gw, err = websocketGateway.New(websocketGateway.Config{
		MaxIdleTime:          time.Second * 3,
		NewConnectionWorkers: 10,
		MaxConcurrency:       1000,
		ListenAddress:        ":81",
	})
	if err != nil {
		b.Fatal(err)
	}

	gw.MessageHandler = func(conn gateway.Conn, streamID int64, data []byte, kvs ...gateway.KeyValue) {
		_ = conn.SendBinary(streamID, []byte{1, 2, 3, 4})
	}
	gw.Run()
	time.Sleep(time.Second)
	var (
		m    []wsutil.Message
		conn net.Conn
	)
	for try := 0; try < 5; try++ {
		conn, _, _, err = ws.Dial(context.Background(), "ws://localhost:81")
		if err == nil {
			break
		}
	}
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = wsutil.WriteMessage(conn, ws.StateClientSide, ws.OpBinary, tools.StrToByte("ABdwdwewefwefwefwewefwefwefweCD"))
		if err != nil {
			b.Fatal(err)
		}
		_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
		m, _ = wsutil.ReadMessage(conn, ws.StateClientSide, m)
		for idx := range m {
			pools.Bytes.Put(m[idx].Payload)
		}
		m = m[:0]
	}
	b.StopTimer()
}

func BenchmarkGatewayParallel(b *testing.B) {
	b.SetParallelism(10)
	var err error
	gw, err = websocketGateway.New(websocketGateway.Config{
		MaxIdleTime:          time.Second * 30,
		NewConnectionWorkers: 1,
		MaxConcurrency:       runtime.NumCPU() * 2,
		ListenAddress:        ":81",
	})
	if err != nil {
		b.Fatal(err)
	}

	gw.MessageHandler = func(conn gateway.Conn, streamID int64, data []byte, kvs ...gateway.KeyValue) {
		_ = conn.SendBinary(streamID, []byte{1, 2, 3, 4})
	}
	gw.Run()
	time.Sleep(time.Second)

	sendMessage := []byte(tools.RandomID(128))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var (
			m    []wsutil.Message
			err  error
			conn net.Conn
		)
		for try := 0; try < 5; try++ {
			conn, _, _, err = ws.Dial(context.Background(), "ws://localhost:81")
			if err == nil {
				break
			}
		}
		if err != nil {
			b.Fatal(err)
		}
		defer conn.Close()
		for pb.Next() {
			err = wsutil.WriteMessage(conn, ws.StateClientSide, ws.OpBinary, sendMessage)
			if err != nil {
				b.Fatal(err)
			}
			_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
			m, err = wsutil.ReadMessage(conn, ws.StateClientSide, m)
			if err != nil {
				b.Fatal(err)
			}

			for idx := range m {
				pools.Bytes.Put(m[idx].Payload)
			}
			m = m[:0]
		}
	})
	b.StopTimer()
}
