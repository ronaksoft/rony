package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/context"
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/memberlist"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/msg"
	raftbadger "github.com/bbva/raft-badger"
	"github.com/dgraph-io/badger/v2"
	"github.com/gobwas/pool/pbytes"
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

//go:generate protoc -I=./vendor -I=.  --gogofaster_out=. node.proto
type Handler func(ctx *context.Context, in *msg.MessageEnvelope)
type GetConstructorNameFunc func(constructor int64) string

// Dispatcher
type Dispatcher interface {
	// All the input arguments are valid in the function context, if you need to pass 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it.
	DispatchUpdate(conn gateway.Conn, streamID int64, authID int64, envelope *msg.UpdateEnvelope)
	// All the input arguments are valid in the function context, if you need to pass 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it.
	DispatchMessage(conn gateway.Conn, streamID int64, authID int64, envelope *msg.MessageEnvelope)
	// All the input arguments are valid in the function context, if you need to pass 'data' or 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it. If 'err' is not nil then envelope will be
	// discarded, it is the user's responsibility to send back appropriate message using 'conn'
	DispatchRequest(conn gateway.Conn, streamID int64, data []byte, envelope *msg.MessageEnvelope) (authID int64, err error)

	DispatchClusterMessage(envelope *msg.MessageEnvelope)
}

// EdgeServer
type EdgeServer struct {
	// General
	bundleID        string
	instanceID      string
	serverID        string
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
	gossipPort    int
	gossip        *memberlist.Memberlist
	cluster       Cluster
	badgerStore   *raftbadger.BadgerStore
}

func NewEdgeServer(bundleID, instanceID string, dispatcher Dispatcher, opts ...Option) *EdgeServer {
	edgeServer := &EdgeServer{
		handlers:   make(map[int64][]Handler),
		bundleID:   bundleID,
		instanceID: instanceID,
		serverID:   fmt.Sprintf("%s.%s", bundleID, instanceID),
		dispatcher: dispatcher,
		getConstructorName: func(constructor int64) string {
			return fmt.Sprintf("%d", constructor)
		},
		dataPath:      ".",
		raftEnabled:   false,
		raftPort:      0,
		raftBootstrap: false,
		gossipPort:    7946,
	}

	for _, opt := range opts {
		opt(edgeServer)
	}

	return edgeServer
}

func (edge *EdgeServer) GetServerID() string {
	return edge.serverID
}

func (edge *EdgeServer) AddHandler(constructor int64, handler ...Handler) {
	edge.handlers[constructor] = handler
}

