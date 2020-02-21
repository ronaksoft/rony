package natsBridge

import "time"

/*
   Creation Time: 2019 - Feb - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// BackGateway Constants
const (
	deliveryReportSubject = "D"
	notifierSubject       = "N"
	messageSubject        = "M"

	sendTimeout     = 5 * time.Second
	sendRetries     = 2
	flushPeriod     = 5 * time.Millisecond
	flushMaxSize    = 250
	flushMaxWorkers = 100

	defaultMaxWorkers = 1
)
