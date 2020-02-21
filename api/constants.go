package api

import "time"

/*
   Creation Time: 2019 - Aug - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	longRequestThreshold = 500 * time.Millisecond
)

// Context Values
const (
	CtxAuthKey   = "AUTH_KEY"
	CtxServerSeq = "S_SEQ"
	CtxClientSeq = "C_SEQ"
	CtxUser      = "USER"
	CtxStreamID  = "SID"
	CtxTemp      = "TEMP"
)
