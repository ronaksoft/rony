package edgec

import (
	"context"

	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2020 - Jul - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Client interface {
	Send(ctx context.Context, req, res *rony.MessageEnvelope) error
	Close() error
	GetRequestID() uint64
}
