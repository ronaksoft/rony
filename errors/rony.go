package errors

import (
	"github.com/pkg/errors"
)

/*
   Creation Time: 2021 - May - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	ErrInvalidRequest     = GenInvalidErr("REQUEST", nil)
	ErrInternalServer     = GenInternalErr("SERVER", nil)
	ErrInvalidHandler     = GenInvalidErr("HANDLER", nil)
	ErrUnavailableRequest = GenUnavailableErr("REQUEST", nil)

	ErrGatewayAlreadyInitialized = errors.New("gateway already initialized")
	ErrRetriesExceeded           = Wrap("maximum retries exceeded")
)
