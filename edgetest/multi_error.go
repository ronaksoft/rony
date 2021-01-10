package edgetest

import (
	"fmt"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2021 - Jan - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type multiError struct {
	tools.SpinLock
	errs []error
}

func (e *multiError) Error() string {
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}
	return fmt.Sprintf("Errors: %d, %v", len(e.errs), e.errs)
}

func (e *multiError) HasError() bool {
	return len(e.errs) > 0
}

func (e *multiError) AddError(err error) {
	if err == nil {
		return
	}
	e.Lock()
	e.errs = append(e.errs, err)
	e.Unlock()
}
