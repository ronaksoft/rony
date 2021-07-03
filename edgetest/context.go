package edgetest

import (
	"time"
)

/*
   Creation Time: 2021 - Jul - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type context interface {
	Run(timeout time.Duration) error
}
