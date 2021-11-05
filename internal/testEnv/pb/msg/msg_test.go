package testmsg_test

import (
	testmsg "github.com/ronaksoft/rony/internal/testEnv/pb/msg"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMsg(t *testing.T) {
	Convey("Marshal/Unmarshal JSON", t, func(c C) {
		e1 := &testmsg.Envelope1{
			E1: "E1Value",
			Embed1: &testmsg.Embed1{
				F1: "Embed1F1Value",
			},
		}
		jsonData, err := e1.MarshalJSON()
		c.So(err, ShouldBeNil)
		c.Println(string(jsonData))
	})
}
