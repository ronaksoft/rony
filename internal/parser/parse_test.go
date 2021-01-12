package parse

import (
	"testing"
)

/*
   Creation Time: 2020 - Aug - 20
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var input = `
{{ @entity cql }}
{{ @tab (x1, x2, x3) }}
{{ @view (x3, x1, x2) }}
{{ @view ((x1, x2), x3) }}
{{ @cnt x1 }}
`

func TestLexer(t *testing.T) {
	l := lex("lex1", input)
	for {
		i := l.nextItem()
		if i.tok == ERROR {
			break
		}
	}
}

func TestParse(t *testing.T) {
	_, err := Parse("PARSER", input)
	if err != nil {
		t.Fatal(err)
	}
	// for idx, n := range tr.Root.Nodes {
	// 	switch n.Type() {
	// 	case NodeText:
	// 		// t.Log(idx, ": ", "Text", []byte(n.String()))
	// 	default:
	// 		// t.Log(idx, ": ", n.String())
	//
	// 	}
	//
	// }
}
