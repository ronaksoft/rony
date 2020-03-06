package rony

import (
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"time"
)

/*
   Creation Time: 2020 - Feb - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	raftApplyTimeout    = time.Second * 3
	gossipUpdateTimeout = time.Second * 5
	gossipLeaveTimeout  = time.Second * 5
)

func init() {
	log.InitLogger(log.DefaultConfig)
}
