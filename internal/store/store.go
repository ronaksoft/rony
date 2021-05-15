package store

import (
	"github.com/hashicorp/raft"
	"github.com/ronaksoft/rony/internal/cluster"
	"io"
)

/*
   Creation Time: 2021 - May - 15
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Store is the finite state machine which will be used when Raft is enabled.
type Store struct {
	c cluster.Cluster
}

func (fsm Store) Apply(raftLog *raft.Log) interface{} {
	// raftCmd := acquireRaftCommand()
	// err := proto.UnmarshalOptions{}.Unmarshal(raftLog.Data, raftCmd)
	// if err != nil {
	// 	return err
	// }
	//
	// // We don't execute the command, if we are the sender server
	// if bytes.Equal(raftCmd.Sender, fsm.c.localServerID) {
	// 	releaseRaftCommand(raftCmd)
	// 	return rony.ErrRaftExecuteOnLeader
	// }
	//
	// err = fsm.c.ReplicaMessageHandler(raftCmd)
	// if err != nil {
	// 	return rony.ErrRaftExecuteOnLeader
	// }
	// releaseRaftCommand(raftCmd)
	return nil
}

func (fsm Store) Snapshot() (raft.FSMSnapshot, error) {
	return &SnapshotBuilder{}, nil
}

func (fsm Store) Restore(rd io.ReadCloser) error {

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
