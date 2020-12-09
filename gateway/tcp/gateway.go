package tcp

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"github.com/ronaksoft/rony/gateway"
	wsutil "github.com/ronaksoft/rony/gateway/tcp/util"
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
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
   Copyright Ronak Software Group 2020
*/

type Protocol int

const (
	Http Protocol = 1 << iota
	Websocket
	Auto = Http | Websocket
)

// Config
type Config struct {
	Concurrency   int
	ListenAddress string
	MaxBodySize   int
	MaxIdleTime   time.Duration
	Protocol      Protocol
	ExternalAddrs []string
}

// Gateway
type Gateway struct {
	gateway.ConnectHandler
	gateway.MessageHandler
	gateway.CloseHandler

	// Internals
	transportMode Protocol
	listenOn      string
	listener      net.Listener
	addrs         []string
	extAddrs      []string
	concurrency   int
	maxBodySize   int

	// Websocket Internals
	upgradeHandler     ws.Upgrader
	maxIdleTime        int64
	conns              map[uint64]*websocketConn
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
		conns:              make(map[uint64]*websocketConn, 100000),
		transportMode:      Auto,
		extAddrs:           config.ExternalAddrs,
	}

	tcpConfig := tcplisten.Config{
		ReusePort:   true,
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
	if config.Protocol != 0 {
		g.transportMode = config.Protocol
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

// Start is non-blocking and call the Run function in background
func (g *Gateway) Start() {
	go g.Run()
}

// Run is blocking and runs the server endless loop until a non-temporary error happens
func (g *Gateway) Run() {
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
	if len(g.extAddrs) > 0 {
		return g.extAddrs
	}
	return g.addrs
}

// GetConn returns the connection identified by connID
func (g *Gateway) GetConn(connID uint64) gateway.Conn {
	return g.getConnection(connID)
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
		if g.transportMode&Websocket == 0 {
			req.SetConnectionClose()
			req.SetStatusCode(http.StatusNotAcceptable)
			return
		}
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
	if g.transportMode&Http == 0 {
		req.SetStatusCode(http.StatusNotAcceptable)
		return
	}
	conn := acquireHttpConn(g, req)
	conn.SetClientIP(tools.StrToByte(clientIP))
	conn.SetClientType(tools.StrToByte(clientType))
	g.ConnectHandler(conn)
	g.MessageHandler(conn, int64(req.ID()), req.PostBody(), kvs...)
	g.CloseHandler(conn)
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
		err error
	)
	wsConn := g.addConnection(c, clientIP, clientType)
	wsConn.desc, err = netpoll.Handle(c,
		netpoll.EventRead|netpoll.EventHup|netpoll.EventOneShot,
	)

	if err != nil {
		log.Warn("Error On NetPoll Description", zap.Error(err))
		g.removeConnection(wsConn.connID)
		return
	}
	err = g.poller.Start(wsConn.desc, wsConn.startEvent)
	if err != nil {
		log.Warn("Error On NetPoll Start", zap.Error(err))
		g.removeConnection(wsConn.connID)
		return
	}
}

func (g *Gateway) addConnection(conn net.Conn, clientIP, clientType string) *websocketConn {
	// Increment total connection counter and connection ID
	totalConns := atomic.AddInt32(&g.connsTotal, 1)
	connID := atomic.AddUint64(&g.connsLastID, 1)
	wsConn := acquireWebsocketConn(g, connID, conn, nil)
	wsConn.SetClientIP(tools.StrToByte(clientIP))

	g.connsMtx.Lock()
	g.conns[connID] = wsConn
	g.connsMtx.Unlock()
	if ce := log.Check(log.DebugLevel, "Websocket Connection Created"); ce != nil {
		ce.Write(
			zap.Uint64("ConnID", connID),
			zap.String("Client", clientType),
			zap.String("IP", wsConn.ClientIP()),
			zap.Int32("Total", totalConns),
		)
	}
	g.ConnectHandler(wsConn)
	g.connGC.monitorConnection(connID)
	return wsConn
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
	releaseWebsocketConn(wsConn)

	if ce := log.Check(log.DebugLevel, "Websocket Connection Removed"); ce != nil {
		ce.Write(
			zap.Uint64("ConnID", wcID),
			zap.Int32("Total", totalConns),
		)
	}
}

func (g *Gateway) getConnection(connID uint64) *websocketConn {
	g.connsMtx.RLock()
	wsConn, ok := g.conns[connID]
	g.connsMtx.RUnlock()
	if ok {
		return wsConn
	}
	return nil
}

func (g *Gateway) readPump(wc *websocketConn, ms []wsutil.Message) {
	defer g.waitGroupReaders.Done()
	var (
		err error
	)

	waitGroup := pools.AcquireWaitGroup()
	_ = wc.conn.SetReadDeadline(time.Now().Add(defaultReadTimout))
	ms = ms[:0]
	ms, err = wsutil.ReadMessage(wc.conn, ws.StateServerSide, ms)
	if err != nil {
		if ce := log.Check(log.DebugLevel, "Error in readPump"); ce != nil {
			ce.Write(
				zap.Uint64("ConnID", wc.connID),
				zap.Error(err),
			)
		}
		// remove the connection from the list
		wc.gateway.removeConnection(wc.connID)
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
			waitGroup.Add(1)
			go func() {
				g.MessageHandler(wc, 0, ms[idx].Payload)
				waitGroup.Done()
			}()
		case ws.OpClose:
			// remove the connection from the list
			wc.gateway.removeConnection(wc.connID)
		default:
			log.Warn("Unknown OpCode")
		}
		pools.Bytes.Put(ms[idx].Payload)
	}
	waitGroup.Wait()
	pools.ReleaseWaitGroup(waitGroup)
}

func (g *Gateway) writePump(wr *writeRequest) {
	defer g.waitGroupWriters.Done()

	if wr.wc.closed {
		return
	}

	switch wr.opCode {
	case ws.OpBinary, ws.OpText:
		_ = wr.wc.conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
		err := wsutil.WriteMessage(wr.wc.conn, ws.StateServerSide, wr.opCode, wr.payload)
		if err != nil {
			if ce := log.Check(log.WarnLevel, "Error in writePump"); ce != nil {
				ce.Write(zap.Error(err), zap.Uint64("ConnID", wr.wc.connID))
			}
			g.removeConnection(wr.wc.connID)
		} else {
			atomic.AddUint64(&g.cntWrites, 1)
		}
	}

}
