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
	NoReplica     Mode = "noReplica"
	SingleReplica Mode = "singleReplica"
	MultiReplica  Mode = "multiReplica"
)

type Cluster interface {
	Raft
	Start()
	Members() []Member
	Join(addr ...string) (int, error)
	Shutdown()
	ReplicaSet() uint64
	SetGatewayAddrs(hostPorts []string) error
	SetTunnelAddrs(hostPorts []string) error
	Addr() string
}

type Raft interface {
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
