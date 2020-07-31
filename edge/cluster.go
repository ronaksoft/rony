package edge

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/memberlist"
	"git.ronaksoftware.com/ronak/rony/pools"
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"
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
func (edge *Server) ClusterMembers() []*ClusterMember {
	return edge.cluster.Members()
}

// ClusterSend sends 'envelope' to the server identified by 'serverID'. It may returns ErrNotFound if the server
// is not in the list. The message will be send with BEST EFFORT and using UDP
func (edge *Server) ClusterSend(serverID []byte, authID int64, envelope *rony.MessageEnvelope) error {
	m := edge.cluster.GetByID(tools.ByteToStr(serverID))
	if m == nil {
		return rony.ErrNotFound
	}

	clusterMessage := acquireClusterMessage()
	clusterMessage.Fill(edge.serverID, authID, envelope)
	b := pools.Bytes.GetLen(clusterMessage.Size())
	_, err := clusterMessage.MarshalToSizedBuffer(b)
	if err != nil {
		return err
	}
	err = edge.gossip.SendBestEffort(m.node, b)
	pools.Bytes.Put(b)
	releaseClusterMessage(clusterMessage)
	return err
}

func (edge *Server) updateCluster(timeout time.Duration) error {
	return edge.gossip.UpdateNode(timeout)
}

// ClusterMember
type ClusterMember struct {
	ServerID    string
	ReplicaSet  uint64
	ShardSet    uint64
	GatewayAddr []string
	Addr        net.IP
	Port        uint16
	RaftPort    int
	RaftState   rony.RaftState
	node        *memberlist.Node
}

func convertMember(sm *memberlist.Node) *ClusterMember {
	edgeNode := rony.EdgeNode{}
	err := edgeNode.Unmarshal(sm.Meta)
	if err != nil {
		log.Warn("Error On ConvertMember",
			zap.Error(err),
			zap.Int("Len", len(sm.Meta)),
		)
		return nil
	}

	return &ClusterMember{
		ServerID:    tools.ByteToStr(edgeNode.ServerID),
		ReplicaSet:  edgeNode.ReplicaSet,
		GatewayAddr: edgeNode.GatewayAddr,
		RaftPort:    int(edgeNode.RaftPort),
		RaftState:   edgeNode.RaftState,
		Addr:        sm.Addr,
		Port:        sm.Port,
		node:        sm,
	}
}

// Cluster
type Cluster struct {
	sync.RWMutex
	byServerID   map[string]*ClusterMember
	byReplicaSet map[uint64][]*ClusterMember
	byShardSet   map[uint64][]*ClusterMember
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
		c.byReplicaSet = make(map[uint64][]*ClusterMember)
		c.byShardSet = make(map[uint64][]*ClusterMember)
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
		c.byReplicaSet = make(map[uint64][]*ClusterMember)
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
	edge *Server
}

func (d delegateEvents) NotifyJoin(n *memberlist.Node) {
	cm := convertMember(n)
	d.edge.cluster.AddMember(cm)
	if cm.ReplicaSet == d.edge.replicaSet {
		_ = joinRaft(d.edge, cm.ServerID, fmt.Sprintf("%s:%d", cm.Addr.String(), cm.RaftPort))
	}
}

func (d delegateEvents) NotifyLeave(n *memberlist.Node) {
	d.edge.cluster.RemoveMember(convertMember(n))
}

func (d delegateEvents) NotifyUpdate(n *memberlist.Node) {
	cm := convertMember(n)
	d.edge.cluster.AddMember(cm)
	if cm.ReplicaSet == d.edge.replicaSet {
		_ = joinRaft(d.edge, cm.ServerID, fmt.Sprintf("%s:%d", cm.Addr.String(), cm.RaftPort))
	}
}

func joinRaft(edge *Server, nodeID, addr string) error {
	if !edge.raftEnabled {
		return rony.ErrRaftNotSet
	}
	if edge.raft.State() != raft.Leader {
		return rony.ErrNotRaftLeader
	}
	futureConfig := edge.raft.GetConfiguration()
	if err := futureConfig.Error(); err != nil {
		return err
	}
	raftConf := futureConfig.Configuration()
	for _, srv := range raftConf.Servers {
		if srv.ID == raft.ServerID(nodeID) && srv.Address == raft.ServerAddress(addr) {
			return nil
		}
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			future := edge.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return err
			}
		}
	}

	future := edge.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if err := future.Error(); err != nil {
		return err
	}
	return nil
}

type delegateNode struct {
	edge *Server
}

func (d delegateNode) NodeMeta(limit int) []byte {
	n := rony.EdgeNode{
		ServerID:    d.edge.serverID,
		ReplicaSet:  d.edge.replicaSet,
		ShardSet:    d.edge.shardSet,
		RaftPort:    uint32(d.edge.raftPort),
		GatewayAddr: d.edge.gateway.Addr(),
		RaftState:   rony.RaftState_None,
	}
	if d.edge.raftEnabled {
		n.RaftState = rony.RaftState(d.edge.raft.State() + 1)
	}

	b, _ := n.Marshal()
	if len(b) > limit {
		log.Warn("Too Large Meta")
		return nil
	}
	return b
}

func (d delegateNode) NotifyMsg(data []byte) {
	if ce := log.Check(log.DebugLevel, "Cluster Message Received"); ce != nil {
		ce.Write(
			zap.ByteString("ServerID", d.edge.serverID),
			zap.Int("Data", len(data)),
		)
	}

	cm := acquireClusterMessage()
	_ = cm.Unmarshal(data)
	dispatchCtx := acquireDispatchCtx(d.edge, nil, 0, cm.AuthID, cm.Sender)
	dispatchCtx.FillEnvelope(cm.Envelope.RequestID, cm.Envelope.Constructor, cm.Envelope.Message)
	releaseClusterMessage(cm)

	d.edge.rateLimitChan <- struct{}{}
	go func() {
		// TODO:: handle error, for instance we might send back an error to the sender
		_ = d.edge.executePrepare(dispatchCtx)
		d.edge.dispatcher.Done(dispatchCtx)
		releaseDispatchCtx(dispatchCtx)
		<-d.edge.rateLimitChan
	}()
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
