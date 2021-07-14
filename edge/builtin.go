package edge

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/internal/msg"
	"github.com/ronaksoft/rony/tools"
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
	pageRepo *msg.PageLocalRepo
	cluster  rony.Cluster
	gateway  rony.Gateway
	store    rony.Store
	serverID string
}

func newBuiltin(edgeServer *Server) *Builtin {
	b := &Builtin{
		cluster:  edgeServer.Cluster(),
		gateway:  edgeServer.Gateway(),
		store:    edgeServer.Store(),
		serverID: edgeServer.GetServerID(),
		pageRepo: msg.NewPageLocalRepo(edgeServer.Store()),
	}
	return b
}

func (pm *Builtin) getNodes(ctx *RequestCtx, in *rony.MessageEnvelope) {
	req := rony.PoolGetNodes.Get()
	defer rony.PoolGetNodes.Put(req)
	res := rony.PoolEdges.Get()
	defer rony.PoolEdges.Put(res)
	err := req.Unmarshal(in.Message)
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
	err := req.Unmarshal(in.Message)
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

func (pm *Builtin) getPage(ctx *RequestCtx, in *rony.MessageEnvelope) {
	if pm.cluster.ReplicaSet() != 1 {
		ctx.PushError(errors.ErrUnavailableRequest)
		return
	}

	req := msg.PoolGetPage.Get()
	defer msg.PoolGetPage.Put(req)
	res := msg.PoolPage.Get()
	defer msg.PoolPage.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	err = pm.store.Update(func(txn *badger.Txn) (err error) {
		_, err = pm.pageRepo.ReadWithTxn(txn, alloc, req.GetPageID(), res)
		if err == nil {
			return
		}
		if req.GetReplicaSet() == 0 {
			return err
		}
		res.ReplicaSet = req.GetReplicaSet()
		res.ID = req.GetPageID()
		return pm.pageRepo.SaveWithTxn(txn, alloc, res)
	})
	if err != nil {
		ctx.PushError(errors.GenInternalErr(err.Error(), err))
		return
	}
	ctx.PushMessage(msg.C_Page, res)
}
