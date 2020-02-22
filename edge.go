package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/context"
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/msg"
	raftbadger "github.com/bbva/raft-badger"
	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"net"
	"os"
	"path/filepath"
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
	TestMode   bool
	BundleID   string
	InstanceID string
	RaftConfig RaftConfig
	// DataPath is a folder which all the internal state data will be written to.
	DataPath           string
	GetConstructorName GetConstructorNameFunc
}

// EdgeServer
type EdgeServer struct {
	// General
	bundleID   string
	instanceID string
	dataPath   string

	// Handlers
	preHandlers        []Handler
	handlers           map[int64][]Handler
	postHandlers       []Handler
	getConstructorName GetConstructorNameFunc

	// Raft Related
	raftEnabled   bool
	raftPort      int
	raftBootstrap bool
	raft          *raft.Raft
}

func NewEdgeServer(config EdgeConfig) (*EdgeServer, error) {
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
		getConstructorName: config.GetConstructorName,
		dataPath:           config.DataPath,
		raftEnabled:        config.RaftConfig.Enabled,
		raftPort:           config.RaftConfig.Port,
		raftBootstrap:      config.RaftConfig.FirstMachine,
	}

	return edgeServer, nil
}

// addHandler
func (edge EdgeServer) AddHandler(constructor int64, handler ...Handler) {
	edge.handlers[constructor] = handler
}

// Execute apply the right handler on the req, the response will be pushed to the clients queue.
func (edge EdgeServer) Execute(authID, userID int64, req *msg.MessageEnvelope) (err error) {
	if edge.raftEnabled {
		if edge.raft.State() != raft.Leader {

			return errors.ErrNotRaftLeader
		}

	}
	switch req.Constructor {
	case msg.C_MessageContainer:
		x := &msg.MessageContainer{}
		err = x.Unmarshal(req.Message)
		if err != nil {
			return err
		}
		xLen := len(x.Envelopes)

		waitGroup := pools.AcquireWaitGroup()
		for i := 0; i < xLen; i++ {
			ctx := pools.AcquireContext(authID, userID, true, false)
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
		ctx := pools.AcquireContext(authID, userID, false, false)
		edge.execute(ctx, req, res)
		pools.ReleaseContext(ctx)
	}

	return
}

// ExecuteWithResult is similar to Execute but it get response filled passed by arguments
func (edge EdgeServer) ExecuteWithResult(authID, userID int64, req, res *msg.MessageEnvelope) (err error) {
	switch req.Constructor {
	case msg.C_MessageContainer:
		x := &msg.MessageContainer{}
		err = x.Unmarshal(req.Message)
		if err != nil {
			return err
		}
		xLen := len(x.Envelopes)
		resContainer := &msg.MessageContainer{
			Envelopes: make([]*msg.MessageEnvelope, 0, xLen),
		}

		waitGroup := pools.AcquireWaitGroup()
		mtxLock := sync.Mutex{}
		for i := 0; i < xLen; i++ {
			ctx := pools.AcquireContext(authID, userID, true, true)
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
		ctx := pools.AcquireContext(authID, userID, false, true)
		edge.execute(ctx, req, res)
		pools.ReleaseContext(ctx)
	}
	return nil
}

func (edge EdgeServer) execute(ctx *context.Context, req, res *msg.MessageEnvelope) {
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
		// TODO:: fix this
		// ctx.PushError(res, errors.ErrCodeInvalid, errors.ErrItemApi)
		return
	}

	// Run the handler
	for idx := range edge.preHandlers {
		edge.preHandlers[idx](ctx, req, res)
		if ctx.Stop {
			break
		}
	}
	if !ctx.Stop {
		for idx := range handlers {
			handlers[idx](ctx, req, res)
			if ctx.Stop {
				break
			}
		}
	}
	if !ctx.Stop {
		for idx := range edge.postHandlers {
			edge.postHandlers[idx](ctx, req, res)
			if ctx.Stop {
				break
			}
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

	return
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
func (edge EdgeServer) Run() (err error) {
	log.Info("Edge Server Started",
		zap.String("BundleID", edge.bundleID),
		zap.String("InstanceID", edge.instanceID),
		zap.String("Gateway", string(gatewayProtocol)),
	)

	if edge.raftEnabled {
		err = edge.runRaft()
		if err != nil {
			return
		}
	}

	err = edge.runGateway()

	return
}
func (edge EdgeServer) runGateway() error {
	switch gatewayProtocol {
	case gateway.Websocket:
		gatewayWebsocket.Run()
	case gateway.HTTP:
		gatewayHTTP.Run()
	case gateway.QUIC:
		gatewayQuic.Run()
	default:
		panic("unknown gateway mode")
	}
	return nil
}
func (edge EdgeServer) runRaft() error {
	// Initialize LogStore for Raft
	badgerOpt := badger.DefaultOptions(filepath.Join(edge.dataPath, "raft"))
	badgerStore, err := raftbadger.New(raftbadger.Options{
		Path:                "",
		BadgerOptions:       &badgerOpt,
		NoSync:              false,
		ValueLogGC:          false,
		GCInterval:          0,
		MandatoryGCInterval: 0,
		GCThreshold:         0,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize Raft
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(edge.GetServerID())
	raftBind := fmt.Sprintf(":%d", edge.raftPort)
	raftAdvertiseAddr, err := net.ResolveIPAddr("tcp", raftBind)
	if err != nil {
		return err
	}

	raftTransport, err := raft.NewTCPTransport(raftBind, raftAdvertiseAddr, 3, 10*time.Second, os.Stdout)
	if err != nil {
		return err
	}

	raftSnapshot, err := raft.NewFileSnapshotStore(raftBind, 3, os.Stdout)
	if err != nil {
		return err
	}

	edge.raft, err = raft.NewRaft(raftConfig, edge, badgerStore, badgerStore, raftSnapshot, raftTransport)
	if err != nil {
		return err
	}

	if edge.raftBootstrap {
		bootConfig := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raftConfig.LocalID,
					Address: raftTransport.LocalAddr(),
				},
			},
		}
		edge.raft.BootstrapCluster(bootConfig)
	}

	return nil
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
	case gateway.QUIC:
		log.Info("Quic Gateway Shutting down")
		// gatewayQuic.Shutdown()
		log.Info("Quic Gateway Shutdown")
	default:
		panic("unknown gateway mode")
	}
	gatewayProtocol = gateway.Undefined
}
