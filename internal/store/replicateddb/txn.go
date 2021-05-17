package replicateddb

import (
	"context"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - May - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func (fsm *Store) newTxn(update bool) *Txn {
	return &Txn{
		ID:     tools.RandomInt64(0),
		store:  fsm,
		update: update,
	}
}

func (fsm *Store) startTxn(txn *Txn) error {
	req := PoolStartTxn.Get()
	defer PoolStartTxn.Put(req)
	req.ID = txn.ID
	req.Update = txn.update
	b := pools.Buffer.FromProto(req)
	ctx, cf := context.WithTimeout(context.TODO(), ProposeTimeout)
	defer cf()
	err := fsm.raft.Propose(ctx, *b.Bytes())
	pools.Buffer.Put(b)
	return err
}

func (fsm *Store) stopTxn(txn *Txn, commit bool) error {
	req := PoolStopTxn.Get()
	defer PoolStopTxn.Put(req)
	req.ID = txn.ID
	req.Commit = commit
	b := pools.Buffer.FromProto(req)
	ctx, cf := context.WithTimeout(context.TODO(), ProposeTimeout)
	defer cf()
	err := fsm.raft.Propose(ctx, *b.Bytes())
	pools.Buffer.Put(b)
	return err
}

type Txn struct {
	ID     int64
	store  *Store
	update bool
}

func (txn *Txn) Delete(alloc *tools.Allocator, keyParts ...interface{}) error {
	me := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(me)
	me.Fill(0, C_Delete, &Delete{
		TxnID: txn.ID,
		Key:   alloc.Gen(keyParts),
	})

	b := pools.Buffer.FromProto(me)
	ctx, cf := context.WithTimeout(context.TODO(), ProposeTimeout)
	defer cf()
	err := txn.store.raft.Propose(ctx, *b.Bytes())
	pools.Buffer.Put(b)
	return err
}

func (txn *Txn) Set(alloc *tools.Allocator, val []byte, keyParts ...interface{}) error {
	me := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(me)
	me.Fill(0, C_Set, &Set{
		TxnID: txn.ID,
		Key:   alloc.Gen(keyParts),
		Value: val,
	})

	b := pools.Buffer.FromProto(me)
	ctx, cf := context.WithTimeout(context.TODO(), ProposeTimeout)
	defer cf()
	err := txn.store.raft.Propose(ctx, *b.Bytes())
	pools.Buffer.Put(b)
	return err
}

func (txn *Txn) Get(alloc *tools.Allocator, keyParts ...interface{}) (v []byte, err error) {
	txn.store.openTxnMtx.Lock()
	ltxn := txn.store.openTxn[txn.ID]
	txn.store.openTxnMtx.Unlock()

	item, err := ltxn.Get(alloc.Gen(keyParts...))
	if err != nil {
		return nil, err
	}
	v, err = item.ValueCopy(v)
	return
}

func (txn *Txn) Exists(alloc *tools.Allocator, keyParts ...interface{}) bool {
	_, err := txn.Get(alloc, keyParts...)
	return err == nil
}
