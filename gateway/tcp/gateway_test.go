package tcpGateway_test

import (
	"context"
	"git.ronaksoftware.com/ronak/rony/gateway"
	tcpGateway "git.ronaksoftware.com/ronak/rony/gateway/tcp"
	wsutil "git.ronaksoftware.com/ronak/rony/gateway/tcp/util"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/ws"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
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

var g *tcpGateway.Gateway

func init() {
	log.InitLogger(log.Config{
		Level: log.DebugLevel,
	})

	g = tcpGateway.MustNew(tcpGateway.Config{
		Concurrency:   10,
		ListenAddress: ":80",
		MaxBodySize:   100,
		MaxIdleTime:   30 * time.Second,
	})
	g.MessageHandler = func(c gateway.Conn, streamID int64, data []byte, kvs ...gateway.KeyValue) {
		log.Debug("Received Message",
			zap.Uint64("ConnID", c.GetConnID()),
			zap.Bool("Persistent", c.Persistent()),
			zap.String("IP", c.GetClientIP()),
		)
	}
	g.ConnectHandler = func(c gateway.Conn) {
		log.Debug("Connection Opened",
			zap.Uint64("ConnID", c.GetConnID()),
			zap.Bool("Persistent", c.Persistent()),
		)
	}
	g.CloseHandler = func(c gateway.Conn) {
		log.Debug("Connection Closed",
			zap.Uint64("ConnID", c.GetConnID()),
			zap.Bool("Persistent", c.Persistent()),
		)
	}
	g.Run()
}

func TestWebsocketConn(t *testing.T) {
	c, _, _, err := ws.Dial(context.Background(), "ws://localhost")
	if err != nil {
		panic(err)
	}
	err = wsutil.WriteMessage(c, ws.StateClientSide, ws.OpBinary, tools.StrToByte(tools.RandomDigit(10)))
	if err != nil {
		log.Warn("Error On Send Websocket Message", zap.Error(err))
	}
	_ = c.Close()
	time.Sleep(time.Second)
}
func TestHttpConn(t *testing.T) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	req.SetHost("127.0.0.1")
	err := fasthttp.Do(req, res)
	if err != nil {
		t.Fatal(err)
	}
}
