package codegen

/*
   Creation Time: 2021 - Jul - 08
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Model interface {
	Name() string
	Table() ModelKey
	Views() []ModelKey
}

type IModelKey interface {
	PartitionKeys() []Prop
	ClusteringKeys() []Prop
	Keys() []Prop
	StringNameTypes(namePrefix string, filter PropFilter) string
	StringNames(prefix string, sep string, nameCase TextCase) string
}
