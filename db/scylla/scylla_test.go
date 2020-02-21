package scylla_test

/*
   Creation Time: 2019 - Nov - 25
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// TODO:: create test_tables before uncommenting the test
// func init() {
// 	testEnv.Init()
// }
//
// type testItem struct {
// 	ID    int64  `db:"id"`
// 	Name  string `db:"name"`
// 	Value string `db:"value"`
// }
//
// func insertUser() {
// 	qp := scylla.NewQueryPool(func() *gocqlx.Queryx {
// 		stmt, names := qb.Insert("test_table").
// 			Columns("id", "name", "value").ToCql()
// 		return gocqlx.Query(scylla.Session().Query(stmt), names)
// 	})
// 	for i := int64(1); i < 1000; i++ {
// 		q := qp.GetQuery()
// 		q.Bind(i, fmt.Sprintf("Name %d", i), fmt.Sprintf("Value %d", i))
// 		err := scylla.Exec(q)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		qp.Put(q)
// 	}
// }

// func TestQueryPool(t *testing.T) {
// 	fmt.Println("Insert To DB")
// 	insertUser()
// 	fmt.Println("Let Get from DB")
// 	qp1 := scylla.NewQueryPool(func() *gocqlx.Queryx {
// 		stmt, names := qb.Select("test_table").
// 			Columns("id", "name", "value").
// 			Where(qb.Eq("id")).
// 			Limit(1).
// 			ToCql()
// 		return gocqlx.Query(scylla.Session().Query(stmt), names)
// 	})
// 	qp2 := scylla.NewQueryPool(func() *gocqlx.Queryx {
// 		stmt, names := qb.Insert("test_table").
// 			Columns("id", "name", "value").ToCql()
// 		return gocqlx.Query(scylla.Session().Query(stmt), names)
// 	})
//
// 	RandomInt := func(n int) int {
// 		b := make([]byte, 100)
// 		_, _ = rand.Read(b)
// 		return int(crc32.ChecksumIEEE(b)) % n
// 	}
// 	cnt := int64(0)
// 	wg := sync.WaitGroup{}
// 	for i := 0; i < 100; i++ {
// 		wg.Add(1)
// 		go func() {
// 			for j := 0; j < 100; j++ {
// 				v := &testItem{}
// 				i := RandomInt(9000) + 1
// 				q := qp1.GetQuery()
// 				q.Bind(i)
// 				err := scylla.Get(q, v)
// 				qp1.Put(q)
// 				if err == nil {
// 					if v.Name != fmt.Sprintf("Name %d", v.ID) || v.Value != fmt.Sprintf("Value %d", v.ID) {
// 						fmt.Println("MisMatched", i, v)
// 					}
// 				}
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go func() {
// 			for j := int64(1); j < 10000; j++ {
// 				q := qp2.GetQuery()
// 				x := tools.RandomInt(0)
// 				q.Bind(x, fmt.Sprintf("Name %d", x), fmt.Sprintf("Value %d", x))
// 				_ = scylla.Exec(q)
// 				qp2.Put(q)
// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// 	fmt.Println("Errors", cnt)
// }
