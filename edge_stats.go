package rony

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gateway"
)

/*
   Creation Time: 2020 - Mar - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// EdgeStats exports some internal metrics data
type EdgeStats struct {
	Address         string
	RaftMembers     int
	RaftState       string
	Members         int
	MembershipScore int
	GatewayProtocol gateway.Protocol
	GatewayAddr     string
}

// Stats exports some internal metrics data packed in 'EdgeStats' struct
func (edge *EdgeServer) Stats() *EdgeStats {
	s := EdgeStats{
		Address:         fmt.Sprintf("%s:%d", edge.gossip.LocalNode().Addr.String(), edge.gossip.LocalNode().Port),
		Members:         len(edge.gossip.Members()),
		MembershipScore: edge.gossip.GetHealthScore(),
		GatewayProtocol: edge.gatewayProtocol,
		GatewayAddr:     edge.gateway.Addr(),
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
