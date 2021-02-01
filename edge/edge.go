package edge

import (
	"errors"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cluster"
	"github.com/ronaksoft/rony/gateway"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
	"github.com/ronaksoft/rony/tunnel"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"os"
	"os/signal"
	"runtime/debug"
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

func init() {
	log.Init(log.DefaultConfig)
}

// Server
type Server struct {
	// General
	serverID []byte

	// Handlers
	preHandlers      []Handler
	handlers         map[int64][]Handler
	postHandlers     []Handler
	readonlyHandlers map[int64]struct{}

	// Gateway's Configs
	gatewayProtocol   gateway.Protocol
	gateway           gateway.Gateway
	gatewayDispatcher Dispatcher

	// Cluster Configs
	cluster cluster.Cluster

	// Tunnel Configs
	tunnel tunnel.Tunnel
}

func NewServer(serverID string, opts ...Option) *Server {
	edgeServer := &Server{
		handlers:          make(map[int64][]Handler),
		readonlyHandlers:  make(map[int64]struct{}),
		serverID:          []byte(serverID),
		gatewayDispatcher: &defaultDispatcher{},
	}

	for _, opt := range opts {
		opt(edgeServer)
	}

	// register default rony handlers
	edgeServer.SetHandlers(rony.C_GetNodes, false, edgeServer.getNodes)

	return edgeServer
}

// GetServerID return this server id, which MUST be unique in the cluster otherwise
// the behaviour is unknown.
func (edge *Server) GetServerID() string {
	return string(edge.serverID)
}

// SetGlobalPreHandlers set the handler which will be called before executing any request. These pre handlers are like middlewares
// which will be automatically triggered for each request. If you want to set pre handler for specific request the your must
// use SetHandlers, PrependHandlers or AppendHandlers
func (edge *Server) SetGlobalPreHandlers(handlers ...Handler) {
	edge.preHandlers = handlers
}

// SetGlobalPostHandlers set the handler which will be called after executing any request. These pre handlers are like middlewares
// which will be automatically triggered for each request. If you want to set post handler for specific request the your must
// use SetHandlers, PrependHandlers or AppendHandlers
func (edge *Server) SetGlobalPostHandlers(handlers ...Handler) {
	edge.postHandlers = handlers
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

// BulkAppendHandlers appends the handlers for all the constructors. This method is useful when you
// a fixed set of pre-handlers for many of your rpc.it appends the handlers for the constructor in order.
// So handlers[n] will be called before handlers[n+1]
func (edge *Server) BulkAppendHandlers(handlers []Handler, constructors ...int64) {
	for _, c := range constructors {
		edge.handlers[c] = append(edge.handlers[c], handlers...)
	}
}

// PrependHandlers prepends the handlers for the constructor in order.
func (edge *Server) PrependHandlers(constructor int64, handlers ...Handler) {
	edge.handlers[constructor] = append(handlers, edge.handlers[constructor]...)
}

// BulkPrependHandlers prepends the handlers for all the constructors. This method is useful when you
// a fixed set of post-handlers for many of your rpc.it appends the handlers for the constructor in order.
func (edge *Server) BulkPrependHandlers(handlers []Handler, constructors ...int64) {
	for _, c := range constructors {
		edge.handlers[c] = append(handlers, edge.handlers[c]...)
	}
}

// Cluster returns a reference to the underlying cluster of the Edge server
func (edge *Server) Cluster() cluster.Cluster {
	return edge.cluster
}

func (edge *Server) executePrepare(dispatchCtx *DispatchCtx) (err error, isLeader bool) {
	// If server is standalone then we are the leader anyway
	isLeader = true
	if edge.cluster == nil || !edge.cluster.RaftEnabled() {
		return
	}

	if edge.cluster.RaftState() == raft.Leader {
		raftCmd := rony.PoolRaftCommand.Get()
		raftCmd.Envelope = rony.PoolMessageEnvelope.Get()
		raftCmd.Fill(edge.serverID, dispatchCtx.req)
		mo := proto.MarshalOptions{UseCachedSize: true}
		buf := pools.Buffer.GetCap(mo.Size(raftCmd))
		var raftCmdBytes []byte
		raftCmdBytes, err = mo.MarshalAppend(*buf.Bytes(), raftCmd)
		if err != nil {
			return
		}
		f := edge.cluster.RaftApply(raftCmdBytes)
		err = f.Error()
		pools.Buffer.Put(buf)
		rony.PoolRaftCommand.Put(raftCmd)
		if err != nil {
			return
		}
	} else {
		isLeader = false
	}
	return
}
func (edge *Server) execute(dispatchCtx *DispatchCtx, isLeader bool) (err error) {
	waitGroup := pools.AcquireWaitGroup()
	switch dispatchCtx.req.GetConstructor() {
	case rony.C_MessageContainer:
		x := &rony.MessageContainer{}
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
				edge.executeFunc(ctx, x.Envelopes[idx], isLeader)
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
		edge.executeFunc(ctx, dispatchCtx.req, isLeader)
		releaseRequestCtx(ctx)
	}
	waitGroup.Wait()
	pools.ReleaseWaitGroup(waitGroup)
	return nil
}
func (edge *Server) executeFunc(requestCtx *RequestCtx, in *rony.MessageEnvelope, isLeader bool) {
	defer edge.recoverPanic(requestCtx, in)

	startTime := tools.CPUTicks()

	// Set the context request
	requestCtx.reqID = in.RequestID

	if !isLeader && requestCtx.dispatchCtx.kind != ReplicaMessage {
		_, ok := edge.readonlyHandlers[in.GetConstructor()]
		if !ok {
			if ce := log.Check(log.DebugLevel, "Redirect To Leader"); ce != nil {
				ce.Write(
					zap.String("RaftLeaderID", edge.cluster.RaftLeaderID()),
					zap.String("State", edge.cluster.RaftState().String()),
				)
			}
			requestCtx.PushRedirectLeader()
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

	if ce := log.Check(log.DebugLevel, "Request Executed"); ce != nil {
		if requestCtx.Kind() != ReplicaMessage {
			ce.Write(
				zap.String("ServerID", edge.GetServerID()),
				zap.String("Kind", requestCtx.Kind().String()),
				zap.Uint64("ReqID", in.GetRequestID()),
				zap.String("C", registry.ConstructorName(in.GetConstructor())),
				zap.Duration("T", time.Duration(tools.CPUTicks()-startTime)),
			)
		}
	}

	return
}
func (edge *Server) recoverPanic(ctx *RequestCtx, in *rony.MessageEnvelope) {
	if r := recover(); r != nil {
		log.Error("Panic Recovered",
			zap.String("C", registry.ConstructorName(in.Constructor)),
			zap.Uint64("ReqID", in.RequestID),
			zap.Uint64("ConnID", ctx.ConnID()),
			zap.Any("Error", r),
		)
		log.Error("Panic Stack Trace", zap.ByteString("Stack", debug.Stack()))
		ctx.PushError(rony.ErrCodeInternal, rony.ErrItemServer)
	}
}

func (edge *Server) onReplicaMessage(raftCmd *rony.RaftCommand) error {
	// _, task := trace.NewTask(context.Background(), "onReplicaMessage")
	// defer task.End()

	dispatchCtx := acquireDispatchCtx(edge, nil, 0, raftCmd.Sender, ReplicaMessage)
	dispatchCtx.FillEnvelope(
		raftCmd.Envelope.GetRequestID(), raftCmd.Envelope.GetConstructor(), raftCmd.Envelope.Message,
		raftCmd.Envelope.Auth, raftCmd.Envelope.Header...,
	)
	err := edge.execute(dispatchCtx, false)
	if err != nil {
		return err
	}
	releaseDispatchCtx(dispatchCtx)
	return nil
}
func (edge *Server) onGatewayMessage(conn rony.Conn, streamID int64, data []byte) {
	// _, task := trace.NewTask(context.Background(), "onGatewayMessage")
	// defer task.End()

	dispatchCtx := acquireDispatchCtx(edge, conn, streamID, edge.serverID, GatewayMessage)
	err := edge.gatewayDispatcher.Interceptor(dispatchCtx, data)
	if err != nil {
		releaseDispatchCtx(dispatchCtx)
		return
	}
	err, isLeader := edge.executePrepare(dispatchCtx)
	if err != nil {
		edge.onError(dispatchCtx, rony.ErrCodeInternal, rony.ErrItemServer)
	} else if err = edge.execute(dispatchCtx, isLeader); err != nil {
		edge.onError(dispatchCtx, rony.ErrCodeInternal, rony.ErrItemServer)
	} else {
		edge.gatewayDispatcher.Done(dispatchCtx)
	}
	releaseDispatchCtx(dispatchCtx)
	return
}
func (edge *Server) onGatewayConnect(conn rony.Conn, kvs ...*rony.KeyValue) {
	edge.gatewayDispatcher.OnOpen(conn, kvs...)
}
func (edge *Server) onGatewayClose(conn rony.Conn) {
	edge.gatewayDispatcher.OnClose(conn)
}
func (edge *Server) onTunnelMessage(conn rony.Conn, tm *rony.TunnelMessage) {
	// _, task := trace.NewTask(context.Background(), "onTunnelMessage")
	// defer task.End()

	dispatchCtx := acquireDispatchCtx(edge, conn, 0, tm.SenderID, TunnelMessage)
	dispatchCtx.FillEnvelope(
		tm.Envelope.GetRequestID(), tm.Envelope.GetConstructor(), tm.Envelope.Message,
		tm.Envelope.Auth, tm.Envelope.Header...,
	)

	err, isLeader := edge.executePrepare(dispatchCtx)
	if err != nil {
		edge.onError(dispatchCtx, rony.ErrCodeInternal, rony.ErrItemServer)
	} else if err = edge.execute(dispatchCtx, isLeader); err != nil {
		edge.onError(dispatchCtx, rony.ErrCodeInternal, rony.ErrItemServer)
	} else {
		edge.onTunnelDone(dispatchCtx)
	}
	releaseDispatchCtx(dispatchCtx)
	return
}
func (edge *Server) onTunnelDone(ctx *DispatchCtx) {
	tm := rony.PoolTunnelMessage.Get()
	defer rony.PoolTunnelMessage.Put(tm)
	switch ctx.BufferSize() {
	case 0:
		return
	case 1:
		tm.SenderReplicaSet = edge.cluster.ReplicaSet()
		tm.SenderID = edge.serverID
		tm.Envelope = ctx.BufferPop()
	default:
		// TODO:: implement it
		panic("not implemented, handle multiple tunnel message")
	}

	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Bytes.GetCap(mo.Size(tm))
	b, _ = mo.MarshalAppend(b, tm)
	_ = ctx.Conn().SendBinary(ctx.streamID, b)
}

func (edge *Server) onError(ctx *DispatchCtx, code, item string) {
	envelope := rony.PoolMessageEnvelope.Get()
	rony.ErrorMessage(envelope, ctx.req.GetRequestID(), code, item)
	switch ctx.kind {
	case GatewayMessage:
		edge.gatewayDispatcher.OnMessage(ctx, envelope)
	case TunnelMessage:
		ctx.BufferPush(envelope.Clone())
	}
	rony.PoolMessageEnvelope.Put(envelope)
}

// StartCluster is non-blocking function which runs the gossip and raft if it is set
func (edge *Server) StartCluster() (err error) {
	if edge.cluster == nil {
		return ErrClusterNotSet
	}

	edge.cluster.Start()
	log.Info("Edge Server:: Cluster Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.String("Cluster", edge.cluster.Addr()),
		zap.Uint64("ReplicaSet", edge.cluster.ReplicaSet()),
	)

	return
}

// StartGateway is non-blocking function runs the gateway in background so we can accept clients requests
func (edge *Server) StartGateway() error {
	if edge.gateway == nil {
		return ErrGatewayNotSet
	}
	edge.gateway.Start()

	log.Info("Edge Server:: Gateway Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.String("Protocol", string(edge.gatewayProtocol)),
		zap.Strings("Addr", edge.gateway.Addr()),
	)

	if edge.cluster != nil {
		return edge.cluster.SetGatewayAddrs(edge.gateway.Addr())
	}

	return nil
}

// StartTunnel is non-blocking function runs the gateway in background so we can accept other servers requests
func (edge *Server) StartTunnel() error {
	if edge.tunnel == nil {
		return ErrTunnelNotSet
	}
	edge.tunnel.Start()

	log.Info("Edge Server:: Tunnel Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.Strings("Addr", edge.tunnel.Addr()),
	)

	if edge.cluster != nil {
		return edge.cluster.SetTunnelAddrs(edge.tunnel.Addr())
	}

	return nil
}

// Start is a helper function which tries to start all three parts of the edge server
// This function does not return any error, if you need to make sure all the parts are started with
// no error, then you must start each section separately.
func (edge *Server) Start() {
	_ = edge.StartCluster()
	_ = edge.StartGateway()
	_ = edge.StartTunnel()
}

// JoinCluster is used to take an existing Cluster and attempt to join a cluster
// by contacting all the given hosts and performing a state sync.
// This returns the number of hosts successfully contacted and an error if
// none could be reached. If an error is returned, the node did not successfully
// join the cluster.
func (edge *Server) JoinCluster(addr ...string) (int, error) {
	return edge.cluster.Join(addr...)
}

// Shutdown gracefully shutdown the services
func (edge *Server) Shutdown() {
	// First shutdown gateway to not accept any more request
	if edge.gateway != nil {
		edge.gateway.Shutdown()
	}

	// Shutdown the cluster
	if edge.cluster != nil {
		edge.cluster.Shutdown()
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

// GetGatewayConn return the gateway connection identified by connID or returns nil if not found.
func (edge *Server) GetGatewayConn(connID uint64) rony.Conn {
	return edge.gateway.GetConn(connID)
}

var (
	ErrClusterNotSet            = errors.New("cluster is not set")
	ErrGatewayNotSet            = errors.New("gateway is not set")
	ErrTunnelNotSet             = errors.New("tunnel is not set")
	ErrUnexpectedTunnelResponse = errors.New("unexpected tunnel response")
	ErrEmptyMemberList          = errors.New("member list is empty")
	ErrMemberNotFound           = errors.New("member not found")
	ErrNoTunnelAddrs            = errors.New("tunnel address does not found")
)
