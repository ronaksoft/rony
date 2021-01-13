package cluster

import (
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2021 - Jan - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type ReplicaMessageHandler func(raftCmd *rony.RaftCommand) error

type Mode string

const (
	// SingleReplica if set then each replica set is only one node. i.e. raft is OFF.
	SingleReplica Mode = "singleReplica"
	// MultiReplica if set then each replica set is a raft cluster
	MultiReplica Mode = "multiReplica"
)

type Cluster interface {
	Raft
	Start()
	Members() []Member
	Join(addr ...string) (int, error)
	Shutdown()
	ReplicaSet() uint64
	TotalReplicas() int
	Addr() string
	SetGatewayAddrs(hostPorts []string) error
	SetTunnelAddrs(hostPorts []string) error
}

type Raft interface {
	RaftEnabled() bool
	RaftMembers(replicaSet uint64) []Member
	RaftState() raft.RaftState
	RaftApply(cmd []byte) raft.ApplyFuture
	RaftLeaderID() string
}

type Member interface {
	Proto(info *rony.NodeInfo) *rony.NodeInfo
	ServerID() string
	RaftState() rony.RaftState
	ReplicaSet() uint64
	GatewayAddr() []string
	TunnelAddr() []string
	RaftPort() int
}
