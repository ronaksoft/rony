package rony

import (
	"bytes"
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

	// We dont execute the command, if we are the sender server
	if bytes.Equal(raftCmd.Sender, fsm.edge.serverID) {
		releaseRaftCommand(raftCmd)
		return ErrRaftExecuteOnLeader
	}

	dispatchCtx := acquireDispatchCtx(fsm.edge, nil, 0, raftCmd.AuthID, raftCmd.Sender)
	dispatchCtx.FillEnvelope(raftCmd.Envelope.RequestID, raftCmd.Envelope.Constructor, raftCmd.Envelope.Message)
	err = fsm.edge.execute(dispatchCtx)
	releaseDispatchCtx(dispatchCtx)
	releaseRaftCommand(raftCmd)
	return nil
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
