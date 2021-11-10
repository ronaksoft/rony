package gopool

import (
	"github.com/panjf2000/ants/v2"
)

/*
   Creation Time: 2020 - Dec - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var goPool *ants.Pool

func Go(f func()) {
	_ = goPool.Submit(f)
}

func init() {
	p, err := ants.NewPool(-1)
	if err != nil {
		panic(err)
	}
	goPool = p
}
