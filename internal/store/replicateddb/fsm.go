package replicateddb

import (
	"context"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/pools"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.uber.org/zap"
	"sync"
	"time"
)

/*
   Creation Time: 2021 - May - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func (fsm *Store) raftLoop() {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			fsm.raft.Tick()
		case rd := <-fsm.raft.Ready():

			err := fsm.wal.Save(&rd.HardState, rd.Entries, &rd.Snapshot)
			if err != nil {
				log.Warn("Error On Storage Save", zap.Error(err))
			}

			if rd.MustSync {
				err = fsm.wal.Sync()
				if err != nil {
					log.Warn("Error On Storage Sync", zap.Error(err))
				}
			}

			wg := &sync.WaitGroup{}
			wg.Add(2)
			go fsm.sendMessages(wg, rd.Messages)
			go fsm.handleCommittedEntries(wg, rd.CommittedEntries)
			wg.Wait()
			fsm.raft.Advance()

		}
	}
}

func (fsm *Store) sendMessages(wg *sync.WaitGroup, msgs []raftpb.Message) {
	defer wg.Done()

	for _, msg := range msgs {
		member := fsm.c.MemberByHash(msg.To)
		if member == nil {
			continue
		}
		c, err := member.Dial()
		if err != nil {
			continue
		}

		b := pools.Bytes.GetLen(msg.Size())
		_, _ = msg.MarshalToSizedBuffer(b)
		me := rony.PoolMessageEnvelope.Get()
		me.Constructor = rony.C_StoreMessage
		me.Message = append(me.Message[:0], b...)
		pools.Bytes.Put(b)
		buf := pools.Buffer.FromProto(me)
		_, err = c.Write(*buf.Bytes())
		pools.Buffer.Put(buf)
		switch msg.Type {
		case raftpb.MsgSnap:
			if err != nil {
				fsm.raft.ReportSnapshot(msg.From, raft.SnapshotFailure)
			} else {
				fsm.raft.ReportSnapshot(msg.From, raft.SnapshotFinish)
			}
		}
	}

}

func (fsm *Store) handleCommittedEntries(wg *sync.WaitGroup, entries []raftpb.Entry) {
	defer wg.Done()
	for _, ce := range entries {
		switch ce.Type {
		case raftpb.EntryNormal:
			var err error
			me := rony.PoolMessageEnvelope.Get()
			_ = me.Unmarshal(ce.Data)
			switch me.Constructor {
			case C_StartTxn:
				err = fsm.applyStartTxn(me.Message)
			case C_StopTxn:
				err = fsm.applyStopTxn(me.Message)
			case C_CommitTxn:
				err = fsm.applyCommitTxn(me.Message)
			case C_Set:
				err = fsm.applySetTxn(me.Message)
			case C_Delete:
				err = fsm.applyDeleteTxn(me.Message)
			}
			rony.PoolMessageEnvelope.Put(me)
			log.Warn("Error On Applying Committed Entry", zap.Error(err))
			// TODO:: handle entry
		case raftpb.EntryConfChange:
			cc := raftpb.ConfChange{}
			_ = cc.Unmarshal(ce.Data)
			fsm.state = fsm.raft.ApplyConfChange(cc)
		}

	}

}

func (fsm *Store) Step(payload []byte) {
	ctx, cf := context.WithTimeout(context.TODO(), ProposeTimeout)
	defer cf()
	m := raftpb.Message{}
	_ = m.Unmarshal(payload)
	err := fsm.raft.Step(ctx, m)
	if err != nil {
		log.Warn("Error On Store RaftApply", zap.Error(err))
	}

	return
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
	if !ok {
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
	if !ok {
		return err
	}

	return txn.Set(x.Key, x.Value)
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
	if !ok {
		return ErrTxnNotFound
	}

	return txn.Delete(x.Key)
}
