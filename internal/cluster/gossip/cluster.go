package gossipCluster

import (
	"fmt"
	"github.com/ronaksoft/memberlist"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/internal/msg"
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
	ServerID       []byte
	Bootstrap      bool
	ReplicaSet     uint64
	GossipIP       string
	GossipPort     int
	AdvertisedIP   string
	AdvertisedPort int
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
	eventChan        chan memberlist.NodeEvent
	gossip           *memberlist.Memberlist
	subscriber       rony.ClusterDelegate
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
		eventChan:        make(chan memberlist.NodeEvent, 1024),
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
	conf.Events = &memberlist.ChannelEventDelegate{
		Ch: c.eventChan,
	}
	conf.Delegate = cd
	conf.LogOutput = ioutil.Discard
	conf.Logger = nil
	conf.BindAddr = c.cfg.GossipIP
	conf.BindPort = c.cfg.GossipPort
	conf.AdvertiseAddr = c.cfg.AdvertisedIP
	conf.AdvertisePort = c.cfg.AdvertisedPort
	if s, err := memberlist.Create(conf); err != nil {
		log.Warn("Error On Creating MemberList", zap.Error(err))

		return err
	} else {
		c.gossip = s
	}

	return c.updateCluster(gossipUpdateTimeout)
}

func (c *Cluster) addMember(n *memberlist.Node) {
	cm, err := newMember(n)
	if err != nil {
		log.Warn("Error On Cluster Node Add", zap.Error(err))
		return
	}

	c.mtx.Lock()
	// add new member to the cluster
	c.membersByID[cm.serverID] = cm
	c.membersByHash[cm.hash] = cm
	if c.membersByReplica[cm.replicaSet] == nil {
		c.membersByReplica[cm.replicaSet] = make(map[string]*Member, 5)
	}
	c.membersByReplica[cm.replicaSet][cm.serverID] = cm
	c.mtx.Unlock()

	if c.subscriber != nil {
		c.subscriber.OnJoin(cm.hash)
	}
}

func (c *Cluster) removeMember(n *memberlist.Node) {
	en := msg.PoolEdgeNode.Get()
	defer msg.PoolEdgeNode.Put(en)
	if err := extractNode(n, en); err != nil {
		log.Warn("Error On Cluster Node Update", zap.Error(err))
		return
	}

	serverID := tools.B2S(en.ServerID)
	c.mtx.Lock()
	delete(c.membersByID, serverID)
	delete(c.membersByHash, en.Hash)
	if c.membersByReplica[en.ReplicaSet] != nil {
		delete(c.membersByReplica[en.ReplicaSet], serverID)
	}
	c.mtx.Unlock()
	if c.subscriber != nil {
		c.subscriber.OnLeave(en.Hash)
	}
}

func (c *Cluster) Start() error {
	if err := c.startGossip(); err != nil {
		return err
	}
	go func() {
		for e := range c.eventChan {
			n := c.gossip.Member(e.Node.Name)
			switch n.State {
			case memberlist.StateLeft:
				c.removeMember(n)
			default:
				c.addMember(n)
			}
		}
	}()
	return nil
}

func (c *Cluster) Join(addr ...string) (int, error) {
	return c.gossip.Join(addr)
}

func (c *Cluster) Leave() error {
	return c.gossip.Leave(gossipLeaveTimeout)
}

func (c *Cluster) Shutdown() {
	// Shutdown gossip
	err := c.gossip.Shutdown()
	if err != nil {
		log.Warn("Error On Shutdown (Gossip)", zap.Error(err))
	}
}

func (c *Cluster) ServerID() string {
	return tools.B2S(c.localServerID)
}

func (c *Cluster) Members() []rony.ClusterMember {
	members := make([]rony.ClusterMember, 0, 16)
	c.mtx.RLock()
	for _, cm := range c.membersByID {
		members = append(members, cm)
	}
	c.mtx.RUnlock()

	return members
}

func (c *Cluster) MembersByReplicaSet(replicaSets ...uint64) []rony.ClusterMember {
	members := make([]rony.ClusterMember, 0, 16)
	c.mtx.RLock()
	for _, rs := range replicaSets {
		for _, cm := range c.membersByReplica[rs] {
			members = append(members, cm)
		}
	}
	c.mtx.RUnlock()

	return members
}

func (c *Cluster) MemberByHash(h uint64) rony.ClusterMember {
	c.mtx.RLock()
	m := c.membersByHash[h]
	c.mtx.RUnlock()
	if m == nil {
		return nil
	}

	return m
}

func (c *Cluster) MemberByID(serverID string) rony.ClusterMember {
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

func (c *Cluster) Subscribe(d rony.ClusterDelegate) {
	c.subscriber = d
}
