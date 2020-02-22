package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/context"
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"git.ronaksoftware.com/ronak/rony/msg"
	raftbadger "github.com/bbva/raft-badger"
	"github.com/dgraph-io/badger/v2"
	"github.com/gobwas/pool/pbytes"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
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

type Handler func(ctx *context.Context, in *msg.MessageEnvelope)
type GetConstructorNameFunc func(constructor int64) string

// EdgeServer
type EdgeServer struct {
	// General
	bundleID          string
	instanceID        string
	dataPath          string
	gatewayProtocol   gateway.Protocol
	gateway           gateway.Gateway
	updateDispatcher  func(authID int64, envelope *msg.UpdateEnvelope)
	messageDispatcher func(authID int64, envelope *msg.MessageEnvelope)

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

func NewEdgeServer(bundleID, instanceID string, opts ...Option) *EdgeServer {
	edgeServer := &EdgeServer{
		handlers:   make(map[int64][]Handler),
		bundleID:   bundleID,
		instanceID: instanceID,
		getConstructorName: func(constructor int64) string {
			return fmt.Sprintf("%d", constructor)
		},
		updateDispatcher: func(authID int64, envelope *msg.UpdateEnvelope) {
			log.Debug("Update Dispatched",
				zap.Int64("AuthID", authID),
				zap.Int64("UpdateID", envelope.UpdateID),
				zap.Int64("C", envelope.Constructor),
			)
		},
		messageDispatcher: func(authID int64, envelope *msg.MessageEnvelope) {
			log.Debug("Message Dispatched",
				zap.Int64("AuthID", authID),
				zap.Uint64("UpdateID", envelope.RequestID),
				zap.Int64("C", envelope.Constructor),
			)
		},
		dataPath:      ".",
		raftEnabled:   false,
		raftPort:      0,
		raftBootstrap: false,
	}

	for _, opt := range opts {
		opt(edgeServer)
	}

	return edgeServer
}

func (edge *EdgeServer) GetServerID() string {
	return fmt.Sprintf("%s.%s", edge.bundleID, edge.instanceID)
}

func (edge *EdgeServer) AddHandler(constructor int64, handler ...Handler) {
	edge.handlers[constructor] = handler
}

// Execute apply the right handler on the req, the response will be pushed to the clients queue.
func (edge *EdgeServer) Execute(authID, userID int64, req *msg.MessageEnvelope) (err error) {
	if edge.raftEnabled {
		if edge.raft.State() != raft.Leader {
			return errors.ErrNotRaftLeader
		}
		raftCmd := pools.AcquireRaftCommand()
		raftCmd.AuthID = authID
		raftCmd.UserID = userID
		raftCmd.Envelope = req
		raftCmdBytes := pbytes.GetLen(raftCmd.Size())
		f := edge.raft.Apply(raftCmdBytes, raftApplyTimeout)
		err = f.Error()
	} else {
		edge.execute(authID, userID, req)
	}
	return
}

func (edge *EdgeServer) execute(authID, userID int64, req *msg.MessageEnvelope) {
	executeFunc := func(ctx *context.Context, req *msg.MessageEnvelope) {
		defer edge.recoverPanic(ctx)
		startTime := time.Now()

		waitGroup := pools.AcquireWaitGroup()
		waitGroup.Add(2)
		go func() {
			for u := range ctx.UpdateChan {
				edge.updateDispatcher(u.AuthID, u.Envelope)
				pools.ReleaseUpdateEnvelope(u.Envelope)
			}
			waitGroup.Done()
		}()
		go func() {
			for m := range ctx.MessageChan {
				edge.messageDispatcher(m.AuthID, m.Envelope)
				pools.ReleaseMessageEnvelope(m.Envelope)
			}
			waitGroup.Done()
		}()

		if ce := log.Check(log.DebugLevel, "Execute (Start)"); ce != nil {
			ce.Write(
				zap.String("Constructor", edge.getConstructorName(req.Constructor)),
				zap.Uint64("RequestID", req.RequestID),
				zap.Int64("AuthID", ctx.AuthID),
			)
		}
		handlers, ok := edge.handlers[req.Constructor]
		if !ok {
			// TODO:: fix this
			// ctx.PushError(res, errors.ErrCodeInvalid, errors.ErrItemApi)
			return
		}

		// Run the handler
		for idx := range edge.preHandlers {
			edge.preHandlers[idx](ctx, req)
			if ctx.Stop {
				break
			}
		}
		if !ctx.Stop {
			for idx := range handlers {
				handlers[idx](ctx, req)
				if ctx.Stop {
					break
				}
			}
		}
		if !ctx.Stop {
			for idx := range edge.postHandlers {
				edge.postHandlers[idx](ctx, req)
				if ctx.Stop {
					break
				}
			}
		}
		if !ctx.Stop {
			ctx.StopExecution()
		}
		waitGroup.Wait()
		duration := time.Now().Sub(startTime)

		if ce := log.Check(log.InfoLevel, "Execute (Finished)"); ce != nil {
			ce.Write(
				zap.Uint64("RequestID", req.RequestID),
				zap.Int64("AuthID", ctx.AuthID),
				zap.Duration("T", duration),
			)
		}

		return
	}
	switch req.Constructor {
	case msg.C_MessageContainer:
		x := &msg.MessageContainer{}
		err := x.Unmarshal(req.Message)
		if err != nil {
			// TODO:: handle error properly
			return
		}
		xLen := len(x.Envelopes)
		waitGroup := pools.AcquireWaitGroup()
		for i := 0; i < xLen; i++ {
			ctx := context.Acquire(authID, userID, true, false)
			nextChan := make(chan struct{}, 1)
			waitGroup.Add(1)
			go func(ctx *context.Context, idx int) {
				executeFunc(ctx, x.Envelopes[idx])
				nextChan <- struct{}{}
				waitGroup.Done()
				context.Release(ctx)
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
		ctx := context.Acquire(authID, userID, false, false)
		executeFunc(ctx, req)
		context.Release(ctx)
	}
	return
}

func (edge *EdgeServer) recoverPanic(ctx *context.Context) {
	if r := recover(); r != nil {
		log.Error("Panic Recovered",
			zap.String("ServerID", edge.GetServerID()),
			zap.Int64("UserID", ctx.UserID),
			zap.Int64("AuthID", ctx.AuthID),
			zap.ByteString("Stack", debug.Stack()),
		)
	}
}

func (edge *EdgeServer) onMessage(c gateway.Conn, streamID int64, data []byte) {
	authID := tools.RandomInt64(0)
	userID := tools.RandomInt64(0)
	err := edge.Execute(authID, userID, &msg.MessageEnvelope{
		Constructor: 200,
		RequestID:   tools.RandomUint64(),
		Message:     data,
	})
	log.WarnOnError("Execute", err)
}

func (edge *EdgeServer) onConnect(connID uint64) {}

func (edge *EdgeServer) onClose(conn gateway.Conn) {}

func (edge *EdgeServer) onFlush(conn gateway.Conn) [][]byte {
	return nil
}

// Run runs the selected gateway, if gateway is not setup it panics
func (edge *EdgeServer) Run() (err error) {
	if edge.gatewayProtocol == gateway.Undefined {
		return errors.ErrGatewayNotInitialized
	}

	log.Info("Edge Server Started",
		zap.String("BundleID", edge.bundleID),
		zap.String("InstanceID", edge.instanceID),
		zap.String("Gateway", string(edge.gatewayProtocol)),
	)

	if edge.raftEnabled {
		err = edge.runRaft()
		if err != nil {
			return
		}
	}

	edge.runGateway()
	return
}
func (edge *EdgeServer) runGateway() {
	edge.gateway.Run()
	return
}
func (edge *EdgeServer) runRaft() error {
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

func (edge *EdgeServer) Shutdown() {
	edge.gateway.Shutdown()
	edge.gatewayProtocol = gateway.Undefined
}
