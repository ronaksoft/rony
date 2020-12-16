package edge

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
   Copyright Ronak Software Group 2020
*/

// ClusterMembers returns a list of all the discovered nodes in the cluster
func (edge *Server) ClusterMembers() []*ClusterMember {
	return edge.cluster.Members()
}

// ClusterSend sends 'envelope' to the server identified by 'serverID'. It may returns ErrNotFound if the server
// is not in the list. The message will be send with BEST EFFORT and using UDP
func (edge *Server) ClusterSend(serverID string, envelope *rony.MessageEnvelope, kvs ...*rony.KeyValue) (err error) {
	m := edge.cluster.GetByID(serverID)
	if m == nil {
		return rony.ErrNotFound
	}

	clusterMessage := acquireClusterMessage()
	clusterMessage.Fill(edge.serverID, envelope, kvs...)

	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Bytes.GetCap(mo.Size(clusterMessage))
	b, err = mo.MarshalAppend(b, clusterMessage)
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
	ShardRange  [2]uint32
	GatewayAddr []string
	ClusterAddr net.IP
	ClusterPort uint16
	RaftPort    int
	RaftState   rony.RaftState
	node        *memberlist.Node
}

func convertMember(sm *memberlist.Node) *ClusterMember {
	edgeNode := &rony.EdgeNode{}
	err := proto.UnmarshalOptions{}.Unmarshal(sm.Meta, edgeNode)
	if err != nil {
		log.Warn("Error On ConvertMember",
			zap.Error(err),
			zap.Int("Len", len(sm.Meta)),
		)
		return nil
	}

	return &ClusterMember{
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
	sync.RWMutex
	byServerID map[string]*ClusterMember
	leaderID   string
}

func (c *Cluster) GetByID(id string) *ClusterMember {
	c.RLock()
	defer c.RUnlock()

	return c.byServerID[id]
}

func (c *Cluster) LeaderID() string {
	return c.leaderID
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

	if m.RaftState == rony.RaftState_Leader {
		c.leaderID = m.ServerID
	}
	c.byServerID[m.ServerID] = m
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
	if m.ServerID == c.leaderID {
		c.leaderID = ""
	}

	delete(c.byServerID, m.ServerID)
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
		_ = joinRaft(d.edge, cm.ServerID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort))
	}
}

func (d delegateEvents) NotifyUpdate(n *memberlist.Node) {
	cm := convertMember(n)
	d.edge.cluster.AddMember(cm)
	if cm.ReplicaSet == d.edge.replicaSet {
		_ = joinRaft(d.edge, cm.ServerID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort))
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

func (d delegateEvents) NotifyLeave(n *memberlist.Node) {
	cm := convertMember(n)
	d.edge.cluster.RemoveMember(cm)
	if cm.ReplicaSet == d.edge.replicaSet {
		_ = leaveRaft(d.edge, cm.ServerID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort))
	}
}

func leaveRaft(edge *Server, nodeID, addr string) error {
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
	return nil
}

type delegateNode struct {
	edge *Server
}

func (d delegateNode) NodeMeta(limit int) []byte {
	n := rony.EdgeNode{
		ServerID:      d.edge.serverID,
		ShardRangeMin: d.edge.shardRange[0],
		ShardRangeMax: d.edge.shardRange[1],
		ReplicaSet:    d.edge.replicaSet,
		RaftPort:      uint32(d.edge.raftPort),
		GatewayAddr:   d.edge.gateway.Addr(),
		RaftState:     *rony.RaftState_None.Enum(),
	}
	if d.edge.raftEnabled {
		n.RaftState = *rony.RaftState(d.edge.raft.State() + 1).Enum()
	}

	b, _ := proto.Marshal(&n)
	if len(b) > limit {
		log.Warn("Too Large Meta", zap.ByteString("ServerID", d.edge.serverID))
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
	_ = proto.Unmarshal(data, cm)
	dispatchCtx := acquireDispatchCtx(d.edge, nil, 0, cm.Sender)
	dispatchCtx.FillEnvelope(
		cm.Envelope.GetRequestID(), cm.Envelope.GetConstructor(), cm.Envelope.Message,
		cm.Envelope.Auth, cm.Envelope.Header...,
	)

	d.edge.rateLimitChan <- struct{}{}
	go func() {
		// TODO:: handle error, for instance we might send back an error to the sender
		d.edge.onClusterMessage(dispatchCtx, cm.Store...)
		d.edge.dispatcher.Done(dispatchCtx)
		releaseDispatchCtx(dispatchCtx)
		releaseClusterMessage(cm)
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
