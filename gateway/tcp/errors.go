package tcp

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
	ErrWriteToClosedConn    = errors.New("write to closed conn")
	ErrUnexpectedSocketRead = errors.New("unexpected socket read")
)
