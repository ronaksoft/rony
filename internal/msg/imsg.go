package msg

import (
	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2021 - Jul - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. -I=../.. --go_out=paths=source_relative:. imsg.proto
//go:generate protoc -I=. -I=../.. --gorony_out=paths=source_relative,rony_opt=no_edge_dep:. imsg.proto
func init() {}

/*
	Extra methods for TunnelMessage
*/

func (x *TunnelMessage) Fill(senderID []byte, senderReplicaSet uint64, e *rony.MessageEnvelope, kvs ...*rony.KeyValue) {
	x.SenderID = append(x.SenderID[:0], senderID...)
	x.SenderReplicaSet = senderReplicaSet
	x.Store = append(x.Store[:0], kvs...)
	if x.Envelope == nil {
		x.Envelope = rony.PoolMessageEnvelope.Get()
	}
	e.DeepCopy(x.Envelope)
}
