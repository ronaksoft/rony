package edge

import (
	"github.com/pkg/errors"
)

/*
   Creation Time: 2021 - Mar - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	ErrClusterNotSet            = errors.New("cluster is not set")
	ErrGatewayNotSet            = errors.New("gateway is not set")
	ErrTunnelNotSet             = errors.New("tunnel is not set")
	ErrUnexpectedTunnelResponse = errors.New("unexpected tunnel response")
	ErrEmptyMemberList          = errors.New("member list is empty")
	ErrMemberNotFound           = errors.New("member not found")
	ErrNoTunnelAddrs            = errors.New("tunnel address does not found")
)
