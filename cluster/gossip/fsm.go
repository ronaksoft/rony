package gossipCluster

import (
	"bytes"
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony"
	"google.golang.org/protobuf/proto"
	"io"
)

/*
   Creation Time: 2020 - Feb - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// raftFSM is the finite state machine which will be used when Raft is enabled.
type raftFSM struct {
	c *Cluster
}

func (fsm raftFSM) Apply(raftLog *raft.Log) interface{} {
	raftCmd := acquireRaftCommand()
	err := proto.UnmarshalOptions{}.Unmarshal(raftLog.Data, raftCmd)
	if err != nil {
		return err
	}

	// We dont execute the command, if we are the sender server
	if bytes.Equal(raftCmd.Sender, fsm.c.localServerID) {
		releaseRaftCommand(raftCmd)
		return rony.ErrRaftExecuteOnLeader
	}

	err = fsm.c.ReplicaMessageHandler(raftCmd)
	if err != nil {
		return rony.ErrRaftExecuteOnLeader
	}
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
