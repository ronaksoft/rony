package badgerRaft

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/raft"
	"io"
)

/*
   Creation Time: 2021 - May - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func (fsm *Store) Apply(raftLog *raft.Log) interface{} {
	storeCmd := PoolStoreCommand.Get()
	defer PoolStoreCommand.Put(storeCmd)
	err := storeCmd.Unmarshal(raftLog.Data)
	if err != nil {
		return err
	}

	switch storeCmd.Type {
	case CommandType_CTStartTxn:
		return fsm.applyStartTxn(storeCmd.Payload)
	case CommandType_CTStopTxn:
		return fsm.applyStopTxn(storeCmd.Payload)
	case CommandType_CTCommitTxn:
		return fsm.applyCommitTxn(storeCmd.Payload)
	case CommandType_CTSet:
		return fsm.applySetTxn(storeCmd.Payload)
	case CommandType_CTDelete:
		return fsm.applyDeleteTxn(storeCmd.Payload)
	case CommandType_CTGet:
		kvs, _ := fsm.applyGetTxn(storeCmd.Payload)
		return kvs
	}

	// TODO:: error err ?
	return nil
}

func (fsm *Store) applyStartTxn(data []byte) error {
	x := PoolStartTxn.Get()
	defer PoolStartTxn.Put(x)

	err := x.Unmarshal(data)
	if err != nil {
		return err
	}

	fsm.openTxnMtx.Lock()
	defer fsm.openTxnMtx.Unlock()

	_, ok := fsm.openTxn[x.ID]
	if ok {
		return ErrDuplicateID
	}

	fsm.openTxn[x.ID] = fsm.db.NewTransaction(x.Update)
	return nil
}

func (fsm *Store) applyStopTxn(data []byte) error {
	x := PoolStopTxn.Get()
	defer PoolStopTxn.Put(x)

	err := x.Unmarshal(data)
	if err != nil {
		return err
	}

	fsm.openTxnMtx.Lock()
	defer fsm.openTxnMtx.Unlock()

	txn, ok := fsm.openTxn[x.ID]
	if !ok {
		return ErrTxnNotFound
	}

	if x.Commit {
		err = txn.Commit()
		if err != nil {
			return err
		}
	}

	txn.Discard()
	delete(fsm.openTxn, x.ID)
	return nil
}

func (fsm *Store) applyCommitTxn(data []byte) error {
	x := PoolCommitTxn.Get()
	defer PoolCommitTxn.Put(x)

	err := x.Unmarshal(data)
	if err != nil {
		return err
	}

	fsm.openTxnMtx.RLock()
	defer fsm.openTxnMtx.RUnlock()

	txn, ok := fsm.openTxn[x.ID]
	if ok {
		return ErrTxnNotFound
	}

	return txn.Commit()
}

func (fsm *Store) applySetTxn(data []byte) error {
	x := PoolSet.Get()
	defer PoolSet.Put(x)

	err := x.Unmarshal(data)
	if err != nil {
		return err
	}

	fsm.openTxnMtx.RLock()
	defer fsm.openTxnMtx.RUnlock()

	txn, ok := fsm.openTxn[x.TxnID]
	if ok {
		return ErrTxnNotFound
	}

	err = txn.Set(x.KV.Key, x.KV.Value)
	if err != nil {
		return err
	}

	return nil
}

func (fsm *Store) applyDeleteTxn(data []byte) error {
	x := PoolDelete.Get()
	defer PoolDelete.Put(x)

	err := x.Unmarshal(data)
	if err != nil {
		return err
	}

	fsm.openTxnMtx.RLock()
	defer fsm.openTxnMtx.RUnlock()

	txn, ok := fsm.openTxn[x.TxnID]
	if ok {
		return ErrTxnNotFound
	}

	err = txn.Delete(x.Key)
	if err != nil {
		return err
	}

	return nil
}

func (fsm *Store) applyGetTxn(data []byte) (*KeyValue, error) {
	req := PoolGet.Get()
	defer PoolGet.Put(req)
	res := &KeyValue{}

	err := req.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	fsm.openTxnMtx.RLock()
	defer fsm.openTxnMtx.RUnlock()

	txn, ok := fsm.openTxn[req.TxnID]
	if ok {
		return nil, ErrTxnNotFound
	}

	item, err := txn.Get(req.Key)
	switch err {
	case nil:
	case badger.ErrKeyNotFound:
		return res, nil
	default:
		return nil, err
	}
	res.Key = append(res.Key, req.Key...)
	_ = item.Value(func(val []byte) error {
		res.Value = append(res.Value, val...)
		return nil
	})
	return res, nil
}

func (fsm *Store) Snapshot() (raft.FSMSnapshot, error) {
	return &SnapshotBuilder{}, nil
}

func (fsm *Store) Restore(rd io.ReadCloser) error {

	return nil
}

// SnapshotBuilder is used for snapshot of Raft logs
type SnapshotBuilder struct{}

func (s SnapshotBuilder) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write([]byte{})
	if err != nil {
		return err
	}
	return sink.Close()
}

func (s SnapshotBuilder) Release() {
}
