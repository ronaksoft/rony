package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/memberlist"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	raftbadger "github.com/bbva/raft-badger"
	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"io/ioutil"
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

//go:generate protoc -I=./vendor -I=.  --gogofaster_out=. msg.proto
type Handler func(ctx *RequestCtx, in *MessageEnvelope)
type GetConstructorNameFunc func(constructor int64) string

// Dispatcher
type Dispatcher interface {
	// All the input arguments are valid in the function context, if you need to pass 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it.
	DispatchUpdate(ctx *DispatchCtx, authID int64, envelope *UpdateEnvelope)
	// All the input arguments are valid in the function context, if you need to pass 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it.
	DispatchMessage(ctx *DispatchCtx, authID int64, envelope *MessageEnvelope)
	// All the input arguments are valid in the function context, if you need to pass 'data' or 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it. If 'err' is not nil then envelope will be
	// discarded, it is the user's responsibility to send back appropriate message using 'conn'
	// Note that conn IS NOT nil in any circumstances.
	DispatchRequest(ctx *DispatchCtx, data []byte) (err error)
}

// EdgeServer
type EdgeServer struct {
	// General
	serverID        []byte
	replicaSet      uint32
	shardSet        uint32
	shardMin        uint32
	shardMax        uint32
	dataPath        string
	gatewayProtocol gateway.Protocol
	gateway         gateway.Gateway
	dispatcher      Dispatcher

	// Handlers
	preHandlers        []Handler
	handlers           map[int64][]Handler
	postHandlers       []Handler
	getConstructorName GetConstructorNameFunc

	// Raft & Gossip
	raftEnabled   bool
	raftPort      int
	raftBootstrap bool
	raftFSM       raftFSM
	raft          *raft.Raft
	rateLimitChan chan struct{}
	gossipPort    int
	gossip        *memberlist.Memberlist
	cluster       Cluster
	badgerStore   *raftbadger.BadgerStore
}

func NewEdgeServer(serverID string, dispatcher Dispatcher, opts ...Option) *EdgeServer {
	edgeServer := &EdgeServer{
		handlers:   make(map[int64][]Handler),
		serverID:   []byte(serverID),
		dispatcher: dispatcher,
		getConstructorName: func(constructor int64) string {
			return fmt.Sprintf("%d", constructor)
		},
		dataPath:      ".",
		raftEnabled:   false,
		raftPort:      0,
		raftBootstrap: false,
		rateLimitChan: make(chan struct{}, clusterMessageRateLimit),
		gossipPort:    7946,
	}

	for _, opt := range opts {
		opt(edgeServer)
	}

	return edgeServer
}

func (edge *EdgeServer) GetServerID() string {
	return string(edge.serverID)
}

func (edge *EdgeServer) AddHandler(constructor int64, handler ...Handler) {
	edge.handlers[constructor] = handler
}

