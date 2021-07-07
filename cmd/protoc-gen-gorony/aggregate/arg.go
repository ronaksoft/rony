package aggregate

import (
	"github.com/jinzhu/inflection"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"hash/crc32"
)

/*
   Creation Time: 2021 - Jul - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type ModelArg struct {
	keys     Key
	Name     string
	Pkg      string
	Type     string
	C        uint32
	HasIndex bool
	Views    []ViewArg
	Fields   []FieldArg
}

func (ma ModelArg) DBKey(prefix string) string {
	return ma.keys.DBKey(prefix)
}
func (ma ModelArg) FuncArgs(prefix string) string {
	return ma.keys.FuncArgs(prefix)
}
func (ma ModelArg) FuncArgsPKs(prefix string) string {
	return ma.keys.FuncArgsPKs(prefix)
}
func (ma ModelArg) FuncArgsCKs(prefix string) string {
	return ma.keys.FuncArgsCKs(prefix)
}
func (ma ModelArg) String(prefix string, sep string, lowerCamel bool) string {
	return ma.keys.String(prefix, sep, lowerCamel)
}
func (ma ModelArg) StringPKs(prefix string, sep string, lowerCamel bool) string {
	return ma.keys.StringPKs(prefix, sep, lowerCamel)
}
func (ma ModelArg) StringCKs(prefix string, sep string, lowerCamel bool) string {
	return ma.keys.StringCKs(prefix, sep, lowerCamel)
}

type ViewArg struct {
	keys Key
	Name string
	Keys []string
}

func (va ViewArg) DBKey(prefix string) string {
	return va.keys.DBKey(prefix)
}
func (va ViewArg) FuncArgs(prefix string) string {
	return va.keys.FuncArgs(prefix)
}
func (va ViewArg) FuncArgsPKs(prefix string) string {
	return va.keys.FuncArgsPKs(prefix)
}
func (va ViewArg) FuncArgsCKs(prefix string) string {
	return va.keys.FuncArgsCKs(prefix)
}
func (va ViewArg) String(prefix string, sep string, lowerCamel bool) string {
	return va.keys.String(prefix, sep, lowerCamel)
}
func (va ViewArg) StringPKs(prefix string, sep string, lowerCamel bool) string {
	return va.keys.StringPKs(prefix, sep, lowerCamel)
}
func (va ViewArg) StringCKs(prefix string, sep string, lowerCamel bool) string {
	return va.keys.StringCKs(prefix, sep, lowerCamel)
}

type FieldArg struct {
	agg       *Aggregate
	Kind      string
	Pkg       string
	Type      string
	Name      string
	ZeroValue string
	HasIndex  bool
}

func (fa FieldArg) DBKey(prefix, postfix string) string {
	return indexKey(fa.agg, fa.Name, prefix, postfix)
}

func (fa FieldArg) DBPrefixKey(prefix, postfix string) string {
	return indexPrefix(fa.agg, fa.Name, tools.ToLowerCamel(inflection.Singular(fa.Name)))
}

func GetArg(g *Generator, m *protogen.Message, agg *Aggregate) ModelArg {
	arg := ModelArg{
		Name:     string(m.Desc.Name()),
		C:        crc32.ChecksumIEEE([]byte(m.Desc.Name())),
		HasIndex: agg.HasIndex,
		keys:     agg.Table,
	}

	arg.Pkg, arg.Type = z.DescParts(g.f, g.g, m.Desc)
	for _, ft := range m.Fields {
		name := string(ft.Desc.Name())
		pkg, typ := z.DescParts(g.f, g.g, ft.Desc.Message())
		opt, _ := ft.Desc.Options().(*descriptorpb.FieldOptions)
		arg.Fields = append(
			arg.Fields,
			FieldArg{
				agg:  agg,
				Kind: arg.kind(ft.Desc),
				Pkg:  pkg, Type: typ,
				Name:      name,
				ZeroValue: z.ZeroValue(ft.Desc),
				HasIndex:  proto.GetExtension(opt, rony.E_RonyIndex).(bool),
			},
		)
	}
	for idx, v := range agg.Views {
		arg.Views = append(arg.Views,
			ViewArg{
				keys: agg.Views[idx],
				Name: agg.Name,
				Keys: v.Keys(),
			},
		)
	}
	return arg
}

func (ma *ModelArg) kind(f protoreflect.FieldDescriptor) string {
	switch f.Cardinality() {
	case protoreflect.Repeated:
		return "repeated"
	default:
		return ""
	}
}
