package rony

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/log"
	"mime/multipart"
	"net"
)

/*
   Creation Time: 2021 - Jan - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Cluster is the component which create and present the whole cluster of Edge nodes.
type Cluster interface {
	Start() error
	Shutdown()
	Join(addr ...string) (int, error)
	Leave() error
	Members() []ClusterMember
	MembersByReplicaSet(replicaSets ...uint64) []ClusterMember
	MemberByID(string) ClusterMember
	MemberByHash(uint64) ClusterMember
	ReplicaSet() uint64
	ServerID() string
	TotalReplicas() int
	Addr() string
	SetGatewayAddrs(hostPorts []string) error
	SetTunnelAddrs(hostPorts []string) error
	Subscribe(d ClusterDelegate)
}

type ClusterMember interface {
	Proto(info *Edge) *Edge
	ServerID() string
	ReplicaSet() uint64
	GatewayAddr() []string
	TunnelAddr() []string
	Dial() (net.Conn, error)
}

type ClusterDelegate interface {
	OnJoin(hash uint64)
	OnLeave(hash uint64)
}

// Gateway defines the gateway interface where clients could connect
// and communicate with the Edge servers
type Gateway interface {
	Start()
	Run()
	Shutdown()
	GetConn(connID uint64) Conn
	Addr() []string
	Protocol() GatewayProtocol
}

type GatewayProtocol int32

const (
	Undefined GatewayProtocol = 0
	Dummy     GatewayProtocol = 1 << iota
	Http
	Websocket
	Quic
	Grpc
	TCP = Http | Websocket // Http & Websocket
)

var protocolNames = map[GatewayProtocol]string{
	Undefined: "Undefined",
	Dummy:     "Dummy",
	Http:      "Http",
	Websocket: "Websocket",
	Quic:      "Quic",
	Grpc:      "Grpc",
	TCP:       "TCP",
}

func (p GatewayProtocol) String() string {
	return protocolNames[p]
}

// HTTP methods were copied from net/http.
const (
	MethodWild    = "*"
	MethodGet     = "GET"     // RFC 7231, 4.3.1
	MethodHead    = "HEAD"    // RFC 7231, 4.3.2
	MethodPost    = "POST"    // RFC 7231, 4.3.3
	MethodPut     = "PUT"     // RFC 7231, 4.3.4
	MethodPatch   = "PATCH"   // RFC 5789
	MethodDelete  = "DELETE"  // RFC 7231, 4.3.5
	MethodConnect = "CONNECT" // RFC 7231, 4.3.6
	MethodOptions = "OPTIONS" // RFC 7231, 4.3.7
	MethodTrace   = "TRACE"   // RFC 7231, 4.3.8
)

// Tunnel provides the communication channel between Edge servers. Tunnel is similar to Gateway in functionalities.
// However, Tunnel should be optimized for inter-communication between Edge servers, and Gateway is optimized for client-server communications.
type Tunnel interface {
	Start()
	Run()
	Shutdown()
	Addr() []string
}

type (
	LocalDB  = badger.DB
	StoreTxn = badger.Txn
)

// Conn defines the Connection interface
type Conn interface {
	ConnID() uint64
	ClientIP() string
	WriteBinary(streamID int64, data []byte) error
	// Persistent returns FALSE if this connection will be closed when edge.DispatchCtx has been done. i.e. HTTP connections
	// It returns TRUE if this connection still alive when edge.DispatchCtx has been done. i.e. WebSocket connections
	Persistent() bool
	Get(key string) interface{}
	Set(key string, val interface{})
}

// RestConn is same as Conn, but it supports REST format apis.
type RestConn interface {
	Conn
	WriteStatus(status int)
	WriteHeader(key, value string)
	MultiPart() (*multipart.Form, error)
	Method() string
	Path() string
	Body() []byte
	Redirect(statusCode int, newHostPort string)
}

type LogLevel = log.Level

// SetLogLevel is used for debugging purpose
func SetLogLevel(l LogLevel) {
	log.SetLevel(l)
}

func RegisterPrometheus(registerer prometheus.Registerer) {
	metrics.Register(registerer)
}

// Router could be used by Edge servers to find entities and redirect clients to the right Edge server.
type Router interface {
	Update(entityID string, replicaSet uint64) error
	Get(entityID string) (replicaSet uint64, err error)
}
