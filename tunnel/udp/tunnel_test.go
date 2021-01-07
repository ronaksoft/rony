package udpTunnel_test

import (
	"bufio"
	"fmt"
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

func TestNewTunnel(t *testing.T) {
	Convey("Tunnel", t, func(c C) {
		t := udpTunnel.New(udpTunnel.Config{
			Concurrency:   0,
			ListenAddress: "127.0.0.1:2374",
			ExternalAddrs: []string{"127.0.0.1:2374"},
		})
		t.Start()
		defer t.Shutdown()

		Convey("Send Data", func(c C) {
			conn, err := net.Dial("udp", "127.0.0.1:2374")
			c.So(err, ShouldBeNil)
			txt := "Hello this is Me"
			n, err := fmt.Fprint(conn, txt)
			c.So(err, ShouldBeNil)
			c.So(n, ShouldEqual, len(txt))
			p := make([]byte, 2048)
			n, err = bufio.NewReader(conn).Read(p)
			c.So(err, ShouldBeNil)
			c.So(n, ShouldEqual, len(txt))
			err = conn.Close()
			c.So(err, ShouldBeNil)
			time.Sleep(time.Second * 2)
		})
	})
}
