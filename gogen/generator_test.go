package gogen_test

import (
	"git.ronaksoftware.com/ronak/rony/gogen"
	"git.ronaksoftware.com/ronak/rony/gogen/plugins/proto2"
	"git.ronaksoftware.com/ronak/rony/gogen/plugins/scylla"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/*
   Creation Time: 2020 - Apr - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func TestGenModel(t *testing.T) {
	Convey("Generate Model", t, func(c C) {
		err := gogen.LoadDescriptors("./testdata")
		c.So(err, ShouldBeNil)

		// 1. Generate ProtoBuffer Files
		err = gogen.Generate(&proto2.ProtoBuffer{}, "./testdata", "samplePackage", "", ".proto")
		c.So(err, ShouldBeNil)

		// 2. Generate Proto Files in Golang

		// 3. Generate Scylla Init
		err = gogen.Generate(&scylla.RepoPlugin{}, "./testdata", "samplePackage", "", "_init.go")
		c.So(err, ShouldBeNil)
	})
}
