package gossipCluster

import (
	"github.com/hashicorp/memberlist"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
)

/*
   Creation Time: 2021 - Jan - 14
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Member
type Member struct {
	serverID    string
	replicaSet  uint64
	ShardRange  [2]uint32
	gatewayAddr []string
	tunnelAddr  []string
	ClusterAddr net.IP
	ClusterPort uint16
	raftPort    int
	raftState   rony.RaftState
	node        *memberlist.Node
}

func (m *Member) ServerID() string {
	return m.serverID
}

func (m *Member) RaftState() rony.RaftState {
	return m.raftState
}

func (m *Member) ReplicaSet() uint64 {
	return m.replicaSet
}

func (m *Member) GatewayAddr() []string {
	return m.gatewayAddr
}

func (m *Member) TunnelAddr() []string {
	return m.tunnelAddr
}

func (m *Member) RaftPort() int {
	return m.raftPort
}

func (m *Member) Proto(p *rony.NodeInfo) *rony.NodeInfo {
	if p == nil {
		p = &rony.NodeInfo{}
	}
	p.ReplicaSet = m.replicaSet
	p.ServerID = m.serverID
	p.HostPorts = append(p.HostPorts, m.gatewayAddr...)
	p.Leader = m.raftState == rony.RaftState_Leader
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
		serverID:    tools.ByteToStr(edgeNode.ServerID),
		replicaSet:  edgeNode.GetReplicaSet(),
		ShardRange:  [2]uint32{edgeNode.ShardRangeMin, edgeNode.ShardRangeMax},
		gatewayAddr: edgeNode.GatewayAddr,
		tunnelAddr:  edgeNode.TunnelAddr,
		raftPort:    int(edgeNode.GetRaftPort()),
		raftState:   edgeNode.GetRaftState(),
		ClusterAddr: sm.Addr,
		ClusterPort: sm.Port,
		node:        sm,
	}
}
