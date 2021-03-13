package tcpGateway

import "errors"

/*
   Creation Time: 2020 - Mar - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	ErrUnsupportedProtocol  = errors.New("gateway protocol is not supported")
	ErrWriteToClosedConn    = errors.New("write to closed conn")
	ErrConnectionClosed     = errors.New("connection closed")
	ErrUnexpectedSocketRead = errors.New("unexpected socket read")
	ErrOpCloseReceived      = errors.New("close operation received")
)
