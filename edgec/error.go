package edgec

import (
	"fmt"
)

/*
   Creation Time: 2020 - Jul - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	ErrTimeout        = fmt.Errorf("time out")
	ErrLeaderRedirect = fmt.Errorf("leader redirect")
)
