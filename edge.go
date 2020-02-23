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
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/serf/serf"
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

	// Raft & Gossip
	raftEnabled   bool
	raftPort      int
	raftBootstrap bool
	raftFSM       raftFSM
	raft          *raft.Raft
	gossipPort    int
	gossip        *serf.Serf
	badgerStore   *raftbadger.BadgerStore
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
		gossipPort:    7946,
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
		_, err = raftCmd.MarshalTo(raftCmdBytes)
		if err != nil {
			return
		}
		f := edge.raft.Apply(raftCmdBytes, raftApplyTimeout)
		err = f.Error()
		pbytes.Put(raftCmdBytes)
		pools.ReleaseRaftCommand(raftCmd)
	} else {
		err = edge.execute(authID, userID, req)
	}
	return
}

func (edge *EdgeServer) execute(authID, userID int64, req *msg.MessageEnvelope) error {
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
			return err
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
	return nil
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

func (edge *EdgeServer) onMessage(conn gateway.Conn, streamID int64, data []byte) {
	authID := tools.RandomInt64(0)
	userID := tools.RandomInt64(0)
	envelope := pools.AcquireMessageEnvelope()
	_ = envelope.Unmarshal(data)
	err := edge.Execute(authID, userID, envelope)
	if err != nil {
		log.Warn("Error OnMessage (Execute)", zap.Error(err))
		envelope := pools.AcquireMessageEnvelope()
		apiErr := errors.NewError(errors.ErrCodeInvalid, errors.ErrItemApi)
		apiErr.ToMessageEnvelope(envelope)
		envelopeBytes := pbytes.GetLen(envelope.Size())
		_, err = envelope.MarshalTo(envelopeBytes)
		if err != nil {
			log.Warn("Error OnMessage (Marshal)", zap.Error(err))
		}
		err = conn.SendBinary(streamID, envelopeBytes)
		if err != nil {
			log.Warn("Error OnMessage (SendBinary)", zap.Error(err))
		}
		pools.ReleaseMessageEnvelope(envelope)
	}
	pools.ReleaseMessageEnvelope(envelope)
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
		time.Sleep(time.Second * 5)
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
	eventChan := make(chan serf.Event, 1000)
	serfConfig := serf.DefaultConfig()
	serfConfig.NodeName = edge.GetServerID()
	serfConfig.Init()

	serfConfig.CoalescePeriod = time.Second
	serfConfig.LogOutput = ioutil.Discard
	serfConfig.EventCh = eventChan
	serfConfig.MemberlistConfig.Delegate = &gossipDelegate{edge: edge}
	serfConfig.MemberlistConfig.Events = &gossipEventDelegate{edge: edge}
	serfConfig.MemberlistConfig.Ping = &gossipPingDelegate{edge: edge}
	serfConfig.MemberlistConfig.Alive = &gossipAliveDelegate{edge: edge}
	serfConfig.MemberlistConfig.LogOutput = ioutil.Discard
	dirPath := filepath.Join(edge.dataPath, "gossip")
	_ = os.MkdirAll(dirPath, os.ModePerm)
	serfConfig.KeyringFile = filepath.Join(dirPath, "keyring")
	serfConfig.SnapshotPath = filepath.Join(dirPath, "snapshot")
	serfConfig.MemberlistConfig.BindPort = edge.gossipPort

	if s, err := serf.Create(serfConfig); err != nil {
		return err
	} else {
		edge.gossip = s
	}
	_ = edge.gossip.SetTags(map[string]string{
		"BundleID":   edge.bundleID,
		"InstanceID": edge.instanceID,
		"RaftNodeID": edge.GetServerID(),
		"RaftPort":   fmt.Sprintf("%d", edge.raftPort),
	})

	go func() {
		for range eventChan {
			for _, m := range edge.gossip.Members() {
				if m.Tags["BundleID"] == edge.bundleID && m.Tags["InstanceID"] != edge.instanceID {
					nodeAddr := fmt.Sprintf("%s:%s", m.Addr.String(), m.Tags["RaftPort"])
					err := edge.joinRaft(m.Tags["RaftNodeID"], nodeAddr)
					if err != nil {
						log.Info("Error On Join Raft",
							zap.String("Addr", nodeAddr),
							zap.Error(err),
						)
					}
				}
			}
		}
	}()

	return nil
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
	}

	return nil
}

func (edge *EdgeServer) Join(addr ...string) error {
	if !edge.raftEnabled {
		return errors.ErrRaftNotSet
	}
	_, err := edge.gossip.Join(addr, false)
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

type EdgeStats struct {
	Address         string
	RaftMembers     int
	RaftState       string
	SerfMembers     int
	SerfState       string
	GatewayProtocol gateway.Protocol
	GatewayAddr     string
}

func (edge *EdgeServer) Stats() EdgeStats {
	if !edge.raftEnabled {
		return EdgeStats{}
	}

	s := EdgeStats{
		Address:         fmt.Sprintf("%s:%d", edge.gossip.LocalMember().Addr.String(), edge.gossip.LocalMember().Port),
		RaftMembers:     0,
		RaftState:       edge.raft.State().String(),
		SerfMembers:     len(edge.gossip.Members()),
		SerfState:       edge.gossip.State().String(),
		GatewayProtocol: edge.gatewayProtocol,
		GatewayAddr:     edge.gateway.Addr(),
	}

	f := edge.raft.GetConfiguration()
	if f.Error() == nil {
		s.RaftMembers = len(f.Configuration().Servers)
	}

	return s
}
