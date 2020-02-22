package httpGateway

import (
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/pool/pbytes"
	"github.com/valyala/fasthttp"
	"github.com/valyala/tcplisten"
	"net"
	"net/http"
	"time"
)

/*
   Creation Time: 2019 - Feb - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Config
type Config struct {
	Concurrency    int
	RequestTimeout time.Duration
	ListenAddress  string
	MaxBodySize    int
	gateway.MessageHandler
	gateway.FlushFunc
}

// Gateway
type Gateway struct {
	gateway.MessageHandler
	gateway.FlushFunc

	// Internal Controlling Params
	listenAddress string
	concurrency   int
	maxBodySize   int
	reqTimeout    time.Duration
}

// New
func New(config Config) *Gateway {
	g := new(Gateway)
	g.listenAddress = config.ListenAddress
	g.reqTimeout = config.RequestTimeout
	g.concurrency = config.Concurrency
	g.maxBodySize = config.MaxBodySize

	if config.MessageHandler == nil {
		g.MessageHandler = func(conn gateway.Conn, streamID int64, data []byte) {}
	} else {
		g.MessageHandler = config.MessageHandler
	}

	return g
}

// Run
func (g *Gateway) Run() {
	tcpConfig := tcplisten.Config{
		ReusePort:   false,
		FastOpen:    true,
		DeferAccept: false,
		Backlog:     8192,
	}
	listener, err := tcpConfig.NewListener("tcp4", g.listenAddress)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		server := fasthttp.Server{
			Name:               "River Edge Server",
			Concurrency:        g.concurrency,
			Handler:            g.requestHandler,
			DisableKeepalive:   true,
			MaxRequestBodySize: g.maxBodySize,
			// MaxConnsPerIP:    500,
		}
		for {
			err := server.Serve(listener)
			if err != nil {
				if nErr, ok := err.(net.Error); ok {
					if !nErr.Temporary() {
						return
					}
				} else {
					return
				}
			}
		}

	}()
}

func (g *Gateway) requestHandler(req *fasthttp.RequestCtx) {
	// ByPass CORS (Cross Origin Resource Sharing) check
	req.Response.Header.Set("Access-Control-Allow-Origin", "*")
	req.Response.Header.Set("Access-Control-Request-Method", "POST, GET, OPTIONS")
	req.Response.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if !req.Request.Header.IsPost() {
		req.SetStatusCode(http.StatusOK)
		return
	}

	var clientIP []byte
	var clientType string
	var detected bool
	req.Request.Header.VisitAll(func(key, value []byte) {
		switch tools.ByteToStr(key) {
		case "Cf-Connecting-Ip":
			clientIP = pbytes.GetLen(len(value))
			copy(clientIP, value)
			detected = true
		case "X-Forwarded-For", "X-Real-Ip", "Forwarded":
			if !detected {
				clientIP = pbytes.GetLen(len(value))
				copy(clientIP, value)
				detected = true
			}
		case "X-Client-Type":
			clientType = string(value)
		}
	})
	if !detected {
		clientIP = req.RemoteIP().To4()
	}
	conn := &Conn{
		gateway:    g,
		req:        req,
		ConnID:     req.ConnID(),
		ClientIP:   tools.ByteToStr(clientIP),
		ClientType: clientType,
	}
	g.MessageHandler(conn, int64(req.ID()), req.PostBody())

}
