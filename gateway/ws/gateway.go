package websocketGateway

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gateway"
	"git.ronaksoftware.com/ronak/rony/gateway/ws/util"
	"git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"github.com/valyala/tcplisten"
	"go.uber.org/zap"
	"net"
	"strings"
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
	NewConnectionWorkers int
	MaxConcurrency       int
	MaxIdleTime          time.Duration
	ListenAddress        string
}

// Gateway
type Gateway struct {
	gateway.ConnectHandler
	gateway.MessageHandler
	gateway.CloseHandler

	// Internal Controlling Params
	newConnWorkers     int
	maxConcurrency     int
	maxIdleTime        int64
	listener           net.Listener
	listenOn           string
	conns              map[uint64]*Conn
	connsMtx           sync.RWMutex
	connsTotal         int32
	connsLastID        uint64
	connGC             *connGC
	poller             netpoll.Poller
	stop               int32
	waitGroupAcceptors *sync.WaitGroup
	waitGroupReaders   *sync.WaitGroup
	waitGroupWriters   *sync.WaitGroup
	cntReads           uint64
	cntWrites          uint64
	addrs              []string
}

func New(config Config) (*Gateway, error) {
	g := new(Gateway)
	g.listenOn = config.ListenAddress
	g.newConnWorkers = config.NewConnectionWorkers
	g.maxConcurrency = config.MaxConcurrency
	g.conns = make(map[uint64]*Conn, 100000)
	g.waitGroupWriters = &sync.WaitGroup{}
	g.waitGroupReaders = &sync.WaitGroup{}
	g.waitGroupAcceptors = &sync.WaitGroup{}
	if config.MaxIdleTime == 0 {
		g.maxIdleTime = int64(defaultConnIdleTime.Seconds())
	} else {
		g.maxIdleTime = int64(config.MaxIdleTime)
	}
	g.connGC = newGC(g)
	g.MessageHandler = func(c gateway.Conn, streamID int64, date []byte) {}
	g.CloseHandler = func(c gateway.Conn) {}
	g.ConnectHandler = func(connID uint64) {}
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

	tcpConfig := tcplisten.Config{
		ReusePort:   true,
		FastOpen:    true,
		DeferAccept: false,
		Backlog:     8192,
	}

	// Setup Listener to listen for TCP connections
	listener, err := tcpConfig.NewListener("tcp4", g.listenOn)
	if err != nil {
		return nil, err
	}

	g.listener = listener

	ta, err := net.ResolveTCPAddr("tcp4", listener.Addr().String())
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

func (g *Gateway) connectionAcceptor() {
	defer g.waitGroupAcceptors.Done()
	var clientIP string
	var clientType string
	var detected bool
	wsUpgrader := ws.DefaultUpgrader
	wsUpgrader.OnHeader = func(key, value []byte) (err error) {
		switch strings.ToLower(tools.ByteToStr(key)) {
		case "cf-connecting-ip":
			clientIP = string(value)
			detected = true
		case "x-forwarded-for", "x-real-ip", "forwarded":
			if !detected {
				clientIP = string(value)
				detected = true
			}
		case "x-client-type":
			clientType = string(value)
		}
		return nil
	}
	for {
		if atomic.LoadInt32(&g.stop) == 1 {
			break
		}
		conn, err := g.listener.Accept()
		if err != nil {
			// log.Warn("Error On Accept", zap.Error(err))
			continue
		}

		if _, err := wsUpgrader.Upgrade(conn); err != nil {
			if ce := log.Check(log.InfoLevel, "Error in Connection Acceptor"); ce != nil {
				ce.Write(
					zap.String("IP", clientIP),
					zap.String("ClientType", clientType),
					zap.Error(err),
				)
			}
			_ = conn.Close()
			continue
		}

		// Set the internal parameters and run two separate go-routines per connection
		// one for reading and one for writing
		var wsConn *Conn
		if detected {
			wsConn = g.addConnection(conn, clientIP, clientType)
		} else {
			wsConn = g.addConnection(conn, strings.Split(conn.RemoteAddr().String(), ":")[0], clientType)
		}

		wsConn.desc, err = netpoll.Handle(conn,
			netpoll.EventRead|netpoll.EventHup|netpoll.EventOneShot,
		)

		if err != nil {
			log.Warn("Error On NetPoll Description", zap.Error(err))
			g.removeConnection(wsConn.ConnID)
			continue
		}
		err = g.poller.Start(wsConn.desc, wsConn.startEvent)
		if err != nil {
			log.Warn("Error On NetPoll Start", zap.Error(err))
			g.removeConnection(wsConn.ConnID)
			continue
		}

		// reset detected flag for next upgrade process
		detected = false
	}
}

func (g *Gateway) addConnection(conn net.Conn, clientIP, clientType string) *Conn {
	// Increment total connection counter and connection ID
	totalConns := atomic.AddInt32(&g.connsTotal, 1)
	connID := atomic.AddUint64(&g.connsLastID, 1)
	wsConn := Conn{
		ClientIP:     clientIP,
		ConnID:       connID,
		conn:         conn,
		gateway:      g,
		lastActivity: time.Now().Unix(),
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
	g.ConnectHandler(connID)
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

func (g *Gateway) readPump(wc *Conn) {
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
	case ws.OpPing:
		_ = wr.wc.conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
		err := wsutil.WriteMessage(wr.wc.conn, ws.StateServerSide, ws.OpPong, nil)
		if err != nil {
			g.removeConnection(wr.wc.ConnID)
		}
	}

}

// Run
func (g *Gateway) Run() {
	// Run multiple connection acceptor
	for i := 0; i < g.newConnWorkers; i++ {
		g.waitGroupAcceptors.Add(1)
		go g.connectionAcceptor()
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
	return g.addrs
}

// GetConnection
func (g *Gateway) GetConnection(connID uint64) *Conn {
	g.connsMtx.RLock()
	wsConn, ok := g.conns[connID]
	g.connsMtx.RUnlock()
	if ok {
		return wsConn
	}
	return nil
}

// TotalConnections
func (g *Gateway) TotalConnections() int32 {
	return atomic.LoadInt32(&g.connsTotal)
}
