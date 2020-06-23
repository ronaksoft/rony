package testEnv

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/gateway"
	httpGateway "git.ronaksoftware.com/ronak/rony/gateway/http"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"go.uber.org/zap"
	"sync/atomic"
)

/*
   Creation Time: 2020 - Apr - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	receivedMessages int32
	receivedUpdates  int32
)

type testDispatcher struct {
}

func (t testDispatcher) OnUpdate(ctx *edge.DispatchCtx, authID int64, envelope *rony.UpdateEnvelope) {
	atomic.AddInt32(&receivedUpdates, 1)
}

func (t testDispatcher) OnMessage(ctx *edge.DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	if ctx.Conn() != nil {
		b := pools.Bytes.GetLen(envelope.Size())
		_, _ = envelope.MarshalToSizedBuffer(b)
		err := ctx.Conn().SendBinary(ctx.StreamID(), b)
		if err != nil {
			log.Warn("Error On SendBinary", zap.Error(err))
		}
		pools.Bytes.Put(b)
	}
	atomic.AddInt32(&receivedMessages, 1)
}

func (t testDispatcher) Prepare(ctx *edge.DispatchCtx, data []byte, kvs ...gateway.KeyValue) (err error) {
	proto := pb.PoolProtoMessage.Get()
	err = proto.Unmarshal(data)
	if err != nil {
		return
	}

	err = ctx.UnmarshalEnvelope(proto.Payload)
	if err != nil {
		return
	}
	ctx.SetAuthID(proto.AuthID)
	pb.PoolProtoMessage.Put(proto)
	return
}

func (t testDispatcher) Done(ctx *edge.DispatchCtx) {}

func initHandlers(edgeServer *edge.Server) {
	edgeServer.AddHandler(100, func(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
		req := pb.PoolReqSimple1.Get()
		defer pb.PoolReqSimple1.Put(req)
		res := pb.PoolResSimple1.Get()
		defer pb.PoolResSimple1.Put(res)
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(in.RequestID, "Invalid", "Proto")
			return
		}
		res.P1 = req.P1
		ctx.PushMessage(ctx.AuthID(), in.RequestID, 201, res)
	})

	edgeServer.AddHandler(101, func(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
		req := pb.PoolReqSimple1.Get()
		defer pb.PoolReqSimple1.Put(req)
		res := pb.PoolResSimple1.Get()
		defer pb.PoolResSimple1.Put(res)
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(in.RequestID, "Invalid", "Proto")
			return
		}
		res.P1 = req.P1

		ts := tools.TimeUnix()
		u := pb.PoolUpdateSimple1.Get()
		defer pb.PoolUpdateSimple1.Put(u)
		ctx.PushMessage(ctx.AuthID(), in.RequestID, 201, res)
		for i := int64(10); i < 20; i++ {
			u.P1 = tools.StrToByte(tools.Int64ToStr(i))
			ctx.PushUpdate(ctx.AuthID(), i, 301, ts, u)
		}
	})
}

func InitEdgeServerWithWebsocket(serverID string, clientPort int, opts ...edge.Option) *edge.Server {
	opts = append(opts,
		edge.WithWebsocketGateway(websocketGateway.Config{
			NewConnectionWorkers: 1,
			MaxConcurrency:       10,
			MaxIdleTime:          0,
			ListenAddress:        fmt.Sprintf(":%d", clientPort),
		}),
	)
	edgeServer := edge.NewServer(serverID, &testDispatcher{}, opts...)
	initHandlers(edgeServer)

	return edgeServer
}

func InitEdgeServerWithHttp(serverID string, clientPort int, opts ...edge.Option) *edge.Server {
	opts = append(opts,
		edge.WithHttpGateway(httpGateway.Config{
			Concurrency:   1 << 20,
			ListenAddress: fmt.Sprintf("127.0.0.1:%d", clientPort),
			MaxBodySize:   1 << 22,
		}),
	)
	edgeServer := edge.NewServer(serverID, &testDispatcher{}, opts...)
	initHandlers(edgeServer)

	return edgeServer
}

func ResetCounters() {
	atomic.StoreInt32(&receivedMessages, 0)
	atomic.StoreInt32(&receivedUpdates, 0)
}
func ReceivedMessages() int32 {
	return atomic.LoadInt32(&receivedMessages)
}

func ReceivedUpdates() int32 {
	return atomic.LoadInt32(&receivedUpdates)
}
