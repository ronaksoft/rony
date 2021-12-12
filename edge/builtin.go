package edge

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/store"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2021 - Jan - 11
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Builtin keep track of pages distribution over Edge servers.
type Builtin struct {
	cluster  rony.Cluster
	gateway  rony.Gateway
	store    *store.Store
	serverID string
}

func newBuiltin(edgeServer *Server) *Builtin {
	b := &Builtin{
		cluster:  edgeServer.Cluster(),
		gateway:  edgeServer.Gateway(),
		serverID: edgeServer.GetServerID(),
	}

	return b
}

func (pm *Builtin) getNodes(ctx *RequestCtx, in *rony.MessageEnvelope) {
	req := rony.PoolGetNodes.Get()
	defer rony.PoolGetNodes.Put(req)
	res := rony.PoolEdges.Get()
	defer rony.PoolEdges.Put(res)

	var err error
	if in.JsonEncoded {
		err = protojson.Unmarshal(in.Message, req)
	} else {
		err = proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	}
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)

		return
	}

	if pm.cluster == nil {
		res.Nodes = append(res.Nodes, &rony.Edge{
			ReplicaSet: 0,
			ServerID:   pm.serverID,
			HostPorts:  pm.gateway.Addr(),
		})
	} else if len(req.ReplicaSet) == 0 {
		members := pm.cluster.MembersByReplicaSet(pm.cluster.ReplicaSet())
		for _, m := range members {
			res.Nodes = append(res.Nodes, m.Proto(nil))
		}
	} else {
		members := pm.cluster.MembersByReplicaSet(req.ReplicaSet...)
		for _, m := range members {
			res.Nodes = append(res.Nodes, m.Proto(nil))
		}
	}

	ctx.PushMessage(rony.C_Edges, res)
}

func (pm *Builtin) getAllNodes(ctx *RequestCtx, in *rony.MessageEnvelope) {
	req := rony.PoolGetNodes.Get()
	defer rony.PoolGetNodes.Put(req)
	res := rony.PoolEdges.Get()
	defer rony.PoolEdges.Put(res)

	var err error
	if in.JsonEncoded {
		err = protojson.Unmarshal(in.Message, req)
	} else {
		err = proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	}
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)

		return
	}

	if pm.cluster == nil {
		res.Nodes = append(res.Nodes, &rony.Edge{
			ReplicaSet: 0,
			ServerID:   pm.serverID,
			HostPorts:  pm.gateway.Addr(),
		})
	} else {
		members := pm.cluster.Members()
		for _, m := range members {
			res.Nodes = append(res.Nodes, m.Proto(nil))
		}
	}
	ctx.PushMessage(rony.C_Edges, res)
}

func (pm *Builtin) ping(ctx *RequestCtx, in *rony.MessageEnvelope) {
	req := rony.PoolPing.Get()
	defer rony.PoolPing.Put(req)
	res := rony.PoolPong.Get()
	defer rony.PoolPong.Put(res)

	var err error
	if in.JsonEncoded {
		err = protojson.Unmarshal(in.Message, req)
	} else {
		err = proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	}
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)

		return
	}

	res.ID = req.ID
	ctx.PushMessage(rony.C_Pong, res)
}
