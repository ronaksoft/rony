package gossipCluster

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
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

func (d *clusterDelegate) NotifyJoin(n *memberlist.Node) {
	cm := convertMember(n)
	d.c.addMember(cm)

	if cm.replicaSet != 0 && cm.replicaSet == d.c.cfg.ReplicaSet {
		err := joinRaft(d.c, cm.serverID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort()))
		if err != nil {
			if ce := log.Check(log.DebugLevel, "Error On Join Raft (NodeJoin)"); ce != nil {
				ce.Write(
					zap.ByteString("This", d.c.localServerID),
					zap.String("NodeID", cm.serverID),
					zap.Error(err),
				)
			}
		} else {
			log.Info("Join Raft (NodeJoin)",
				zap.ByteString("This", d.c.localServerID),
				zap.String("NodeID", cm.serverID),
			)
		}
	}
}

func (d *clusterDelegate) NotifyUpdate(n *memberlist.Node) {
	cm := convertMember(n)
	d.c.addMember(cm)
	if cm.replicaSet != 0 && cm.replicaSet == d.c.cfg.ReplicaSet {
		err := joinRaft(d.c, cm.serverID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort()))
		if err != nil {
			if ce := log.Check(log.DebugLevel, "Error On Join Raft (NodeUpdate)"); ce != nil {
				ce.Write(
					zap.ByteString("This", d.c.localServerID),
					zap.String("NodeID", cm.serverID),
					zap.Error(err),
				)
			}
		} else {
			log.Info("Join Raft (NodeUpdate)",
				zap.ByteString("This", d.c.localServerID),
				zap.String("NodeID", cm.serverID),
			)
		}
	}
}

func (d *clusterDelegate) NotifyAlive(n *memberlist.Node) error {
	cm := convertMember(n)
	d.c.addMember(cm)
	if cm.replicaSet != 0 && cm.replicaSet == d.c.cfg.ReplicaSet {
		err := joinRaft(d.c, cm.serverID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort()))
		if err == nil {
			if ce := log.Check(log.DebugLevel, "Join Raft (NodeAlive)"); ce != nil {
				ce.Write(
					zap.ByteString("This", d.c.localServerID),
					zap.String("NodeID", cm.serverID),
				)
			}
		}
	}
	return nil
}

func joinRaft(c *Cluster, nodeID, addr string) error {
	if c.raft == nil {
		return nil
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
			return rony.ErrRaftAlreadyJoined
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

func (d *clusterDelegate) NotifyLeave(n *memberlist.Node) {
	cm := convertMember(n)
	d.c.removeMember(cm)
	if cm.replicaSet == d.c.cfg.ReplicaSet {
		_ = leaveRaft(d.c, cm.serverID, fmt.Sprintf("%s:%d", cm.ClusterAddr.String(), cm.RaftPort()))
	}
}

func leaveRaft(c *Cluster, nodeID, addr string) error {
	if c.raft == nil {
		return nil
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

func (d *clusterDelegate) NodeMeta(limit int) []byte {
	n := rony.EdgeNode{
		ServerID:      d.c.localServerID,
		ShardRangeMin: d.c.localShardRange[0],
		ShardRangeMax: d.c.localShardRange[1],
		ReplicaSet:    d.c.cfg.ReplicaSet,
		RaftPort:      uint32(d.c.cfg.RaftPort),
		GatewayAddr:   d.c.localGatewayAddr,
		TunnelAddr:    d.c.localTunnelAddr,
		RaftState:     *rony.RaftState_Leader.Enum(),
	}

	if d.c.raft != nil {
		n.RaftState = *rony.RaftState(d.c.raft.State() + 1).Enum()
	}

	b, _ := proto.Marshal(&n)
	if len(b) > limit {
		log.Warn("Too Large Meta", zap.ByteString("ServerID", d.c.localServerID))
		return nil
	}
	return b
}

func (d *clusterDelegate) NotifyMsg(data []byte) {}

func (d *clusterDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (d *clusterDelegate) LocalState(join bool) []byte {
	return nil
}

func (d *clusterDelegate) MergeRemoteState(buf []byte, join bool) {
	return
}
