package edge

import (
	"fmt"
	raftbadger "github.com/bbva/raft-badger"
	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/gateway"
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

/*
   Creation Time: 2020 - Feb - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Handler func(ctx *RequestCtx, in *rony.MessageEnvelope)
type GetConstructorNameFunc func(constructor int64) string

// Dispatcher
type Dispatcher interface {
	// All the input arguments are valid in the function context, if you need to pass 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it.
	OnMessage(ctx *DispatchCtx, authID int64, envelope *rony.MessageEnvelope)
	// All the input arguments are valid in the function context, if you need to pass 'data' or 'envelope' to other
	// async functions, make sure to hard copy (clone) it before sending it. If 'err' is not nil then envelope will be
	// discarded, it is the user's responsibility to send back appropriate message using 'conn'
	// Note that conn IS NOT nil in any circumstances.
	Prepare(ctx *DispatchCtx, data []byte, kvs ...gateway.KeyValue) (err error)
	// This will be called when the context has been finished, this lets cleaning up, or in case you need to flush the
	// messages and updates in one go.
	Done(ctx *DispatchCtx)
	// This will be called when a new connection has been opened
	OnOpen(conn gateway.Conn)
	// This will be called when a connection is closed
	OnClose(conn gateway.Conn)
}

// Server
type Server struct {
	// General
	serverID        []byte
	replicaSet      uint64
	dataPath        string
	gatewayProtocol gateway.Protocol
	gateway         gateway.Gateway
	dispatcher      Dispatcher

	// Handlers
	preHandlers        []Handler
	handlers           map[int64][]Handler
	postHandlers       []Handler
	readonlyHandlers   map[int64]struct{}
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

func NewServer(serverID string, dispatcher Dispatcher, opts ...Option) *Server {
	edgeServer := &Server{
		handlers:         make(map[int64][]Handler),
		readonlyHandlers: make(map[int64]struct{}),
		serverID:         []byte(serverID),
		dispatcher:       dispatcher,
		getConstructorName: func(constructor int64) string {
			return fmt.Sprintf("%d", constructor)
		},
		dataPath:      ".",
		raftEnabled:   false,
		raftPort:      0,
		raftBootstrap: false,
		rateLimitChan: make(chan struct{}, clusterMessageRateLimit),
		gossipPort:    7946,
		cluster: Cluster{
			byServerID: make(map[string]*ClusterMember, 100),
		},
	}

	for _, opt := range opts {
		opt(edgeServer)
	}

	return edgeServer
}

// GetServerID return this server id, which MUST be unique in the cluster otherwise
// the behaviour is unknown.
func (edge *Server) GetServerID() string {
	return string(edge.serverID)
}

// SetHandlers set the handlers for the constructor. 'leaderOnly' is applicable ONLY if the cluster is run
// with Raft support. If cluster is a Raft enabled cluster, then by setting 'leaderOnly' to TRUE, requests sent
// to a follower server will return redirect error to the client. For standalone servers 'leaderOnly' does not
// affect.
func (edge *Server) SetHandlers(constructor int64, leaderOnly bool, handlers ...Handler) {
	if !leaderOnly {
		edge.readonlyHandlers[constructor] = struct{}{}
	}
	edge.handlers[constructor] = handlers
}

// AppendHandlers appends the handlers for the constructor in order. So handlers[n] will be called before
// handlers[n+1].
func (edge *Server) AppendHandlers(constructor int64, handlers ...Handler) {
	edge.handlers[constructor] = append(edge.handlers[constructor], handlers...)
}

// PrependHandlers prepends the handlers for the constructor in order.
func (edge *Server) PrependHandlers(constructor int64, handlers ...Handler) {
	edge.handlers[constructor] = append(handlers, edge.handlers[constructor]...)
}

func (edge *Server) executePrepare(dispatchCtx *DispatchCtx) (err error) {
	// If server is standalone then we are the leader anyway
	isLeader := true
	if edge.raftEnabled {
		if edge.raft.State() == raft.Leader {
			raftCmd := acquireRaftCommand()
			raftCmd.Fill(edge.serverID, dispatchCtx.authID, dispatchCtx.req)
			mo := proto.MarshalOptions{}
			raftCmdBytes := pools.Bytes.GetCap(mo.Size(raftCmd))
			raftCmdBytes, err = mo.MarshalAppend(raftCmdBytes, raftCmd)
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
		} else {
			isLeader = false
		}
	}
	err = edge.execute(dispatchCtx, isLeader)
	return err
}
func (edge *Server) execute(dispatchCtx *DispatchCtx, isLeader bool) (err error) {
	waitGroup := acquireWaitGroup()
	switch dispatchCtx.req.GetConstructor() {
	case rony.C_MessageContainer:
		x := &rony.MessageContainer{}
		mo := proto.UnmarshalOptions{
			Merge: true,
		}
		err = mo.Unmarshal(dispatchCtx.req.Message, x)
		if err != nil {
			return
		}
		xLen := len(x.Envelopes)
		for i := 0; i < xLen; i++ {
			ctx := acquireRequestCtx(dispatchCtx, true)
			nextChan := make(chan struct{}, 1)
			waitGroup.Add(1)
			go func(ctx *RequestCtx, idx int) {
				edge.executeFunc(dispatchCtx, ctx, x.Envelopes[idx], isLeader)
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
		edge.executeFunc(dispatchCtx, ctx, dispatchCtx.req, isLeader)
		releaseRequestCtx(ctx)
	}
	waitGroup.Wait()
	releaseWaitGroup(waitGroup)
	return nil
}
func (edge *Server) executeFunc(dispatchCtx *DispatchCtx, requestCtx *RequestCtx, in *rony.MessageEnvelope, isLeader bool) {
	defer edge.recoverPanic(requestCtx, in)

	var startTime time.Time

	if ce := log.Check(log.DebugLevel, "Execute (Start)"); ce != nil {
		startTime = time.Now()
		ce.Write(
			zap.String("Constructor", edge.getConstructorName(in.GetConstructor())),
			zap.Uint64("RequestID", in.GetRequestID()),
			zap.Int64("AuthID", dispatchCtx.authID),
		)
	}
	if !isLeader {
		_, ok := edge.readonlyHandlers[in.GetConstructor()]
		if !ok {
			if ce := log.Check(log.DebugLevel, "Redirect To Leader"); ce != nil {
				ce.Write(
					zap.String("LeaderID", edge.cluster.leaderID),
					zap.String("State", edge.raft.State().String()),
				)
			}
			if leaderID := edge.cluster.leaderID; leaderID == "" {
				requestCtx.PushError(rony.ErrCodeUnavailable, rony.ErrItemRaftLeader)
			} else {
				requestCtx.PushMessage(
					rony.C_Redirect,
					&rony.Redirect{
						LeaderHostPort: edge.cluster.GetByID(leaderID).GatewayAddr,
						ServerID:       leaderID,
					},
				)
			}
			return
		}
	}

	handlers, ok := edge.handlers[in.GetConstructor()]
	if !ok {
		requestCtx.PushError(rony.ErrCodeInvalid, rony.ErrItemHandler)
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
			zap.Uint64("RequestID", in.GetRequestID()),
			zap.Int64("AuthID", dispatchCtx.authID),
			zap.Duration("T", time.Now().Sub(startTime)),
		)
	}

	return
}
func (edge *Server) recoverPanic(ctx *RequestCtx, in *rony.MessageEnvelope) {
	if r := recover(); r != nil {
		log.Error("Panic Recovered",
			zap.ByteString("ServerID", edge.serverID),
			zap.Uint64("ConnID", ctx.ConnID()),
			zap.Int64("AuthID", ctx.AuthID()),
			zap.Any("Error", r),
		)
		ctx.PushError(rony.ErrCodeInternal, rony.ErrItemServer)
	}
}

func (edge *Server) HandleGatewayMessage(conn gateway.Conn, streamID int64, data []byte, kvs ...gateway.KeyValue) {
	// _, task := trace.NewTask(context.Background(), "Handle Gateway Message")
	// defer task.End()

	dispatchCtx := acquireDispatchCtx(edge, conn, streamID, 0, edge.serverID)
	err := edge.dispatcher.Prepare(dispatchCtx, data, kvs...)
	if err != nil {
		releaseDispatchCtx(dispatchCtx)
		return
	}
	err = edge.executePrepare(dispatchCtx)
	if err != nil {
		edge.onError(dispatchCtx, rony.ErrCodeInternal, rony.ErrItemServer)
	}
	edge.dispatcher.Done(dispatchCtx)
	releaseDispatchCtx(dispatchCtx)
}
func (edge *Server) onError(dispatchCtx *DispatchCtx, code, item string) {
	envelope := acquireMessageEnvelope()
	rony.ErrorMessage(envelope, dispatchCtx.req.GetRequestID(), code, item)
	edge.dispatcher.OnMessage(dispatchCtx, dispatchCtx.authID, envelope)
	releaseMessageEnvelope(envelope)
}
func (edge *Server) onConnect(conn gateway.Conn) {
	edge.dispatcher.OnOpen(conn)
}
func (edge *Server) onClose(conn gateway.Conn) {
	edge.dispatcher.OnClose(conn)
}

// StartCluster is non-blocking function which runs the gossip and raft if it is set
func (edge *Server) StartCluster() (err error) {
	log.Info("Edge Server Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.String("Gateway", string(edge.gatewayProtocol)),
		zap.Bool("Raft", edge.raftEnabled),
		zap.Int("GossipPort", edge.gossipPort),
	)

	notifyChan := make(chan bool, 1)
	if edge.raftEnabled {
		err = edge.startRaft(notifyChan)
		if err != nil {
			return
		}
	}

	err = edge.startGossip()
	if err != nil {
		return
	}
	go func() {
		for range notifyChan {
			err := tools.Try(10, time.Millisecond, func() error {
				return edge.updateCluster(gossipUpdateTimeout)
			})
			if err != nil {
				log.Warn("Rony got error on updating the cluster",
					zap.Error(err),
					zap.ByteString("ID", edge.serverID),
				)
			}
		}
	}()

	return
}
func (edge *Server) startGossip() error {
	dirPath := filepath.Join(edge.dataPath, "gossip")
	_ = os.MkdirAll(dirPath, os.ModePerm)

	conf := memberlist.DefaultWANConfig()
	conf.Name = string(edge.serverID)
	conf.Events = &delegateEvents{edge: edge}
	conf.Delegate = &delegateNode{edge: edge}
	conf.LogOutput = ioutil.Discard
	conf.Logger = nil
	conf.BindPort = edge.gossipPort
	if s, err := memberlist.Create(conf); err != nil {
		return err
	} else {
		edge.gossip = s
	}

	return edge.updateCluster(gossipUpdateTimeout)
}
func (edge *Server) startRaft(notifyChan chan bool) (err error) {
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
	// raftConfig.LogLevel = "DEBUG"
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
			if err == raft.ErrCantBootstrap {
				log.Info("Error On Raft Bootstrap", zap.Error(err))
			} else {
				log.Warn("Error On Raft Bootstrap", zap.Error(err))
			}

		}
	}

	return nil
}

// StartGateway is non-blocking function runs the gateway in background so we can accept clients requests
func (edge *Server) StartGateway() {
	edge.gateway.Run()
}

// JoinCluster is used to take an existing Cluster and attempt to join a cluster
// by contacting all the given hosts and performing a state sync.
// This returns the number of hosts successfully contacted and an error if
// none could be reached. If an error is returned, the node did not successfully
// join the cluster.
func (edge *Server) JoinCluster(addr ...string) (int, error) {
	return edge.gossip.Join(addr)
}

// Shutdown gracefully shutdown the services
func (edge *Server) Shutdown() {
	// First shutdown gateway to not accept any more request
	if edge.gateway != nil {
		edge.gateway.Shutdown()
	}

	// Second shutdown raft, if it is enabled
	if edge.raftEnabled {
		if f := edge.raft.Snapshot(); f.Error() != nil {
			if f.Error() != raft.ErrNothingNewToSnapshot {
				log.Warn("Error On Shutdown (Raft Snapshot)",
					zap.Error(f.Error()),
					zap.String("ServerID", edge.GetServerID()),
				)
			} else {
				log.Info("Error On Shutdown (Raft Snapshot)", zap.Error(f.Error()))
			}

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
	if edge.gossip != nil {
		err := edge.gossip.Leave(gossipLeaveTimeout)
		if err != nil {
			log.Warn("Error On Leaving Cluster, but we shutdown anyway", zap.Error(err))
		}
		err = edge.gossip.Shutdown()
		if err != nil {
			log.Warn("Error On Shutdown (Gossip)", zap.Error(err))
		}
	}

	edge.gatewayProtocol = gateway.Undefined
	log.Info("Server Shutdown!", zap.ByteString("ID", edge.serverID))
}

// Shutdown blocks until any of the signals has been called
func (edge *Server) ShutdownWithSignal(signals ...os.Signal) {
	ch := make(chan os.Signal)
	signal.Notify(ch, signals...)

	// Wait for signal
	<-ch
	edge.Shutdown()
}
