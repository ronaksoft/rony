package quicGateway

import (
	"context"
	"crypto/tls"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/logger"
	"github.com/lucas-clemente/quic-go"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"sync"
	"sync/atomic"
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

// Config
type Config struct {
	gateway.CloseHandler
	// MessageHandler must not keep the caller busy, and preferably it must do the heavy work,
	// on background or other routines. Because our gateway use a constant number of readPump
	// routines.
	gateway.MessageHandler
	gateway.ConnectHandler
	gateway.FailedWriteHandler
	gateway.SuccessWriteHandler
	gateway.FlushFunc

	NewConnectionWorkers int
	MaxConcurrency       int
	ListenAddress        string
	Certificates         []tls.Certificate
}

// Gateway
type Gateway struct {
	gateway.ConnectHandler
	gateway.MessageHandler
	gateway.CloseHandler
	gateway.FailedWriteHandler
	gateway.SuccessWriteHandler
	gateway.FlushFunc

	// Internal Controlling Params
	newConnWorkers int
	maxConcurrency int
	listenOn       string
	conns          map[uint64]*Conn
	connsMtx       sync.RWMutex
	connsLastID    uint64
	connsTotal     int32
	certs          []tls.Certificate
}

func New(config Config) (*Gateway, error) {
	g := new(Gateway)
	g.listenOn = config.ListenAddress
	g.newConnWorkers = config.NewConnectionWorkers
	g.maxConcurrency = config.MaxConcurrency
	g.conns = make(map[uint64]*Conn, 100000)
	g.certs = config.Certificates

	if config.MessageHandler == nil {
		g.MessageHandler = func(c gateway.Conn, streamID int64, date []byte) {}
	} else {
		g.MessageHandler = config.MessageHandler
	}
	if config.CloseHandler == nil {
		g.CloseHandler = func(c gateway.Conn) {}
	} else {
		g.CloseHandler = config.CloseHandler
	}
	if config.ConnectHandler == nil {
		g.ConnectHandler = func(connID uint64) {}
	} else {
		g.ConnectHandler = config.ConnectHandler
	}
	if config.FailedWriteHandler == nil {
		g.FailedWriteHandler = func(c gateway.Conn, data []byte, err error) {}
	} else {
		g.FailedWriteHandler = config.FailedWriteHandler
	}
	if config.SuccessWriteHandler == nil {
		g.SuccessWriteHandler = func(c gateway.Conn) {}
	} else {
		g.SuccessWriteHandler = config.SuccessWriteHandler
	}
	if config.FlushFunc == nil {
		g.FlushFunc = func(c gateway.Conn) [][]byte { return nil }
	} else {
		g.FlushFunc = config.FlushFunc
	}

	return g, nil
}

func (g *Gateway) connectionAcceptor(l quic.Listener) {
	for {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
		conn, err := l.Accept(ctx)
		if err != nil {
			log.Warn("Error on Listener Accept",
				zap.Error(err),
			)
			cancelFunc()
			continue
		}

		go g.connHandler(g.addConnection(conn))
	}
}

// addConnection
func (g *Gateway) addConnection(conn quic.Session) *Conn {
	// Increment total connection counter and connection ID
	totalConns := atomic.AddInt32(&g.connsTotal, 1)
	connID := atomic.AddUint64(&g.connsLastID, 1)
	qConn := Conn{
		ClientIP:  conn.RemoteAddr().String(),
		ConnID:    connID,
		conn:      conn,
		gateway:   g,
		streams:   make(map[quic.StreamID]quic.Stream),
		flushChan: make(chan bool, 1),
	}
	g.connsMtx.Lock()
	g.conns[connID] = &qConn
	g.connsMtx.Unlock()
	if ce := log.Check(log.DebugLevel, "Quic Connection connected"); ce != nil {
		ce.Write(
			zap.Uint64("ConnectionID", connID),
			zap.Int32("Total", totalConns),
		)
	}
	g.ConnectHandler(connID)
	return &qConn
}

// removeConnection
func (g *Gateway) removeConnection(connID uint64) {
	g.connsMtx.Lock()
	qConn, ok := g.conns[connID]
	if !ok {
		g.connsMtx.Unlock()
		return
	}
	delete(g.conns, connID)
	g.connsMtx.Unlock()

	totalConns := atomic.AddInt32(&g.connsTotal, -1)
	_ = qConn.conn.Close()

	qConn.Lock()
	if !qConn.closed {
		g.CloseHandler(qConn)

		qConn.closed = true
		qConn.flushChan = nil
	}
	qConn.Unlock()
	if ce := log.Check(log.DebugLevel, "Quic Connection Removed"); ce != nil {
		ce.Write(
			zap.Uint64("ConnectionID", connID),
			zap.Int32("Total", totalConns),
		)
	}

}

func (g *Gateway) connHandler(qConn *Conn) {
	if ce := log.Check(log.DebugLevel, "New Conn"); ce != nil {
		ce.Write(
			zap.Uint64("ConnID", qConn.ConnID),
			zap.String("ClientIP", qConn.ClientIP),
		)
	}
	waitGroup := waitGroupPool.Get().(*sync.WaitGroup)
	for {
		stream, err := qConn.conn.AcceptStream(context.Background())
		if err != nil {
			if nErr, ok := err.(net.Error); ok {
				if !nErr.Temporary() {
					g.removeConnection(qConn.ConnID)
					break
				}
			}
			continue
		}
		waitGroup.Add(1)
		atomic.AddInt32(&qConn.noStreams, 1)
		qConn.Lock()
		qConn.streams[stream.StreamID()] = stream
		qConn.Unlock()
		go g.streamHandler(waitGroup, qConn, stream)
	}
	waitGroup.Wait()
	waitGroupPool.Put(waitGroup)
}

func (g *Gateway) streamHandler(waitGroup *sync.WaitGroup, qConn *Conn, stream quic.Stream) {
	defer waitGroup.Done()
	defer atomic.AddInt32(&qConn.noStreams, -1)
	if ce := log.Check(log.DebugLevel, "New Stream"); ce != nil {
		ce.Write(
			zap.Any("StreamID", stream.StreamID()),
			zap.Uint64("ConnID", qConn.ConnID),
			zap.Int32("OpenStreams", qConn.noStreams),
		)
	}

	_ = stream.SetReadDeadline(time.Now().Add(readWait))
	bytes, _ := ioutil.ReadAll(stream)

	atomic.AddUint64(&qConn.Counters.IncomingCount, 1)
	atomic.AddUint64(&qConn.Counters.IncomingBytes, uint64(len(bytes)))

	g.MessageHandler(qConn, int64(stream.StreamID()), bytes)
}

// Run
func (g *Gateway) Run() {
	listener, err := quic.ListenAddr(
		g.listenOn,
		&tls.Config{
			Certificates:       g.certs,
			InsecureSkipVerify: true,
			NextProtos:         []string{"RonakProto"},
		},
		&quic.Config{
			ConnectionIDLength: 8,
			MaxIncomingStreams: maxStreamsPerConn,
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Run multiple connection acceptor
	for i := 0; i < g.newConnWorkers; i++ {
		go g.connectionAcceptor(listener)
	}

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
