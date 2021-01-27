package kvgen_test

import (
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb/model"
	"github.com/ronaksoft/rony/repo/kv"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/*
   Creation Time: 2021 - Jan - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func init() {
	testEnv.Init()
}

func TestGenerate(t *testing.T) {
	Convey("Generate func", t, func(c C) {
		Convey("Save", save)
	})
}

func save(c C) {
	m1 := &model.Model1{
		ID:       int32(tools.RandomInt64(1000)),
		ShardKey: 100,
		P1:       tools.RandomID(10),
		P2:       []string{"1", "2", "3"},
		P5:       tools.RandomUint64(0),
		Enum:     model.Enum_Something,
	}
	m2 := &model.Model1{
		ID:       int32(tools.RandomInt64(1000)),
		ShardKey: 1001,
		P1:       tools.RandomID(10),
		P2:       []string{"3", "4"},
		P5:       tools.RandomUint64(0),
		Enum:     model.Enum_Something,
	}
	m3 := &model.Model1{
		ID:       int32(tools.RandomInt64(1000)),
		ShardKey: 1001,
		P1:       tools.RandomID(10),
		P2:       []string{"3", "4", "5"},
		P5:       tools.RandomUint64(0),
		Enum:     model.Enum_Else,
	}

	err := model.SaveModel1(m1)
	c.So(err, ShouldBeNil)
	err = model.SaveModel1(m2)
	c.So(err, ShouldBeNil)
	err = model.SaveModel1(m3)
	c.So(err, ShouldBeNil)
	res, err := model.ListModel1ByEnum(model.Enum_Else, kv.NewListOption().SetLimit(10))
	c.So(err, ShouldBeNil)
	c.So(res, ShouldHaveLength, 1)
	c.So(res[0].ID, ShouldEqual, m3.ID)

	res, err = model.ListModel1ByP2("3", kv.NewListOption().SetLimit(10))
	c.So(err, ShouldBeNil)
	c.So(res, ShouldHaveLength, 3)
	res, err = model.ListModel1ByP2("4", kv.NewListOption().SetLimit(10))
	c.So(err, ShouldBeNil)
	c.So(res, ShouldHaveLength, 2)
	res, err = model.ListModel1ByP2("5", kv.NewListOption().SetLimit(10))
	c.So(err, ShouldBeNil)
	c.So(res, ShouldHaveLength, 1)
	res, err = model.ListModel1ByP2("6", kv.NewListOption().SetLimit(10))
	c.So(err, ShouldBeNil)
	c.So(res, ShouldHaveLength, 0)

	m1p, err := model.ReadModel1ByEnumAndShardKeyAndID(m1.Enum, m1.ShardKey, m1.ID, nil)
	c.So(err, ShouldBeNil)
	c.So(m1p.P1, ShouldEqual, m1.P1)
	c.So(m1p.P2, ShouldResemble, m1.P2)

	err = model.DeleteModel1(m1.ID, m1.ShardKey)
	c.So(err, ShouldBeNil)

	m1p, err = model.ReadModel1ByEnumAndShardKeyAndID(m1.Enum, m1.ShardKey, m1.ID, nil)
	c.So(err, ShouldNotBeNil)
	c.So(m1p, ShouldBeNil)
}
