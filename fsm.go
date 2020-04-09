package rony

import (
	"github.com/hashicorp/raft"
	"io"
)

/*
   Creation Time: 2020 - Feb - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type raftFSM struct {
	edge *EdgeServer
}

func (fsm raftFSM) Apply(raftLog *raft.Log) interface{} {
	raftCmd := acquireRaftCommand()
	err := raftCmd.Unmarshal(raftLog.Data)
	if err != nil {
		return err
	}
	dispatchCtx := acquireDispatchCtx(fsm.edge, nil, 0, raftCmd.AuthID, nil)
	dispatchCtx.FillEnvelope(raftCmd.Envelope.RequestID, raftCmd.Envelope.Constructor, raftCmd.Envelope.Message)
	err = fsm.edge.execute(dispatchCtx)
	releaseDispatchCtx(dispatchCtx)
	releaseRaftCommand(raftCmd)
	return err
}

func (fsm raftFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &snapshot{}, nil
}

func (fsm raftFSM) Restore(io.ReadCloser) error {
	return nil
}

type snapshot struct{}

func (s snapshot) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write([]byte{})
	if err != nil {
		return err
	}
	return sink.Close()
}

func (s snapshot) Release() {
}
