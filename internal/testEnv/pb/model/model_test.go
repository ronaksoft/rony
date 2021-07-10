package model_test

import (
	"github.com/ronaksoft/rony/internal/store/localdb"
	"github.com/ronaksoft/rony/internal/testEnv/pb/model"
	"github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
	"os"
	"testing"
)

/*
   Creation Time: 2021 - Jul - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func TestMain(m *testing.M) {
	s, err := localdb.New(localdb.DefaultConfig("./_hdd"))
	if err != nil {
		panic(err)
	}

	store.Init(store.Config{
		DB: s.DB(),
	})
	code := m.Run()
	_ = os.RemoveAll("./_hdd")
	os.Exit(code)
}

func TestModel(t *testing.T) {
	Convey("Model - Auto-Generated Code Tests", t, func(c C) {
		Convey("Create/Read/Update/Read", func(c C) {
			m1 := &model.Model1{
				ID:       int32(tools.FastRand()),
				ShardKey: int32(tools.FastRand()),
				P1:       tools.RandomID(32),
				P2:       tools.RandomIDs(32, 43, 12, 10),
				P5:       tools.RandomUint64(0),
				Enum:     model.Enum_Something,
			}
			err := model.CreateModel1(m1)
			c.So(err, ShouldBeNil)
			m2, err := model.ReadModel1(m1.ID, m1.ShardKey, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m2), ShouldBeTrue)
			m3, err := model.ReadModel1ByEnumAndShardKeyAndID(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m3), ShouldBeTrue)
			m1.P1 = tools.RandomID(32)
			m1.Enum = model.Enum_Else
			err = model.UpdateModel1(m1.ID, m1.ShardKey, m1)
			c.So(err, ShouldBeNil)

			m5, err := model.ReadModel1(m1.ID, m1.ShardKey, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m5), ShouldBeTrue)

			m6, err := model.ReadModel1ByEnumAndShardKeyAndID(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m6), ShouldBeTrue)

			_, err = model.ReadModel1ByEnumAndShardKeyAndID(m2.Enum, m2.ShardKey, m2.ID, nil)
			c.So(err, ShouldNotBeNil)

			ms1, err := model.ListModel1ByP1(m1.P1, store.NewListOption(), nil)
			c.So(err, ShouldBeNil)
			c.So(ms1, ShouldHaveLength, 1)

			ms1, err = model.ListModel1ByP1(m2.P1, store.NewListOption(), nil)
			c.So(err, ShouldBeNil)
			c.So(ms1, ShouldHaveLength, 0)

			err = model.DeleteModel1(m2.ID, m2.ShardKey)
			c.So(err, ShouldBeNil)

			_, err = model.ReadModel1(m1.ID, m1.ShardKey, nil)
			c.So(err, ShouldNotBeNil)

			ms1, err = model.ListModel1ByP1(m1.P1, store.NewListOption(), nil)
			c.So(err, ShouldBeNil)
			c.So(ms1, ShouldHaveLength, 0)

			ms1, err = model.ListModel1ByP1(m5.P1, store.NewListOption(), nil)
			c.So(err, ShouldBeNil)
			c.So(ms1, ShouldHaveLength, 0)
		})
		Convey("List/Iter", func(c C) {
			total := int32(100)
			for i := int32(0) ; i < total ;i++ {
				m1 := &model.Model1{
					ID:       int32(tools.FastRand()),
					ShardKey: int32(tools.FastRand()),
					P1:       tools.RandomID(32),
					P2:       tools.RandomIDs(32, 43, 12, 10),
					P5:       tools.RandomUint64(0),
					Enum:     model.Enum_Something,
				}
				err := model.CreateModel1(m1)
				c.So(err, ShouldBeNil)
			}
			m1s, err := model.ListModel1(0, 0, store.NewListOption().SetLimit(50), nil, model.Model1OrderByEnum)
			c.So(err, ShouldBeNil)
			c.So(m1s, ShouldHaveLength, 50)
			m1s, err = model.ListModel1(0, 0, store.NewListOption().SetLimit(total * 2), nil, model.Model1OrderByEnum)
			c.So(err, ShouldBeNil)
			c.So(m1s, ShouldHaveLength, total)
		})
	})
}
