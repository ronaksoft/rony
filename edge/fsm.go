package edge

import (
	"bytes"
	"git.ronaksoftware.com/ronak/rony"
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

// raftFSM is the finite state machine which will be used when Raft is enabled.
type raftFSM struct {
	edge *Server
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
		return rony.ErrRaftExecuteOnLeader
	}

	dispatchCtx := acquireDispatchCtx(fsm.edge, nil, 0, raftCmd.AuthID, raftCmd.Sender)
	dispatchCtx.FillEnvelope(raftCmd.Envelope.RequestID, raftCmd.Envelope.Constructor, raftCmd.Envelope.Message)

	err = fsm.edge.execute(dispatchCtx, false)
	fsm.edge.dispatcher.Done(dispatchCtx)
	releaseDispatchCtx(dispatchCtx)
	releaseRaftCommand(raftCmd)
	return nil
}

func (fsm raftFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &raftSnapshot{}, nil
}

func (fsm raftFSM) Restore(rd io.ReadCloser) error {

	return nil
}

// raftSnapshot is used for snapshot of Raft logs
type raftSnapshot struct{}

func (s raftSnapshot) Persist(sink raft.SnapshotSink) error {
	_, err := sink.Write([]byte{})
	if err != nil {
		return err
	}
	return sink.Close()
}

func (s raftSnapshot) Release() {
}
