package tcpGateway_test

import (
	"context"
	"git.ronaksoft.com/ronak/rony/gateway"
	tcpGateway "git.ronaksoft.com/ronak/rony/gateway/tcp"
	wsutil "git.ronaksoft.com/ronak/rony/gateway/tcp/util"
	log "git.ronaksoft.com/ronak/rony/internal/logger"
	"git.ronaksoft.com/ronak/rony/internal/testEnv"
	"git.ronaksoft.com/ronak/rony/pools"
	"git.ronaksoft.com/ronak/rony/tools"
	"github.com/gobwas/ws"
	"net"
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

var gw *tcpGateway.Gateway

func init() {
	testEnv.Init()
}

// func TestWebsocketConn(t *testing.T) {
//
// }
//
// func TestHttpConn(t *testing.T) {
//
// }
//
// func BenchmarkHttpConn(b *testing.B) {
// 	log.SetLevel(log.InfoLevel)
//
// }

func BenchmarkWebsocketConn(b *testing.B) {
	log.SetLevel(log.InfoLevel)
	b.SetParallelism(10)
	var err error
	gw, err = tcpGateway.New(tcpGateway.Config{
		Concurrency:   1000,
		ListenAddress: ":82",
		MaxBodySize:   0,
		MaxIdleTime:   time.Second * 30,
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
			conn, _, _, err = ws.Dial(context.Background(), "ws://localhost:82")
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
