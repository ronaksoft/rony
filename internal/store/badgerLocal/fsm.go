package badgerLocal

import (
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

func (s *Store) Apply(log *raft.Log) interface{} {
	panic("BUG! not supported")
}

func (s *Store) Snapshot() (raft.FSMSnapshot, error) {
	panic("BUG! not supported")
}

func (s *Store) Restore(closer io.ReadCloser) error {
	panic("BUG! not supported")
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
