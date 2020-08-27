package testEnv

import (
	"fmt"
	"git.ronaksoft.com/ronak/rony"
	"git.ronaksoft.com/ronak/rony/edge"
	"git.ronaksoft.com/ronak/rony/gateway"
	tcpGateway "git.ronaksoft.com/ronak/rony/gateway/tcp"
	log "git.ronaksoft.com/ronak/rony/internal/logger"
	"git.ronaksoft.com/ronak/rony/pools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"sync/atomic"
)

/*
   Creation Time: 2020 - Apr - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	receivedMessages int32
	receivedUpdates  int32
)

type testDispatcher struct {
}

func (t testDispatcher) OnOpen(conn gateway.Conn) {

}

func (t testDispatcher) OnClose(conn gateway.Conn) {

}

func (t testDispatcher) OnMessage(ctx *edge.DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	if ctx.Conn() != nil {
		mo := proto.MarshalOptions{
			UseCachedSize: true,
		}

		b := pools.Bytes.GetCap(mo.Size(envelope))
		b, _ = mo.MarshalAppend(b, envelope)
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
		edge.WithTcpGateway(tcpGateway.Config{
			Protocol:      tcpGateway.Websocket,
			Concurrency:   10,
			MaxIdleTime:   0,
			ListenAddress: fmt.Sprintf(":%d", clientPort),
		}),
	)
	edgeServer := edge.NewServer(serverID, &testDispatcher{}, opts...)

	return edgeServer
}

func InitEdgeServerWithHttp(serverID string, clientPort int, opts ...edge.Option) *edge.Server {
	opts = append(opts,
		edge.WithTcpGateway(tcpGateway.Config{
			Protocol:      tcpGateway.Http,
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
