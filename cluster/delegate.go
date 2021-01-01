package cluster

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	log "github.com/ronaksoft/rony/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2021 - Jan - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type clusterDelegate struct {
	c *Cluster
}

func (d clusterDelegate) NotifyJoin(n *memberlist.Node) {
	cm := convertMember(n)
	d.c.AddMember(cm)
	if cm.ReplicaSet == d.c.replicaSet {
		_ = joinRaft(d.c, cm.ServerID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort))
	}
}

func (d clusterDelegate) NotifyUpdate(n *memberlist.Node) {
	cm := convertMember(n)
	d.c.AddMember(cm)
	if cm.ReplicaSet == d.c.replicaSet {
		_ = joinRaft(d.c, cm.ServerID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort))
	}
}

func joinRaft(c *Cluster, nodeID, addr string) error {
	if !c.raftEnabled {
		return rony.ErrRaftNotSet
	}
	if c.raft.State() != raft.Leader {
		return rony.ErrNotRaftLeader
	}
	futureConfig := c.raft.GetConfiguration()
	if err := futureConfig.Error(); err != nil {
		return err
	}
	raftConf := futureConfig.Configuration()
	for _, srv := range raftConf.Servers {
		if srv.ID == raft.ServerID(nodeID) && srv.Address == raft.ServerAddress(addr) {
			return nil
		}
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			future := c.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return err
			}
		}
	}

	future := c.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if err := future.Error(); err != nil {
		return err
	}
	return nil
}

func (d clusterDelegate) NotifyLeave(n *memberlist.Node) {
	cm := convertMember(n)
	d.c.RemoveMember(cm)
	if cm.ReplicaSet == d.c.replicaSet {
		_ = leaveRaft(d.c, cm.ServerID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort))
	}
}

func leaveRaft(c *Cluster, nodeID, addr string) error {
	if !c.raftEnabled {
		return rony.ErrRaftNotSet
	}
	if c.raft.State() != raft.Leader {
		return rony.ErrNotRaftLeader
	}
	futureConfig := c.raft.GetConfiguration()
	if err := futureConfig.Error(); err != nil {
		return err
	}
	raftConf := futureConfig.Configuration()
	for _, srv := range raftConf.Servers {
		if srv.ID == raft.ServerID(nodeID) && srv.Address == raft.ServerAddress(addr) {
			return nil
		}
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			future := c.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d clusterDelegate) NodeMeta(limit int) []byte {
	n := rony.EdgeNode{
		ServerID:      d.c.serverID,
		ShardRangeMin: d.c.shardRange[0],
		ShardRangeMax: d.c.shardRange[1],
		ReplicaSet:    d.c.replicaSet,
		RaftPort:      uint32(d.c.raftPort),
		GatewayAddr:   d.c.gatewayAddr,
		RaftState:     *rony.RaftState_None.Enum(),
	}
	if d.c.raftEnabled {
		n.RaftState = *rony.RaftState(d.c.raft.State() + 1).Enum()
	}

	b, _ := proto.Marshal(&n)
	if len(b) > limit {
		log.Warn("Too Large Meta", zap.ByteString("ServerID", d.c.serverID))
		return nil
	}
	return b
}

func (d clusterDelegate) NotifyMsg(data []byte) {
	if ce := log.Check(log.DebugLevel, "Cluster Message Received"); ce != nil {
		ce.Write(
			zap.ByteString("ServerID", d.c.serverID),
			zap.Int("Data", len(data)),
		)
	}
	cm := acquireClusterMessage()
	_ = proto.Unmarshal(data, cm)
	d.c.onClusterMessage(cm)
	releaseClusterMessage(cm)
}

func (d clusterDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d clusterDelegate) LocalState(join bool) []byte {
	return nil
}

func (d clusterDelegate) MergeRemoteState(buf []byte, join bool) {
	return
}
