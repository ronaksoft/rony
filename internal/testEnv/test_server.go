package testEnv

import (
	"fmt"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgetest"
	"github.com/ronaksoft/rony/log"
	"github.com/ronaksoft/rony/pools"
	"go.uber.org/zap"
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
	l log.Logger
}

func (t testDispatcher) Encoder(me *rony.MessageEnvelope, buf *pools.ByteBuffer) error {
	mo := proto.MarshalOptions{UseCachedSize: true}
	bb, _ := mo.MarshalAppend(*buf.Bytes(), me)
	buf.SetBytes(&bb)

	atomic.AddInt32(&receivedMessages, 1)

	return nil
}

func (t testDispatcher) Decoder(data []byte, me *rony.MessageEnvelope) error {
	return me.Unmarshal(data)
}

func (t testDispatcher) OnOpen(conn rony.Conn, kvs ...*rony.KeyValue) {

}

func (t testDispatcher) OnClose(conn rony.Conn) {

}

func (t testDispatcher) Done(ctx *edge.DispatchCtx) {
	ctx.BufferPopAll(func(envelope *rony.MessageEnvelope) {
		buf := pools.Buffer.FromProto(envelope)
		err := ctx.Conn().WriteBinary(ctx.StreamID(), *buf.Bytes())
		if err != nil {
			t.l.Warn("Error On WriteBinary", zap.Error(err))
		}
		pools.Buffer.Put(buf)

		atomic.AddInt32(&receivedMessages, 1)
	})
}

func EdgeServer(serverID string, listenPort int, concurrency int, opts ...edge.Option) *edge.Server {
	opts = append(opts,
		edge.WithCustomDispatcher(&testDispatcher{
			l: log.DefaultLogger,
		}),
		edge.WithTcpGateway(edge.TcpGatewayConfig{
			Concurrency:   concurrency,
			ListenAddress: fmt.Sprintf(":%d", listenPort),
			MaxIdleTime:   time.Second,
			Protocol:      rony.TCP,
			ExternalAddrs: []string{fmt.Sprintf("127.0.0.1:%d", listenPort)},
		}),
	)
	edgeServer := edge.NewServer(serverID, opts...)

	return edgeServer
}

func TestServer(serverID string) *edgetest.Server {
	return edgetest.NewServer(serverID, &testDispatcher{
		l: log.DefaultLogger,
	})
}
