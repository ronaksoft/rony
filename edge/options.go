package edge

import (
	"runtime"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	gossipCluster "github.com/ronaksoft/rony/internal/cluster/gossip"
	dummyGateway "github.com/ronaksoft/rony/internal/gateway/dummy"
	tcpGateway "github.com/ronaksoft/rony/internal/gateway/tcp"
	scyllaRouter "github.com/ronaksoft/rony/internal/router/scylla"
	sqlRouter "github.com/ronaksoft/rony/internal/router/sql"
	udpTunnel "github.com/ronaksoft/rony/internal/tunnel/udp"
	"github.com/ronaksoft/rony/log"
	"github.com/scylladb/gocqlx/v2"
	"gorm.io/gorm"
)

/*
   Creation Time: 2020 - Feb - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Option func(edge *Server)

func WithTracer(tracer trace.Tracer) Option {
	return func(edge *Server) {
		if tracer == nil {
			return
		}
		edge.tracer = tracer
	}
}

func WithJSONDispatcher() Option {
	return func(edge *Server) {
		edge.dispatcher = &JSONDispatcher{}
	}
}

// WithCustomDispatcher enables custom dispatcher to write your specific event handlers.
func WithCustomDispatcher(d Dispatcher) Option {
	return func(edge *Server) {
		if d != nil {
			edge.dispatcher = d
		}
	}
}

type GossipClusterConfig struct {
	Bootstrap      bool
	ReplicaSet     uint64
	GossipIP       string
	GossipPort     int
	AdvertisedIP   string
	AdvertisedPort int
}

// WithGossipCluster enables the cluster in gossip mode. This mod is eventually consistent mode but there is
// no need to a central key-value store or any other 3rd party service to run the cluster
func WithGossipCluster(clusterConfig GossipClusterConfig) Option {
	return func(edge *Server) {
		edge.cluster = gossipCluster.New(
			gossipCluster.Config{
				ServerID:       edge.serverID,
				Bootstrap:      clusterConfig.Bootstrap,
				ReplicaSet:     clusterConfig.ReplicaSet,
				GossipPort:     clusterConfig.GossipPort,
				GossipIP:       clusterConfig.GossipIP,
				AdvertisedIP:   clusterConfig.AdvertisedIP,
				AdvertisedPort: clusterConfig.AdvertisedPort,
				Logger:         log.DefaultLogger.With("GossipCluster"),
			},
		)
	}
}

func WithCustomerCluster(c rony.Cluster) Option {
	return func(edge *Server) {
		edge.cluster = c
	}
}

type gatewayDelegate struct {
	edge *Server
}

func (g gatewayDelegate) OnConnect(c rony.Conn, kvs ...*rony.KeyValue) {
	g.edge.onGatewayConnect(c, kvs...)
}

func (g gatewayDelegate) OnMessage(c rony.Conn, streamID int64, data []byte) {
	g.edge.onGatewayMessage(c, streamID, data)
}

func (g gatewayDelegate) OnClose(c rony.Conn) {
	g.edge.onGatewayClose(c)
}

type TcpGatewayConfig struct {
	Concurrency   int
	ListenAddress string
	MaxBodySize   int
	MaxIdleTime   time.Duration
	Protocol      rony.GatewayProtocol
	ExternalAddrs []string
}

// WithTcpGateway set the gateway to tcp which can support http and/or websocket
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithTcpGateway(gatewayConfig TcpGatewayConfig) Option {
	return func(edge *Server) {
		if edge.gateway != nil {
			panic(errors.ErrGatewayAlreadyInitialized)
		}
		if gatewayConfig.Protocol == rony.Undefined {
			gatewayConfig.Protocol = rony.TCP
		}
		if gatewayConfig.Concurrency == 0 {
			gatewayConfig.Concurrency = runtime.NumCPU() * 100
		}
		gatewayTcp, err := tcpGateway.New(tcpGateway.Config{
			Concurrency:   gatewayConfig.Concurrency,
			ListenAddress: gatewayConfig.ListenAddress,
			MaxBodySize:   gatewayConfig.MaxBodySize,
			MaxIdleTime:   gatewayConfig.MaxIdleTime,
			Protocol:      gatewayConfig.Protocol,
			ExternalAddrs: gatewayConfig.ExternalAddrs,
			Logger:        log.DefaultLogger.With("TCPGateway"),
		})
		if err != nil {
			panic(err)
		}

		gatewayTcp.Subscribe(&gatewayDelegate{edge: edge})
		edge.gateway = gatewayTcp
	}
}

type DummyGatewayConfig = dummyGateway.Config

// WithTestGateway set the gateway to a dummy gateway which is useful for writing tests.
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithTestGateway(gatewayConfig DummyGatewayConfig) Option {
	return func(edge *Server) {
		if edge.gateway != nil {
			panic(errors.ErrGatewayAlreadyInitialized)
		}
		gatewayDummy, err := dummyGateway.New(gatewayConfig)
		if err != nil {
			panic(err)
		}
		gatewayDummy.Subscribe(&gatewayDelegate{edge: edge})
		edge.gateway = gatewayDummy
	}
}

func WithCustomGateway(gateway rony.Gateway) Option {
	return func(edge *Server) {
		if edge.gateway != nil {
			panic(errors.ErrGatewayAlreadyInitialized)
		}

		gateway.Subscribe(&gatewayDelegate{edge: edge})
		edge.gateway = gateway
	}
}

type UdpTunnelConfig struct {
	ListenAddress string
	MaxBodySize   int
	ExternalAddrs []string
}

// WithUdpTunnel set the tunnel to a udp based tunnel which provides communication channel between
// edge servers.
func WithUdpTunnel(config UdpTunnelConfig) Option {
	return func(edge *Server) {
		tunnelUDP, err := udpTunnel.New(udpTunnel.Config{
			ServerID:      edge.GetServerID(),
			ListenAddress: config.ListenAddress,
			MaxBodySize:   config.MaxBodySize,
			ExternalAddrs: config.ExternalAddrs,
			Logger:        log.DefaultLogger.With("Tunnel"),
		})
		if err != nil {
			panic(err)
		}
		tunnelUDP.MessageHandler = edge.onTunnelMessage
		edge.tunnel = tunnelUDP
	}
}

func WithCustomRouter(r rony.Router) Option {
	return func(edge *Server) {
		edge.router = r
	}
}

type ScyllaRouterConfig struct {
	DbSession gocqlx.Session
}

func WithScyllaRouter(config ScyllaRouterConfig) Option {
	return func(edge *Server) {
		edge.router = scyllaRouter.New(config.DbSession)
	}
}

type SqlRouterConfig struct {
	DB *gorm.DB
}

func WithSqlRouter(config SqlRouterConfig) Option {
	return func(edge *Server) {
		edge.router = sqlRouter.New(config.DB)
	}
}
