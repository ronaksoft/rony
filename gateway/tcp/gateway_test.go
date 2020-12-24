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
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
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

func TestWebsocketConn(t *testing.T) {
	rony.SetLogLevel(0)
	gw, err := tcpGateway.New(tcpGateway.Config{
		Concurrency:   10,
		ListenAddress: "0.0.0.0:88",
		MaxBodySize:   0,
		MaxIdleTime:   0,
		Protocol:      tcpGateway.Auto,
		ExternalAddrs: []string{"127.0.0.1:88"},
	})

	if err != nil {
		t.Fatal(err)
	}
	gw.MessageHandler = func(c gateway.Conn, streamID int64, data []byte) {
		e := &rony.MessageEnvelope{}
		_ = e.Unmarshal(data)
		out, _ := e.Marshal()
		err := c.SendBinary(streamID, out)
		if err != nil {
			fmt.Println("MessageHandler:", err.Error())
		}
	}
	gw.Start()
	defer gw.Shutdown()
	Convey("Websocket Server Test", t, func(c C) {
		for i := 0; i < 10000; i++ {
			wsc := edgec.NewWebsocket(edgec.Config{
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
			ctx, cf := context.WithTimeout(context.TODO(), time.Second)
			err = wsc.SendWithContext(ctx, req, res)
			c.So(err, ShouldBeNil)
			cf()
			err = wsc.Close()
			c.So(err, ShouldBeNil)
		}
	})
}

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

	gw.MessageHandler = func(conn gateway.Conn, streamID int64, data []byte) {
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
