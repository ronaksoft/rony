package api

import (
	"crypto/rsa"
	"fmt"
	natsBridge "git.ronaksoftware.com/ronak/rony/bridge/nats"
	"git.ronaksoftware.com/ronak/rony/config"
	"git.ronaksoftware.com/ronak/rony/context"
	"git.ronaksoftware.com/ronak/rony/db/redis"
	"git.ronaksoftware.com/ronak/rony/gateway"
	grpcGateway "git.ronaksoftware.com/ronak/rony/gateway/grpc"
	httpGateway "git.ronaksoftware.com/ronak/rony/gateway/http"
	quicGateway "git.ronaksoftware.com/ronak/rony/gateway/quic"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	log "git.ronaksoftware.com/ronak/rony/logger"
	"git.ronaksoftware.com/ronak/rony/metrics"
	"git.ronaksoftware.com/ronak/rony/msg"
	"git.ronaksoftware.com/ronak/rony/pools"
	"github.com/gobwas/pool/pbytes"
	"github.com/monnand/dhkx"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"runtime/debug"
	"sync"
	"time"
)

/*
   Creation Time: 2019 - Aug - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	TestMode bool

	// internals
	gatewayProtocol  gateway.Protocol
	gatewayWebsocket *websocketGateway.Gateway
	gatewayHTTP      *httpGateway.Gateway
	gatewayGRPC      *grpcGateway.Gateway
	gatewayQuic      *quicGateway.Gateway
	bridge           *natsBridge.Bridge
	handlersMap      map[int64][]Handler
)

type Handler func(ctx *context.Context, in, out *msg.MessageEnvelope)

func InitWebsocketGateway(config websocketGateway.Config) {
	if gatewayProtocol != "" {
		panic("bug: only one gateway could be setup")
	}
	initialize()
	var err error
	gatewayWebsocket, err = websocketGateway.New(config)
	log.PanicOnError("Error On Initializing Websocket Gateway", err)
	go func() {
		for {
			gatewayStats := gatewayWebsocket.Stats()
			metrics.Gauge(metrics.GaugeActiveConnections).Set(float64(gatewayWebsocket.TotalConnections()))
			metrics.Gauge(metrics.GaugeGatewayWebsocketInQueue).Set(float64(gatewayStats.InQueue))
			metrics.Gauge(metrics.GaugeGatewayWebsocketOutQueue).Set(float64(gatewayStats.OutQueue))
			w, d, wb, db := CountUniqueVisits()
			metrics.Gauge(metrics.GaugeActiveUsers24h).Set(float64(d))
			metrics.Gauge(metrics.GaugeActiveUsers1w).Set(float64(w))
			metrics.Gauge(metrics.GaugeActiveUsersLast24h).Set(float64(db))
			metrics.Gauge(metrics.GaugeActiveUsersLast1w).Set(float64(wb))
			time.Sleep(time.Second * 5)
		}
	}()
	gatewayProtocol = gateway.Websocket
}

func InitHttpGateway(config httpGateway.Config) {
	if gatewayProtocol != "" {
		panic("bug: only one gateway could be setup")
	}
	initialize()
	gatewayHTTP = httpGateway.New(config)
	gatewayProtocol = gateway.HTTP
}

func InitGrpcGateway(config grpcGateway.Config) {
	if gatewayProtocol != "" {
		panic("bug: only one gateway could be setup")
	}
	initialize()
	gatewayGRPC = grpcGateway.New(config)
	gatewayProtocol = gateway.GRPC
}

func InitQuicGateway(config quicGateway.Config) {
	if gatewayProtocol != "" {
		panic("bug: only one gateway could be setup")
	}
	initialize()
	var err error
	gatewayQuic, err = quicGateway.New(config)
	log.PanicOnError("Error On Initializing Quic Gateway", err)
	go func() {
		// for {
		// 	gatewayStats := quicGateway.Stats()
		// 	metrics.Gauge(metrics.GaugeActiveConnections).Set(float64(gatewayWebsocket.TotalConnections()))
		// 	metrics.Gauge(metrics.GaugeGatewayWebsocketInQueue).Set(float64(gatewayStats.InQueue))
		// 	metrics.Gauge(metrics.GaugeGatewayWebsocketOutQueue).Set(float64(gatewayStats.OutQueue))
		// 	time.Sleep(time.Second * 5)
		// }
	}()
	gatewayProtocol = gateway.QUIC
}

func initialize() {
	handlersMap = make(map[int64][]Handler)
	TestMode = config.GetBool(config.TestMode)

	initBridge()

	return
}
func initBridge() {
	// Initialize bridge
	natsConfig := nats.GetDefaultOptions()
	natsConfig.Name = fmt.Sprintf("%s.%s", config.BundleID, config.InstanceID)
	natsConfig.Url = config.GetString(config.NatsURL)
	var err error
	bridge, err = natsBridge.NewBridge(natsBridge.Config{
		BundleID:   config.GetString(config.BundleID),
		InstanceID: config.GetString(config.InstanceID),
		Options:    natsConfig,
		Timeout:    config.GetDuration(config.NatsTimeout),
		Retries:    config.GetInt(config.NatsRetries),
		ErrHandler: func(err error) {
			log.Warn("Error On bridge", zap.Error(err))
		},
		MessageHandler: nil,
		NotifyHandler:  bridgeNotifyHandler,
	})
	log.PanicOnError("Error On bridge Setup", err)
}
func bridgeNotifyHandler(connIDs []uint64) {
	metrics.Counter(metrics.CntBridgeNotify).Add(1)
	if ce := log.Check(log.DebugLevel, "Notify Arrived"); ce != nil {
		ce.Write(
			zap.Uint64s("ConnIDs", connIDs),
		)
	}
	switch gatewayProtocol {
	case gateway.Websocket:
		for idx := range connIDs {
			wsConn := gatewayWebsocket.GetConnection(connIDs[idx])
			if wsConn == nil {
				if ce := log.Check(log.DebugLevel, "Flush Signal For Nil Websocket"); ce != nil {
					ce.Write(zap.Uint64("ConnID", connIDs[idx]))
				}
				continue
			}
			wsConn.Flush()
		}
	case gateway.GRPC:
		for idx := range connIDs {
			conn := gatewayGRPC.GetConnection(connIDs[idx])
			if conn == nil {
				log.Warn("Flush Signal For Nil GRPC stream", zap.Uint64("ConnID", connIDs[idx]))
				continue
			}
			conn.Flush()
		}
	default:

	}
}

// addHandler
func AddHandler(constructor int64, handler ...Handler) {
	handlersMap[constructor] = handler
}

// Execute apply the right handler on the req, the response will be pushed to the clients queue.
func Execute(authID, userID int64, connID uint64, req *msg.MessageEnvelope) {
	switch req.Constructor {
	case msg.C_MessageContainer:
		x := &msg.MessageContainer{}
		_ = x.Unmarshal(req.Message)
		xLen := len(x.Envelopes)

		waitGroup := pools.AcquireWaitGroup()
		for i := 0; i < xLen; i++ {
			ctx := pools.AcquireContext(authID, userID, connID, true, false)
			nextChan := make(chan struct{}, 1)
			waitGroup.Add(1)
			go func(ctx *context.Context, idx int) {
				res := &msg.MessageEnvelope{}
				execute(ctx, x.Envelopes[idx], res)
				nextChan <- struct{}{}
				waitGroup.Done()
				pools.ReleaseContext(ctx)
			}(ctx, i)
			select {
			case <-ctx.NextChan:
				// The handler supported quick return
			case <-nextChan:
			}
		}
		waitGroup.Wait()
		pools.ReleaseWaitGroup(waitGroup)
	default:
		res := &msg.MessageEnvelope{}
		ctx := pools.AcquireContext(authID, userID, connID, false, false)
		execute(ctx, req, res)
		pools.ReleaseContext(ctx)
	}

	return
}

// ExecuteWithResult is similar to Execute but it get response filled passed by arguments
func ExecuteWithResult(authID, userID int64, connID uint64, req, res *msg.MessageEnvelope) {
	ctx := pools.AcquireContext(authID, userID, connID, false, true)
	defer pools.ReleaseContext(ctx)

	switch req.Constructor {
	case msg.C_MessageContainer:
		x := &msg.MessageContainer{}
		_ = x.Unmarshal(req.Message)
		xLen := len(x.Envelopes)
		resContainer := &msg.MessageContainer{}
		resContainer.Envelopes = make([]*msg.MessageEnvelope, 0, xLen)

		waitGroup := pools.AcquireWaitGroup()
		mtxLock := sync.Mutex{}
		for i := 0; i < xLen; i++ {
			ctx := pools.AcquireContext(authID, userID, connID, true, true)
			nextChan := make(chan struct{}, 1)
			waitGroup.Add(1)
			go func(ctx *context.Context, idx int) {
				localRes := &msg.MessageEnvelope{}
				execute(ctx, x.Envelopes[idx], localRes)
				nextChan <- struct{}{}
				mtxLock.Lock()
				resContainer.Envelopes = append(resContainer.Envelopes, localRes)
				mtxLock.Unlock()
				waitGroup.Done()
				pools.ReleaseContext(ctx)
			}(ctx, i)
			select {
			case <-ctx.NextChan:
				// The handler supported quick return
			case <-nextChan:
			}
		}
		resContainer.Length = int32(len(resContainer.Envelopes))
		res.RequestID = req.RequestID
		msg.ResultMessageContainer(res, resContainer)
		waitGroup.Wait()
		pools.ReleaseWaitGroup(waitGroup)
	default:
		ctx := pools.AcquireContext(authID, userID, connID, false, true)
		execute(ctx, req, res)
	}
}

func execute(ctx *context.Context, req, res *msg.MessageEnvelope) {
	defer recoverPanic(ctx)
	startTime := time.Now()

	if ce := log.Check(log.DebugLevel, "Execute"); ce != nil {
		ce.Write(
			zap.String("Constructor", msg.ConstructorNames[req.Constructor]),
			zap.Uint64("RequestID", req.RequestID),
			zap.Int64("AuthID", ctx.AuthID),
		)
	}
	res.RequestID = req.RequestID
	handlers, ok := handlersMap[req.Constructor]
	if !ok {
		PushError(ctx, res, msg.ErrCodeInvalid, msg.ErrItemApi)
		return
	}

	// Run the handler
	for idx := range handlers {
		handlers[idx](ctx, req, res)
		if ctx.Stop {
			break
		}
	}

	if ce := log.Check(log.DebugLevel, "Response"); ce != nil {
		ce.Write(
			zap.String("Constructor", msg.ConstructorNames[res.Constructor]),
			zap.Int64("C", res.Constructor),
			zap.Uint64("RequestID", res.RequestID),
			zap.Int64("AuthID", ctx.AuthID),
			zap.Int64("UserID", ctx.UserID),
			zap.Int("MessageSize", len(res.Message)),
		)
	}
	duration := time.Now().Sub(startTime)
	if duration > longRequestThreshold {
		metrics.CounterVec(metrics.CntSlowFunctions).WithLabelValues(msg.ConstructorNames[req.Constructor]).Add(1)
		if ce := log.Check(log.InfoLevel, "Execute Too Long"); ce != nil {
			ce.Write(
				zap.Int64("AuthID", ctx.AuthID),
				zap.String("Constructor", msg.ConstructorNames[req.Constructor]),
				zap.Uint64("RequestID", req.RequestID),
				zap.Duration("Duration", duration),
			)
		}
	}
	metrics.Histogram(metrics.HistRequestTimeMS).Observe(float64(duration / time.Millisecond))
	return

}

func recoverPanic(ctx *context.Context) {
	if r := recover(); r != nil {
		log.Error("Panic Recovered",
			zap.String("ServerID", fmt.Sprintf("%s.%s", config.GetString(config.BundleID), config.GetString(config.InstanceID))),
			zap.Int64("UserID", ctx.UserID),
			zap.Int64("AuthID", ctx.AuthID),
			zap.ByteString("Stack", debug.Stack()),
		)
	}
}

// Run runs the selected gateway, if gateway is not setup it panics
func Run() {
	// Initialize API
	log.Info("API Server Started",
		zap.String("bundleID", config.GetString(config.BundleID)),
		zap.String("instanceID", config.GetString(config.InstanceID)),
		zap.String("Gateway", string(gatewayProtocol)),
	)

	switch gatewayProtocol {
	case gateway.Websocket:
		gatewayWebsocket.Run()
		log.Info("Websocket Gateway Started", zap.String("Addr", config.GetString(config.GatewayListenAddr)))
	case gateway.HTTP:
		gatewayHTTP.Run()
		log.Info("Http Gateway Started", zap.String("Addr", config.GetString(config.GatewayListenAddr)))
	case gateway.QUIC:
		gatewayQuic.Run()
		log.Info("Quic Gateway Started", zap.String("Addr", config.GetString(config.GatewayListenAddr)))
	case gateway.GRPC:
		gatewayGRPC.Run()
		log.Info("GRPC Gateway Started", zap.String("Addr", config.GetString(config.GatewayListenAddr)))
	default:
		panic("unknown gateway mode")
	}
}

func Shutdown() {
	switch gatewayProtocol {
	case gateway.Websocket:
		log.Info("Websocket Gateway Shutting down")
		gatewayWebsocket.Shutdown()
		log.Info("Websocket Gateway Shutdown")
	case gateway.HTTP:
		log.Info("Http Gateway Shutting down")
		// gatewayHTTP.Shutdown()
		log.Info("Http Gateway Shutdown")
	case gateway.QUIC:
		log.Info("Quic Gateway Shutting down")
		// gatewayQuic.Shutdown()
		log.Info("Quic Gateway Shutdown")
	case gateway.GRPC:
		log.Info("GRPC Gateway Shutting down")
		// gatewayGRPC.Shutdown()
		log.Info("GRPC Gateway Shutdown")
	default:
		panic("unknown gateway mode")
	}
	gatewayProtocol = ""
}
