package rony

import (
	"hash/crc64"
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

var AccessHashTable = crc64.MakeTable(0x23740630002374)
