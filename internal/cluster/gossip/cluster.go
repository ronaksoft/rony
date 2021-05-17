package gossipCluster

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"hash/crc64"
	"io/ioutil"
	"sync"
	"time"
)

/*
   Creation Time: 2021 - Jan - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/
var (
	crcTable = crc64.MakeTable(crc64.ISO)
)

type Config struct {
	ServerID   []byte
	Bootstrap  bool
	ReplicaSet uint64
	GossipPort int
}

type Cluster struct {
	dataPath         string
	cfg              Config
	mtx              sync.RWMutex
	localServerID    []byte
	localGatewayAddr []string
	localTunnelAddr  []string
	membersByReplica map[uint64]map[string]*Member
	membersByID      map[string]*Member
	membersByHash    map[uint64]*Member
	gossip           *memberlist.Memberlist
	rateLimitChan    chan struct{}
	subscriber       cluster.Delegate
}

func New(dataPath string, cfg Config) *Cluster {
	if cfg.GossipPort == 0 {
		cfg.GossipPort = 7946
	}

	c := &Cluster{
		dataPath:         dataPath,
		cfg:              cfg,
		localServerID:    cfg.ServerID,
		membersByID:      make(map[string]*Member, 4096),
		membersByReplica: make(map[uint64]map[string]*Member, 1024),
		membersByHash:    make(map[uint64]*Member, 4096),
		rateLimitChan:    make(chan struct{}, clusterMessageRateLimit),
	}

	return c
}

func (c *Cluster) updateCluster(timeout time.Duration) error {
	return c.gossip.UpdateNode(timeout)
}

func (c *Cluster) startGossip() error {
	cd := &clusterDelegate{
		c: c,
	}
	conf := memberlist.DefaultLANConfig()
	conf.Name = string(c.localServerID)
	conf.Events = cd
	conf.Delegate = cd
	conf.Alive = cd
	conf.LogOutput = ioutil.Discard
	conf.Logger = nil
	conf.BindPort = c.cfg.GossipPort
	if s, err := memberlist.Create(conf); err != nil {
		log.Warn("Error On Creating MemberList", zap.Error(err))
		return err
	} else {
		c.gossip = s
	}

	return c.updateCluster(gossipUpdateTimeout)
}

func (c *Cluster) addMember(m *Member) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// add new member to the cluster
	c.membersByID[m.serverID] = m
	c.membersByHash[m.hash] = m
	if c.membersByReplica[m.replicaSet] == nil {
		c.membersByReplica[m.replicaSet] = make(map[string]*Member, 5)
	}
	c.membersByReplica[m.replicaSet][m.serverID] = m
}

func (c *Cluster) removeMember(en *rony.EdgeNode) {
	serverID := tools.B2S(en.ServerID)
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.membersByID, serverID)
	delete(c.membersByHash, en.Hash)
	if c.membersByReplica[en.ReplicaSet] != nil {
		delete(c.membersByReplica[en.ReplicaSet], serverID)
	}
}

func (c *Cluster) Start() {
	err := c.startGossip()
	if err != nil {
		return
	}
}

func (c *Cluster) Join(addr ...string) (int, error) {
	return c.gossip.Join(addr)
}

func (c *Cluster) Shutdown() {
	// Shutdown gossip
	err := c.gossip.Leave(gossipLeaveTimeout)
	if err != nil {
		log.Warn("Error On Leaving Cluster, but we shutdown anyway", zap.Error(err))
	}
	err = c.gossip.Shutdown()
	if err != nil {
		log.Warn("Error On Shutdown (Gossip)", zap.Error(err))
	}

}

func (c *Cluster) ServerID() []byte {
	return c.localServerID
}

func (c *Cluster) Members() []cluster.Member {
	members := make([]cluster.Member, 0, 16)
	c.mtx.RLock()
	for _, cm := range c.membersByID {
		members = append(members, cm)
	}
	c.mtx.RUnlock()
	return members
}

func (c *Cluster) MembersByReplicaSet(replicaSets ...uint64) []cluster.Member {
	members := make([]cluster.Member, 0, 16)
	c.mtx.RLock()
	for _, rs := range replicaSets {
		for _, cm := range c.membersByReplica[rs] {
			members = append(members, cm)
		}
	}
	c.mtx.RUnlock()
	return members
}

func (c *Cluster) MemberByHash(h uint64) cluster.Member {
	c.mtx.RLock()
	m := c.membersByHash[h]
	c.mtx.RUnlock()
	if m == nil {
		return nil
	}
	return m
}

func (c *Cluster) MemberByID(serverID string) cluster.Member {
	c.mtx.RLock()
	m := c.membersByID[serverID]
	c.mtx.RUnlock()
	if m == nil {
		return nil
	}
	return m
}

func (c *Cluster) TotalReplicas() int {
	return len(c.membersByReplica)
}

func (c *Cluster) SetGatewayAddrs(hostPorts []string) error {
	c.localGatewayAddr = hostPorts
	return c.updateCluster(gossipUpdateTimeout)
}

func (c *Cluster) SetTunnelAddrs(hostPorts []string) error {
	c.localTunnelAddr = hostPorts
	return c.updateCluster(gossipUpdateTimeout)
}

func (c *Cluster) Addr() string {
	return fmt.Sprintf("%s:%d", c.gossip.LocalNode().Addr.String(), c.gossip.LocalNode().Port)
}

func (c *Cluster) ReplicaSet() uint64 {
	return c.cfg.ReplicaSet
}

func (c *Cluster) Subscribe(d cluster.Delegate) {
	c.subscriber = d
}
