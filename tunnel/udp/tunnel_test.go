package udpTunnel_test

import (
	"bufio"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/tools"
	udpTunnel "github.com/ronaksoft/rony/tunnel/udp"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
	"time"
)

/*
   Creation Time: 2021 - Jan - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func init() {
	testEnv.Init()
	log.SetLevel(log.DebugLevel)
}

func TestNewTunnel(t *testing.T) {
	Convey("Tunnel", t, func(c C) {
		t := udpTunnel.New(udpTunnel.Config{
			Concurrency:   0,
			ListenAddress: "127.0.0.1:2374",
			ExternalAddrs: []string{"127.0.0.1:2374"},
		})
		t.MessageHandler = func(conn rony.Conn, tm *rony.TunnelMessage) {
			b, _ := tm.Marshal()
			err := conn.SendBinary(0, b)
			c.So(err, ShouldBeNil)
		}
		t.Start()
		defer t.Shutdown()

		Convey("Send Data", func(c C) {
			conn, err := net.Dial("udp", "127.0.0.1:2374")
			c.So(err, ShouldBeNil)

			req := &rony.TunnelMessage{
				SenderID:         []byte("SomeSender"),
				SenderReplicaSet: tools.RandomUint64(1000),
				Store:            nil,
				Envelope:         nil,
			}
			res := &rony.TunnelMessage{}
			out, _ := req.Marshal()
			n, err := conn.Write(out)
			c.So(err, ShouldBeNil)
			c.So(n, ShouldEqual, len(out))
			p := make([]byte, 2048)
			n, err = bufio.NewReader(conn).Read(p)
			c.So(err, ShouldBeNil)
			c.So(n, ShouldEqual, len(out))
			err = res.Unmarshal(p[:n])
			c.So(err, ShouldBeNil)
			c.So(res.SenderID, ShouldResemble, req.SenderID)
			c.So(res.SenderReplicaSet, ShouldEqual, req.SenderReplicaSet)
			err = conn.Close()
			c.So(err, ShouldBeNil)
			time.Sleep(time.Second * 2)
		})
	})
}
