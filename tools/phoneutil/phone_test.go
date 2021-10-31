package phoneutil_test

import (
	"github.com/ronaksoft/rony/tools/phoneutil"
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
		"989121228718":   "989121228718",
		"+989121228718":  "989121228718",
		"9121228718":     "989121228718",
		"00989121228718": "989121228718",
	}

	for ph, cph := range phones {
		sph, err := phoneutil.SanitizePhone(ph, "IR")
		if sph != cph || err != nil {
			t.Fatal(err)
		}
	}

	phones = map[string]string{
		"989121228718":   "989121228718",
		"+989121228718":  "989121228718",
		"00989121228718": "989121228718",
	}

	for ph, cph := range phones {
		sph, err := phoneutil.SanitizePhone(ph, "")
		if sph != cph || err != nil {
			t.Log(ph, "->", sph)
			t.Fatal(err)
		}
	}
}
