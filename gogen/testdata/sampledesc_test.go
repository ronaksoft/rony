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

var (
	s *gocql.Session
)

func init() {
	var (
		err error
	)
	cfg := gocql.NewCluster("river.ronaksoftware.com:9042")
	cfg.CQLVersion = "3.4.4"
	s, err = cfg.CreateSession()
	if err != nil {
		panic(err)
	}

	sampleDesc.Init(s)
}

func TestCreateTables(t *testing.T) {
	err := sampleDesc.CreateTables(s, "rony")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveModel1(t *testing.T) {
	err := sampleDesc.SaveModel1(&sampleDesc.Model1{
		P1: "Ehsan",
		P2: []int32{1, 2, 3},
		P3: "Moosa",
	})
	if err != nil {
		t.Fatal(err)
	}

	m, err := sampleDesc.GetModel1("Ehsan", []int32{1, 2, 3}, "Moosa")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(m.P1, m.P2, m.P3)
}
