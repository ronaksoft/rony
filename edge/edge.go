package edge

import (
	"bufio"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/gateway"
	tcpGateway "github.com/ronaksoft/rony/internal/gateway/tcp"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/internal/store/localdb"
	"github.com/ronaksoft/rony/internal/tunnel"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/rest"
	"github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
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

type Server struct {
	// General
	dataDir  string
	serverID []byte

	// Handlers
	preHandlers  []Handler
	handlers     map[int64]*HandlerOption
	postHandlers []Handler

	// Edge components
	cluster    cluster.Cluster
	tunnel     tunnel.Tunnel
	store      store.Store
	gateway    gateway.Gateway
	dispatcher Dispatcher
}

func NewServer(serverID string, opts ...Option) *Server {
	// Initialize metrics
	metrics.Init(map[string]string{
		"ServerID": serverID,
	})

	edgeServer := &Server{
		dataDir:    "./_hdd",
		handlers:   make(map[int64]*HandlerOption),
		serverID:   []byte(serverID),
		dispatcher: &defaultDispatcher{},
	}

	for _, opt := range opts {
		opt(edgeServer)
	}

	if edgeServer.store == nil {
		edgeServer.store, _ = localdb.New(localdb.DefaultConfig(edgeServer.dataDir))
	}

	// register builtin rony handlers
	builtin := newBuiltin(edgeServer.GetServerID(), edgeServer.Gateway(), edgeServer.Cluster())
	edgeServer.SetHandler(NewHandlerOptions().SetConstructor(rony.C_GetNodes).Append(builtin.GetNodes))
	edgeServer.SetHandler(NewHandlerOptions().SetConstructor(rony.C_GetAllNodes).Append(builtin.GetAllNodes))
	edgeServer.SetHandler(NewHandlerOptions().SetConstructor(rony.C_GetPage).Append(builtin.GetPage).setBuiltin())

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

// SetHandler set the handlers for the constructor. 'leaderOnly' is applicable ONLY if the cluster is run
// with Raft support. If cluster is a Raft enabled cluster, then by setting 'leaderOnly' to TRUE, requests sent
// to a follower server will return redirect error to the client. For standalone servers 'leaderOnly' does not
// affect.
func (edge *Server) SetHandler(ho *HandlerOption) {
	for c := range ho.constructors {
		edge.handlers[c] = ho
	}

}

// GetHandler returns the handlers of the constructor
func (edge *Server) GetHandler(constructor int64) *HandlerOption {
	return edge.handlers[constructor]
}

// SetRestWrapper set a REST wrapper to expose RPCs in REST (Representational State Transfer) format
func (edge *Server) SetRestWrapper(method string, path string, f *rest.Factory) {
	switch gw := edge.gateway.(type) {
	case *tcpGateway.Gateway:
		if !gw.Support(gateway.Http) {
			panic("tcp gateway does not support http protocol")
		}
		gw.SetProxy(method, path, f)
	default:
		panic("only works on tcp gateway")
	}
}

// Cluster returns a reference to the underlying cluster of the Edge server
func (edge *Server) Cluster() cluster.Cluster {
	return edge.cluster
}

func (edge *Server) Store() store.Store {
	return edge.store
}

// Gateway returns a reference to the underlying gateway of the Edge server
func (edge *Server) Gateway() gateway.Gateway {
	return edge.gateway
}

func (edge *Server) execute(dispatchCtx *DispatchCtx) (err error) {
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
				edge.executeFunc(ctx, x.Envelopes[idx])
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
		edge.executeFunc(ctx, dispatchCtx.req)
		releaseRequestCtx(ctx)
	}
	waitGroup.Wait()
	pools.ReleaseWaitGroup(waitGroup)
	return nil
}
func (edge *Server) executeFunc(requestCtx *RequestCtx, in *rony.MessageEnvelope) {
	defer edge.recoverPanic(requestCtx, in)

	startTime := tools.CPUTicks()

	// Set the context request
	requestCtx.reqID = in.RequestID

	ho, ok := edge.handlers[in.GetConstructor()]
	if !ok {
		requestCtx.PushError(errors.ErrInvalidHandler)
		return
	}

	// Run the handler
	if !ho.builtin {
		for idx := range edge.preHandlers {
			edge.preHandlers[idx](requestCtx, in)
			if requestCtx.stop {
				break
			}
		}
	}
	if !requestCtx.stop {
		for idx := range ho.handlers {
			ho.handlers[idx](requestCtx, in)
			if requestCtx.stop {
				break
			}
		}
	}
	if !requestCtx.stop && !ho.builtin {
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
		ce.Write(
			zap.String("ServerID", edge.GetServerID()),
			zap.String("Kind", requestCtx.Kind().String()),
			zap.Uint64("ReqID", in.GetRequestID()),
			zap.String("C", registry.ConstructorName(in.GetConstructor())),
			zap.Duration("T", time.Duration(tools.CPUTicks()-startTime)),
		)
	}
	switch requestCtx.Kind() {
	case GatewayMessage:
		metrics.ObserveHistogram(metrics.HistGatewayRequestTime, float64(time.Duration(tools.CPUTicks()-startTime)/time.Millisecond))
	case TunnelMessage:
		metrics.ObserveHistogram(metrics.HistTunnelRequestTime, float64(time.Duration(tools.CPUTicks()-startTime)/time.Millisecond))
	}
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
		ctx.PushError(errors.ErrInternalServer)
	}
}

