package httpGateway

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/valyala/fasthttp"
	"github.com/valyala/tcplisten"
	"net"
	"net/http"
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
	Concurrency   int
	ListenAddress string
	MaxBodySize   int
}

// Gateway
type Gateway struct {
	gateway.MessageHandler

	// Internal Controlling Params
	listenOn    string
	concurrency int
	maxBodySize int
}

// New
func New(config Config) *Gateway {
	g := new(Gateway)
	g.listenOn = config.ListenAddress
	g.concurrency = config.Concurrency
	g.maxBodySize = config.MaxBodySize
	g.MessageHandler = func(conn gateway.Conn, streamID int64, data []byte) {
		fmt.Println("Request Received")
	}
	return g
}

// Run
func (g *Gateway) Run() {
	tcpConfig := tcplisten.Config{
		ReusePort:   false,
		FastOpen:    true,
		DeferAccept: true,
		Backlog:     2048,
	}
	listener, err := tcpConfig.NewListener("tcp4", g.listenOn)
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
	req.SetConnectionClose()
	// ByPass CORS (Cross Origin Resource Sharing) check
	req.Response.Header.Set("Access-Control-Allow-Origin", "*")
	req.Response.Header.Set("Access-Control-Request-Method", "POST, GET, OPTIONS")
	req.Response.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if !req.Request.Header.IsPost() {
		req.SetStatusCode(http.StatusOK)
		return
	}

	conn := acquireConn(g, req)

	var detected bool
	req.Request.Header.VisitAll(func(key, value []byte) {
		switch tools.ByteToStr(key) {
		case "Cf-Connecting-Ip":
			conn.ClientIP = append(conn.ClientIP[:0], value...)
			detected = true
		case "X-Forwarded-For", "X-Real-Ip", "Forwarded":
			if !detected {
				conn.ClientIP = append(conn.ClientIP[:0], value...)
				detected = true
			}
		case "X-Client-Type":
			conn.ClientType = append(conn.ClientType[:0], value...)
		}
	})
	if !detected {
		conn.ClientIP = append(conn.ClientIP, req.RemoteIP().To4()...)
	}

	g.MessageHandler(conn, int64(req.ID()), req.PostBody())
	releaseConn(conn)
}

func (g *Gateway) Shutdown() {}

// Addr return the address which gateway is listen on
func (g *Gateway) Addr() string {
	return g.listenOn
}
