package scylla

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gen"
	"github.com/iancoleman/strcase"
)

/*
   Creation Time: 2020 - Apr - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func snakeCaseAll(s []string) []string {
	o := make([]string, len(s))
	for idx := range s {
		o[idx] = strcase.ToSnake(s[idx])
	}
	return o
}

func validateModel(m *gen.Model) {
	primaryKeys := append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...)
	for _, fk := range m.FilterKeys {
		filterKeys := append(fk.PartitionKeys, fk.ClusteringKeys...)
		if len(filterKeys) != len(primaryKeys) {
			panic("filter keys must have all the primary keys")
		}
		for _, k := range filterKeys {
			p, err := m.GetProperty(k)
			if err != nil {
				panic("invalid filter key")
			}
			if p.CheckOption(gen.Slice) {
				panic("primary and/or filter keys must not be sliced type")
			}
			if !inArray(primaryKeys, k) {
				panic(fmt.Sprintf("key %s(%s) does not exists in the PrimaryKey for model (%s)", fk.Name, k, m))
			}
		}
	}
}

func inArray(arr []string, v string) bool {
	for i := range arr {
		if arr[i] == v {
			return true
		}
	}
	return false
}
