package errors

import (
	"fmt"
	"sync"
)

/*
   Creation Time: 2019 - Oct - 20
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Multi struct {
	sync.Mutex
	errs []error
}

func (e *Multi) Error() string {
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}
	return fmt.Sprintf("Errors: %d, %v", len(e.errs), e.errs)
}

func (e *Multi) HasError() bool {
	return len(e.errs) > 0
}

func (e *Multi) AddError(err error) {
	e.Lock()
	e.errs = append(e.errs, err)
	e.Unlock()
}
