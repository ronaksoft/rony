package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/errors"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/memberlist"
	"git.ronaksoftware.com/ronak/rony/msg"
	"github.com/gobwas/pool/pbytes"
	"go.uber.org/zap"
	"net"
	"sync"
)

/*
   Creation Time: 2020 - Feb - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// ClusterMembers returns a list of all the discovered nodes in the cluster
func (edge *EdgeServer) ClusterMembers() []*ClusterMember {
	return edge.cluster.Members()
}

// ClusterSend sends 'envelope' to the server identified by 'serverID'. It may returns ErrNotFound if the server
// is not in the list. The message will be send with BEST EFFORT and using UDP
func (edge *EdgeServer) ClusterSend(serverID string, envelope *msg.MessageEnvelope) error {
	m := edge.cluster.GetByID(serverID)
	if m == nil {
		return errors.ErrNotFound
	}
	b := pbytes.GetLen(envelope.Size())
	_, err := envelope.MarshalTo(b)
	if err != nil {
		return err
	}
	err = edge.gossip.SendBestEffort(m.node, b)
	pbytes.Put(b)
	return err
}

func (edge *EdgeServer) updateCluster() error {
	return edge.gossip.UpdateNode(0)
}

// ClusterMember
type ClusterMember struct {
	ServerID    string
	ReplicaSet  uint32
	ShardSet    uint32
	ShardMin    uint32
	ShardMax    uint32
	GatewayAddr string
	Addr        net.IP
	Port        uint16
	RaftPort    int
	RaftState   RaftState
	node        *memberlist.Node
}

func convertMember(sm *memberlist.Node) *ClusterMember {
	edgeNode := EdgeNode{}
	err := edgeNode.Unmarshal(sm.Meta)
	if err != nil {
		log.Warn("Error On ConvertMember",
			zap.Error(err),
			zap.Int("Len", len(sm.Meta)),
		)
		return nil
	}

	return &ClusterMember{
		ServerID:    edgeNode.ServerID,
		ReplicaSet:  edgeNode.ReplicaSet,
		ShardMin:    edgeNode.ShardMin,
		ShardMax:    edgeNode.ShardMax,
		GatewayAddr: edgeNode.GatewayAddr,
		RaftPort:    int(edgeNode.RaftPort),
		RaftState:   edgeNode.RaftState,
		Addr:        sm.Addr,
		Port:        sm.Port,
		node:        sm,
	}
}

type Cluster struct {
	sync.RWMutex
	byServerID   map[string]*ClusterMember
	byReplicaSet map[uint32][]*ClusterMember
	byShardSet   map[uint32][]*ClusterMember
}

func (c *Cluster) GetByID(id string) *ClusterMember {
	c.RLock()
	defer c.RUnlock()

	return c.byServerID[id]
}

func (c *Cluster) AddMember(m *ClusterMember) {
	c.Lock()
	defer c.Unlock()

	if m == nil {
		return
	}
	if len(m.ServerID) == 0 {
		return
	}

	if c.byReplicaSet == nil {
		c.byReplicaSet = make(map[uint32][]*ClusterMember)
		c.byServerID = make(map[string]*ClusterMember)
		c.byServerID[m.ServerID] = m
		c.byReplicaSet[m.ReplicaSet] = append(c.byReplicaSet[m.ReplicaSet], m)
		return
	}

	c.byServerID[m.ServerID] = m
	for idx := range c.byReplicaSet[m.ReplicaSet] {
		if c.byReplicaSet[m.ReplicaSet][idx].ServerID == m.ServerID {
			c.byReplicaSet[m.ReplicaSet][idx] = m
			return
		}
	}
	c.byReplicaSet[m.ReplicaSet] = append(c.byReplicaSet[m.ReplicaSet], m)

	for idx := range c.byShardSet[m.ShardSet] {
		if c.byShardSet[m.ShardSet][idx].ServerID == m.ServerID {
			c.byShardSet[m.ShardSet][idx] = m
			return
		}
	}
	c.byShardSet[m.ShardSet] = append(c.byShardSet[m.ShardSet], m)
}

func (c *Cluster) RemoveMember(m *ClusterMember) {
	c.Lock()
	defer c.Unlock()

	if m == nil {
		return
	}
	if len(m.ServerID) == 0 {
		return
	}

	if c.byReplicaSet == nil {
		c.byReplicaSet = make(map[uint32][]*ClusterMember)
		c.byServerID = make(map[string]*ClusterMember)
		return
	}
	delete(c.byServerID, m.ServerID)
	for idx := range c.byReplicaSet[m.ReplicaSet] {
		if c.byReplicaSet[m.ReplicaSet][idx].ServerID == m.ServerID {
			c.byReplicaSet[m.ReplicaSet][idx] = c.byReplicaSet[m.ReplicaSet][len(c.byReplicaSet[m.ReplicaSet])-1]
			c.byReplicaSet[m.ReplicaSet] = c.byReplicaSet[m.ReplicaSet][:len(c.byReplicaSet[m.ReplicaSet])-1]
			return
		}
	}
}

func (c *Cluster) Members() []*ClusterMember {
	members := make([]*ClusterMember, 0, 10)
	c.RLock()
	for _, cm := range c.byServerID {
		members = append(members, cm)
	}
	c.RUnlock()
	return members
}

type delegateEvents struct {
	edge *EdgeServer
}

func (d delegateEvents) NotifyJoin(n *memberlist.Node) {
	cm := convertMember(n)
	d.edge.cluster.AddMember(cm)
	if cm.ReplicaSet == d.edge.replicaSet {
		_ = d.edge.joinRaft(cm.ServerID, fmt.Sprintf("%s:%d", cm.Addr.String(), cm.RaftPort))
	}
}

func (d delegateEvents) NotifyLeave(n *memberlist.Node) {
	d.edge.cluster.RemoveMember(convertMember(n))
}

func (d delegateEvents) NotifyUpdate(n *memberlist.Node) {
	cm := convertMember(n)
	d.edge.cluster.AddMember(cm)
	_ = d.edge.joinRaft(cm.ServerID, fmt.Sprintf("%s:%d", cm.Addr.String(), cm.RaftPort))
}

type delegateNode struct {
	edge *EdgeServer
}

func (d delegateNode) NodeMeta(limit int) []byte {
	n := EdgeNode{
		ServerID:    d.edge.serverID,
		ReplicaSet:  d.edge.replicaSet,
		ShardSet:    d.edge.shardSet,
		ShardMax:    d.edge.shardMax,
		ShardMin:    d.edge.shardMin,
		RaftPort:    uint32(d.edge.raftPort),
		GatewayAddr: d.edge.gateway.Addr(),
		RaftState:   RaftState_None,
	}
	if d.edge.raftEnabled {
		n.RaftState = RaftState(d.edge.raft.State() + 1)
	}

	b, _ := n.Marshal()
	if len(b) > limit {
		log.Warn("Too Large Meta")
		return nil
	}
	return b
}

func (d delegateNode) NotifyMsg(data []byte) {
	log.Info("Message Received",
		zap.String("ServerID", d.edge.GetServerID()),
		zap.Int("Data", len(data)),
	)
}

func (d delegateNode) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d delegateNode) LocalState(join bool) []byte {
	return nil
}

func (d delegateNode) MergeRemoteState(buf []byte, join bool) {
	return
}