func (edge *EdgeServer) executePrepare(dispatchCtx *DispatchCtx) (err error) {
	if edge.raftEnabled {
		if edge.raft.State() != raft.Leader {
			return ErrNotRaftLeader
		}
		raftCmd := acquireRaftCommand()
		raftCmd.AuthID = dispatchCtx.authID
		raftCmd.Envelope = dispatchCtx.req
		raftCmdBytes := pools.Bytes.GetLen(raftCmd.Size())
		_, err = raftCmd.MarshalTo(raftCmdBytes)
		if err != nil {
			return
		}
		f := edge.raft.Apply(raftCmdBytes, raftApplyTimeout)
		err = f.Error()
		pools.Bytes.Put(raftCmdBytes)
		releaseRaftCommand(raftCmd)
		if err != nil {
			return
		}
	}
	err = edge.execute(dispatchCtx)
	return err
}
func (edge *EdgeServer) execute(dispatchCtx *DispatchCtx) (err error) {
	waitGroup := acquireWaitGroup()
	switch dispatchCtx.req.Constructor {
	case C_MessageContainer:
		x := &MessageContainer{}
		err = x.Unmarshal(dispatchCtx.req.Message)
		if err != nil {
			return
		}
		xLen := len(x.Envelopes)
		for i := 0; i < xLen; i++ {
			ctx := acquireRequestCtx(dispatchCtx, true)
			nextChan := make(chan struct{}, 1)
			waitGroup.Add(1)
			go func(ctx *RequestCtx, idx int) {
				edge.executeFunc(dispatchCtx, ctx, x.Envelopes[idx])
				nextChan <- struct{}{}
				waitGroup.Done()
				releaseRequestCtx(ctx)
			}(ctx, i)
			select {
			case <-ctx.nextChan:
				// The handler supported quick return
			case <-nextChan:
			}
		}
	default:
		ctx := acquireRequestCtx(dispatchCtx, false)
		edge.executeFunc(dispatchCtx, ctx, dispatchCtx.req)
		releaseRequestCtx(ctx)
	}
	waitGroup.Wait()
	releaseWaitGroup(waitGroup)
	return nil
}
func (edge *EdgeServer) executeFunc(dispatchCtx *DispatchCtx, requestCtx *RequestCtx, in *MessageEnvelope) {
	defer edge.recoverPanic(dispatchCtx)

	startTime := time.Now()
	if ce := log.Check(log.DebugLevel, "Execute (Start)"); ce != nil {
		ce.Write(
			zap.String("Constructor", edge.getConstructorName(in.Constructor)),
			zap.Uint64("RequestID", in.RequestID),
			zap.Int64("AuthID", dispatchCtx.authID),
		)
	}
	handlers, ok := edge.handlers[in.Constructor]
	if !ok {
		requestCtx.PushError(in.RequestID, ErrCodeInvalid, ErrItemHandler)
		return
	}

	// Run the handler
	for idx := range edge.preHandlers {
		edge.preHandlers[idx](requestCtx, in)
		if requestCtx.stop {
			break
		}
	}
	if !requestCtx.stop {
		for idx := range handlers {
			handlers[idx](requestCtx, in)
			if requestCtx.stop {
				break
			}
		}
	}
	if !requestCtx.stop {
		for idx := range edge.postHandlers {
			edge.postHandlers[idx](requestCtx, in)
			if requestCtx.stop {
				break
			}
		}
	}
	if !requestCtx.stop {
		requestCtx.StopExecution()
	}

	if ce := log.Check(log.DebugLevel, "Execute (Finished)"); ce != nil {
		ce.Write(
			zap.Uint64("RequestID", in.RequestID),
			zap.Int64("AuthID", dispatchCtx.authID),
			zap.Duration("T", time.Now().Sub(startTime)),
		)
	}

	return
}
func (edge *EdgeServer) recoverPanic(dispatchCtx *DispatchCtx) {
	if r := recover(); r != nil {
		log.Error("Panic Recovered",
			zap.ByteString("ServerID", edge.serverID),
			zap.Uint64("ConnID", dispatchCtx.conn.GetConnID()),
			zap.Int64("AuthID", dispatchCtx.authID),
			zap.ByteString("Stack", debug.Stack()),
		)
	}
}

func (edge *EdgeServer) onGatewayMessage(conn gateway.Conn, streamID int64, data []byte) {
	dispatchCtx := acquireDispatchCtx(edge, conn, streamID, 0, edge.serverID)

	err := edge.dispatcher.DispatchRequest(dispatchCtx, data)
	if err != nil {
		releaseDispatchCtx(dispatchCtx)
		return
	}
	err = edge.executePrepare(dispatchCtx)
	switch err {
	case nil:
	case ErrNotRaftLeader:
		edge.onError(dispatchCtx, ErrCodeUnavailable, ErrItemRaftLeader)
	default:
		edge.onError(dispatchCtx, ErrCodeInternal, ErrItemServer)
	}
	releaseDispatchCtx(dispatchCtx)
}
func (edge *EdgeServer) onError(dispatchCtx *DispatchCtx, code, item string) {
	envelope := acquireMessageEnvelope()
	ErrorMessage(envelope, code, item)
	edge.dispatcher.DispatchMessage(dispatchCtx, dispatchCtx.authID, envelope)
	releaseMessageEnvelope(envelope)
}
func (edge *EdgeServer) onConnect(connID uint64)   {}
func (edge *EdgeServer) onClose(conn gateway.Conn) {}
func (edge *EdgeServer) onFlush(conn gateway.Conn) [][]byte {
	return nil
}

