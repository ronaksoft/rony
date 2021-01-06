package testEnv

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgetest"
	"github.com/ronaksoft/rony/gateway"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/pools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"sync/atomic"
	"time"
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

func (t testDispatcher) OnOpen(conn gateway.Conn, kvs ...gateway.KeyValue) {

}

func (t testDispatcher) OnClose(conn gateway.Conn) {

}

func (t testDispatcher) OnMessage(ctx *edge.DispatchCtx, envelope *rony.MessageEnvelope) {
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

func (t testDispatcher) Interceptor(ctx *edge.DispatchCtx, data []byte) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (t testDispatcher) Done(ctx *edge.DispatchCtx) {}

func InitEdgeServerWithWebsocket(serverID string, listenPort int, concurrency int, opts ...edge.Option) *edge.Server {
	opts = append(opts,
		edge.WithTcpGateway(edge.TcpGatewayConfig{
			Concurrency:   concurrency,
			MaxIdleTime:   time.Second,
			ListenAddress: fmt.Sprintf(":%d", listenPort),
			ExternalAddrs: []string{fmt.Sprintf("127.0.0.1:%d", listenPort)},
		}),
	)
	edgeServer := edge.NewServer(serverID, &testDispatcher{}, opts...)

	return edgeServer
}

func InitTestServer(serverID string) *edgetest.Server {
	return edgetest.NewServer(serverID, &testDispatcher{})
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
