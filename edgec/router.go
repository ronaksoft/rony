package edgec

import (
	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2021 - Jan - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Router interface {
	UpdateRoute(req *rony.MessageEnvelope, replicaSet uint64)
	GetRoute(req *rony.MessageEnvelope) (replicaSet uint64)
}
