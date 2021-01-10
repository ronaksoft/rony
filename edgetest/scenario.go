package edgetest

import (
	"time"
)

/*
   Creation Time: 2021 - Jan - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type scenario []*context

func newScenario() *scenario {
	return &scenario{}
}

func (s *scenario) Append(c *context) {
	*s = append(*s, c)
}

func (s scenario) Run(timeout time.Duration) error {
	for _, ctx := range s {
		err := ctx.Run(timeout)
		if err != nil {
			return err
		}
	}
	return nil
}
