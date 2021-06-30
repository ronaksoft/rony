package tools

import (
	"fmt"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/pools"
)

/*
   Creation Time: 2021 - Apr - 05
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/


func RunInSequence(fns ...func() error) error {
	for idx, fn := range fns {
		if err := fn(); err != nil {
			return errors.Wrap(fmt.Sprintf("index=%d", idx))(err)
		}
	}
	return nil
}

func RunInParallel(fns ...func() error) error {
	mErr := &errors.Multi{}
	waitGroup := pools.AcquireWaitGroup()
	for idx, fn := range fns {
		waitGroup.Add(1)
		go func(idx int, fn func() error) {
			mErr.AddError(fn())
			waitGroup.Done()
		}(idx, fn)
	}
	waitGroup.Wait()
	pools.ReleaseWaitGroup(waitGroup)
	if mErr.HasError() {
		return mErr
	}
	return nil
}
