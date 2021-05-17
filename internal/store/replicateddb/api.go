package replicateddb

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - May - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func (fsm *Store) Apply(req, res *rony.StoreMessage) {
	switch req.Constructor {
	case C_StartTxn:
		fsm.applyStartTxn(req, res)
	case C_StopTxn:
		fsm.applyStopTxn(req, res)
	case C_CommitTxn:
		fsm.applyCommitTxn(req, res)
	case C_Set:
		fsm.applySetTxn(req, res)
	case C_Delete:
		fsm.applyDeleteTxn(req, res)
	case C_Get:
		fsm.applyGetTxn(req, res)

	}
	return
}

func fillError(x *rony.StoreMessage, code, item string) {
	e := rony.PoolError.Get()
	e.Code = code
	e.Items = item
	buf := pools.Buffer.FromProto(e)
	x.Constructor = rony.C_Error
	x.Payload = append(x.Payload[:0], *buf.Bytes()...)
	rony.PoolError.Put(e)
}

func fillKeyValue(x *rony.StoreMessage, key, val string) {
	e := rony.PoolKeyValue.Get()
	e.Key = key
	e.Value = val
	buf := pools.Buffer.FromProto(e)
	x.Constructor = rony.C_KeyValue
	x.Payload = append(x.Payload[:0], *buf.Bytes()...)
	rony.PoolKeyValue.Put(e)
}

func (fsm *Store) applyStartTxn(req, res *rony.StoreMessage) {
	x := PoolStartTxn.Get()
	defer PoolStartTxn.Put(x)

	err := x.Unmarshal(req.Payload)
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
		return
	}

	fsm.openTxnMtx.Lock()
	defer fsm.openTxnMtx.Unlock()

	_, ok := fsm.openTxn[x.ID]
	if ok {
		fillError(res, rony.ErrCodeAlreadyExists, "TXN_ID")
		return
	}

	fsm.openTxn[x.ID] = fsm.db.NewTransaction(x.Update)

	return
}

func (fsm *Store) applyStopTxn(req, res *rony.StoreMessage) {
	x := PoolStopTxn.Get()
	defer PoolStopTxn.Put(x)

	err := x.Unmarshal(req.Payload)
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
		return
	}

	fsm.openTxnMtx.Lock()
	defer fsm.openTxnMtx.Unlock()

	txn, ok := fsm.openTxn[x.ID]
	if !ok {
		fillError(res, rony.ErrCodeUnavailable, "TXN_ID")
		return
	}

	if x.Commit {
		err = txn.Commit()
		if err != nil {
			fillError(res, rony.ErrCodeInternal, err.Error())
			return
		}
	}

	txn.Discard()
	delete(fsm.openTxn, x.ID)
	return
}

func (fsm *Store) applyCommitTxn(req, res *rony.StoreMessage) {
	x := PoolCommitTxn.Get()
	defer PoolCommitTxn.Put(x)

	err := x.Unmarshal(req.Payload)
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
		return
	}

	fsm.openTxnMtx.RLock()
	defer fsm.openTxnMtx.RUnlock()

	txn, ok := fsm.openTxn[x.ID]
	if !ok {
		fillError(res, rony.ErrCodeUnavailable, "TXN_ID")
		return
	}

	err = txn.Commit()
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
	}
}

func (fsm *Store) applySetTxn(req, res *rony.StoreMessage) {
	x := PoolSet.Get()
	defer PoolSet.Put(x)

	err := x.Unmarshal(req.Payload)
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
		return
	}

	fsm.openTxnMtx.RLock()
	defer fsm.openTxnMtx.RUnlock()

	txn, ok := fsm.openTxn[x.TxnID]
	if !ok {
		fillError(res, rony.ErrCodeUnavailable, "TXN_ID")
		return
	}

	err = txn.Set(x.Key, x.Value)
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
		return
	}

	return
}

func (fsm *Store) applyDeleteTxn(req, res *rony.StoreMessage) {
	x := PoolDelete.Get()
	defer PoolDelete.Put(x)

	err := x.Unmarshal(req.Payload)
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
		return
	}

	fsm.openTxnMtx.RLock()
	defer fsm.openTxnMtx.RUnlock()

	txn, ok := fsm.openTxn[x.TxnID]
	if !ok {
		fillError(res, rony.ErrCodeUnavailable, "TXN_ID")
		return
	}

	err = txn.Delete(x.Key)
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
		return

	}

	return
}

func (fsm *Store) applyGetTxn(req, res *rony.StoreMessage) {
	x := PoolGet.Get()
	defer PoolGet.Put(x)

	err := req.Unmarshal(req.Payload)
	if err != nil {
		fillError(res, rony.ErrCodeInternal, err.Error())
		return
	}

	fsm.openTxnMtx.RLock()
	defer fsm.openTxnMtx.RUnlock()

	txn, ok := fsm.openTxn[x.TxnID]
	if !ok {
		fillError(res, rony.ErrCodeUnavailable, "TXN_ID")
		return
	}

	item, err := txn.Get(x.Key)
	switch err {
	case nil:
	case store.ErrKeyNotFound:
		fillError(res, rony.ErrCodeUnavailable, "KEY")
		return
	default:
		fillError(res, rony.ErrCodeInternal, err.Error())
		return
	}

	_ = item.Value(func(val []byte) error {
		fillKeyValue(res, tools.B2S(x.Key), tools.B2S(val))
		return nil
	})
	return
}
