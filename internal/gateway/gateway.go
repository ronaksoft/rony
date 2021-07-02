package gateway

import (
	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2019 - Aug - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type (
	ConnectHandler = func(c rony.Conn, kvs ...*rony.KeyValue)
	MessageHandler = func(c rony.Conn, streamID int64, data []byte)
	CloseHandler   = func(c rony.Conn)
)
