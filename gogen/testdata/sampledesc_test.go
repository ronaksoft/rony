package sampleDesc_test

import (
	sampleDesc "git.ronaksoftware.com/ronak/rony/gogen/testdata"
	"github.com/gocql/gocql"
	"testing"
)

/*
   Creation Time: 2020 - Apr - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func TestCreateTables(t *testing.T) {
	cfg := gocql.NewCluster("river.ronaksoftware.com:9042")
	cfg.CQLVersion = "3.4.4"
	s, err := cfg.CreateSession()
	if err != nil {
		t.Fatal(err)
	}
	sampleDesc.Init(s)
	err = sampleDesc.CreateTables(s, "rony")
	if err != nil {
		t.Fatal(err)
	}
}
