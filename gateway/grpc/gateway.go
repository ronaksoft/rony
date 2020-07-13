package grpcGateway

import (
	"context"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

var _ rony.RonyServer = (*Gateway)(nil)

type Config struct {
	ConnectHandler gateway.ConnectHandler
	MessageHandler gateway.MessageHandler
	CloseHandler   gateway.CloseHandler
	ListenAddress  string
}

type Gateway struct {
	gateway.ConnectHandler
	gateway.MessageHandler
	gateway.CloseHandler

	listenAddress string
	listener      net.Listener
	server        *grpc.Server

	conns    map[uint64]*Conn
	connsMtx sync.RWMutex

	connsTotal  int32
	connsLastID uint64
}

func New(config Config) *Gateway {
	g := &Gateway{
		listenAddress: config.ListenAddress,
		conns:         make(map[uint64]*Conn),
	}

	if g.ConnectHandler == nil {
		g.ConnectHandler = func(connID uint64) {}

	}
	if g.MessageHandler == nil {
		g.MessageHandler = func(c gateway.Conn, streamID int64, data []byte, kvs ...gateway.KeyValue) {
			log.Debug("Using no-op message handler.")
		}
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
	rony.RegisterRonyServer(s, g)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal("GRPC listen error.", zap.Error(err))
		}
	}()
}

func (g *Gateway) Pipeline(srv rony.Rony_PipelineServer) error {
	authKey, err := authKeyFromMetadata(srv.Context())
	if err != nil {
		return err
	}

	conn := g.addConnection(srv, authKey)
	defer g.removeConnection(conn.ConnID)

	g.ConnectHandler(conn.ConnID)

	log.Info("New Pipeline stream.",
		zap.Uint64("connID", conn.ConnID),
	)

	for {
		envelope, err := srv.Recv()
		switch err {
		case nil:
		case io.EOF:
			log.Info("End of stream.",
				zap.Uint64("connID", conn.ConnID),
			)
			return nil
		default:
			log.Warn("Error On Received Envelope", zap.Error(err))
			continue
		}

		if envelope == nil {
			log.Info("Received nil message.",
				zap.Uint64("connID", conn.ConnID),
			)
			return nil
		}

		envelopeBytes := pools.Bytes.GetLen(envelope.Size())
		_, _ = envelope.MarshalToSizedBuffer(envelopeBytes)
		g.MessageHandler(conn, 0, envelopeBytes)
		pools.Bytes.Put(envelopeBytes)
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

func (g *Gateway) addConnection(srv rony.Rony_PipelineServer, authKey []byte) *Conn {
	atomic.AddInt32(&g.connsTotal, 1)
	connID := atomic.AddUint64(&g.connsLastID, 1)

	conn := Conn{
		ConnID: connID,
		srv:    srv,
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
