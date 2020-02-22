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

type RaftConfig struct {
	Enabled      bool
	Port         int
	FirstMachine bool
}

func (edge EdgeServer) Join(nodeID, addr string) error {
	futureConfig := edge.raft.GetConfiguration()
	if err := futureConfig.Error(); err != nil {
		return err
	}

	for _, srv := range futureConfig.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeID) && srv.Address == raft.ServerAddress(addr) {
			return nil
		}
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			future := edge.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return err
			}
		}
	}

	future := edge.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if err := future.Error(); err != nil {
		return err
	}
	return nil
}

func (edge EdgeServer) Apply(*raft.Log) interface{} {
	panic("implement me")
}

func (edge EdgeServer) Snapshot() (raft.FSMSnapshot, error) {
	panic("implement me")
}

func (edge EdgeServer) Restore(io.ReadCloser) error {
	panic("implement me")
}
