package model_test

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/testEnv/pb/model"
	"github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/migrate"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
	"os"
	"testing"
	"time"
)

/*
   Creation Time: 2021 - Jul - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	sampleStore rony.Store
	scyllaDB    gocqlx.Session
)

func TestMain(m *testing.M) {
	s, err := store.New(store.DefaultConfig("./_hdd"))
	if err != nil {
		panic(err)
	}
	sampleStore = s

	cfg := gocql.NewCluster("127.0.0.1:9042")
	cfg.Timeout = time.Second * 3
	cfg.Keyspace = "rony"
	scyllaDB, err = gocqlx.WrapSession(cfg.CreateSession())
	if err != nil {
		panic(err)
	}

	err = migrate.FromFS(context.TODO(), scyllaDB, os.DirFS("./"))
	if err != nil {
		panic(err)
	}

	code := m.Run()
	_ = os.RemoveAll("./_hdd")
	os.Exit(code)
}

func TestModelLocalRepo(t *testing.T) {
	Convey("Model - LocalRepo Auto-Generated Code Tests", t, func(c C) {
		Convey("Create/Read/Update/Read", func(c C) {
			repo := model.NewModel1LocalRepo(sampleStore)
			m1 := &model.Model1{
				ID:       int32(tools.FastRand()),
				ShardKey: int32(tools.FastRand()),
				P1:       tools.RandomID(32),
				P2:       tools.RandomIDs(32, 43, 12, 10),
				P5:       tools.RandomUint64(0),
				Enum:     model.Enum_Something,
			}
			err := repo.Create(m1)
			c.So(err, ShouldBeNil)
			m2, err := repo.Read(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m2), ShouldBeTrue)
			m3, err := repo.ReadByCustomerSort(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m3), ShouldBeTrue)
			m1.P1 = tools.RandomID(32)
			err = repo.Update(m1.ID, m1.ShardKey, m1.Enum, m1)
			c.So(err, ShouldBeNil)

			m5, err := repo.Read(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m5), ShouldBeTrue)

			m6, err := repo.ReadByCustomerSort(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m6), ShouldBeTrue)

			_, err = repo.ReadByCustomerSort(m2.Enum, m2.ShardKey, m2.ID, nil)
			c.So(err, ShouldBeNil)

			err = repo.Delete(m2.ID, m2.ShardKey, m2.Enum)
			c.So(err, ShouldBeNil)

			_, err = repo.Read(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldNotBeNil)
		})
		Convey("List/Iter", func(c C) {
			repo := model.NewModel3LocalRepo(sampleStore)
			total := int32(100)
			start := int32(10)
			for i := int32(0); i < total; i++ {
				err := repo.Create(&model.Model3{
					ID:       int64(i + start),
					ShardKey: tools.RandomInt32(10),
					P1:       tools.S2B(tools.RandomID(32)),
					P2:       tools.RandomIDs(32, 43, 12, 10),
					P5:       nil,
				})
				c.So(err, ShouldBeNil)
			}

			res, err := repo.List(
				model.Model3PK{
					ID:       0,
					ShardKey: 0,
				},
				store.NewListOption().SetLimit(total/2),
				nil,
			)
			c.So(err, ShouldBeNil)
			c.So(res, ShouldHaveLength, total/2)
			c.So(res[0].ID, ShouldEqual, start)

			res, err = repo.List(
				model.Model3PK{
					ID:       int64(start + 10),
					ShardKey: 0,
				},
				store.NewListOption().SetLimit(total/2),
				nil,
			)
			c.So(err, ShouldBeNil)
			c.So(res, ShouldHaveLength, total/2)
			c.So(res[0].ID, ShouldEqual, start+10)

			res, err = repo.List(
				model.Model3PK{
					ID:       int64(start + 10),
					ShardKey: 4,
				},
				store.NewListOption().SetLimit(total/2),
				nil,
			)
			c.So(err, ShouldBeNil)
			c.So(res, ShouldHaveLength, total/2)
			c.So(res[0].ID, ShouldBeGreaterThanOrEqualTo, start+10)
			c.So(res[0].ID, ShouldBeLessThanOrEqualTo, start+11)

			res, err = repo.List(
				model.Model3PK{
					ID:       int64(start + 10),
					ShardKey: 4,
				},
				store.NewListOption().SetLimit(total/2).SetBackward(),
				nil,
			)
			c.So(err, ShouldBeNil)
			c.So(res[0].ID, ShouldBeGreaterThanOrEqualTo, start+9)
			c.So(res[0].ID, ShouldBeLessThanOrEqualTo, start+10)

			res, err = repo.List(
				nil,
				store.NewListOption().SetLimit(total*2).SetSkip(10),
				nil,
			)
			c.So(err, ShouldBeNil)
			c.So(res, ShouldHaveLength, total-10)
		})
		Convey("ListByIndex", func(c C) {
			repo := model.NewModel1LocalRepo(sampleStore)
			total := int32(100)
			start := int32(10)
			possibleValues := []string{"A", "B", "C", "D", "E", "F"}
			for i := int32(0); i < total; i++ {
				err := repo.Create(&model.Model1{
					ID:       i + start,
					ShardKey: tools.RandomInt32(10),
					P1:       possibleValues[tools.RandomInt(len(possibleValues))],
					P2: []string{
						possibleValues[tools.RandomInt(len(possibleValues))],
						possibleValues[tools.RandomInt(len(possibleValues))],
						possibleValues[tools.RandomInt(len(possibleValues))],
					},
					P5: tools.RandomUint64(0),
				})
				c.So(err, ShouldBeNil)
			}

			res, err := repo.ListByP1("A", store.NewListOption().SetLimit(total), nil)
			c.So(err, ShouldBeNil)
			for _, x := range res {
				c.So(x.P1, ShouldEqual, "A")
			}

			res, err = repo.ListByP2("A", store.NewListOption().SetLimit(total), nil)
			c.So(err, ShouldBeNil)
			for _, x := range res {
				c.So("A", ShouldBeIn, x.P2)
			}
		})
	})
}

func TestModelGlobalRepo(t *testing.T) {
	Convey("Model - GlobalRepo Auto-Generated Code Tests", t, func(c C) {
		Convey("Create/Read/Update/Read", func(c C) {
			repo := model.NewModel1GlobalRepo(scyllaDB)
			m1 := &model.Model1{
				ID:       int32(tools.FastRand()),
				ShardKey: int32(tools.FastRand()),
				P1:       tools.RandomID(32),
				P2:       tools.RandomIDs(32, 43, 12, 10),
				P5:       tools.RandomUint64(0),
				Enum:     model.Enum_Something,
			}
			err := repo.Insert(m1, false)
			c.So(err, ShouldBeNil)
			m2, err := repo.Get(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m2), ShouldBeTrue)
			m3, err := repo.GetByCustomerSort(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m3), ShouldBeTrue)
			m1.P1 = tools.RandomID(32)
			err = repo.Update(m1)
			c.So(err, ShouldBeNil)

			m5, err := repo.Get(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m5), ShouldBeTrue)

			m6, err := repo.GetByCustomerSort(m1.Enum, m1.ShardKey, m1.ID, nil)
			c.So(err, ShouldBeNil)
			c.So(proto.Equal(m1, m6), ShouldBeTrue)

			_, err = repo.GetByCustomerSort(m2.Enum, m2.ShardKey, m2.ID, nil)
			c.So(err, ShouldBeNil)

			err = repo.Delete(m2.ID, m2.ShardKey, m2.Enum)
			c.So(err, ShouldBeNil)

			_, err = repo.Get(m1.ID, m1.ShardKey, m1.Enum, nil)
			c.So(err, ShouldNotBeNil)
		})
		Convey("List/Iter", func(c C) {
			repo := model.NewModel3GlobalRepo(scyllaDB)
			total := int32(100)
			start := int32(10)
			for i := int32(0); i < total; i++ {
				err := repo.Insert(&model.Model3{
					ID:       int64(i + start),
					ShardKey: tools.RandomInt32(3),
					P1:       tools.S2B(tools.RandomID(32)),
					P2:       tools.RandomIDs(32, 43, 12, 10),
					P5:       nil,
				}, false)
				c.So(err, ShouldBeNil)
			}

			res, err := repo.List(
				model.Model3PK{
					ID:       int64(start),
					ShardKey: 1,
				},
				uint(total/2),
			)
			c.So(err, ShouldBeNil)
			for _, x := range res {
				c.So(x.ID, ShouldEqual, start)
				c.So(x.ShardKey, ShouldEqual, 1)
			}
		})
	})
}
