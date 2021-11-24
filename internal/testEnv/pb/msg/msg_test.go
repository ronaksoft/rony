package testmsg_test

import (
	"testing"

	testmsg "github.com/ronaksoft/rony/internal/testEnv/pb/msg"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMsg(t *testing.T) {
	Convey("Marshal/Unmarshal JSON", t, func(c C) {
		Convey("NormalEnvelope", func(c C) {
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
		Convey("RonyEnvelope", func(c C) {
			e1 := &testmsg.Envelope1{
				E1: "E1Value",
				Embed1: &testmsg.Embed1{
					F1: "Embed1F1Value",
				},
			}
			e1Bytes, err := e1.Marshal()
			c.So(err, ShouldBeNil)
			e2 := &testmsg.Envelope2{
				Constructor: testmsg.C_Envelope1,
				Message:     e1Bytes,
			}

			jsonData, err := e2.MarshalJSON()
			c.So(err, ShouldBeNil)
			c.Println(string(jsonData))

			e22 := &testmsg.Envelope2{}
			err = e22.UnmarshalJSON(jsonData)
			c.So(err, ShouldBeNil)
			c.So(e22.Constructor, ShouldEqual, e2.Constructor)
			c.So(e22.Message, ShouldEqual, e2.Message)
		})
	})
}
