package httpGateway_test

import (
	"fmt"
	"net"
	"testing"
)

/*
   Creation Time: 2020 - Apr - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func TestGateway_Addr(t *testing.T) {

	sts := []string{
		":2374", "0.0.0.0:2374", "127.0.0.1:2374", "10.2.3.43",
	}
	for _, s := range sts {
		ta, err := net.ResolveTCPAddr("tcp4", s)
		if err != nil {
			t.Fatal(s, err)
		}
		fmt.Println(ta.IP, ta.Port)
	}
}
