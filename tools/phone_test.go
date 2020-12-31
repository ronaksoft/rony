package tools

import (
	"testing"
)

/*
   Creation Time: 2021 - Jan - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func TestSanitizePhone(t *testing.T) {
	phones := map[string]string{
		"989121228718":  "989121228718",
		"+989121228718": "989121228718",
		"9121228718":    "989121228718",
	}

	for ph, cph := range phones {
		sph := SanitizePhone(ph, "IR")
		if sph != cph {
			t.Fatal()
		}
	}
}
