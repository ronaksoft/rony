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
	ErrAlreadyExists           = errors.New("already exists")
	ErrInvalidMediaType        = errors.New("invalid media type")
	ErrInvalidProvider         = errors.New("invalid provider")
	ErrPeerIsNil               = errors.New("peer is nil")
	ErrWrongCounterValue       = errors.New("wrong counter value")
	ErrEmpty                   = errors.New("empty key")
	ErrInvalidData             = errors.New("invalid data")
	ErrMaxPinnedDialogsReached = errors.New("max pinned dialogs reached")
	ErrNilUpdate               = errors.New("empty update")
	ErrMD5MisMatched           = errors.New("md5 mismatch")
	ErrInvalidClusterID        = errors.New("invalid cluster id")
	ErrFileTooLarge            = errors.New("file too large")
	ErrPartTooLarge            = errors.New("part is too large")
	ErrInvalidRange            = errors.New("invalid range")
	ErrNoStorageClient         = errors.New("no storage client")
	ErrWriteToClosedConn       = errors.New("write to closed conn")
	ErrWriteToFullBufferedConn = errors.New("write to full buffer conn")
	ErrUnexpectedSocketRead    = errors.New("unexpected read from socket")
	ErrSrpExpired              = errors.New("srp expired")
)

func Wrap(txt string, err error) error {
	return errors.New(fmt.Sprintf("%s: %v", txt, err))
}
