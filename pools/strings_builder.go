package pools

import (
	"strings"
	"sync"
)

/*
   Creation Time: 2020 - May - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var stringsBuilderPool sync.Pool

func AcquireStringsBuilder() *strings.Builder {
	sb := stringsBuilderPool.Get()
	if sb == nil {
		return &strings.Builder{}
	}

	return sb.(*strings.Builder)
}

func ReleaseStringsBuilder(sb *strings.Builder) {
	sb.Reset()
	stringsBuilderPool.Put(sb)
}
