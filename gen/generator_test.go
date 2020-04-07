package gen_test

import (
	"git.ronaksoftware.com/ronak/rony/gen"
	"git.ronaksoftware.com/ronak/rony/gen/plugins/proto2"
	"git.ronaksoftware.com/ronak/rony/gen/plugins/scylla"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"os/exec"
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
		err := gen.LoadDescriptors("./testdata")
		c.So(err, ShouldBeNil)

		// 1. Generate ProtoBuffer Files
		err = gen.Generate(&proto2.ProtoBuffer{}, "./testdata", "samplePackage", "", ".proto")
		c.So(err, ShouldBeNil)

		// 2. Generate Proto Files in Golang
		cmd := exec.Command("/usr/local/bin/protoc", "-I", "./testdata", "-I", "../", "--gogofaster_out=./testdata", "sampledesc.proto")
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		err = cmd.Run()
		c.So(err, ShouldBeNil)

		// 3. Generate Scylla Init
		err = gen.Generate(&scylla.RepoPlugin{}, "./testdata", "samplePackage", "", "_init.go")
		c.So(err, ShouldBeNil)
	})
}
