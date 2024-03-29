package tcpGateway

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"github.com/panjf2000/ants/v2"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/internal/gateway/tcp/cors"
	wsutil "github.com/ronaksoft/rony/internal/gateway/tcp/util"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
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

// Config holds all the configuration for Gateway
type Config struct {
	Concurrency   int
	ListenAddress string
	MaxBodySize   int
	MaxIdleTime   time.Duration
	Protocol      rony.GatewayProtocol
	ExternalAddrs []string
	Logger        log.Logger
	// TextDataFrame if is set to TRUE then websocket data frames use OpText otherwise use OpBinary
	TextDataFrame bool

	// CORS
	AllowedHeaders []string // Default Allow All
	AllowedOrigins []string // Default Allow All
	AllowedMethods []string // Default Allow All
}

// Gateway is one of the main components of the Rony framework. Basically Gateway is the component
// that connects edge.Server with the external world. Clients which are not part of our cluster MUST
// connect to our edge servers through Gateway.
// This is an implementation of gateway.Gateway interface with support for **Http** and **Websocket** connections.
type Gateway struct {
	// Internals
	cfg                Config
	transportMode      rony.GatewayProtocol
	listener           *wrapListener
	listenerAddressMtx sync.RWMutex
	listenerAddresses  []string
	poller             netpoll.Poller
	stop               int32
	waitGroupAcceptors *sync.WaitGroup
	waitGroupReaders   *sync.WaitGroup
	waitGroupWriters   *sync.WaitGroup
	cntReads           uint64
	cntWrites          uint64
	cors               *cors.CORS
	delegate           rony.GatewayDelegate

	// Websocket Internals
	upgradeHandler ws.Upgrader
	connGC         *websocketConnGC
	maxIdleTime    int64
	conns          map[uint64]*websocketConn
	connsMtx       sync.RWMutex
	connsTotal     int32
	connsLastID    uint64
}

func New(config Config) (*Gateway, error) {
	var err error

	if config.Logger == nil {
		config.Logger = log.DefaultLogger
	}

	g := &Gateway{
		cfg:                config,
		maxIdleTime:        int64(defaultConnIdleTime),
		waitGroupReaders:   &sync.WaitGroup{},
		waitGroupWriters:   &sync.WaitGroup{},
		waitGroupAcceptors: &sync.WaitGroup{},
		conns:              make(map[uint64]*websocketConn, 100000),
		transportMode:      rony.TCP,
		cors: cors.New(cors.Config{
			AllowedHeaders: config.AllowedHeaders,
			AllowedMethods: config.AllowedMethods,
			AllowedOrigins: config.AllowedOrigins,
		}),
	}

	g.listener, err = newWrapListener(g.cfg.ListenAddress)
	if err != nil {
		return nil, err
	}

	if config.MaxIdleTime != 0 {
		g.maxIdleTime = int64(config.MaxIdleTime)
	}
	if config.Protocol != rony.Undefined {
		g.transportMode = config.Protocol
	}

	switch g.transportMode {
	case rony.Websocket, rony.Http, rony.TCP:
	default:
		return nil, ErrUnsupportedProtocol
	}

	// initialize websocket upgrade handler
	g.upgradeHandler = ws.DefaultUpgrader

	// initialize idle websocket garbage collector
	g.connGC = newWebsocketConnGC(g)

	// set handlers
	if poller, err := netpoll.New(&netpoll.Config{
		OnWaitError: func(e error) {
			g.cfg.Logger.Warn("Error On NetPoller Wait",
				zap.Error(e),
			)
		},
	}); err != nil {
		return nil, err
	} else {
		g.poller = poller
	}

	// try to detect the ip address of the listener
	err = g.detectListenerAddress()
	if err != nil {
		g.cfg.Logger.Warn("Rony:: Gateway got error on detecting listener addresses", zap.Error(err))

		return nil, err
	}

	goPoolB, err = ants.NewPool(g.cfg.Concurrency,
		ants.WithNonblocking(false),
		ants.WithPreAlloc(true),
	)
	if err != nil {
		return nil, err
	}

	goPoolNB, err = ants.NewPool(g.cfg.Concurrency,
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
		err := g.detectListenerAddress()
		if err != nil {
			g.cfg.Logger.Warn("Gateway got error on detecting listener address", zap.Error(err))
		}
		time.Sleep(time.Second * 15)
	}
}

