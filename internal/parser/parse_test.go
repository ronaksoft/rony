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

func TestParse(t *testing.T) {
	tr, err := Parse("PARSER", input)
	if err != nil {
		t.Fatal(err)
	}
	for idx, n := range tr.Root.Nodes {
		switch n.Type() {
		case NodeText:
			t.Log(idx, ": ", "Text", []byte(n.String()))
		default:
			t.Log(idx, ": ", n.String())

		}

	}
}