func (edge *Server) onGatewayMessage(conn rony.Conn, streamID int64, data []byte, byPassDispatcher bool) {
	// _, task := trace.NewTask(context.Background(), "onGatewayMessage")
	// defer task.End()

	var (
		err error
	)

	// Fill dispatch context with data. We use the GatewayDispatcher or consume data directly based on the
	// byPassDispatcher argument
	dispatchCtx := acquireDispatchCtx(edge, conn, streamID, edge.serverID, GatewayMessage)
	if byPassDispatcher {
		err = dispatchCtx.UnmarshalEnvelope(data)
	} else {
		err = edge.dispatcher.Interceptor(dispatchCtx, data)
	}
	if err != nil {
		releaseDispatchCtx(dispatchCtx)
		return
	}

	if err := edge.execute(dispatchCtx); err != nil {
		edge.onError(dispatchCtx, errors.ErrInternalServer)
	} else {
		edge.dispatcher.Done(dispatchCtx)
	}

	// Release the dispatch context
	releaseDispatchCtx(dispatchCtx)
}
func (edge *Server) onGatewayConnect(conn rony.Conn, kvs ...*rony.KeyValue) {
	edge.dispatcher.OnOpen(conn, kvs...)
}
func (edge *Server) onGatewayClose(conn rony.Conn) {
	edge.dispatcher.OnClose(conn)
}
func (edge *Server) onTunnelMessage(conn rony.Conn, tm *rony.TunnelMessage) {
	// _, task := trace.NewTask(context.Background(), "onTunnelMessage")
	// defer task.End()

	// Fill the dispatch context envelope from the received tunnel message
	dispatchCtx := acquireDispatchCtx(edge, conn, 0, tm.SenderID, TunnelMessage)
	dispatchCtx.FillEnvelope(
		tm.Envelope.GetRequestID(), tm.Envelope.GetConstructor(), tm.Envelope.Message,
		tm.Envelope.Auth, tm.Envelope.Header...,
	)

	if err := edge.execute(dispatchCtx); err != nil {
		edge.onError(dispatchCtx, errors.ErrInternalServer)
	} else {
		edge.onTunnelDone(dispatchCtx)
	}

	// Release the dispatch context
	releaseDispatchCtx(dispatchCtx)
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
	buf := pools.Buffer.FromProto(tm)
	_ = ctx.Conn().SendBinary(ctx.streamID, *buf.Bytes())
	pools.Buffer.Put(buf)
}
func (edge *Server) onError(ctx *DispatchCtx, err *rony.Error) {
	envelope := rony.PoolMessageEnvelope.Get()
	err.ToEnvelope(envelope)
	switch ctx.kind {
	case GatewayMessage:
		edge.dispatcher.OnMessage(ctx, envelope)
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
		zap.String("Protocol", edge.gateway.Protocol().String()),
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

	// Shutdown the store
	if edge.store != nil {
		edge.store.Shutdown()
	}

	log.Info("Server Shutdown!", zap.ByteString("ID", edge.serverID))
}

// ShutdownWithSignal blocks until any of the signals has been called
func (edge *Server) ShutdownWithSignal(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)

	// Wait for signal
	<-ch
	edge.Shutdown()
}

// GetGatewayConn return the gateway connection identified by connID or returns nil if not found.
func (edge *Server) GetGatewayConn(connID uint64) rony.Conn {
	if edge.gateway == nil {
		return nil
	}
	return edge.gateway.GetConn(connID)
}

// ExecuteRemote sends and receives a request through the Tunnel interface of the receiver Edge node.
func (edge *Server) ExecuteRemote(replicaSet uint64, req, res *rony.MessageEnvelope) error {
	return edge.TryExecuteRemote(1, 0, replicaSet, req, res)
}
func (edge *Server) TryExecuteRemote(attempts int, retryWait time.Duration, replicaSet uint64, req, res *rony.MessageEnvelope) error {
	startTime := tools.CPUTicks()
	err := tools.Try(attempts, retryWait, func() error {
		target := edge.getReplicaMember(replicaSet)
		if target == nil {
			return ErrMemberNotFound
		}

		return edge.sendRemoteCommand(target, req, res)
	})
	metrics.ObserveHistogram(metrics.HistTunnelRoundtripTime, float64(time.Duration(tools.CPUTicks()-startTime)/time.Millisecond))
	return err
}
func (edge *Server) getReplicaMember(replicaSet uint64) (target cluster.Member) {
	members := edge.cluster.MembersByReplicaSet(replicaSet)
	if len(members) == 0 {
		return nil
	}
	for idx := range members {
		target = members[idx]
		break
	}
	return
}
func (edge *Server) sendRemoteCommand(target cluster.Member, req, res *rony.MessageEnvelope) error {
	conn, err := target.Dial()
	if err != nil {
		return err
	}

	// Get a rony.TunnelMessage from pool and put it back into the pool when we are done
	tmOut := rony.PoolTunnelMessage.Get()
	defer rony.PoolTunnelMessage.Put(tmOut)
	tmIn := rony.PoolTunnelMessage.Get()
	defer rony.PoolTunnelMessage.Put(tmIn)
	tmOut.Fill(edge.serverID, edge.cluster.ReplicaSet(), req)

	// Marshal and send over the wire
	buf := pools.Buffer.FromProto(tmOut)
	_, err = conn.Write(*buf.Bytes())
	pools.Buffer.Put(buf)
	if err != nil {
		return err
	}

	// Wait for response and unmarshal it
	buf = pools.Buffer.GetLen(4096)
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	n, err := bufio.NewReader(conn).Read(*buf.Bytes())
	if err != nil || n == 0 {
		return err
	}
	err = tmIn.Unmarshal((*buf.Bytes())[:n])
	if err != nil {
		return err
	}
	pools.Buffer.Put(buf)

	// deep copy
	tmIn.Envelope.DeepCopy(res)
	return nil
}
