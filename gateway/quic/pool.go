package quicGateway

import "sync"

/*
   Creation Time: 2019 - Aug - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var waitGroupPool = sync.Pool{
	New: func() interface{} {
		return &sync.WaitGroup{}
	},
}
