package grpcGateway

import (
	"context"
	"git.ronaksoftware.com/ronak/rony/gateway"
	"io"
	"net"
	"sync"
	"sync/atomic"

	log "git.ronaksoftware.com/ronak/rony/logger"
	"git.ronaksoftware.com/ronak/rony/msg"
	"git.ronaksoftware.com/ronak/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var _ msg.BotProxyServer = (*Gateway)(nil)

type ConnectHandler func(gateway.Conn) error
type MessageHandler func(gateway.Conn, *msg.MessageEnvelope)
type FlushFunc func(gateway.Conn) *msg.MessageEnvelope

type Config struct {
	ConnectHandler ConnectHandler
	MessageHandler MessageHandler
	FlushFunc      FlushFunc
	ListenAddress  string
}

type Gateway struct {
	listenAddress  string
	listener       net.Listener
	server         *grpc.Server
	connectHandler ConnectHandler
	messageHandler MessageHandler
	flushFunc      FlushFunc

	conns    map[uint64]*Conn
	connsMtx sync.RWMutex

	connsTotal  int32
	connsLastID uint64
}

func New(config Config) *Gateway {
	g := &Gateway{
		connectHandler: config.ConnectHandler,
		messageHandler: config.MessageHandler,
		flushFunc:      config.FlushFunc,
		listenAddress:  config.ListenAddress,
		conns:          make(map[uint64]*Conn),
	}

	if g.connectHandler == nil {
		g.connectHandler = func(gateway.Conn) error { log.Debug("Using no-op connect handler."); return nil }
	}
	if g.messageHandler == nil {
		g.messageHandler = func(gateway.Conn, *msg.MessageEnvelope) { log.Debug("Using no-op message handler.") }
	}
	if g.flushFunc == nil {
		g.flushFunc = func(gateway.Conn) *msg.MessageEnvelope { log.Debug("Using no-op flush func."); return nil }
	}

	return g
}

func (g *Gateway) Run() {
	lis, err := net.Listen("tcp", g.listenAddress)
	if err != nil {
		log.Fatal(err.Error())
	}
	g.listener = lis

	s := grpc.NewServer()
	g.server = s

	msg.RegisterBotProxyServer(s, g)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal("GRPC listen error.", zap.Error(err))
		}
	}()
}

func (g *Gateway) ServeBotRequest(srv msg.BotProxy_ServeBotRequestServer) error {
	authKey, err := authKeyFromMetadata(srv.Context())
	if err != nil {
		return err
	}

	conn := g.addConnection(srv, authKey)
	defer g.removeConnection(conn.ConnID)

	if err := g.connectHandler(conn); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	log.Info("New ServeBotRequest stream.",
		zap.Uint64("connID", conn.ConnID),
	)

	for {
		msg, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				log.Info("End of stream.",
					zap.Uint64("connID", conn.ConnID),
				)
				return nil
			}
		}
		if msg == nil {
			log.Info("Received nil message.",
				zap.Uint64("connID", conn.ConnID),
			)
			return nil
		}

		g.messageHandler(conn, msg)
	}
}

func authKeyFromMetadata(ctx context.Context) ([]byte, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Info("Failed to get stream metadata.")
		return nil, status.Errorf(codes.DataLoss, "failed to get stream metadata")
	}

	authKeys, ok := md["auth-key"]
	if !ok {
		log.Info("No auth key was found in metadata.")
		return nil, status.Error(codes.Unauthenticated, "no auth key was found in metadata")
	}
	if authKeys[0] == "" {
		log.Info("Received empty auth key.")
		return nil, status.Error(codes.Unauthenticated, "received empty auth key")
	}

	return tools.StrToByte(authKeys[0]), nil
}

func (g *Gateway) addConnection(srv msg.BotProxy_ServeBotRequestServer, authKey []byte) *Conn {
	atomic.AddInt32(&g.connsTotal, 1)
	connID := atomic.AddUint64(&g.connsLastID, 1)

	conn := Conn{
		ConnID:    connID,
		srv:       srv,
		AuthKey:   authKey,
		flushFunc: g.flushFunc,
	}

	p, ok := peer.FromContext(srv.Context())
	if ok {
		conn.ClientIP = p.Addr.String()
	}

	g.connsMtx.Lock()
	g.conns[connID] = &conn
	g.connsMtx.Unlock()

	return &conn
}

func (g *Gateway) removeConnection(connID uint64) {
	atomic.AddInt32(&g.connsTotal, -1)

	g.connsMtx.Lock()
	delete(g.conns, connID)
	g.connsMtx.Unlock()
}

func (g *Gateway) GetConnection(connID uint64) *Conn {
	g.connsMtx.RLock()
	defer g.connsMtx.RUnlock()

	return g.conns[connID]
}
