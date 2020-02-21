package quicGateway_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gateway"
	quicGateway "git.ronaksoftware.com/ronak/rony/gateway/quic"
	"git.ronaksoftware.com/ronak/rony/gateway/quic/testdata"
	log "git.ronaksoftware.com/ronak/rony/logger"
	"git.ronaksoftware.com/ronak/rony/metrics"
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/lucas-clemente/quic-go"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)

/*
   Creation Time: 2019 - Aug - 26
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	cntServerReceived int32
	cntServerFlush    int32
	gw                *quicGateway.Gateway
	cert              tls.Certificate
)

func init() {
	log.InitLogger(log.WarnLevel, "")
	metrics.Run("Test", "Test", 8000)
	testdata.GenerateSelfSignedCerts("./_key.pem", "./_cert.pem")
	cert = testdata.GetCertificate("./_key.pem", "./_cert.pem")
}

func TestGateway(t *testing.T) {
	Convey("Test QuicGateway", t, func(c C) {
		Convey("Run Server", func(c C) {
			var err error
			gw, err = quicGateway.New(quicGateway.Config{
				CloseHandler: func(c gateway.Conn) {},
				MessageHandler: func(conn gateway.Conn, streamID int64, data []byte) {
					c.So(data, ShouldResemble, tools.StrToByte(fmt.Sprintf("%v", streamID)))
					_ = conn.(*quicGateway.Conn).SendBinary(quic.StreamID(streamID), []byte{1, 2, 3, 4})
				},
				ConnectHandler: func(connID uint64) {
					// _Log.Info("New Connection", zap.Uint64("ConnID", connID))
				},
				FlushFunc: func(c gateway.Conn) [][]byte {
					// return [][]byte{{1}, {2}}
					return nil
				},
				NewConnectionWorkers: 1,
				MaxConcurrency:       1,
				ListenAddress:        ":82",
				Certificates:         []tls.Certificate{cert},
			})
			if err != nil {
				c.So(err, ShouldBeNil)
			}
			go gw.Run()
			time.Sleep(time.Second)
		})

		Convey("Run Clients", func(c C) {
			wg := sync.WaitGroup{}
			for j := 0; j < 10; j++ {
				wg.Add(1)
				go func() {
					client, err := quic.DialAddr("localhost:82", &tls.Config{
						Certificates:       []tls.Certificate{cert},
						InsecureSkipVerify: true,
						NextProtos:         []string{"RonakProto"},
					}, &quic.Config{
						KeepAlive: true,
					})
					c.So(err, ShouldBeNil)
					waitGroup := sync.WaitGroup{}
					for i := 0; i < 10; i++ {
						waitGroup.Add(1)
						go func() {
							s, err := client.OpenStreamSync(context.Background())
							c.So(err, ShouldBeNil)
							_, err = s.Write(tools.StrToByte(fmt.Sprintf("%v", s.StreamID())))
							c.So(err, ShouldBeNil)

							// br, _ := ioutil.ReadAll(s)
							// c.So(br, ShouldHaveLength, 4)
							// _ = br
							time.Sleep(200 * time.Millisecond)
							_ = s.Close()
							waitGroup.Done()
						}()
					}
					waitGroup.Wait()
					_ = client.Close()
					wg.Done()
				}()
			}
			wg.Wait()
		})
	})
}
