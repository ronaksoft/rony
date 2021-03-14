package gossipCluster

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
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

// Config
type Config struct {
	ServerID   []byte
	Bootstrap  bool
	RaftPort   int
	ReplicaSet uint64
	Mode       cluster.Mode
	GossipPort int
}

// Cluster
type Cluster struct {
	cluster.ReplicaMessageHandler
	dataPath         string
	cfg              Config
	mtx              sync.RWMutex
	localServerID    []byte
	localGatewayAddr []string
	localTunnelAddr  []string
	replicaLeaderID  string
	replicaMembers   map[uint64]map[string]*Member
	clusterMembers   map[string]*Member

	// Raft & Gossip
	raftFSM raftFSM
	raft    *raft.Raft
	gossip  *memberlist.Memberlist
	// badgerStore   *raftbadger.BadgerStore
	rateLimitChan chan struct{}
}

func New(dataPath string, cfg Config) *Cluster {
	if cfg.GossipPort == 0 {
		cfg.GossipPort = 7946
	}

	switch cfg.Mode {
	case cluster.MultiReplica, cluster.SingleReplica:
	default:
		panic("only singleReplica and multiReplica supported")
	}
	c := &Cluster{
		dataPath:       dataPath,
		cfg:            cfg,
		localServerID:  cfg.ServerID,
		clusterMembers: make(map[string]*Member, 100),
		replicaMembers: make(map[uint64]map[string]*Member, 100),
		rateLimitChan:  make(chan struct{}, clusterMessageRateLimit),
	}

	c.raftFSM = raftFSM{c: c}
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

func (c *Cluster) startRaft(notifyChan chan bool) (err error) {
	// Initialize LogStore for Raft
	dirPath := filepath.Join(c.dataPath, "raft")
	_ = os.MkdirAll(dirPath, os.ModePerm)

	// Initialize Raft
	raftConfig := raft.DefaultConfig()
	raftConfig.LogLevel = "WARN"
	raftConfig.NotifyCh = notifyChan
	raftConfig.Logger = hclog.NewNullLogger()
	raftConfig.LocalID = raft.ServerID(c.localServerID)

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	var raftAdvertiseAddr *net.TCPAddr
	var raftBind string
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok {
			if ipNet.IP == nil || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
				continue
			}
			raftBind = fmt.Sprintf("%s:%d", ipNet.IP.String(), c.cfg.RaftPort)
			raftAdvertiseAddr, err = net.ResolveTCPAddr("tcp", raftBind)
			if err != nil {
				return err
			}
			break
		}
	}

	raftTransport, err := raft.NewTCPTransport(raftBind, raftAdvertiseAddr, 3, 10*time.Second, os.Stdout)
	if err != nil {
		return err
	}

	raftSnapshot, err := raft.NewFileSnapshotStore(dirPath, 3, os.Stdout)
	if err != nil {
		return err
	}

	badgerStore := &BadgerStore{}
	c.raft, err = raft.NewRaft(raftConfig, c.raftFSM, badgerStore, badgerStore, raftSnapshot, raftTransport)
	if err != nil {
		return err
	}

	if c.cfg.Bootstrap {
		bootConfig := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raftConfig.LocalID,
					Address: raftTransport.LocalAddr(),
				},
			},
		}
		f := c.raft.BootstrapCluster(bootConfig)
		if err := f.Error(); err != nil {
			if err == raft.ErrCantBootstrap {
				log.Info("Error On Raft Bootstrap", zap.Error(err))
			} else {
				log.Warn("Error On Raft Bootstrap", zap.Error(err))
			}

		}
	}

	return nil
}

func (c *Cluster) addMember(m *Member) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// set the raft leader if the newly added member is the leader and is in our replica set
	if m.replicaSet == c.cfg.ReplicaSet && m.raftState == rony.RaftState_Leader {
		c.replicaLeaderID = m.serverID
	}

	// add new member to the cluster
	c.clusterMembers[m.serverID] = m
	if c.replicaMembers[m.replicaSet] == nil {
		c.replicaMembers[m.replicaSet] = make(map[string]*Member, 5)
	}
	c.replicaMembers[m.replicaSet][m.serverID] = m
}

