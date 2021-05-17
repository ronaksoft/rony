package edge

import (
	"github.com/ronaksoft/rony/internal/gateway"
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
	ReplicaSet      uint64
	Members         int
	MembershipScore int
	TunnelAddr      []string
	GatewayProtocol gateway.Protocol
	GatewayAddr     []string
}

// Stats exports some internal metrics data packed in 'Stats' struct
func (edge *Server) Stats() *Stats {
	s := Stats{}

	if edge.gateway != nil {
		s.GatewayAddr = edge.gateway.Addr()
		s.GatewayProtocol = edge.gateway.Protocol()
	}

	if edge.tunnel != nil {
		s.TunnelAddr = edge.tunnel.Addr()
	}

	if edge.cluster != nil {
		s.ReplicaSet = edge.cluster.ReplicaSet()
		s.Address = edge.cluster.Addr()
		s.Members = len(edge.cluster.Members())
	}

	return &s
}
