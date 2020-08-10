package tcpGateway

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gateway"
	wsutil "git.ronaksoftware.com/ronak/rony/gateway/tcp/util"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"github.com/valyala/fasthttp"
	"github.com/valyala/tcplisten"
	"go.uber.org/zap"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
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
	Concurrency   int
	ListenAddress string
	MaxBodySize   int
	MaxIdleTime   time.Duration
}

// Gateway
type Gateway struct {
	gateway.ConnectHandler
	gateway.MessageHandler
	gateway.CloseHandler

	// Internals
	listenOn    string
	listener    net.Listener
	addrs       []string
	concurrency int
	maxBodySize int

	// Websocket Internals
	upgradeHandler     ws.Upgrader
	maxIdleTime        int64
	conns              map[uint64]*WebsocketConn
	connsMtx           sync.RWMutex
	connsTotal         int32
	connsLastID        uint64
	connGC             *websocketConnGC
	poller             netpoll.Poller
	stop               int32
	waitGroupAcceptors *sync.WaitGroup
	waitGroupReaders   *sync.WaitGroup
	waitGroupWriters   *sync.WaitGroup
	cntReads           uint64
	cntWrites          uint64
}

// New
func New(config Config) (*Gateway, error) {
	var (
		err error
	)
	g := &Gateway{
		listenOn:           config.ListenAddress,
		concurrency:        config.Concurrency,
		maxBodySize:        config.MaxBodySize,
		maxIdleTime:        int64(defaultConnIdleTime.Seconds()),
		waitGroupReaders:   &sync.WaitGroup{},
		waitGroupWriters:   &sync.WaitGroup{},
		waitGroupAcceptors: &sync.WaitGroup{},
		conns:              make(map[uint64]*WebsocketConn, 100000),
	}

	tcpConfig := tcplisten.Config{
		ReusePort:   false,
		FastOpen:    true,
		DeferAccept: true,
		Backlog:     2048,
	}

	g.listener, err = tcpConfig.NewListener("tcp4", g.listenOn)
	if err != nil {
		return nil, err
	}

	if config.MaxIdleTime != 0 {
		g.maxIdleTime = int64(config.MaxIdleTime)
	}

	// initialize websocket upgrade handler
	g.upgradeHandler = ws.DefaultUpgrader

	// initialize idle websocket garbage collector
	g.connGC = newWebsocketConnGC(g)

	// set handlers
	g.MessageHandler = func(c gateway.Conn, streamID int64, date []byte, kvs ...gateway.KeyValue) {}
	g.CloseHandler = func(c gateway.Conn) {}
	g.ConnectHandler = func(c gateway.Conn) {}
	if poller, err := netpoll.New(&netpoll.Config{
		OnWaitError: func(e error) {
			log.Warn("Error On NetPoller Wait",
				zap.Error(e),
			)
		},
	}); err != nil {
		return nil, err
	} else {
		g.poller = poller
	}

	// try to detect the ip address of the listener
	ta, err := net.ResolveTCPAddr("tcp4", g.listener.Addr().String())
	if err != nil {
		return nil, err
	}
	if ta.IP.IsUnspecified() {
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, a := range addrs {
				switch x := a.(type) {
				case *net.IPNet:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					g.addrs = append(g.addrs, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				case *net.IPAddr:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					g.addrs = append(g.addrs, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				case *net.TCPAddr:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					g.addrs = append(g.addrs, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				}
			}
		}
	} else {
		g.addrs = append(g.addrs, fmt.Sprintf("%s:%d", ta.IP, ta.Port))
	}

	return g, nil
}

func MustNew(config Config) *Gateway {
	g, err := New(config)
	if err != nil {
		panic(err)
	}
	return g
}

// Run
func (g *Gateway) Run() {
	go func() {
		server := fasthttp.Server{
			Name:               "Rony TCP Gateway",
			Concurrency:        g.concurrency,
			Handler:            g.requestHandler,
			KeepHijackedConns:  true,
			MaxRequestBodySize: g.maxBodySize,
		}
		for {
			conn, err := g.listener.Accept()
			if err != nil {
				// log.Warn("Error On Accept", zap.Error(err))
				continue
			}

			wc := newWrapConn(conn)
			err = server.ServeConn(wc)
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

// Shutdown
func (g *Gateway) Shutdown() {
	// 1. Stop Accepting New Connections, i.e. Stop ConnectionAcceptor routines
	log.Info("Connection Acceptors are closing...")
	atomic.StoreInt32(&g.stop, 1)
	_ = g.listener.Close()
	g.waitGroupAcceptors.Wait()
	log.Info("Connection Acceptors all closed")

	// TODO:: why wait ?!
	time.Sleep(time.Second * 3)

	// 2. Close all readPumps
	log.Info("Read Pumpers are closing")
	g.waitGroupReaders.Wait()
	log.Info("Read Pumpers all closed")

	// 3. Close all writePumps
	log.Info("Write Pumpers are closing")
	g.waitGroupWriters.Wait()
	log.Info("Write Pumpers all closed")

	log.Info("Stats",
		zap.Uint64("Reads", g.cntReads),
		zap.Uint64("Writes", g.cntWrites),
	)
}

// Addr return the address which gateway is listen on
func (g *Gateway) Addr() []string {
	return g.addrs
}

func (g *Gateway) requestHandler(req *fasthttp.RequestCtx) {
	// ByPass CORS (Cross Origin Resource Sharing) check
	req.Response.Header.Set("Access-Control-Allow-Origin", "*")
	req.Response.Header.Set("Access-Control-Request-Method", "POST, GET, OPTIONS")
	req.Response.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if req.Request.Header.IsOptions() {
		req.SetStatusCode(http.StatusOK)
		return
	}

	var (
		kvs            = make([]gateway.KeyValue, 0, 4)
		clientIP       string
		clientType     string
		clientDetected bool
	)

	req.Request.Header.VisitAll(func(key, value []byte) {
		switch tools.ByteToStr(key) {
		case "Cf-Connecting-Ip":
			clientIP = string(value)
			clientDetected = true
		case "X-Forwarded-For", "X-Real-Ip", "Forwarded":
			if clientDetected {
				clientIP = string(value)
				clientDetected = true
			}
		case "X-Client-Type":
			clientType = string(value)
		default:
			kvs = append(kvs, gateway.KeyValue{
				Key:   string(key),
				Value: string(value),
			})
		}
	})
	if !clientDetected {
		clientIP = string(req.RemoteIP().To4())
	}

	if req.Request.Header.ConnectionUpgrade() {
		wc := req.Conn().(*wrapConn)
		wc.ReadyForUpgrade()
		req.HijackSetNoResponse(true)
		req.Hijack(func(c net.Conn) {
			g.waitGroupAcceptors.Add(1)
			g.connectionAcceptor(wc, clientIP, clientType)
		})
		return
	}

	req.SetConnectionClose()
	conn := acquireHttpConn(g, req)
	conn.ClientIP = append(conn.ClientIP[:0], tools.StrToByte(clientIP)...)
	conn.ClientType = append(conn.ClientType[:0], tools.StrToByte(clientType)...)
	g.MessageHandler(conn, int64(req.ID()), req.PostBody(), kvs...)
	releaseHttpConn(conn)
}

func (g *Gateway) connectionAcceptor(c net.Conn, clientIP, clientType string) {
	defer g.waitGroupAcceptors.Done()
	if atomic.LoadInt32(&g.stop) == 1 {
		return
	}
	if _, err := g.upgradeHandler.Upgrade(c); err != nil {
		if ce := log.Check(log.InfoLevel, "Error in Connection Acceptor"); ce != nil {
			ce.Write(
				zap.String("IP", clientIP),
				zap.String("ClientType", clientType),
				zap.Error(err),
			)
		}
		_ = c.Close()
		return
	}

	var (
		wsConn *WebsocketConn
		err    error
	)
	wsConn = g.addConnection(c, clientIP, clientType)
	wsConn.desc, err = netpoll.Handle(c,
		netpoll.EventRead|netpoll.EventHup|netpoll.EventOneShot,
	)

	if err != nil {
		log.Warn("Error On NetPoll Description", zap.Error(err))
		g.removeConnection(wsConn.ConnID)
		return
	}
	err = g.poller.Start(wsConn.desc, wsConn.startEvent)
	if err != nil {
		log.Warn("Error On NetPoll Start", zap.Error(err))
		g.removeConnection(wsConn.ConnID)
		return
	}
}

func (g *Gateway) addConnection(conn net.Conn, clientIP, clientType string) *WebsocketConn {
	// Increment total connection counter and connection ID
	totalConns := atomic.AddInt32(&g.connsTotal, 1)
	connID := atomic.AddUint64(&g.connsLastID, 1)
	wsConn := WebsocketConn{
		ClientIP:     clientIP,
		ConnID:       connID,
		conn:         conn,
		gateway:      g,
		lastActivity: tools.TimeUnix(),
		buf:          tools.NewLinkedList(),
	}
	g.connsMtx.Lock()
	g.conns[connID] = &wsConn
	g.connsMtx.Unlock()
	if ce := log.Check(log.DebugLevel, "Websocket Connection Created"); ce != nil {
		ce.Write(
			zap.Uint64("ConnID", connID),
			zap.String("Client", clientType),
			zap.String("IP", clientIP),
			zap.Int32("Total", totalConns),
		)
	}
	g.ConnectHandler(&wsConn)
	g.connGC.monitorConnection(connID)
	return &wsConn
}

func (g *Gateway) removeConnection(wcID uint64) {
	g.connsMtx.Lock()
	wsConn, ok := g.conns[wcID]
	if !ok {
		g.connsMtx.Unlock()
		return
	}
	delete(g.conns, wcID)
	g.connsMtx.Unlock()
	totalConns := atomic.AddInt32(&g.connsTotal, -1)
	if wsConn.desc != nil {
		_ = g.poller.Stop(wsConn.desc)
		_ = wsConn.desc.Close()
	}
	_ = wsConn.conn.Close()
	wsConn.Lock()
	if !wsConn.closed {
		g.CloseHandler(wsConn)
		wsConn.closed = true
	}
	wsConn.Unlock()
	if ce := log.Check(log.DebugLevel, "Websocket Connection Removed"); ce != nil {
		ce.Write(
			zap.Uint64("ConnID", wcID),
			zap.Int32("Total", totalConns),
		)
	}
}

func (g *Gateway) getConnection(connID uint64) *WebsocketConn {
	g.connsMtx.RLock()
	wsConn, ok := g.conns[connID]
	g.connsMtx.RUnlock()
	if ok {
		return wsConn
	}
	return nil
}

func (g *Gateway) readPump(wc *WebsocketConn) {
	defer g.waitGroupReaders.Done()
	var (
		err error
		ms  []wsutil.Message
	)

	_ = wc.conn.SetReadDeadline(time.Now().Add(defaultReadTimout))
	ms = ms[:0]
	ms, err = wsutil.ReadMessage(wc.conn, ws.StateServerSide, ms)
	if err != nil {
		if ce := log.Check(log.DebugLevel, "Error in readPump"); ce != nil {
			ce.Write(
				zap.Int64("AuthID", wc.AuthID),
				zap.Error(err),
			)
		}
		// remove the connection from the list
		wc.gateway.removeConnection(wc.ConnID)
		return
	}
	atomic.AddUint64(&g.cntReads, 1)
	_ = wc.gateway.poller.Resume(wc.desc)

	// handle messages
	for idx := range ms {
		switch ms[idx].OpCode {
		case ws.OpPong:
		case ws.OpPing:
			_ = wc.conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
			err = wsutil.WriteMessage(wc.conn, ws.StateServerSide, ws.OpPong, ms[idx].Payload)
			if err != nil {
				log.Warn("Error On Write OpPing", zap.Error(err))
			}
		case ws.OpBinary:
			g.MessageHandler(wc, 0, ms[idx].Payload)
		case ws.OpClose:
			// remove the connection from the list
			wc.gateway.removeConnection(wc.ConnID)
		default:
			log.Warn("Unknown OpCode")
		}
		pools.Bytes.Put(ms[idx].Payload)
	}
}

func (g *Gateway) writePump(wr writeRequest) {
	defer g.waitGroupWriters.Done()

	if wr.wc.closed {
		return
	}

	switch wr.opCode {
	case ws.OpBinary:
		// Try to write to the wire in WEBSOCKET_WRITE_SHORT_WAIT time
		_ = wr.wc.conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
		err := wsutil.WriteMessage(wr.wc.conn, ws.StateServerSide, ws.OpBinary, wr.payload)
		if err != nil {
			if ce := log.Check(log.WarnLevel, "Error in writePump"); ce != nil {
				ce.Write(zap.Error(err), zap.Uint64("ConnID", wr.wc.ConnID))
			}
			g.removeConnection(wr.wc.ConnID)
		} else {
			atomic.AddUint64(&g.cntWrites, 1)
		}
		// Put the bytes into the pool to be reused again
		pools.Bytes.Put(wr.payload)
	}

}
