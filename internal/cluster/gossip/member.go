package gossipCluster

import (
	"github.com/hashicorp/memberlist"
	"github.com/pkg/errors"
	"github.com/ronaksoft/rony"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
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
	mtx         sync.RWMutex
	idx         int
	serverID    string
	replicaSet  uint64
	ShardRange  [2]uint32
	gatewayAddr []string
	tunnelAddr  []string
	ClusterAddr net.IP
	ClusterPort uint16
	raftPort    int
	raftState   rony.RaftState
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

func (m *Member) Proto(p *rony.Edge) *rony.Edge {
	if p == nil {
		p = &rony.Edge{}
	}
	p.ReplicaSet = m.replicaSet
	p.ServerID = m.serverID
	p.HostPorts = append(p.HostPorts, m.gatewayAddr...)
	p.Leader = m.raftState == rony.RaftState_Leader
	return p
}

func (m *Member) Merge(en *rony.EdgeNode) {
	m.replicaSet = en.GetReplicaSet()
	m.gatewayAddr = append(m.gatewayAddr[:0], en.GetGatewayAddr()...)
	m.tunnelAddr = append(m.tunnelAddr[:0], en.GetTunnelAddr()...)

	// merge raft
	m.raftPort = int(en.GetRaftPort())
	m.raftState = en.GetRaftState()

}

func (m *Member) TunnelConn() (net.Conn, error) {
	if len(m.tunnelAddr) == 0 {
		return nil, ErrNoTunnelAddrs
	}

	idx := m.idx
	for {
		conn, err := net.Dial("udp", m.tunnelAddr[idx])
		if err == nil {
			m.idx = idx
			return conn, nil
		}
		idx = (idx + 1) % len(m.tunnelAddr)
		if idx == m.idx {
			return nil, err
		}
	}
}

func newMember(sm *memberlist.Node) (*Member, error) {
	en := rony.PoolEdgeNode.Get()
	defer rony.PoolEdgeNode.Put(en)

	err := extractNode(sm, en)
	if err != nil {
		return nil, err
	}

	m := &Member{
		serverID:    string(en.ServerID),
		ClusterAddr: sm.Addr,
		ClusterPort: sm.Port,
	}
	m.Merge(en)

	return m, nil
}

func extractNode(n *memberlist.Node, en *rony.EdgeNode) error {
	return proto.UnmarshalOptions{}.Unmarshal(n.Meta, en)
}

var (
	ErrNoTunnelAddrs = errors.New("tunnel address does not found")
)
