package replicateddb

import (
	"fmt"
)

/*
   Creation Time: 2021 - May - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	ErrDuplicateID = fmt.Errorf("duplicate id")
	ErrTxnNotFound = fmt.Errorf("txn not found")
	ErrUnknown     = fmt.Errorf("unknown")
)
