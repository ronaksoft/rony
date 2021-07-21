package rony

import (
	"net"
)

/*
   Creation Time: 2021 - Jul - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

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
