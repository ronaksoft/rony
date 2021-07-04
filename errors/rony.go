package errors

import (
	"fmt"
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
)

var (
	ErrClusterNotSet             = fmt.Errorf("cluster is not set")
	ErrGatewayNotSet             = fmt.Errorf("gateway is not set")
	ErrTunnelNotSet              = fmt.Errorf("tunnel is not set")
	ErrUnexpectedResponse        = fmt.Errorf("unexpected response")
	ErrUnexpectedTunnelResponse  = fmt.Errorf("unexpected tunnel response")
	ErrMemberNotFound            = fmt.Errorf("member not found")
	ErrGatewayAlreadyInitialized = fmt.Errorf("gateway already initialized")
	ErrNoTunnelAddrs             = fmt.Errorf("tunnel address does not found")
	ErrRetriesExceeded           = Wrap("maximum retries exceeded")
	ErrConnectionNotExists       = fmt.Errorf("connection does not exists")
)
