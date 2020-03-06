package rony

import (
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
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
		log.Fatal("Error On Raft Apply",
			zap.Int("Len", len(raftLog.Data)),
			zap.Any("LogType", raftLog.Type),
			zap.Uint64("Index", raftLog.Index),
			zap.Uint64("Term", raftLog.Term),
		)
	}
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
