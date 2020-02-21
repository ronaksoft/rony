package natsBridge

import "sync"

/*
   Creation Time: 2019 - Feb - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	_MessagePool = sync.Pool{
		New: func() interface{} {
			return new(Message)
		},
	}
)
