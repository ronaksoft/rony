package model_test

import (
	"github.com/ronaksoft/rony/internal/store/localdb"
	"github.com/ronaksoft/rony/internal/testEnv/pb/model"
	"github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
	"github.com/scylladb/gocqlx/v2"
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

func TestModelLocalRepo(t *testing.T) {
	Convey("Model - LocalRepo Auto-Generated Code Tests", t, func(c C) {
		SkipConvey("Create/Read/Update/Read", func(c C) {
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
			m2, err := model.ReadModel1(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m2), ShouldBeTrue)
			m3, err := model.ReadModel1ByEnumAndShardKeyAndID(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m3), ShouldBeTrue)
			m1.P1 = tools.RandomID(32)
			err = model.UpdateModel1(m1.ID, m1.ShardKey, m1.Enum, m1)
			c.So(err, ShouldBeNil)

			m5, err := model.ReadModel1(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m5), ShouldBeTrue)

			m6, err := model.ReadModel1ByEnumAndShardKeyAndID(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m6), ShouldBeTrue)

			_, err = model.ReadModel1ByEnumAndShardKeyAndID(m2.Enum, m2.ShardKey, m2.ID, nil)
			c.So(err, ShouldBeNil)

			err = model.DeleteModel1(m2.ID, m2.ShardKey, m2.Enum)
			c.So(err, ShouldBeNil)

			_, err = model.ReadModel1(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldNotBeNil)
		})
	})
}

func TestModelRemoteRepo(t *testing.T) {
	Convey("Model - RemoteRep Auto-Generated Code Tests", t, func(c C) {
		SkipConvey("Create/Read/Update/Read", func(c C) {
			m1 := &model.Model1{
				ID:       int32(tools.FastRand()),
				ShardKey: int32(tools.FastRand()),
				P1:       tools.RandomID(32),
				P2:       tools.RandomIDs(32, 43, 12, 10),
				P5:       tools.RandomUint64(0),
				Enum:     model.Enum_Something,
			}
			repo1 := model.NewModel1RemoteRepo(gocqlx.Session{})
			err := repo1.Insert(m1, false)

			c.So(err, ShouldBeNil)
			m2, err := model.ReadModel1(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m2), ShouldBeTrue)
			m3, err := model.ReadModel1ByEnumAndShardKeyAndID(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m3), ShouldBeTrue)
			m1.P1 = tools.RandomID(32)
			err = model.UpdateModel1(m1.ID, m1.ShardKey, m1.Enum, m1)
			c.So(err, ShouldBeNil)

			m5, err := model.ReadModel1(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m5), ShouldBeTrue)

			m6, err := model.ReadModel1ByEnumAndShardKeyAndID(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m6), ShouldBeTrue)

			_, err = model.ReadModel1ByEnumAndShardKeyAndID(m2.Enum, m2.ShardKey, m2.ID, nil)
			c.So(err, ShouldBeNil)

			err = model.DeleteModel1(m2.ID, m2.ShardKey, m2.Enum)
			c.So(err, ShouldBeNil)

			_, err = model.ReadModel1(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldNotBeNil)
		})
	})
}
