package edge

import (
	"github.com/ronaksoft/rony/internal/gateway"
)

/*
   Creation Time: 2021 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Proxy interface {
	OnReceive(ctx *gateway.RequestCtx)
	OnSend(d []byte) []byte
}
