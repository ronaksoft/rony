package cluster

import (
	"github.com/ronaksoft/rony"
	"net"
)

/*
   Creation Time: 2021 - Jan - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Cluster interface {
	Start()
	Members() []Member
	MembersByReplicaSet(replicaSets ...uint64) []Member
	Join(addr ...string) (int, error)
	Shutdown()
	ReplicaSet() uint64
	TotalReplicas() int
	Addr() string
	SetGatewayAddrs(hostPorts []string) error
	SetTunnelAddrs(hostPorts []string) error
}

type Member interface {
	Proto(info *rony.Edge) *rony.Edge
	ServerID() string
	ReplicaSet() uint64
	GatewayAddr() []string
	TunnelAddr() []string
	TunnelConn() (net.Conn, error)
}