func (edge *EdgeServer) execute(conn gateway.Conn, streamID, authID int64, req *msg.MessageEnvelope) (err error) {
	if edge.raftEnabled {
		if edge.raft.State() != raft.Leader {
			return errors.ErrNotRaftLeader
		}
		raftCmd := pools.AcquireRaftCommand()
		raftCmd.AuthID = authID
		raftCmd.Envelope = req
		raftCmdBytes := pbytes.GetLen(raftCmd.Size())
		_, err = raftCmd.MarshalTo(raftCmdBytes)
		if err != nil {
			return
		}
		f := edge.raft.Apply(raftCmdBytes, raftApplyTimeout)
		err = f.Error()
		pbytes.Put(raftCmdBytes)
		pools.ReleaseRaftCommand(raftCmd)
		if err != nil {
			return
		}
	}
	executeFunc := func(ctx *context.Context, in *msg.MessageEnvelope) {
		defer edge.recoverPanic(ctx)
		startTime := time.Now()

		waitGroup := pools.AcquireWaitGroup()
		waitGroup.Add(2)
		go func() {
			for u := range ctx.UpdateChan {
				edge.dispatcher.DispatchUpdate(conn, streamID, u.AuthID, u.Envelope)
				if ce := log.Check(log.DebugLevel, "Update Dispatched"); ce != nil {
					ce.Write(
						zap.Uint64("ConnID", conn.GetConnID()),
						zap.Int64("AuthID", u.AuthID),
						zap.Int64("UpdateID", u.Envelope.UpdateID),
						zap.String("C", edge.getConstructorName(u.Envelope.Constructor)),
					)
				}

				pools.ReleaseUpdateEnvelope(u.Envelope)
			}
			waitGroup.Done()
		}()
		go func() {
			for m := range ctx.MessageChan {
				edge.dispatcher.DispatchMessage(conn, streamID, m.AuthID, m.Envelope)
				if ce := log.Check(log.DebugLevel, "Message Dispatched"); ce != nil {
					ce.Write(
						zap.Uint64("ConnID", conn.GetConnID()),
						zap.Int64("AuthID", m.AuthID),
						zap.Uint64("RequestID", m.Envelope.RequestID),
						zap.String("C", edge.getConstructorName(m.Envelope.Constructor)),
					)
				}
				pools.ReleaseMessageEnvelope(m.Envelope)
			}
			waitGroup.Done()
		}()

		if ce := log.Check(log.DebugLevel, "Execute (Start)"); ce != nil {
			ce.Write(
				zap.String("Constructor", edge.getConstructorName(in.Constructor)),
				zap.Uint64("RequestID", in.RequestID),
				zap.Int64("AuthID", ctx.AuthID),
			)
		}
		handlers, ok := edge.handlers[in.Constructor]
		if !ok {
			ctx.PushError(in.RequestID, errors.ErrCodeInvalid, errors.ErrItemHandler)
			return
		}

		// Run the handler
		for idx := range edge.preHandlers {
			edge.preHandlers[idx](ctx, in)
			if ctx.Stop {
				break
			}
		}
		if !ctx.Stop {
			for idx := range handlers {
				handlers[idx](ctx, in)
				if ctx.Stop {
					break
				}
			}
		}
		if !ctx.Stop {
			for idx := range edge.postHandlers {
				edge.postHandlers[idx](ctx, in)
				if ctx.Stop {
					break
				}
			}
		}
		if !ctx.Stop {
			ctx.StopExecution()
		}
		waitGroup.Wait()
		pools.ReleaseWaitGroup(waitGroup)

		duration := time.Now().Sub(startTime)
		if ce := log.Check(log.DebugLevel, "Execute (Finished)"); ce != nil {
			ce.Write(
				zap.Uint64("RequestID", in.RequestID),
				zap.Int64("AuthID", ctx.AuthID),
				zap.Duration("T", duration),
			)
		}

		return
	}
	switch req.Constructor {
	case msg.C_MessageContainer:
		x := &msg.MessageContainer{}
		err = x.Unmarshal(req.Message)
		if err != nil {
			return
		}
		xLen := len(x.Envelopes)
		waitGroup := pools.AcquireWaitGroup()
		for i := 0; i < xLen; i++ {
			ctx := context.Acquire(conn.GetConnID(), authID, true, false)
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
		ctx := context.Acquire(conn.GetConnID(), authID, false, false)
		executeFunc(ctx, req)
		context.Release(ctx)
	}
	return nil
}
func (edge *EdgeServer) recoverPanic(ctx *context.Context) {
	if r := recover(); r != nil {
		log.Error("Panic Recovered",
			zap.String("ServerID", edge.GetServerID()),
			zap.Uint64("ConnID", ctx.ConnID),
			zap.Int64("AuthID", ctx.AuthID),
			zap.ByteString("Stack", debug.Stack()),
		)
	}
}

func (edge *EdgeServer) onMessage(conn gateway.Conn, streamID int64, data []byte) {
	envelope := pools.AcquireMessageEnvelope()
	authID, err := edge.dispatcher.DispatchRequest(conn, streamID, data, envelope)
	if err != nil {
		// pools.ReleaseMessageEnvelope(envelope)
		return
	}
	startTime := time.Now()
	err = edge.execute(conn, streamID, authID, envelope)
	if ce := log.Check(log.DebugLevel, "Execute Time"); ce != nil {
		ce.Write(zap.Duration("D", time.Now().Sub(startTime)))
	}
	switch err {
	case nil:
	case errors.ErrNotRaftLeader:
		edge.onError(conn, streamID, authID, errors.ErrCodeUnavailable, errors.ErrItemRaftLeader)
	default:
		log.Warn("Error On Execute", zap.Error(err))
		edge.onError(conn, streamID, authID, errors.ErrCodeInternal, errors.ErrItemServer)

	}
	pools.ReleaseMessageEnvelope(envelope)
}
func (edge *EdgeServer) onError(conn gateway.Conn, streamID, authID int64, code, item string) {
	envelope := pools.AcquireMessageEnvelope()
	errors.NewError(code, item).ToMessageEnvelope(envelope)
	edge.dispatcher.DispatchMessage(conn, streamID, authID, envelope)
	pools.ReleaseMessageEnvelope(envelope)
}
func (edge *EdgeServer) onConnect(connID uint64)   {}
func (edge *EdgeServer) onClose(conn gateway.Conn) {}
func (edge *EdgeServer) onFlush(conn gateway.Conn) [][]byte {
	return nil
}

// Run runs the selected gateway, if gateway is not setup it returns error
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
		err = edge.runGossip()
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
func (edge *EdgeServer) runGossip() error {
	dirPath := filepath.Join(edge.dataPath, "gossip")
	_ = os.MkdirAll(dirPath, os.ModePerm)

	conf := memberlist.DefaultLANConfig()
	conf.Name = edge.GetServerID()
	conf.Events = &delegateEvents{edge: edge}
	conf.Delegate = &delegateNode{edge: edge}
	conf.LogOutput = ioutil.Discard
	conf.BindPort = edge.gossipPort
	if s, err := memberlist.Create(conf); err != nil {
		return err
	} else {
		edge.gossip = s
	}

	return edge.updateCluster()
}
func (edge *EdgeServer) runRaft() (err error) {
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
	raftConfig.Logger = hclog.NewNullLogger()
	raftConfig.LocalID = raft.ServerID(edge.GetServerID())
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

func (edge *EdgeServer) JoinCluster(addr ...string) error {
	if !edge.raftEnabled {
		return errors.ErrRaftNotSet
	}
	_, err := edge.gossip.Join(addr)

	return err
}
func (edge *EdgeServer) joinRaft(nodeID, addr string) error {
	if !edge.raftEnabled {
		return errors.ErrRaftNotSet
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

func (edge *EdgeServer) Shutdown() {
	// First shutdown gateway to not accept
	edge.gateway.Shutdown()

	if edge.raftEnabled {
		err := edge.gossip.Shutdown()
		if err != nil {
			log.Warn("Error On Shutdown (Gossip)", zap.Error(err))
		}

		sf := edge.raft.Snapshot()
		if err := sf.Error(); err != nil {
			log.Warn("Error On Shutdown (Raft Snapshot)", zap.Error(err))
		}
		f := edge.raft.Shutdown()
		if err := f.Error(); err != nil {
			log.Warn("Error On Shutdown (Raft Shutdown)", zap.Error(err))
		}

		err = edge.badgerStore.Close()
		if err != nil {
			log.Warn("Error On Shutdown (Close Badger", zap.Error(err))
		}
	}

	edge.gatewayProtocol = gateway.Undefined
	log.Info("Server Shutdown!", zap.String("ID", edge.GetServerID()))
}

// EdgeStats exports some internal metrics data
type EdgeStats struct {
	Address         string
	RaftMembers     int
	RaftState       string
	Members         int
	MembershipScore int
	GatewayProtocol gateway.Protocol
	GatewayAddr     string
}

// Stats exports some internal metrics data packed in 'EdgeStats' struct
func (edge *EdgeServer) Stats() EdgeStats {
	if !edge.raftEnabled {
		return EdgeStats{}
	}

	s := EdgeStats{
		Address:         fmt.Sprintf("%s:%d", edge.gossip.LocalNode().Addr.String(), edge.gossip.LocalNode().Port),
		RaftMembers:     0,
		RaftState:       edge.raft.State().String(),
		Members:         len(edge.gossip.Members()),
		MembershipScore: edge.gossip.GetHealthScore(),
		GatewayProtocol: edge.gatewayProtocol,
		GatewayAddr:     edge.gateway.Addr(),
	}

	f := edge.raft.GetConfiguration()
	if f.Error() == nil {
		s.RaftMembers = len(f.Configuration().Servers)
	}

	return s
}
