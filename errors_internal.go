package rony

import (
	"errors"
	"fmt"
)

/*
   Creation Time: 2019 - Oct - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

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
