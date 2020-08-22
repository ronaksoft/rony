package edge

import (
	"fmt"
	"git.ronaksoft.com/ronak/rony/gateway"
)

/*
   Creation Time: 2020 - Mar - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Stats exports some internal metrics data
type Stats struct {
	Address         string
	RaftMembers     int
	RaftState       string
	ReplicaSet      uint64
	Members         int
	MembershipScore int
	GatewayProtocol gateway.Protocol
	GatewayAddr     []string
}

// Stats exports some internal metrics data packed in 'Stats' struct
func (edge *Server) Stats() *Stats {
	s := Stats{
		Address:         fmt.Sprintf("%s:%d", edge.gossip.LocalNode().Addr.String(), edge.gossip.LocalNode().Port),
		Members:         len(edge.gossip.Members()),
		MembershipScore: edge.gossip.GetHealthScore(),
		GatewayProtocol: edge.gatewayProtocol,
		GatewayAddr:     edge.gateway.Addr(),
		ReplicaSet:      edge.replicaSet,
	}

	if edge.raftEnabled {
		s.RaftState = edge.raft.State().String()
		f := edge.raft.GetConfiguration()
		if f.Error() == nil {
			s.RaftMembers = len(f.Configuration().Servers)
		}
	}

	return &s
}
