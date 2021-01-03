package gossipCluster

import (
	"github.com/ronaksoft/rony"
	"sync"
)

/*
   Creation Time: 2019 - Oct - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var raftCommandPool sync.Pool

func acquireRaftCommand() *rony.RaftCommand {
	v := raftCommandPool.Get()
	if v == nil {
		return &rony.RaftCommand{
			Envelope: &rony.MessageEnvelope{},
		}
	}
	return v.(*rony.RaftCommand)
}

func releaseRaftCommand(x *rony.RaftCommand) {
	x.Sender = x.Sender[:0]
	x.Store = x.Store[:0]
	x.Envelope.Message = x.Envelope.Message[:0]
	x.Envelope.RequestID = 0
	x.Envelope.Constructor = 0
	raftCommandPool.Put(x)
}
