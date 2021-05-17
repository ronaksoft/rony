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
	Shutdown()
	Join(addr ...string) (int, error)
	Members() []Member
	MembersByReplicaSet(replicaSets ...uint64) []Member
	MemberByID(string) Member
	MemberByHash(uint64) Member
	ReplicaSet() uint64
	ServerID() []byte
	TotalReplicas() int
	Addr() string
	SetGatewayAddrs(hostPorts []string) error
	SetTunnelAddrs(hostPorts []string) error
	Subscribe(d Delegate)
}

type Member interface {
	Proto(info *rony.Edge) *rony.Edge
	ServerID() string
	ReplicaSet() uint64
	GatewayAddr() []string
	TunnelAddr() []string
	Dial() (net.Conn, error)
}

type Delegate interface {
	OnJoin(hash uint64)
	OnLeave(hash uint64)
}
