package gogen_test

import (
	"git.ronaksoftware.com/ronak/rony/gogen"
	"git.ronaksoftware.com/ronak/rony/gogen/plugins/proto2"
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
		err := gogen.LoadDescriptors("./extra/testdata")
		c.So(err, ShouldBeNil)

		err = gogen.Generate(&proto2.ProtoBuffer{}, "./extra/_output", "samplePackage", "", ".proto")
		c.So(err, ShouldBeNil)
	})
}