// Run runs the selected gateway, if gateway is not setup it returns error
func (edge *EdgeServer) Run() (err error) {
	if edge.gatewayProtocol == gateway.Undefined {
		return ErrGatewayNotInitialized
	}

	log.Info("Edge Server Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.String("Gateway", string(edge.gatewayProtocol)),
	)

	// We must run the gateway before gossip since some information about the gateway will be spread to other nodes
	// by gossip protocol.
	edge.runGateway()

	notifyChan := make(chan bool, 1)
	if edge.raftEnabled {
		err = edge.runRaft(notifyChan)
		if err != nil {
			return
		}
	}
	err = edge.runGossip()
	if err != nil {
		return
	}

	go func() {
		for range notifyChan {
			err := tools.Try(10, time.Millisecond, func() error {
				return edge.updateCluster(gossipUpdateTimeout)
			})
			if err != nil {
				log.Warn("Error On Update Cluster", zap.Error(err))
			}
		}
	}()
	return
}
func (edge *EdgeServer) runGateway() {
	edge.gateway.Run()
	return
}
func (edge *EdgeServer) runGossip() error {
	dirPath := filepath.Join(edge.dataPath, "gossip")
	_ = os.MkdirAll(dirPath, os.ModePerm)

	conf := memberlist.DefaultWANConfig()
	conf.Name = string(edge.serverID)
	conf.Events = &delegateEvents{edge: edge}
	conf.Delegate = &delegateNode{edge: edge}
	conf.LogOutput = ioutil.Discard
	conf.BindPort = edge.gossipPort
	if s, err := memberlist.Create(conf); err != nil {
		return err
	} else {
		edge.gossip = s
	}

	return edge.updateCluster(gossipUpdateTimeout)
}
func (edge *EdgeServer) runRaft(notifyChan chan bool) (err error) {
	// Initialize LogStore for Raft
	dirPath := filepath.Join(edge.dataPath, "raft")
	_ = os.MkdirAll(dirPath, os.ModePerm)
	badgerOpt := badger.DefaultOptions(dirPath).WithLogger(nil)
	edge.badgerStore, err = raftbadger.New(raftbadger.Options{
		Path:                dirPath,
		BadgerOptions:       &badgerOpt,
		NoSync:              false,
		ValueLogGC:          false,
		GCInterval:          0,
		MandatoryGCInterval: 0,
		GCThreshold:         0,
	})
	if err != nil {
		return
	}

	// Initialize Raft
	raftConfig := raft.DefaultConfig()
	raftConfig.LogLevel = "WARN"
	raftConfig.NotifyCh = notifyChan
	raftConfig.Logger = hclog.NewNullLogger()
	raftConfig.LocalID = raft.ServerID(edge.serverID)
	raftBind := fmt.Sprintf(":%d", edge.raftPort)
	raftAdvertiseAddr, err := net.ResolveTCPAddr("tcp", raftBind)
	if err != nil {
		return err
	}

	raftTransport, err := raft.NewTCPTransport(raftBind, raftAdvertiseAddr, 3, 10*time.Second, os.Stdout)
	if err != nil {
		return err
	}

	raftSnapshot, err := raft.NewFileSnapshotStore(dirPath, 3, os.Stdout)
	if err != nil {
		return err
	}

	edge.raft, err = raft.NewRaft(raftConfig, edge.raftFSM, edge.badgerStore, edge.badgerStore, raftSnapshot, raftTransport)
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
		f := edge.raft.BootstrapCluster(bootConfig)
		if err := f.Error(); err != nil {
			log.Warn("Error On Raft Bootstrap", zap.Error(err))
		}
		time.Sleep(time.Second * 3)
	}

	return nil
}

// JoinCluster joins this node to one or more cluster members
func (edge *EdgeServer) JoinCluster(addr ...string) error {
	_, err := edge.gossip.Join(addr)
	return err
}
func (edge *EdgeServer) joinRaft(nodeID, addr string) error {
	if !edge.raftEnabled {
		return ErrRaftNotSet
	}
	futureConfig := edge.raft.GetConfiguration()
	if err := futureConfig.Error(); err != nil {
		return err
	}
	for _, srv := range futureConfig.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) && srv.Address == raft.ServerAddress(addr) {
			return nil
		}
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			future := edge.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return err
			}
		}
	}

	future := edge.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if err := future.Error(); err != nil {
		return err
	}
	return nil
}

// Shutdown gracefully shutdown the services
func (edge *EdgeServer) Shutdown() {
	// First shutdown gateway to not accept
	edge.gateway.Shutdown()

	// Second shutdown raft, if it is enabled
	if edge.raftEnabled {
		if f := edge.raft.Snapshot(); f.Error() != nil {
			log.Warn("Error On Shutdown (Raft Snapshot)", zap.Error(f.Error()))
		}
		if f := edge.raft.Shutdown(); f.Error() != nil {
			log.Warn("Error On Shutdown (Raft Shutdown)", zap.Error(f.Error()))
		}

		err := edge.badgerStore.Close()
		if err != nil {
			log.Warn("Error On Shutdown (Close Raft Badger)", zap.Error(err))
		}
	}

	// Shutdown gossip
	err := edge.gossip.Leave(gossipLeaveTimeout)
	if err != nil {
		log.Warn("Error On Leaving Cluster, but we shutdown anyway", zap.Error(err))
	}
	err = edge.gossip.Shutdown()
	if err != nil {
		log.Warn("Error On Shutdown (Gossip)", zap.Error(err))
	}

	edge.gatewayProtocol = gateway.Undefined
	log.Info("Server Shutdown!", zap.ByteString("ID", edge.serverID))
}
