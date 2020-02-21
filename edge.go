package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/bridge"
	"git.ronaksoftware.com/ronak/rony/context"
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/logger"
	"git.ronaksoftware.com/ronak/rony/msg"
	"git.ronaksoftware.com/ronak/rony/pools"
	"go.uber.org/zap"
	"runtime/debug"
	"sync"
	"time"
)

/*
   Creation Time: 2020 - Feb - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Handler func(ctx *context.Context, in, out *msg.MessageEnvelope)
type GetConstructorNameFunc func(constructor int64) string

// EdgeConfig
type EdgeConfig struct {
	TestMode           bool
	BundleID           string
	InstanceID         string
	GetConstructorName GetConstructorNameFunc
}

// EdgeServer
type EdgeServer struct {
	handlers           map[int64][]Handler
	bundleID           string
	instanceID         string
	getConstructorName GetConstructorNameFunc
	bridge             bridge.Bridge
}

func NewEdgeServer(config EdgeConfig, bridge bridge.Bridge) (*EdgeServer, error) {
	if gatewayProtocol == gateway.Undefined {
		return nil, errors.ErrGatewayNotInitialized
	}
	if config.GetConstructorName == nil {
		config.GetConstructorName = func(constructor int64) string {
			return ""
		}
	}
	edgeServer := &EdgeServer{
		handlers:           make(map[int64][]Handler),
		bundleID:           config.BundleID,
		instanceID:         config.InstanceID,
		bridge:             bridge,
		getConstructorName: config.GetConstructorName,
	}

	return edgeServer, nil
}

// addHandler
func (edge EdgeServer) AddHandler(constructor int64, handler ...Handler) {
	edge.handlers[constructor] = handler
}

// Execute apply the right handler on the req, the response will be pushed to the clients queue.
func (edge EdgeServer) Execute(authID, userID int64, connID uint64, req *msg.MessageEnvelope) {
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
				edge.execute(ctx, x.Envelopes[idx], res)
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
		edge.execute(ctx, req, res)
		pools.ReleaseContext(ctx)
	}

	return
}

// ExecuteWithResult is similar to Execute but it get response filled passed by arguments
func (edge EdgeServer) ExecuteWithResult(authID, userID int64, connID uint64, req, res *msg.MessageEnvelope) error {
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
				edge.execute(ctx, x.Envelopes[idx], localRes)
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
		edge.execute(ctx, req, res)
		pools.ReleaseContext(ctx)
	}
}

func (edge EdgeServer) execute(ctx *context.Context, req, res *msg.MessageEnvelope) error {
	defer edge.recoverPanic(ctx)
	startTime := time.Now()

	if ce := log.Check(log.DebugLevel, "Execute"); ce != nil {
		ce.Write(
			zap.String("Constructor", edge.getConstructorName(req.Constructor)),
			zap.Uint64("RequestID", req.RequestID),
			zap.Int64("AuthID", ctx.AuthID),
		)
	}
	res.RequestID = req.RequestID
	handlers, ok := edge.handlers[req.Constructor]
	if !ok {
		return errors.ErrConstructorNotHandled
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
			zap.String("Constructor", edge.getConstructorName(req.Constructor)),
			zap.Int64("C", res.Constructor),
			zap.Uint64("RequestID", res.RequestID),
			zap.Int64("AuthID", ctx.AuthID),
			zap.Int64("UserID", ctx.UserID),
			zap.Int("MessageSize", len(res.Message)),
		)
	}
	duration := time.Now().Sub(startTime)
	if duration > longRequestThreshold {
		if ce := log.Check(log.InfoLevel, "Execute Too Long"); ce != nil {
			ce.Write(
				zap.Int64("AuthID", ctx.AuthID),
				zap.String("Constructor", edge.getConstructorName(req.Constructor)),
				zap.Uint64("RequestID", req.RequestID),
				zap.Duration("Duration", duration),
			)
		}
	}

	return nil

}

func (edge EdgeServer) recoverPanic(ctx *context.Context) {
	if r := recover(); r != nil {
		log.Error("Panic Recovered",
			zap.String("ServerID", edge.GetServerID()),
			zap.Int64("UserID", ctx.UserID),
			zap.Int64("AuthID", ctx.AuthID),
			zap.ByteString("Stack", debug.Stack()),
		)
	}
}

func (edge EdgeServer) GetServerID() string {
	return fmt.Sprintf("%s.%s", edge.bundleID, edge.instanceID)
}

// Run runs the selected gateway, if gateway is not setup it panics
func (edge EdgeServer) Run() {
	// Initialize API
	log.Info("API Server Started",
		zap.String("BundleID", edge.bundleID),
		zap.String("InstanceID", edge.instanceID),
		zap.String("Gateway", string(gatewayProtocol)),
	)

	switch gatewayProtocol {
	case gateway.Websocket:
		gatewayWebsocket.Run()
	case gateway.HTTP:
		gatewayHTTP.Run()
	case gateway.Quic:
		gatewayQuic.Run()
	case gateway.Grpc:
		gatewayGRPC.Run()
	default:
		panic("unknown gateway mode")
	}
}

func (edge EdgeServer) Shutdown() {
	switch gatewayProtocol {
	case gateway.Websocket:
		log.Info("Websocket Gateway Shutting down")
		gatewayWebsocket.Shutdown()
		log.Info("Websocket Gateway Shutdown")
	case gateway.HTTP:
		log.Info("Http Gateway Shutting down")
		// gatewayHTTP.Shutdown()
		log.Info("Http Gateway Shutdown")
	case gateway.Quic:
		log.Info("Quic Gateway Shutting down")
		// gatewayQuic.Shutdown()
		log.Info("Quic Gateway Shutdown")
	case gateway.Grpc:
		log.Info("Grpc Gateway Shutting down")
		// gatewayGRPC.Shutdown()
		log.Info("Grpc Gateway Shutdown")
	default:
		panic("unknown gateway mode")
	}
	gatewayProtocol = gateway.Undefined
}
