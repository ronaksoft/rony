package edge

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"go.opentelemetry.io/otel/codes"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/internal/msg"
	"github.com/ronaksoft/rony/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

/*
   Creation Time: 2020 - Feb - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type MessageKind byte

const (
	_ MessageKind = iota
	GatewayMessage
	TunnelMessage
)

var (
	messageKindNames = map[MessageKind]string{
		GatewayMessage: "GatewayMessage",
		TunnelMessage:  "TunnelMessage",
	}
)

func (c MessageKind) String() string {
	return messageKindNames[c]
}

// Server represents Edge serve. Edge server sticks all other parts of the system together.
type Server struct {
	// General
	name     string
	serverID []byte
	logger   log.Logger

	// Handlers
	preHandlers  []Handler
	handlers     map[uint64]*HandlerOption
	postHandlers []Handler

	// Edge components
	router     rony.Router
	cluster    rony.Cluster
	tunnel     rony.Tunnel
	gateway    rony.Gateway
	dispatcher Dispatcher
	restMux    *restMux

	// Tracing
	tracer     trace.Tracer
	propagator propagation.TraceContext
}

func NewServer(serverID string, opts ...Option) *Server {
	// Initialize metrics
	metrics.Init(map[string]string{
		"ServerID": serverID,
	})

	edgeServer := &Server{
		handlers:   make(map[uint64]*HandlerOption),
		serverID:   []byte(serverID),
		dispatcher: &defaultDispatcher{},
		logger:     log.With("EDGE"),
	}

	for _, opt := range opts {
		opt(edgeServer)
	}

	// register builtin rony handlers
	builtin := newBuiltin(edgeServer)
	edgeServer.SetHandler(NewHandlerOptions().SetConstructor(rony.C_GetNodes).Append(builtin.getNodes))
	edgeServer.SetHandler(NewHandlerOptions().SetConstructor(rony.C_GetAllNodes).Append(builtin.getAllNodes))
	edgeServer.SetHandler(NewHandlerOptions().SetConstructor(rony.C_Ping).Append(builtin.ping))

	return edgeServer
}

// GetServerID return this server id, which MUST be unique in the cluster otherwise
// the behaviour is unknown.
func (edge *Server) GetServerID() string {
	return string(edge.serverID)
}

// SetGlobalPreHandlers set the handler which will be called before executing any request.
// These pre handlers are like middlewares which will be automatically triggered for each request.
// If you want to set pre handler for specific request then you must use SetHandlers,
// PrependHandlers or AppendHandlers
func (edge *Server) SetGlobalPreHandlers(handlers ...Handler) {
	edge.preHandlers = handlers
}

// SetGlobalPostHandlers set the handler which will be called after executing any request.
// These pre handlers are like middlewares which will be automatically triggered for each request.
// If you want to set post handler for specific request then you must use SetHandlers,
// PrependHandlers or AppendHandlers
func (edge *Server) SetGlobalPostHandlers(handlers ...Handler) {
	edge.postHandlers = handlers
}

// SetHandler set the handlers for the constructor.
func (edge *Server) SetHandler(ho *HandlerOption) {
	for c := range ho.constructors {
		edge.handlers[c] = ho
	}
}

// GetHandler returns the handlers of the constructor
func (edge *Server) GetHandler(constructor uint64) *HandlerOption {
	return edge.handlers[constructor]
}

// SetRestProxy set a REST wrapper to expose RPCs in REST (Representational State Transfer) format
func (edge *Server) SetRestProxy(method string, path string, p RestProxy) {
	if edge.restMux == nil {
		edge.restMux = &restMux{
			routes: map[string]*trie{},
		}
	}
	edge.restMux.Set(method, path, p)
}

// SetCustomDispatcher replace the default dispatcher with your custom dispatcher. Usually you don't need to
// use custom dispatcher, but in some rare cases you can implement your own custom dispatcher
func (edge *Server) SetCustomDispatcher(d Dispatcher) {
	edge.dispatcher = d
}

// Cluster returns a reference to the underlying cluster of the Edge server
func (edge *Server) Cluster() rony.Cluster {
	return edge.cluster
}

// Gateway returns a reference to the underlying gateway of the Edge server
func (edge *Server) Gateway() rony.Gateway {
	return edge.gateway
}

// Tunnel returns a reference to the underlying tunnel of the Edge server
func (edge *Server) Tunnel() rony.Tunnel {
	return edge.tunnel
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
			ctx := acquireRequestCtx(dispatchCtx, x.SyncRun)
			waitGroup.Add(1)
			go func(ctx *RequestCtx, idx int) {
				edge.executeFunc(ctx, x.Envelopes[idx])
				if x.SyncRun {
					select {
					case ctx.nextChan <- struct{}{}:
					default:
					}
				}
				waitGroup.Done()
				releaseRequestCtx(ctx)
			}(ctx, i)

			if x.SyncRun {
				// wait until we allowed to go to next
				<-ctx.nextChan
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

	if edge.tracer != nil {
		switch in.Constructor {
		case rony.C_Ping, rony.C_GetAllNodes, rony.C_GetNodes:
		default:
			requestCtx.ctx = edge.propagator.Extract(requestCtx.ctx, in.Carrier())
			var span trace.Span
			requestCtx.ctx, span = edge.tracer.
				Start(
					requestCtx.ctx,
					fmt.Sprintf("%s/%s", ho.serviceName, ho.methodName),
					trace.WithAttributes(
						semconv.RPCServiceKey.String(ho.serviceName),
						semconv.RPCMethodKey.String(ho.methodName),
					),
					trace.WithSpanKind(trace.SpanKindServer),
				)
			defer span.End()
		}
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
	if !ho.builtin {
		for idx := range edge.postHandlers {
			edge.postHandlers[idx](requestCtx, in)
		}
	}
	if !requestCtx.stop {
		requestCtx.StopExecution()
	}

	if ce := edge.logger.Check(log.DebugLevel, "Request Executed"); ce != nil {
		ce.Write(
			zap.String("ServerID", edge.GetServerID()),
			zap.String("Kind", requestCtx.Kind().String()),
			zap.Uint64("ReqID", in.GetRequestID()),
			zap.String("C", registry.C(in.GetConstructor())),
			zap.Duration("T", time.Duration(tools.CPUTicks()-startTime)),
		)
	}

	metrics.AddCounterVec(metrics.CntRPC, 1, registry.C(in.Constructor))

	if requestCtx.err == nil {
		requestCtx.Span().SetStatus(codes.Ok, "")
	} else {
		requestCtx.Span().SetStatus(codes.Error, requestCtx.err.Error())
	}

	switch requestCtx.Kind() {
	case GatewayMessage:
		metrics.ObserveHistogram(
			metrics.HistGatewayRequestTime,
			float64(time.Duration(tools.CPUTicks()-startTime)/time.Millisecond),
		)
	case TunnelMessage:
		metrics.ObserveHistogram(
			metrics.HistTunnelRequestTime,
			float64(time.Duration(tools.CPUTicks()-startTime)/time.Millisecond),
		)
	}
}
func (edge *Server) recoverPanic(ctx *RequestCtx, in *rony.MessageEnvelope) {
	if r := recover(); r != nil {
		edge.logger.Error("Panic Recovered",
			zap.String("C", registry.C(in.Constructor)),
			zap.Uint64("ReqID", in.RequestID),
			zap.Uint64("ConnID", ctx.ConnID()),
			zap.Any("Error", r),
		)
		edge.logger.Error("Panic Stack Trace", zap.ByteString("Stack", debug.Stack()))
		ctx.PushError(errors.ErrInternalServer)
	}
}

func (edge *Server) onGatewayMessage(conn rony.Conn, streamID int64, data []byte) {
	// Fill dispatch context with data. We use the GatewayDispatcher or consume data directly based on the
	// byPassDispatcher argument
	dispatchCtx := acquireDispatchCtx(edge, conn, streamID, edge.serverID, GatewayMessage)
	defer releaseDispatchCtx(dispatchCtx)

	// if it is REST API then we take a different approach.
	if !conn.Persistent() && edge.restMux != nil {
		if rc, ok := conn.(rony.RestConn); ok {
			path, p := edge.restMux.Search(rc)
			if p != nil {
				if edge.tracer != nil {
					var span trace.Span
					dispatchCtx.ctx, span = edge.tracer.
						Start(
							dispatchCtx.ctx,
							fmt.Sprintf("%s %s", rc.Method(), path),
							trace.WithAttributes(
								semconv.HTTPMethodKey.String(rc.Method()),
							),
							trace.WithSpanKind(trace.SpanKindServer),
						)
					defer span.End()
				}

				edge.onGatewayRest(dispatchCtx, rc, p)

				return
			}
		}
	}

	err := edge.dispatcher.Decoder(data, dispatchCtx.req)
	if err != nil {
		return
	}

	err = edge.execute(dispatchCtx)
	if err != nil {
		edge.onError(dispatchCtx, errors.ErrInternalServer)
	} else {
		edge.dispatcher.Done(dispatchCtx)
	}
}
func (edge *Server) onGatewayRest(ctx *DispatchCtx, conn rony.RestConn, proxy RestProxy) {
	// apply the transformation on the client message before execute it
	err := proxy.ClientMessage(conn, ctx)
	if err != nil {
		conn.WriteStatus(http.StatusInternalServerError)
		b, _ := errors.New(errors.Internal, err.Error()).MarshalJSON()
		_ = conn.WriteBinary(0, b)

		return
	}

	err = edge.execute(ctx)
	if err != nil {
		conn.WriteStatus(http.StatusInternalServerError)
		b, _ := errors.New(errors.Internal, err.Error()).MarshalJSON()
		_ = conn.WriteBinary(0, b)

		return
	}

	// apply the transformation on the server message before sending it to the client
	err = proxy.ServerMessage(conn, ctx)
	if err != nil {
		conn.WriteStatus(http.StatusInternalServerError)
		b, _ := errors.New(errors.Internal, err.Error()).MarshalJSON()
		_ = conn.WriteBinary(0, b)
	}
}
func (edge *Server) onGatewayConnect(conn rony.Conn, kvs ...*rony.KeyValue) {
	edge.dispatcher.OnOpen(conn, kvs...)
}
func (edge *Server) onGatewayClose(conn rony.Conn) {
	edge.dispatcher.OnClose(conn)
}
func (edge *Server) onTunnelMessage(conn rony.Conn, tm *msg.TunnelMessage) {
	// Fill the dispatch context envelope from the received tunnel message
	dispatchCtx := acquireDispatchCtx(edge, conn, 0, tm.SenderID, TunnelMessage)
	tm.Envelope.DeepCopy(dispatchCtx.req)

	if err := edge.execute(dispatchCtx); err != nil {
		edge.onError(dispatchCtx, errors.ErrInternalServer)
	} else {
		edge.onTunnelDone(dispatchCtx)
	}

	// Release the dispatch context
	releaseDispatchCtx(dispatchCtx)
}
func (edge *Server) onTunnelDone(ctx *DispatchCtx) {
	tm := msg.PoolTunnelMessage.Get()
	defer msg.PoolTunnelMessage.Put(tm)
	switch ctx.BufferSize() {
	case 0:
		return
	case 1:
		tm.SenderReplicaSet = edge.cluster.ReplicaSet()
		tm.SenderID = edge.serverID
		ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
			envelope.DeepCopy(tm.Envelope)
			buf := pools.Buffer.FromProto(tm)
			_ = ctx.Conn().WriteBinary(ctx.streamID, *buf.Bytes())
			pools.Buffer.Put(buf)
		})
	default:
		// TODO:: implement it
		panic("not implemented, handle multiple tunnel message")
	}
}
func (edge *Server) onError(ctx *DispatchCtx, err *rony.Error) {
	envelope := rony.PoolMessageEnvelope.Get()
	err.ToEnvelope(envelope)
	switch ctx.kind {
	case GatewayMessage:
		buf := pools.Buffer.GetCap(1024)
		_ = edge.dispatcher.Encoder(envelope, buf)
		_ = ctx.conn.WriteBinary(ctx.streamID, *buf.Bytes())
		pools.Buffer.Put(buf)
		rony.PoolMessageEnvelope.Put(envelope)
	case TunnelMessage:
		ctx.BufferPush(envelope)
	}
}

// StartCluster is non-blocking function which runs the cluster component of the Edge server.
func (edge *Server) StartCluster() (err error) {
	if edge.cluster == nil {
		edge.logger.Warn("Cluster is NOT set",
			zap.ByteString("ServerID", edge.serverID),
		)

		return errors.ErrClusterNotSet
	}

	err = edge.cluster.Start()
	if err != nil {
		return
	}

	edge.logger.Info("Cluster Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.String("Cluster", edge.cluster.Addr()),
		zap.Uint64("ReplicaSet", edge.cluster.ReplicaSet()),
	)

	return
}

// StartGateway is non-blocking function runs the gateway in background so we can accept clients requests
func (edge *Server) StartGateway() error {
	if edge.gateway == nil {
		edge.logger.Warn("Gateway is NOT set",
			zap.ByteString("ServerID", edge.serverID),
		)

		return errors.ErrGatewayNotSet
	}
	edge.gateway.Start()

	edge.logger.Info("Gateway Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.String("Protocol", edge.gateway.Protocol().String()),
		zap.Strings("Addr", edge.gateway.Addr()),
	)

	if edge.cluster != nil {
		return errors.WrapText("cluster:")(edge.cluster.SetGatewayAddrs(edge.gateway.Addr()))
	}

	return nil
}

// StartTunnel is non-blocking function runs the gateway in background so we can accept other servers requests
func (edge *Server) StartTunnel() error {
	if edge.tunnel == nil {
		edge.logger.Warn("Tunnel is NOT set",
			zap.ByteString("ServerID", edge.serverID),
		)

		return errors.ErrTunnelNotSet
	}
	edge.tunnel.Start()

	edge.logger.Info("Tunnel Started",
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
// no error, then you must start each section separately. i.e. use StartGateway, StartCluster and StartTunnel
// functions.
func (edge *Server) Start() {
	if err := edge.StartCluster(); err != nil && err != errors.ErrClusterNotSet {
		panic(err)
	}
	if err := edge.StartGateway(); err != nil && err != errors.ErrGatewayNotSet {
		panic(err)
	}
	if err := edge.StartTunnel(); err != nil && err != errors.ErrTunnelNotSet {
		panic(err)
	}
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

	// Shutdown the tunnel
	if edge.tunnel != nil {
		edge.tunnel.Shutdown()
	}

	edge.logger.Info("Server Shutdown!", zap.ByteString("ID", edge.serverID))
}

// ShutdownWithSignal blocks until any of the signals has been called
func (edge *Server) ShutdownWithSignal(signals ...os.Signal) error {
	edge.WaitForSignal(signals...)
	if edge.cluster != nil {
		if err := edge.cluster.Leave(); err != nil {
			return err
		}
	}
	edge.Shutdown()

	return nil
}

func (edge *Server) WaitForSignal(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)

	// Wait for signal
	<-ch
}

// GetGatewayConn return the gateway connection identified by connID or returns nil if not found.
func (edge *Server) GetGatewayConn(connID uint64) rony.Conn {
	if edge.gateway == nil {
		return nil
	}

	conn := edge.gateway.GetConn(connID)
	if conn == nil {
		return nil
	}

	return conn
}

// TunnelRequest sends and receives a request through the Tunnel interface of the receiver Edge node.
func (edge *Server) TunnelRequest(replicaSet uint64, req, res *rony.MessageEnvelope) error {
	return edge.TryTunnelRequest(1, 0, replicaSet, req, res)
}
func (edge *Server) TryTunnelRequest(
	attempts int, retryWait time.Duration, replicaSet uint64,
	req, res *rony.MessageEnvelope,
) error {
	if edge.tunnel == nil {
		return errors.ErrTunnelNotSet
	}
	startTime := tools.CPUTicks()
	err := tools.Try(attempts, retryWait, func() error {
		target := edge.getReplicaMember(replicaSet)
		if target == nil {
			return errors.ErrMemberNotFound
		}

		return edge.sendRemoteCommand(target, req, res)
	})
	metrics.ObserveHistogram(
		metrics.HistTunnelRoundTripTime,
		float64(time.Duration(tools.CPUTicks()-startTime)/time.Millisecond),
	)

	return err
}
func (edge *Server) getReplicaMember(replicaSet uint64) (target rony.ClusterMember) {
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
func (edge *Server) sendRemoteCommand(target rony.ClusterMember, req, res *rony.MessageEnvelope) error {
	conn, err := target.Dial()
	if err != nil {
		return err
	}

	// Get a rony.TunnelMessage from pool and put it back into the pool when we are done
	tmOut := msg.PoolTunnelMessage.Get()
	defer msg.PoolTunnelMessage.Put(tmOut)
	tmIn := msg.PoolTunnelMessage.Get()
	defer msg.PoolTunnelMessage.Put(tmIn)
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
