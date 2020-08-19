package parser

import (
	"fmt"
	"testing"
)

/*
   Creation Time: 2020 - Aug - 18
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var input = `
{{@model}}
{{@tab ((x1, x2), x3)}}
{{@view (x3, x1, x2)}}
`

func TestLexer(t *testing.T) {
	l := lex("lex1", input, "", "")
	for {
		i := l.nextItem()
		if i.typ == ERROR {
			break
		}
		fmt.Println(i.String())
	}

}
