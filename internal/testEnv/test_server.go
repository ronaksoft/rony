package testEnv

import (
	"fmt"
	"time"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgetest"
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

func EdgeServer(serverID string, listenPort int, concurrency int, opts ...edge.Option) *edge.Server {
	opts = append(opts,
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
	return edgetest.NewServer(
		serverID, &edge.DefaultDispatcher{},
	)
}

func TestJSONServer(serverID string) *edgetest.Server {
	return edgetest.NewServer(
		serverID, &edge.JSONDispatcher{},
	)
}
