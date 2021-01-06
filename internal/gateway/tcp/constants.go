package tcpGateway

import "time"

/*
   Creation Time: 2019 - Feb - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	defaultReadTimout   = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
	defaultConnIdleTime = 5 * time.Minute
)
