package edge

import (
	"fmt"
	"git.ronaksoft.com/ronak/rony"
	"git.ronaksoft.com/ronak/rony/gateway"
	log "git.ronaksoft.com/ronak/rony/internal/logger"
	"git.ronaksoft.com/ronak/rony/internal/memberlist"
	"git.ronaksoft.com/ronak/rony/pools"
	"git.ronaksoft.com/ronak/rony/tools"
	raftbadger "github.com/bbva/raft-badger"
	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
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
   Copyright Ronak Software Group 2018
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

func (edge *Server) GetServerID() string {
	return string(edge.serverID)
}

func (edge *Server) AddHandler(constructor int64, handler ...Handler) {
	edge.handlers[constructor] = handler
}

func (edge *Server) AddBeforeHandler(constructor int64, handlers ...Handler) {

}

func (edge *Server) AddReadOnlyHandler(constructor int64, handler ...Handler) {
	edge.readonlyHandlers[constructor] = struct{}{}
	edge.AddHandler(constructor, handler...)
}

func (edge *Server) executePrepare(dispatchCtx *DispatchCtx) (err error) {
	readyOnly := false
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
			readyOnly = true
		}
	}
	err = edge.execute(dispatchCtx, readyOnly)
	return err
}
func (edge *Server) execute(dispatchCtx *DispatchCtx, readyOnly bool) (err error) {
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
				edge.executeFunc(dispatchCtx, ctx, x.Envelopes[idx], readyOnly)
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
		edge.executeFunc(dispatchCtx, ctx, dispatchCtx.req, readyOnly)
		releaseRequestCtx(ctx)
	}
	waitGroup.Wait()
	releaseWaitGroup(waitGroup)
	return nil
}
func (edge *Server) executeFunc(dispatchCtx *DispatchCtx, requestCtx *RequestCtx, in *rony.MessageEnvelope, readOnly bool) {
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
	if readOnly {
		_, ok := edge.readonlyHandlers[in.GetConstructor()]
		if !ok {
			leaderID := edge.cluster.leaderID
			if ce := log.Check(log.DebugLevel, "Redirect To Leader"); ce != nil {
				ce.Write(
					zap.String("LeaderID", leaderID),
					zap.String("State", edge.raft.State().String()),
				)
			}
			if leaderID != "" {
				requestCtx.PushMessage(
					rony.C_Redirect,
					&rony.Redirect{
						LeaderHostPort: edge.cluster.GetByID(leaderID).GatewayAddr,
						ServerID:       leaderID,
					},
				)
			} else {
				requestCtx.PushError(rony.ErrCodeUnavailable, rony.ErrItemRaftLeader)
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

// RunCluster runs the gossip and raft if it is set
func (edge *Server) RunCluster() (err error) {
	log.Info("Edge Server Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.String("Gateway", string(edge.gatewayProtocol)),
	)

	notifyChan := make(chan bool, 1)
	if edge.raftEnabled {
		err = edge.runRaft(notifyChan)
		if err != nil {
			return
		}
	}

	if edge.gatewayProtocol != gateway.Undefined {
		err = edge.runGossip()
		if err != nil {
			return
		}
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
func (edge *Server) runGossip() error {
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
func (edge *Server) runRaft(notifyChan chan bool) (err error) {
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

// RunGateway runs the gateway then we can accept clients requests
func (edge *Server) RunGateway() {
	edge.gateway.Run()
}

// JoinCluster joins this node to one or more cluster members
func (edge *Server) JoinCluster(addr ...string) error {
	_, err := edge.gossip.Join(addr)
	return err
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
