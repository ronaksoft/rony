package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ronaksoft/rony/pools"
	"strings"
	"sync"
)

/*
   Creation Time: 2021 - Apr - 05
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type mErrors struct {
	mtx  sync.Mutex
	errs []error
}

func (errs *mErrors) Error() string {
	sb := strings.Builder{}
	for _, err := range errs.errs {
		sb.WriteString(err.Error())
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (errs *mErrors) Add(err error) {
	if err == nil {
		return
	}
	errs.mtx.Lock()
	errs.errs = append(errs.errs, err)
	errs.mtx.Unlock()
}

func (errs *mErrors) HasError() bool {
	if len(errs.errs) > 0 {
		return true
	}
	return false
}

func RunInSequence(fns ...func() error) error {
	for idx, fn := range fns {
		if err := fn(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("index=%d", idx))
		}
	}
	return nil
}

func RunInParallel(fns ...func() error) error {
	mErr := &mErrors{}
	waitGroup := pools.AcquireWaitGroup()
	for idx, fn := range fns {
		waitGroup.Add(1)
		go func(idx int, fn func() error) {
			mErr.Add(fn())
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
