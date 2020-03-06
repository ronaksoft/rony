package rony

import (
	"errors"
	"fmt"
)

/*
   Creation Time: 2019 - Nov - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Error Codes
const (
	ErrCodeInternal         = "E00"
	ErrCodeInvalid          = "E01"
	ErrCodeUnavailable      = "E02"
	ErrCodeTooMany          = "E03"
	ErrCodeTooFew           = "E04"
	ErrCodeIncomplete       = "E05"
	ErrCodeTimeout          = "E06"
	ErrCodeAccess           = "E07"
	ErrCodeAlreadyExists    = "E08"
	ErrCodeBusy             = "E09"
	ErrCodeOutOfRange       = "E10"
	ErrCodePartiallyApplied = "E11"
	ErrCodeExpired          = "E12"
	ErrCodeNotImplemented   = "E13"
)

// Error Items
const (
	ErrItemServer     = "SERVER"
	ErrItemRaftLeader = "RAFT_LEADER"
	ErrItemHandler    = "HANDLER"
	ErrItemRequest    = "REQUEST"
)


var (
	ErrGatewayAlreadyInitialized = errors.New("gateway already initialized")
	ErrGatewayNotInitialized     = errors.New("there is no gateway defined")
	ErrNotFound                  = errors.New("not found")
	ErrNotRaftLeader             = errors.New("not raft leader")
	ErrRaftNotSet                = errors.New("raft not set")
)

func Wrap(txt string, err error) error {
	return errors.New(fmt.Sprintf("%s: %v", txt, err))
}