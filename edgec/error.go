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
	ErrNoConnection      = fmt.Errorf("no connection")
	ErrNoHost            = fmt.Errorf("no host")
	ErrTimeout           = fmt.Errorf("time out")
	ErrReplicaMaster     = fmt.Errorf("leader redirect")
	ErrReplicaSetSession = fmt.Errorf("replica redirect session")
	ErrReplicaSetRequest = fmt.Errorf("replica redirect request")
	ErrUnknownResponse   = fmt.Errorf("unknown response")
)
