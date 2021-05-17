package gossipCluster

import (
	"time"
)

/*
   Creation Time: 2021 - Jan - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	gossipUpdateTimeout     = time.Second * 5
	gossipLeaveTimeout      = time.Second * 5
	clusterMessageRateLimit = 100
)
