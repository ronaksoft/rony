package store

import (
	"fmt"
)

/*
   Creation Time: 2021 - Mar - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	ErrAlreadyExists = fmt.Errorf("already exists")
	ErrEmptyObject   = fmt.Errorf("empty object")
)
