package sampleDesc_test

import (
	sampleDesc "git.ronaksoftware.com/ronak/rony/gen/testdata"
	"github.com/gocql/gocql"
	. "github.com/smartystreets/goconvey/convey"
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
	Convey("CreateTables", t, func(c C) {
		cfg := gocql.NewCluster("river.ronaksoftware.com:9042")
		cfg.CQLVersion = "3.4.4"
		s, err := cfg.CreateSession()
		c.So(err, ShouldBeNil)

		err = sampleDesc.CreateTables(s, "rony")
		c.So(err, ShouldBeNil)
		s.Close()
	})

}

func TestSaveModel1(t *testing.T) {
	Convey("Save/Get Model", t, func(c C) {
		cfg := gocql.NewCluster("river.ronaksoftware.com:9042")
		cfg.CQLVersion = "3.4.4"
		cfg.Keyspace = "rony"
		s, err := cfg.CreateSession()
		c.So(err, ShouldBeNil)

		sampleDesc.Init(s)
		err = sampleDesc.SaveModel1(&sampleDesc.Model1{
			P1: "Ehsan",
			P2: []int32{1, 2, 3},
			P3: "Moosa",
		})
		c.So(err, ShouldBeNil)

		m, err := sampleDesc.GetModel1("Ehsan", "Moosa")
		c.So(err, ShouldBeNil)
		c.So(m.P1, ShouldEqual, "Ehsan")
		c.So(m.P2, ShouldHaveLength, 3)
		c.So(m.P3, ShouldEqual, "Moosa")

		m, err = sampleDesc.GetModel1ByP3("Moosa", "Ehsan")
		c.So(err, ShouldBeNil)
		c.So(m.P1, ShouldEqual, "Ehsan")
		c.So(m.P2, ShouldHaveLength, 3)
		c.So(m.P3, ShouldEqual, "Moosa")
		s.Close()
	})
}
