package quicGateway

import "time"

/*
   Creation Time: 2019 - Aug - 26
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	flushDebounceTime = 10 * time.Millisecond
	flushDebounceMax  = 3
	writeWait         = time.Second
	readWait          = time.Second
	maxStreamsPerConn = 10
)
