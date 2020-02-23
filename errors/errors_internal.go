package errors

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
	ErrBridgeAlreadyInitialized  = errors.New("bridge already initialized")
	ErrEmpty                     = errors.New("empty key")
	ErrInvalidData               = errors.New("invalid data")
	ErrConstructorNotHandled     = errors.New("constructor not handled")
	ErrNotRaftLeader             = errors.New("not raft leader")
	ErrRaftNotSet                = errors.New("raft not set")
	ErrGossipNotSet              = errors.New("gossip not set")
	ErrWriteToClosedConn         = errors.New("write to closed conn")
	ErrWriteToFullBufferedConn   = errors.New("write to full buffer conn")
	ErrUnexpectedSocketRead      = errors.New("unexpected read from socket")
)

func Wrap(txt string, err error) error {
	return errors.New(fmt.Sprintf("%s: %v", txt, err))
}
