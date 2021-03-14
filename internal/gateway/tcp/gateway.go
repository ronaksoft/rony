package tcpGateway

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"github.com/panjf2000/ants/v2"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
	wsutil "github.com/ronaksoft/rony/internal/gateway/tcp/util"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/tools"
	"github.com/valyala/fasthttp"
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

type UnsafeConn interface {
	net.Conn
	UnsafeConn() net.Conn
}

// Config
type Config struct {
	Concurrency   int
	ListenAddress string
	MaxBodySize   int
	MaxIdleTime   time.Duration
	Protocol      gateway.Protocol
	ExternalAddrs []string
	ProxyEnabled  bool
}

// Gateway
type Gateway struct {
	gateway.ConnectHandler
	gateway.MessageHandler
	gateway.CloseHandler

	// Internals
	transportMode gateway.Protocol
	listenOn      string
	listener      *wrapListener
	addrsMtx      sync.RWMutex
	addrs         []string
	extAddrs      []string
	concurrency   int
	maxBodySize   int

	// Http Internals
	httpProxy *HttpProxy

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
		maxIdleTime:        int64(defaultConnIdleTime),
		waitGroupReaders:   &sync.WaitGroup{},
		waitGroupWriters:   &sync.WaitGroup{},
		waitGroupAcceptors: &sync.WaitGroup{},
		conns:              make(map[uint64]*websocketConn, 100000),
		transportMode:      gateway.TCP,
		extAddrs:           config.ExternalAddrs,
	}

	if config.ProxyEnabled {
		g.httpProxy = &HttpProxy{}
	}

	g.listener, err = newWrapListener(g.listenOn)
	if err != nil {
		return nil, err
	}

	if config.MaxIdleTime != 0 {
		g.maxIdleTime = int64(config.MaxIdleTime)
	}
	if config.Protocol != gateway.Undefined {
		g.transportMode = config.Protocol
	}

	switch g.transportMode {
	case gateway.Websocket, gateway.Http, gateway.TCP:
	default:
		return nil, ErrUnsupportedProtocol
	}

	// initialize websocket upgrade handler
	g.upgradeHandler = ws.DefaultUpgrader

	// initialize idle websocket garbage collector
	g.connGC = newWebsocketConnGC(g)

	// set handlers
	g.MessageHandler = func(c rony.Conn, streamID int64, data []byte) {}
	g.CloseHandler = func(c rony.Conn) {}
	g.ConnectHandler = func(c rony.Conn, kvs ...*rony.KeyValue) {}
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
	err = g.detectAddrs()
	if err != nil {
		log.Warn("Rony:: Gateway got error on detecting addrs", zap.Error(err))
		return nil, err
	}

	goPoolB, err = ants.NewPool(g.concurrency,
		ants.WithNonblocking(false),
		ants.WithPreAlloc(true),
	)
	if err != nil {
		return nil, err
	}

	goPoolNB, err = ants.NewPool(g.concurrency,
		ants.WithNonblocking(true),
		ants.WithPreAlloc(true),
	)
	if err != nil {
		return nil, err
	}

	// run the watchdog in background
	go g.watchdog()

	return g, nil
}

func MustNew(config Config) *Gateway {
	g, err := New(config)
	if err != nil {
		panic(err)
	}
	return g
}

func (g *Gateway) watchdog() {
	for {
		metrics.SetGauge(metrics.GaugeActiveWebsocketConnections, float64(g.TotalConnections()))
		err := g.detectAddrs()
		if err != nil {
			log.Warn("Rony:: Gateway got error on detecting addrs", zap.Error(err))
		}
		time.Sleep(time.Second * 15)
	}

}

