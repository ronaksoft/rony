package edge

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cluster"
	"github.com/ronaksoft/rony/store"
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
	cluster   cluster.Cluster
	serverID  string
	hostPorts []string
	rs        uint64
}

func newBuiltin(serverID string, hostPorts []string, c cluster.Cluster) *Builtin {
	return &Builtin{
		cluster:   c,
		serverID:  serverID,
		hostPorts: hostPorts,
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
			HostPorts:  pm.hostPorts,
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
	panic("implement me")
}

// Get returns the replica set which owns pageID. If this pageID was not exists in our database
// and newReplicaSet is non-zero, then database updates with new value, and the new value
// will be returned.
func (pm *Builtin) Get(ctx *RequestCtx, pageID uint32, newReplicaSet uint64) (uint64, error) {
	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()

	page := rony.Page{}
	err := store.Update(func(txn *badger.Txn) (err error) {
		_, err = rony.ReadPageWithTxn(txn, alloc, pageID, &page)
		if err == nil {
			return
		}
		if newReplicaSet == 0 {
			return
		}
		page.ReplicaSet = newReplicaSet
		page.ID = pageID
		return rony.SavePageWithTxn(txn, alloc, &page)
	})
	if err != nil {
		return 0, err
	}
	return page.ReplicaSet, nil
}

// Set updates the database with new entry. If replace is not set, then it will returns
// error if there is already a value in the database.
func (pm *Builtin) Set(ctx *RequestCtx, pageID uint32, replicaSet uint64, replace bool) error {
	alloc := store.NewAllocator()
	defer alloc.ReleaseAll()

	page := &rony.Page{
		ID:         pageID,
		ReplicaSet: replicaSet,
	}
	return store.Update(func(txn *badger.Txn) (err error) {
		_, err = rony.ReadPageWithTxn(txn, alloc, pageID, page)
		if err == nil && !replace {
			err = rony.ErrAlreadyExists
			return
		}
		page.ReplicaSet = replicaSet
		page.ID = pageID
		return rony.SavePageWithTxn(txn, alloc, page)
	})
}
