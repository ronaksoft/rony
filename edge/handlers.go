package edge

import (
	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2021 - Jan - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func (edge *Server) getNodes(ctx *RequestCtx, in *rony.MessageEnvelope) {
	req := rony.PoolGetNodes.Get()
	defer rony.PoolGetNodes.Put(req)
	res := rony.PoolNodeInfoMany.Get()
	defer rony.PoolNodeInfoMany.Put(res)
	err := req.Unmarshal(in.Message)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	if len(req.ReplicaSet) == 0 {
		members := edge.cluster.RaftMembers(edge.cluster.ReplicaSet())
		for _, m := range members {
			res.Nodes = append(res.Nodes, m.Proto(nil))
		}
	} else {
		for _, rs := range req.ReplicaSet {
			members := edge.cluster.RaftMembers(rs)
			for _, m := range members {
				res.Nodes = append(res.Nodes, m.Proto(nil))
			}
		}
	}

	ctx.PushMessage(rony.C_NodeInfoMany, res)
	return

}