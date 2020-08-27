package registry

import (
	"fmt"
)

/*
   Creation Time: 2020 - Aug - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var constructors = map[int64]string{}

func RegisterConstructor(c int64, n string) {
	if old, ok := constructors[c]; ok {
		panic(fmt.Sprintf("constructor already exists %s:%s", old, n))
	} else {
		constructors[c] = n
	}
}

func ConstructorName(c int64) string {
	return constructors[c]
}