func (c *Cluster) updateMember(en *rony.EdgeNode) (*Member, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	m := c.clusterMembers[tools.B2S(en.ServerID)]
	if m == nil {
		return nil, fmt.Errorf("update a node not exists: %s", tools.B2S(en.ServerID))
	}

	m.Merge(en)
	return m, nil
}

func (c *Cluster) removeMember(en *rony.EdgeNode) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	serverID := tools.B2S(en.ServerID)
	if serverID == c.replicaLeaderID {
		c.replicaLeaderID = ""
	}

	delete(c.clusterMembers, serverID)

	if c.replicaMembers[en.ReplicaSet] != nil {
		delete(c.replicaMembers[en.ReplicaSet], serverID)
	}
}

func (c *Cluster) Start() {
	err := c.startGossip()
	if err != nil {
		return
	}

	if c.cfg.Mode == cluster.MultiReplica {
		notifyChan := make(chan bool, 1)
		err := c.startRaft(notifyChan)
		if err != nil {
			log.Warn("Error On Starting Raft", zap.Error(err))
			return
		}
		go func() {
			for range notifyChan {
				err := tools.Try(10, time.Millisecond, func() error {
					return c.updateCluster(gossipUpdateTimeout)
				})
				if err != nil {
					log.Warn("Rony got error on updating the cluster",
						zap.Error(err),
						zap.ByteString("ID", c.localServerID),
					)
				}
			}
		}()
	}

}

func (c *Cluster) Join(addr ...string) (int, error) {
	return c.gossip.Join(addr)
}

func (c *Cluster) Shutdown() {
	//  shutdown raft, if it is enabled
	if c.cfg.Mode == cluster.MultiReplica {
		if f := c.raft.Snapshot(); f.Error() != nil {
			if f.Error() != raft.ErrNothingNewToSnapshot {
				log.Warn("Error On Shutdown (Raft Snapshot)",
					zap.Error(f.Error()),
					zap.ByteString("ServerID", c.localServerID),
				)
			} else {
				log.Info("Error On Shutdown (Raft Snapshot)", zap.Error(f.Error()))
			}

		}
		if f := c.raft.Shutdown(); f.Error() != nil {
			log.Warn("Error On Shutdown (Raft Shutdown)", zap.Error(f.Error()))
		}
	}

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

func (c *Cluster) Members() []cluster.Member {
	members := make([]cluster.Member, 0, 16)
	c.mtx.RLock()
	for _, cm := range c.clusterMembers {
		members = append(members, cm)
	}
	c.mtx.RUnlock()
	return members
}

func (c *Cluster) TotalReplicas() int {
	return len(c.replicaMembers)
}

func (c *Cluster) RaftEnabled() bool {
	return c.raft != nil
}

func (c *Cluster) RaftState() raft.RaftState {
	switch c.cfg.Mode {
	case cluster.SingleReplica:
		return raft.Leader
	case cluster.MultiReplica:
		return c.raft.State()
	default:
		panic("unsupported replica mode")
	}
}

func (c *Cluster) RaftApply(cmd []byte) raft.ApplyFuture {
	switch c.cfg.Mode {
	case cluster.SingleReplica:
		return nil
	case cluster.MultiReplica:
		return c.raft.Apply(cmd, raftApplyTimeout)
	default:
		panic("unsupported replica mode")
	}
}

func (c *Cluster) RaftMembers(replicaSet uint64) []cluster.Member {
	members := make([]cluster.Member, 0, 8)
	c.mtx.RLock()
	for _, cm := range c.replicaMembers[replicaSet] {
		members = append(members, cm)
	}
	c.mtx.RUnlock()
	return members
}

func (c *Cluster) RaftLeaderID() string {
	return c.replicaLeaderID
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
