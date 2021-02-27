package edge

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/store"
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
	cluster  cluster.Cluster
	gateway  gateway.Gateway
	serverID string
	rs       uint64
}

func newBuiltin(serverID string, gw gateway.Gateway, c cluster.Cluster) *Builtin {
	return &Builtin{
		cluster:  c,
		gateway:  gw,
		serverID: serverID,
	}
}

func (pm *Builtin) GetNodes(ctx *RequestCtx, in *rony.MessageEnvelope) {
	req := rony.PoolGetNodes.Get()
	defer rony.PoolGetNodes.Put(req)
	res := rony.PoolEdges.Get()
	defer rony.PoolEdges.Put(res)
	err := req.Unmarshal(in.Message)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	if pm.cluster == nil {
		res.Nodes = append(res.Nodes, &rony.Edge{
			ReplicaSet: 0,
			ServerID:   pm.serverID,
			HostPorts:  pm.gateway.Addr(),
			Leader:     true,
		})
	} else if len(req.ReplicaSet) == 0 {
		members := pm.cluster.RaftMembers(pm.cluster.ReplicaSet())
		for _, m := range members {
			res.Nodes = append(res.Nodes, m.Proto(nil))
		}
	} else {
		for _, rs := range req.ReplicaSet {
			members := pm.cluster.RaftMembers(rs)
			for _, m := range members {
				res.Nodes = append(res.Nodes, m.Proto(nil))
			}
		}
	}
	ctx.PushMessage(rony.C_Edges, res)
	return
}

func (pm *Builtin) GetPage(ctx *RequestCtx, in *rony.MessageEnvelope) {
	if pm.cluster.ReplicaSet() != 1 {
		ctx.PushError(rony.ErrCodeUnavailable, rony.ErrItemRequest)
		return
	}

	req := rony.PoolGetPage.Get()
	defer rony.PoolGetPage.Put(req)
	res := rony.PoolPage.Get()
	defer rony.PoolPage.Put(res)
	err := proto.UnmarshalOptions{Merge: true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()

	err = store.Update(func(txn *badger.Txn) (err error) {
		_, err = rony.ReadPageWithTxn(txn, alloc, req.GetPageID(), res)
		if err == nil {
			return
		}
		if req.GetReplicaSet() == 0 {
			return err
		}
		res.ReplicaSet = req.GetReplicaSet()
		res.ID = req.GetPageID()
		return rony.SavePageWithTxn(txn, alloc, res)
	})
	if err != nil {
		ctx.PushError(rony.ErrCodeInternal, err.Error())
		return
	}
	ctx.PushMessage(rony.C_Page, res)
}
