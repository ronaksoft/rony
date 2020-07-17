package pools

import "github.com/gobwas/pool/pbytes"

/*
   Creation Time: 2020 - Mar - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var TinyBytes = pbytes.New(4, 128)
var Bytes = pbytes.New(32, 64*(1<<10))
var LargeBytes = pbytes.New(32*(1<<10), 1<<20)
