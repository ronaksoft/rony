package edge

import (
	"github.com/ronaksoft/rony/gateway"
)

/*
   Creation Time: 2020 - Mar - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Stats exports some internal metrics data
type Stats struct {
	Address         string
	RaftMembers     int
	RaftState       string
	ReplicaSet      uint64
	Members         int
	MembershipScore int
	TunnelAddr      []string
	GatewayProtocol gateway.Protocol
	GatewayAddr     []string
}

// Stats exports some internal metrics data packed in 'Stats' struct
func (edge *Server) Stats() *Stats {
	s := Stats{
		GatewayProtocol: edge.gatewayProtocol,
	}

	if edge.gateway != nil {
		s.GatewayAddr = edge.gateway.Addr()
	}

	if edge.tunnel != nil {
		s.TunnelAddr = edge.tunnel.Addr()
	}

	if edge.cluster != nil {
		s.ReplicaSet = edge.cluster.ReplicaSet()
		s.Address = edge.cluster.Addr()
		s.Members = len(edge.cluster.Members())
		s.RaftMembers = len(edge.cluster.RaftMembers(s.ReplicaSet))
		s.RaftState = edge.cluster.RaftState().String()
	}

	return &s
}
