package cluster

import (
	"fmt"
	raftbadger "github.com/bbva/raft-badger"
	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

type MessageHandler func(cm *rony.ClusterMessage)
type ReplicaMessageHandler func(raftCmd *rony.RaftCommand) error

// Member
type Member struct {
	ServerID    string
	ReplicaSet  uint64
	ShardRange  [2]uint32
	GatewayAddr []string
	ClusterAddr net.IP
	ClusterPort uint16
	RaftPort    int
	RaftState   rony.RaftState
	node        *memberlist.Node
}

func (m *Member) Proto(p *rony.NodeInfo) *rony.NodeInfo {
	if p == nil {
		p = &rony.NodeInfo{}
	}
	p.ReplicaSet = m.ReplicaSet
	p.ServerID = m.ServerID
	p.HostPorts = append(p.HostPorts, m.GatewayAddr...)
	p.Leader = m.RaftState == rony.RaftState_Leader
	return p
}

func convertMember(sm *memberlist.Node) *Member {
	edgeNode := &rony.EdgeNode{}
	err := proto.UnmarshalOptions{}.Unmarshal(sm.Meta, edgeNode)
	if err != nil {
		log.Warn("Error On ConvertMember",
			zap.Error(err),
			zap.Int("Len", len(sm.Meta)),
		)
		return nil
	}

	return &Member{
		ServerID:    tools.ByteToStr(edgeNode.ServerID),
		ReplicaSet:  edgeNode.GetReplicaSet(),
		ShardRange:  [2]uint32{edgeNode.ShardRangeMin, edgeNode.ShardRangeMax},
		GatewayAddr: edgeNode.GatewayAddr,
		RaftPort:    int(edgeNode.GetRaftPort()),
		RaftState:   edgeNode.GetRaftState(),
		ClusterAddr: sm.Addr,
		ClusterPort: sm.Port,
		node:        sm,
	}
}

// Cluster
type Cluster struct {
	mtx              sync.RWMutex
	mode             Mode
	dataPath         string
	localServerID    []byte
	localGatewayAddr []string
	localShardRange  [2]uint32
	replicaLeaderID  string
	replicaSet       uint64
	replicaMembers   map[uint64]map[string]*Member
	clusterMembers   map[string]*Member

	// External Handlers
	onReplicaMessage ReplicaMessageHandler

	// Raft
	raftEnabled   bool
	raftPort      int
	raftBootstrap bool
	raftFSM       raftFSM
	raft          *raft.Raft
	// Gossip
	gossipPort    int
	gossip        *memberlist.Memberlist
	badgerStore   *raftbadger.BadgerStore
	rateLimitChan chan struct{}
}

func New(serverID []byte, replicaHandler ReplicaMessageHandler) *Cluster {
	return &Cluster{
		mode:             NoReplica,
		onReplicaMessage: replicaHandler,
		localServerID:    serverID,
		clusterMembers:   make(map[string]*Member, 100),
		replicaMembers:   make(map[uint64]map[string]*Member, 100),
		dataPath:         ".",
		raftEnabled:      false,
		raftPort:         0,
		raftBootstrap:    false,
		rateLimitChan:    make(chan struct{}, clusterMessageRateLimit),
		gossipPort:       7946,
	}
}

func (c *Cluster) updateCluster(timeout time.Duration) error {
	return c.gossip.UpdateNode(timeout)
}

func (c *Cluster) startGossip() error {
	dirPath := filepath.Join(c.dataPath, "gossip")
	_ = os.MkdirAll(dirPath, os.ModePerm)

	cd := &clusterDelegate{
		c: c,
	}
	conf := memberlist.DefaultLANConfig()
	conf.Name = string(c.localServerID)
	conf.Events = cd
	conf.Delegate = cd
	conf.LogOutput = ioutil.Discard
	conf.Logger = nil
	conf.BindPort = c.gossipPort
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
	badgerOpt := badger.DefaultOptions(dirPath).WithLogger(nil)
	c.badgerStore, err = raftbadger.New(raftbadger.Options{
		Path:                dirPath,
		BadgerOptions:       &badgerOpt,
		NoSync:              false,
		ValueLogGC:          false,
		GCInterval:          0,
		MandatoryGCInterval: 0,
		GCThreshold:         0,
	})
	if err != nil {
		return
	}

	// Initialize Raft
	raftConfig := raft.DefaultConfig()
	// raftConfig.LogLevel = "DEBUG"
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
			raftBind = fmt.Sprintf("%s:%d", ipNet.IP.String(), c.raftPort)
			raftAdvertiseAddr, err = net.ResolveTCPAddr("tcp", raftBind)
			if err != nil {
				return err
			}
			break
		}
	}

	log.Info("Raft",
		zap.String("Bind", raftBind),
		zap.String("Advertise", raftAdvertiseAddr.String()),
	)

	raftTransport, err := raft.NewTCPTransport(raftBind, raftAdvertiseAddr, 3, 10*time.Second, os.Stdout)
	if err != nil {
		return err
	}

	raftSnapshot, err := raft.NewFileSnapshotStore(dirPath, 3, os.Stdout)
	if err != nil {
		return err
	}

	c.raft, err = raft.NewRaft(raftConfig, c.raftFSM, c.badgerStore, c.badgerStore, raftSnapshot, raftTransport)
	if err != nil {
		return err
	}

	if c.raftBootstrap {
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

func (c *Cluster) Start() {
	notifyChan := make(chan bool, 1)
	if c.raftEnabled {
		err := c.startRaft(notifyChan)
		if err != nil {
			log.Warn("Error On Starting Raft", zap.Error(err))
			return
		}
	}

	err := c.startGossip()
	if err != nil {
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

func (c *Cluster) Join(addr ...string) (int, error) {
	return c.gossip.Join(addr)
}

func (c *Cluster) Shutdown() {
	// Second shutdown raft, if it is enabled
	if c.raftEnabled {
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

		err := c.badgerStore.Close()
		if err != nil {
			log.Warn("Error On Shutdown (Close Raft Badger)", zap.Error(err))
		}
	}

	// Shutdown gossip
	if c.gossip != nil {
		err := c.gossip.Leave(gossipLeaveTimeout)
		if err != nil {
			log.Warn("Error On Leaving Cluster, but we shutdown anyway", zap.Error(err))
		}
		err = c.gossip.Shutdown()
		if err != nil {
			log.Warn("Error On Shutdown (Gossip)", zap.Error(err))
		}
	}

}

func (c *Cluster) Send(serverID string, envelope *rony.MessageEnvelope, kvs ...*rony.KeyValue) (err error) {
	m := c.GetByID(serverID)
	if m == nil {
		return rony.ErrNotFound
	}

	cm := acquireClusterMessage()
	cm.Fill(c.localServerID, envelope, kvs...)

	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Bytes.GetCap(mo.Size(cm))
	b, err = mo.MarshalAppend(b, cm)
	if err != nil {
		return err
	}
	err = c.gossip.SendBestEffort(m.node, b)
	pools.Bytes.Put(b)
	releaseClusterMessage(cm)
	return err
}

func (c *Cluster) SetRaft(replicaSet uint64, bindPort int, bootstrap bool, mod Mode) {
	switch mod {
	case MultiReplica, SingleReplica:
		c.mode = mod
	default:
		panic("only singleReplica and multiReplica supported")
	}
	c.raftFSM = raftFSM{c: c}
	c.replicaSet = replicaSet
	c.raftEnabled = true
	c.raftPort = bindPort
	c.raftBootstrap = bootstrap

}

func (c *Cluster) SetGossipPort(port int) {
	c.gossipPort = port
}

func (c *Cluster) SetDataPath(path string) {
	c.dataPath = path
}

func (c *Cluster) SetGatewayAddrs(addrs []string) error {
	c.localGatewayAddr = addrs
	return c.updateCluster(gossipUpdateTimeout)
}

func (c *Cluster) RaftEnabled() bool {
	return c.raftEnabled
}

func (c *Cluster) RaftState() raft.RaftState {
	return c.raft.State()
}

func (c *Cluster) RaftApply(cmd []byte) raft.ApplyFuture {
	return c.raft.Apply(cmd, raftApplyTimeout)
}

func (c *Cluster) RaftConfigs() raft.ConfigurationFuture {
	return c.raft.GetConfiguration()
}

func (c *Cluster) GossipPort() int {
	return c.gossipPort
}

func (c *Cluster) GetByID(id string) *Member {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.clusterMembers[id]
}

func (c *Cluster) LeaderID() string {
	return c.replicaLeaderID
}

func (c *Cluster) AddMember(m *Member) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if m == nil {
		return
	}
	if len(m.ServerID) == 0 {
		return
	}

	if m.ReplicaSet == c.replicaSet && m.RaftState == rony.RaftState_Leader {
		c.replicaLeaderID = m.ServerID
	}
	c.clusterMembers[m.ServerID] = m
	if c.replicaMembers[m.ReplicaSet] == nil {
		c.replicaMembers[m.ReplicaSet] = make(map[string]*Member, 100)
	}
	c.replicaMembers[m.ReplicaSet][m.ServerID] = m
}

func (c *Cluster) RemoveMember(m *Member) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if m == nil {
		return
	}

	if len(m.ServerID) == 0 {
		return
	}

	if m.ServerID == c.replicaLeaderID {
		c.replicaLeaderID = ""
	}

	delete(c.clusterMembers, m.ServerID)

	if c.replicaMembers[m.ReplicaSet] != nil {
		delete(c.replicaMembers[m.ReplicaSet], m.ServerID)
	}
}

func (c *Cluster) Members() []*Member {
	members := make([]*Member, 0, 16)
	c.mtx.RLock()
	for _, cm := range c.clusterMembers {
		members = append(members, cm)
	}
	c.mtx.RUnlock()
	return members
}

func (c *Cluster) ReplicaMembers(replicaSet uint64) []*Member {
	members := make([]*Member, 0, 8)
	c.mtx.RLock()
	for _, cm := range c.replicaMembers[replicaSet] {
		members = append(members, cm)
	}
	c.mtx.RUnlock()
	return members
}

func (c *Cluster) Addr() string {
	return fmt.Sprintf("%s:%d", c.gossip.LocalNode().Addr.String(), c.gossip.LocalNode().Port)
}

func (c *Cluster) ReplicaSet() uint64 {
	return c.replicaSet
}
