package gossipCluster

import (
	"github.com/ronaksoft/memberlist"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/internal/msg"
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

type Member struct {
	idx         int
	serverID    string
	replicaSet  uint64
	hash        uint64
	gatewayAddr []string
	tunnelAddr  []string
	ClusterAddr net.IP
	ClusterPort uint16
}

func (m *Member) ServerID() string {
	return m.serverID
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

func (m *Member) Proto(p *rony.Edge) *rony.Edge {
	if p == nil {
		p = &rony.Edge{}
	}
	p.ReplicaSet = m.replicaSet
	p.ServerID = m.serverID
	p.HostPorts = append(p.HostPorts, m.gatewayAddr...)
	return p
}

func (m *Member) Merge(en *msg.EdgeNode) {
	m.replicaSet = en.GetReplicaSet()
	m.gatewayAddr = append(m.gatewayAddr[:0], en.GetGatewayAddr()...)
	m.tunnelAddr = append(m.tunnelAddr[:0], en.GetTunnelAddr()...)
}

func (m *Member) Dial() (net.Conn, error) {
	if len(m.tunnelAddr) == 0 {
		return nil, errors.ErrNoTunnelAddrs
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
	en := msg.PoolEdgeNode.Get()
	defer msg.PoolEdgeNode.Put(en)

	if err := extractNode(sm, en); err != nil {
		return nil, err
	}

	m := &Member{
		serverID:    string(en.ServerID),
		hash:        en.Hash,
		ClusterAddr: sm.Addr,
		ClusterPort: sm.Port,
	}
	m.Merge(en)

	return m, nil
}

func extractNode(n *memberlist.Node, en *msg.EdgeNode) error {
	return proto.UnmarshalOptions{}.Unmarshal(n.Meta, en)
}
