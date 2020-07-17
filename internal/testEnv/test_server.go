package testEnv

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/gateway"
	httpGateway "git.ronaksoftware.com/ronak/rony/gateway/http"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/pools"
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
	return ctx.UnmarshalEnvelope(data)
}

func (t testDispatcher) Done(ctx *edge.DispatchCtx) {}

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
