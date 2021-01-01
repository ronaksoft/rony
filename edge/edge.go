package edge

import (
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cluster"
	"github.com/ronaksoft/rony/gateway"
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/registry"
	"github.com/ronaksoft/rony/tools"
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

type Handler func(ctx *RequestCtx, in *rony.MessageEnvelope)

// Server
type Server struct {
	// General
	serverID        []byte
	gatewayProtocol gateway.Protocol
	gateway         gateway.Gateway
	dispatcher      Dispatcher

	// Handlers
	preHandlers      []Handler
	handlers         map[int64][]Handler
	postHandlers     []Handler
	readonlyHandlers map[int64]struct{}

	// Raft & Gossip
	cluster *cluster.Cluster
}

func NewServer(serverID string, dispatcher Dispatcher, opts ...Option) *Server {
	edgeServer := &Server{
		handlers:         make(map[int64][]Handler),
		readonlyHandlers: make(map[int64]struct{}),
		serverID:         []byte(serverID),
		dispatcher:       dispatcher,
	}
	edgeServer.cluster = cluster.New(
		edgeServer.serverID,
		edgeServer.onReplicaMessage,
		edgeServer.onClusterMessage,
	)

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

// SetPreHandlers set the handler which will be called before executing any request. These pre handlers are like middlewares
// which will be automatically triggered for each request. If you want to set pre handler for specific request the your must
// use SetHandlers, PrependHandlers or AppendHandlers
func (edge *Server) SetPreHandlers(handlers ...Handler) {
	edge.preHandlers = handlers
}

// SetPreHandlers set the handler which will be called after executing any request. These pre handlers are like middlewares
// which will be automatically triggered for each request. If you want to set post handler for specific request the your must
// use SetHandlers, PrependHandlers or AppendHandlers
func (edge *Server) SetPostHandlers(handlers ...Handler) {
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

// PrependHandlers prepends the handlers for the constructor in order.
func (edge *Server) PrependHandlers(constructor int64, handlers ...Handler) {
	edge.handlers[constructor] = append(handlers, edge.handlers[constructor]...)
}

// ClusterMembers returns a list of all the discovered nodes in the cluster
func (edge *Server) ClusterMembers() []*cluster.Member {
	return edge.cluster.Members()
}

// ClusterSend sends 'envelope' to the server identified by 'serverID'. It may returns ErrNotFound if the server
// is not in the list. The message will be send with BEST EFFORT and using UDP
func (edge *Server) ClusterSend(serverID string, envelope *rony.MessageEnvelope, kvs ...*rony.KeyValue) (err error) {
	return edge.cluster.Send(serverID, envelope, kvs...)
}

func (edge *Server) executePrepare(dispatchCtx *DispatchCtx) (err error, isLeader bool) {
	// If server is standalone then we are the leader anyway
	isLeader = true
	if !edge.cluster.RaftEnabled() {
		return
	}

	if edge.cluster.RaftState() == raft.Leader {
		raftCmd := acquireRaftCommand()
		raftCmd.Fill(edge.serverID, dispatchCtx.req)
		mo := proto.MarshalOptions{UseCachedSize: true}
		raftCmdBytes := pools.Bytes.GetCap(mo.Size(raftCmd))
		raftCmdBytes, err = mo.MarshalAppend(raftCmdBytes, raftCmd)
		if err != nil {
			return
		}
		f := edge.cluster.RaftApply(raftCmdBytes, raftApplyTimeout)
		err = f.Error()
		pools.Bytes.Put(raftCmdBytes)
		releaseRaftCommand(raftCmd)
		if err != nil {
			return
		}
	} else {
		isLeader = false
	}
	return
}
func (edge *Server) execute(dispatchCtx *DispatchCtx, isLeader bool) (err error) {
	waitGroup := acquireWaitGroup()
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
	releaseWaitGroup(waitGroup)
	return nil
}
func (edge *Server) executeFunc(requestCtx *RequestCtx, in *rony.MessageEnvelope, isLeader bool) {
	defer edge.recoverPanic(requestCtx, in)

	var startTime int64
	if ce := log.Check(log.DebugLevel, "Execute (Start)"); ce != nil {
		startTime = tools.CPUTicks()
		ce.Write(
			zap.String("C", registry.ConstructorName(in.GetConstructor())),
			zap.Uint64("ReqID", in.GetRequestID()),
		)
	}

	// Set the context request
	requestCtx.reqID = in.RequestID

	if !isLeader {
		_, ok := edge.readonlyHandlers[in.GetConstructor()]
		if !ok {
			if ce := log.Check(log.DebugLevel, "Redirect To Leader"); ce != nil {
				ce.Write(
					zap.String("LeaderID", edge.cluster.LeaderID()),
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

	if ce := log.Check(log.DebugLevel, "Execute (Finished)"); ce != nil {
		ce.Write(
			zap.Uint64("ReqID", in.GetRequestID()),
			zap.Duration("T", time.Duration(tools.CPUTicks()-startTime)),
		)
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

func (edge *Server) onGatewayMessage(conn gateway.Conn, streamID int64, data []byte) {
	// _, task := trace.NewTask(context.Background(), "onGatewayMessage")
	// defer task.End()

	dispatchCtx := acquireDispatchCtx(edge, conn, streamID, edge.serverID)
	err := edge.dispatcher.Interceptor(dispatchCtx, data)
	if err != nil {
		releaseDispatchCtx(dispatchCtx)
		return
	}
	err, isLeader := edge.executePrepare(dispatchCtx)
	if err != nil {
		edge.onError(dispatchCtx, rony.ErrCodeInternal, rony.ErrItemServer)
	}
	err = edge.execute(dispatchCtx, isLeader)
	if err != nil {
		edge.onError(dispatchCtx, rony.ErrCodeInternal, rony.ErrItemServer)
	}
	edge.dispatcher.Done(dispatchCtx)
	releaseDispatchCtx(dispatchCtx)
}
func (edge *Server) onError(dispatchCtx *DispatchCtx, code, item string) {
	envelope := acquireMessageEnvelope()
	rony.ErrorMessage(envelope, dispatchCtx.req.GetRequestID(), code, item)
	edge.dispatcher.OnMessage(dispatchCtx, envelope)
	releaseMessageEnvelope(envelope)
}
func (edge *Server) onConnect(conn gateway.Conn, kvs ...gateway.KeyValue) {
	edge.dispatcher.OnOpen(conn, kvs...)
}
func (edge *Server) onClose(conn gateway.Conn) {
	edge.dispatcher.OnClose(conn)
}

func (edge *Server) onClusterMessage(cm *rony.ClusterMessage) {
	// _, task := trace.NewTask(context.Background(), "onClusterMessage")
	// defer task.End()
}
func (edge *Server) onReplicaMessage(raftCmd *rony.RaftCommand) error {
	// _, task := trace.NewTask(context.Background(), "onReplicaMessage")
	// defer task.End()

	dispatchCtx := acquireDispatchCtx(edge, nil, 0, raftCmd.Sender)
	dispatchCtx.FillEnvelope(
		raftCmd.Envelope.GetRequestID(), raftCmd.Envelope.GetConstructor(), raftCmd.Envelope.Message,
		raftCmd.Envelope.Auth, raftCmd.Envelope.Header...,
	)

	err := edge.execute(dispatchCtx, false)
	if err != nil {
		return err
	}
	edge.dispatcher.Done(dispatchCtx)
	releaseDispatchCtx(dispatchCtx)
	return nil
}

// StartCluster is non-blocking function which runs the gossip and raft if it is set
func (edge *Server) StartCluster() (err error) {
	log.Info("Edge Server Started",
		zap.ByteString("ServerID", edge.serverID),
		zap.String("Gateway", string(edge.gatewayProtocol)),
		zap.Bool("Raft", edge.cluster.RaftEnabled()),
		zap.Int("GossipPort", edge.cluster.GossipPort()),
	)

	edge.cluster.Start()
	return
}

// StartGateway is non-blocking function runs the gateway in background so we can accept clients requests
func (edge *Server) StartGateway() {
	edge.gateway.Start()
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
	edge.cluster.Shutdown()

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
func (edge *Server) GetGatewayConn(connID uint64) gateway.Conn {
	return edge.gateway.GetConn(connID)
}