func (g *Gateway) detectAddrs() error {
	// try to detect the ip address of the listener
	ta, err := net.ResolveTCPAddr("tcp4", g.listener.Addr().String())
	if err != nil {
		return err
	}
	lAddrs := make([]string, 0, 10)
	if ta.IP.IsUnspecified() {
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, a := range addrs {
				switch x := a.(type) {
				case *net.IPNet:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					lAddrs = append(lAddrs, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				case *net.IPAddr:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					lAddrs = append(lAddrs, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				case *net.TCPAddr:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					lAddrs = append(lAddrs, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				}
			}
		}
	} else {

		lAddrs = append(lAddrs, fmt.Sprintf("%s:%d", ta.IP, ta.Port))

	}
	g.addrsMtx.Lock()
	g.addrs = append(g.addrs[:0], lAddrs...)
	g.addrsMtx.Unlock()
	return nil
}

// Start is non-blocking and call the Run function in background
func (g *Gateway) Start() {
	go g.Run()
}

// Run is blocking and runs the server endless loop until a non-temporary error happens
func (g *Gateway) Run() {
	// if http proxy is set then, HttpProxy is responsible for handling the request
	if g.httpProxy != nil {
		g.httpProxy.handler = g.MessageHandler
	}

	// initialize the fasthttp server.
	server := fasthttp.Server{
		Name:               "Rony TCP-Gateway",
		Handler:            g.requestHandler,
		Concurrency:        g.concurrency,
		KeepHijackedConns:  true,
		MaxRequestBodySize: g.maxBodySize,
		DisableKeepalive:   true,
		CloseOnShutdown:    true,
	}

	// start serving in blocking mode
	err := server.Serve(g.listener)
	if err != nil {
		log.Warn("Error On Serve", zap.Error(err))
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

	g.connsMtx.Lock()
	for id, c := range g.conns {
		fmt.Println("Conn(", id, ") Stalled",
			time.Duration(tools.CPUTicks()-atomic.LoadInt64(&c.startTime)),
			time.Duration(tools.CPUTicks()-(atomic.LoadInt64(&c.lastActivity))),
		)
	}
	g.connsMtx.Unlock()
}

// Addr return the address which gateway is listen on
func (g *Gateway) Addr() []string {
	if len(g.extAddrs) > 0 {
		return g.extAddrs
	}
	g.addrsMtx.RLock()
	addrs := g.addrs
	g.addrsMtx.RUnlock()
	return addrs
}

// GetConn returns the connection identified by connID
func (g *Gateway) GetConn(connID uint64) rony.Conn {
	c := g.getConnection(connID)
	if c == nil {
		return nil
	}
	return c
}

func (g *Gateway) TotalConnections() int {
	g.connsMtx.RLock()
	n := len(g.conns)
	g.connsMtx.RUnlock()
	return n
}

func (g *Gateway) requestHandler(reqCtx *gateway.RequestCtx) {
	// ByPass CORS (Cross Origin Resource Sharing) check
	// TODO:: let developer choose
	reqCtx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	reqCtx.Response.Header.Set("Access-Control-Request-Method", "*")
	reqCtx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	if reqCtx.Request.Header.IsOptions() {
		reqCtx.SetStatusCode(http.StatusOK)
		reqCtx.SetConnectionClose()
		return
	}

	// extract required information from the header of the RequestCtx
	meta := acquireConnInfo(reqCtx)

	// If this is a Http Upgrade then we handle websocket
	if reqCtx.Request.Header.ConnectionUpgrade() {
		if g.transportMode == gateway.Http {
			reqCtx.SetConnectionClose()
			reqCtx.SetStatusCode(http.StatusNotAcceptable)
			return
		}
		reqCtx.HijackSetNoResponse(true)
		reqCtx.Hijack(func(c net.Conn) {
			wc, _ := c.(UnsafeConn).UnsafeConn().(*wrapConn)
			wc.ReadyForUpgrade()
			g.waitGroupAcceptors.Add(1)
			g.websocketHandler(wc, meta)
			releaseConnInfo(meta)
		})
		return
	}

	// This is going to be an HTTP request
	reqCtx.SetConnectionClose()
	if g.transportMode == gateway.Websocket {
		reqCtx.SetStatusCode(http.StatusNotAcceptable)
		return
	}
	conn := acquireHttpConn(g, reqCtx)
	conn.SetClientIP(meta.clientIP)
	conn.SetClientType(meta.clientType)

	metrics.IncCounter(metrics.CntGatewayIncomingHttpMessage)

	g.ConnectHandler(conn, meta.kvs...)
	releaseConnInfo(meta)

	if g.httpProxy != nil {
		conn.proxy = g.httpProxy.search(tools.B2S(reqCtx.Method()), tools.B2S(reqCtx.Request.URI().PathOriginal()), conn)
		if conn.proxy == nil {
			g.MessageHandler(conn, int64(reqCtx.ID()), reqCtx.PostBody())
		} else {
			g.httpProxy.handle(conn, reqCtx)
		}
	} else {
		g.MessageHandler(conn, int64(reqCtx.ID()), reqCtx.PostBody())
	}

	g.CloseHandler(conn)
	releaseHttpConn(conn)
}

func (g *Gateway) websocketHandler(c net.Conn, meta *connInfo) {
	defer g.waitGroupAcceptors.Done()
	if atomic.LoadInt32(&g.stop) == 1 {
		return
	}
	if _, err := g.upgradeHandler.Upgrade(c); err != nil {
		if ce := log.Check(log.InfoLevel, "Error in Connection Acceptor"); ce != nil {
			ce.Write(
				zap.String("IP", tools.B2S(meta.clientIP)),
				zap.String("ClientType", tools.B2S(meta.clientType)),
				zap.Error(err),
			)
		}
		_ = c.Close()
		return
	}

	var (
		err error
	)

	wsConn, err := newWebsocketConn(g, c, meta.clientIP)
	if err != nil {
		log.Warn("Error On NetPoll Description", zap.Error(err), zap.Int("Total", g.TotalConnections()))
		return
	}

	g.ConnectHandler(wsConn, meta.kvs...)

	err = wsConn.registerDesc()
	if err != nil {
		log.Warn("Error On RegisterDesc", zap.Error(err))
	}
}

func (g *Gateway) websocketReadPump(wc *websocketConn, wg *sync.WaitGroup, ms []wsutil.Message) (err error) {
	ms = ms[:0]
	ms, err = wc.read(ms)
	if err != nil {
		if ce := log.Check(log.DebugLevel, "Error in websocketReadPump"); ce != nil {
			ce.Write(
				zap.Uint64("ConnID", wc.connID),
				zap.Error(err),
			)
		}
		return ErrUnexpectedSocketRead
	}
	atomic.AddUint64(&g.cntReads, 1)

	// handle messages
	for idx := range ms {
		switch ms[idx].OpCode {
		case ws.OpPong:
		case ws.OpPing:
			err = wc.write(ws.OpPing, ms[idx].Payload)
			if err != nil {
				log.Warn("Error On Write OpPing", zap.Error(err))
			}
		case ws.OpBinary:
			wg.Add(1)
			_ = goPoolB.Submit(func() {
				metrics.IncCounter(metrics.CntGatewayIncomingWebsocketMessage)
				g.MessageHandler(wc, 0, ms[idx].Payload)
				wg.Done()
			})

		case ws.OpClose:
			// remove the connection from the list
			err = ErrOpCloseReceived
		default:
			log.Warn("Unknown OpCode")
		}
	}

	return err
}

func (g *Gateway) websocketWritePump(wr *writeRequest) (err error) {
	defer g.waitGroupWriters.Done()

	switch wr.opCode {
	case ws.OpBinary, ws.OpText:
		err = wr.wc.write(wr.opCode, wr.payload)
		if err != nil {
			if ce := log.Check(log.DebugLevel, "Error in websocketWritePump"); ce != nil {
				ce.Write(zap.Error(err), zap.Uint64("ConnID", wr.wc.connID))
			}
		} else {
			atomic.AddUint64(&g.cntWrites, 1)
		}

	}
	return
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

func (g *Gateway) SetProxy(
	method, path string, handle gateway.ProxyHandle,
) {
	g.httpProxy.Set(method, path, handle)
}
