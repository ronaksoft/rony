package edge

import (
	log "github.com/ronaksoft/rony/internal/logger"
	"time"
)

/*
   Creation Time: 2020 - Feb - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	raftApplyTimeout        = time.Second * 3
	gossipUpdateTimeout     = time.Second * 5
	gossipLeaveTimeout      = time.Second * 5
	clusterMessageRateLimit = 100
)

// Bridge Parameters
const (
	deliveryReportSubject = "D"
	notifierSubject       = "N"
	messageSubject        = "M"
	sendTimeout           = 5 * time.Second
	sendRetries           = 2
	flushPeriod           = time.Millisecond
	flushMaxSize          = 1000
	flushMaxWorkers       = 25

	defaultMaxWorkers = 1
)

func init() {
	log.Init(log.DefaultConfig)
}