func (g *Gateway) detectListenerAddress() error {
	// try to detect the ip address of the listener
	ta, err := net.ResolveTCPAddr("tcp4", g.listener.Addr().String())
	if err != nil {
		return err
	}
	listenerAddresses := make([]string, 0, 10)
	if ta.IP.IsUnspecified() {
		interfaceAddresses, err := net.InterfaceAddrs()
		if err == nil {
			for _, a := range interfaceAddresses {
				switch x := a.(type) {
				case *net.IPNet:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					listenerAddresses = append(listenerAddresses, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				case *net.IPAddr:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					listenerAddresses = append(listenerAddresses, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				case *net.TCPAddr:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					listenerAddresses = append(listenerAddresses, fmt.Sprintf("%s:%d", x.IP.String(), ta.Port))
				}
			}
		}
	} else {
		listenerAddresses = append(listenerAddresses, fmt.Sprintf("%s:%d", ta.IP, ta.Port))
	}
	g.listenerAddressMtx.Lock()
	g.listenerAddresses = append(g.listenerAddresses[:0], listenerAddresses...)
	g.listenerAddressMtx.Unlock()

	return nil
}

func (g *Gateway) Subscribe(d rony.GatewayDelegate) {
	g.delegate = d
}

// Start is non-blocking and call the Run function in background
func (g *Gateway) Start() {
	go g.Run()
}

// Run is blocking and runs the server endless loop until a non-temporary error happens
func (g *Gateway) Run() {
	// initialize the fasthttp server.
	server := fasthttp.Server{
		Name:               "Rony TCP-Gateway",
		Handler:            g.requestHandler,
		Concurrency:        g.cfg.Concurrency,
		KeepHijackedConns:  true,
		MaxRequestBodySize: g.cfg.MaxBodySize,
		DisableKeepalive:   true,
		CloseOnShutdown:    true,
	}

	// start serving in blocking mode
	err := server.Serve(g.listener)
	if err != nil {
		g.cfg.Logger.Warn("Error On Serve", zap.Error(err))
	}
}

// Shutdown closes the server by stopping services in sequence, in a way that all the flying request
// will be served before server shutdown.
func (g *Gateway) Shutdown() {
	// 1. Stop Accepting New Connections, i.e. Stop ConnectionAcceptor routines
	g.cfg.Logger.Info("Connection Acceptors are closing...")
	atomic.StoreInt32(&g.stop, 1)
	_ = g.listener.Close()
	g.waitGroupAcceptors.Wait()
	g.cfg.Logger.Info("Connection Acceptors all closed")

	// 2. Close all readPumps
	g.cfg.Logger.Info("Read Pumpers are closing")
	g.waitGroupReaders.Wait()
	g.cfg.Logger.Info("Read Pumpers all closed")

	// 3. Close all writePumps
	g.cfg.Logger.Info("Write Pumpers are closing")
	g.waitGroupWriters.Wait()
	g.cfg.Logger.Info("Write Pumpers all closed")

	g.cfg.Logger.Info("Stats",
		zap.Uint64("Reads", g.cntReads),
		zap.Uint64("Writes", g.cntWrites),
	)

	g.connsMtx.Lock()
	for id, c := range g.conns {
		g.cfg.Logger.Info("Conn Stalled",
			zap.Uint64("ID", id),
			zap.Duration("SinceStart", time.Duration(tools.CPUTicks()-atomic.LoadInt64(&c.startTime))),
			zap.Duration("SinceLastActivity", time.Duration(tools.CPUTicks()-(atomic.LoadInt64(&c.lastActivity)))),
		)
	}
	g.connsMtx.Unlock()
}

// Addr return the address which gateway is listen on
func (g *Gateway) Addr() []string {
	if len(g.cfg.ExternalAddrs) > 0 {
		return g.cfg.ExternalAddrs
	}
	g.listenerAddressMtx.RLock()
	addrs := g.listenerAddresses
	g.listenerAddressMtx.RUnlock()

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

func (g *Gateway) Support(p rony.GatewayProtocol) bool {
	return g.transportMode&p == p
}

func (g *Gateway) TotalConnections() int {
	g.connsMtx.RLock()
	n := len(g.conns)
	g.connsMtx.RUnlock()

	return n
}

func (g *Gateway) Protocol() rony.GatewayProtocol {
	return g.transportMode
}

func (g *Gateway) requestHandler(reqCtx *fasthttp.RequestCtx) {
	if g.cors.Handle(reqCtx) {
		return
	}

	// extract required information from the header of the RequestCtx
	connInfo := acquireConnInfo(reqCtx)

	// If this is a Http Upgrade then we Handle websocket
	if connInfo.Upgrade() {
		if !g.Support(rony.Websocket) {
			reqCtx.SetConnectionClose()
			reqCtx.SetStatusCode(http.StatusNotAcceptable)

			return
		}
		reqCtx.HijackSetNoResponse(true)
		reqCtx.Hijack(
			func(c net.Conn) {
				wc, _ := c.(UnsafeConn).UnsafeConn().(*wrapConn)
				wc.ReadyForUpgrade()
				g.waitGroupAcceptors.Add(1)
				g.websocketHandler(wc, connInfo)
				releaseConnInfo(connInfo)
			},
		)

		return
	}

	// This is going to be an HTTP request
	reqCtx.SetConnectionClose()
	if !g.Support(rony.Http) {
		reqCtx.SetStatusCode(http.StatusNotAcceptable)

		return
	}

	conn := acquireHttpConn(g, reqCtx)
	conn.SetClientIP(connInfo.clientIP)
	conn.SetClientType(connInfo.clientType)
	for k, v := range connInfo.kvs {
		conn.Set(k, v)
	}

	metrics.IncCounter(metrics.CntGatewayIncomingHttpMessage)

	g.delegate.OnConnect(conn)

	g.delegate.OnMessage(conn, int64(reqCtx.ID()), reqCtx.PostBody())

	g.delegate.OnClose(conn)

	releaseConnInfo(connInfo)
	releaseHttpConn(conn)
}

func (g *Gateway) websocketHandler(c net.Conn, meta *connInfo) {
	defer g.waitGroupAcceptors.Done()
	if atomic.LoadInt32(&g.stop) == 1 {
		return
	}
	if _, err := g.upgradeHandler.Upgrade(c); err != nil {
		if ce := g.cfg.Logger.Check(log.InfoLevel, "got error in websocket upgrade"); ce != nil {
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
		g.cfg.Logger.Warn("got error on creating websocket connection",
			zap.Error(err),
			zap.Int("Total", g.TotalConnections()),
		)

		return
	}
	for k, v := range meta.kvs {
		wsConn.Set(k, v)
	}

	g.delegate.OnConnect(wsConn)

	err = wsConn.registerDesc()
	if err != nil {
		g.cfg.Logger.Warn("got error in registering conn desc",
			zap.Error(err),
			zap.Any("Conn", wsConn.conn),
		)
	}
}

func (g *Gateway) websocketReadPump(wc *websocketConn, wg *sync.WaitGroup) (err error) {
	var ms []wsutil.Message
	ms, err = wc.read(ms)
	if err != nil {
		if ce := g.cfg.Logger.Check(log.DebugLevel, "got error in websocket read pump"); ce != nil {
			ce.Write(
				zap.Uint64("ConnID", wc.connID),
				zap.Error(err),
			)
		}

		return errors.Wrap(ErrUnexpectedSocketRead)(err)
	}
	atomic.AddUint64(&g.cntReads, 1)

	// Handle messages
	for idx := range ms {
		switch ms[idx].OpCode {
		case ws.OpPong:
		case ws.OpPing:
			err = wc.write(ws.OpPong, ms[idx].Payload)
			pools.Bytes.Put(ms[idx].Payload)
		case ws.OpBinary, ws.OpText:
			wg.Add(1)
			_ = goPoolB.Submit(
				func(idx int) func() {
					return func() {
						metrics.IncCounter(metrics.CntGatewayIncomingWebsocketMessage)
						g.delegate.OnMessage(wc, 0, ms[idx].Payload)
						pools.Bytes.Put(ms[idx].Payload)
						wg.Done()
					}
				}(idx),
			)
		case ws.OpClose:
			// remove the connection from the list
			err = ErrOpCloseReceived
		default:
			g.cfg.Logger.Warn("Unknown OpCode")
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
			if ce := g.cfg.Logger.Check(log.DebugLevel, "Error in websocketWritePump"); ce != nil {
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
