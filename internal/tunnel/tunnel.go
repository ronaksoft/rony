package tunnel

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/msg"
)

/*
   Creation Time: 2021 - Jan - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type MessageHandler = func(conn rony.Conn, tm *msg.TunnelMessage)
