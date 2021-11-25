package errors

import (
	"fmt"
)

/*
   Creation Time: 2021 - May - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func WrapText(txt string) func(err error) error {
	return func(err error) error {
		if err == nil {
			return nil
		}

		return fmt.Errorf("%s: %v", txt, err)
	}
}

func Wrap(wrapper error) func(err error) error {
	return func(err error) error {
		return fmt.Errorf("%s: %v", wrapper.Error(), err)
	}
}
