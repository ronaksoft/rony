package storeCluster

import (
	raftbadger "github.com/bbva/raft-badger"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony/cluster"
	"sync"
)

/*
   Creation Time: 2021 - Jan - 14
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Config
type Config struct {
	ServerID   []byte
	Bootstrap  bool
	RaftPort   int
	ReplicaSet uint64
	Mode       cluster.Mode
	GossipPort int
	DataPath   string
}

type Cluster struct {
	cluster.ReplicaMessageHandler
	cfg              Config
	mtx              sync.RWMutex
	localServerID    []byte
	localGatewayAddr []string
	localTunnelAddr  []string
	localShardRange  [2]uint32
	replicaLeaderID  string
	replicaMembers   map[uint64]map[string]*Member
	clusterMembers   map[string]*Member

	// Raft & Gossip
	raftFSM       raftFSM
	raft          *raft.Raft
	gossip        *memberlist.Memberlist
	badgerStore   *raftbadger.BadgerStore
	rateLimitChan chan struct{}
}

func (c *Cluster) RaftEnabled() bool {
	panic("implement me")
}

func (c *Cluster) RaftMembers(replicaSet uint64) []cluster.Member {
	panic("implement me")
}

func (c *Cluster) RaftState() raft.RaftState {
	panic("implement me")
}

func (c *Cluster) RaftApply(cmd []byte) raft.ApplyFuture {
	panic("implement me")
}

func (c *Cluster) RaftLeaderID() string {
	panic("implement me")
}

func (c *Cluster) Start() {
	panic("implement me")
}

func (c *Cluster) Members() []cluster.Member {
	panic("implement me")
}

func (c *Cluster) Join(addr ...string) (int, error) {
	panic("implement me")
}

func (c *Cluster) Shutdown() {
	panic("implement me")
}

func (c *Cluster) ReplicaSet() uint64 {
	panic("implement me")
}

func (c *Cluster) TotalReplicas() int {
	panic("implement me")
}

func (c *Cluster) Addr() string {
	panic("implement me")
}

func (c *Cluster) SetGatewayAddrs(hostPorts []string) error {
	panic("implement me")
}

func (c *Cluster) SetTunnelAddrs(hostPorts []string) error {
	panic("implement me")
}
