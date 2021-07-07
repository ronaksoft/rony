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
	Name        string
	Pkg         string
	Type        string
	C           uint32
	HasIndex    bool
	DBKey       func(prefix string) string
	FuncArgs    func(prefix string) string
	FuncArgsPKs func(prefix string) string
	FuncArgsCKs func(prefix string) string
	String      func(prefix string, sep string, lowerCamel bool) string
	StringPKs   func(prefix string, sep string, lowerCamel bool) string
	StringCKs   func(prefix string, sep string, lowerCamel bool) string
	Views       []ViewArg
	Fields      []FieldArg
}
type ViewArg struct {
	Name        string
	Keys        []string
	DBKey       func(prefix string) string
	FuncArgs    func(prefix string) string
	FuncArgsPKs func(prefix string) string
	FuncArgsCKs func(prefix string) string
	String      func(prefix string, sep string, lowerCamel bool) string
	StringPKs   func(prefix string, sep string, lowerCamel bool) string
	StringCKs   func(prefix string, sep string, lowerCamel bool) string
}
type FieldArg struct {
	Kind      string
	Pkg       string
	Type      string
	Name      string
	ZeroValue string
	DBKey     func(pre, post string) string
	DBPrefix  func(pre, post string) string
	HasIndex  bool
}

func GetArg(g *Generator, m *protogen.Message, agg *Aggregate) ModelArg {
	arg := ModelArg{
		Name:     string(m.Desc.Name()),
		C:        crc32.ChecksumIEEE([]byte(m.Desc.Name())),
		HasIndex: agg.HasIndex,
		DBKey: func(prefix string) string {
			return agg.Table.DBKey(prefix)
		},
		FuncArgs: func(prefix string) string {
			return agg.Table.FuncArgs(prefix)
		},
		FuncArgsPKs: func(prefix string) string {
			return agg.Table.FuncArgsPKs(prefix)
		},
		FuncArgsCKs: func(prefix string) string {
			return agg.Table.FuncArgsCKs(prefix)
		},
		String: func(prefix, sep string, lowerCamel bool) string {
			return agg.Table.String(prefix, sep, lowerCamel)
		},
		StringPKs: func(prefix, sep string, lowerCamel bool) string {
			return agg.Table.StringPKs(prefix, sep, lowerCamel)
		},
		StringCKs: func(prefix, sep string, lowerCamel bool) string {
			return agg.Table.StringCKs(prefix, sep, lowerCamel)
		},
	}

	arg.Pkg, arg.Type = z.DescParts(g.f, g.g, m.Desc)
	for _, ft := range m.Fields {
		name := string(ft.Desc.Name())
		pkg, typ := z.DescParts(g.f, g.g, ft.Desc.Message())
		opt, _ := ft.Desc.Options().(*descriptorpb.FieldOptions)
		arg.Fields = append(
			arg.Fields,
			FieldArg{
				Kind: arg.kind(ft.Desc),
				Pkg:  pkg, Type: typ,
				Name:      name,
				ZeroValue: z.ZeroValue(ft.Desc),
				HasIndex:  proto.GetExtension(opt, rony.E_RonyIndex).(bool),
				DBKey: func(pre, post string) string {
					return indexKey(agg, name, pre, post)
				},
				DBPrefix: func(pre, post string) string {
					return indexPrefix(agg, name, tools.ToLowerCamel(inflection.Singular(name)))
				},
			},
		)
	}
	for idx, v := range agg.Views {
		arg.Views = append(arg.Views,
			ViewArg{
				Name: agg.Name,
				Keys: v.Keys(),
				DBKey: func(prefix string) string {
					return agg.Views[idx].DBKey(prefix)
				},
				FuncArgs: func(prefix string) string {
					return agg.Views[idx].FuncArgs(prefix)
				},
				FuncArgsPKs: func(prefix string) string {
					return agg.Views[idx].FuncArgsPKs(prefix)
				},
				FuncArgsCKs: func(prefix string) string {
					return agg.Views[idx].FuncArgsCKs(prefix)
				},
				String: func(prefix, sep string, lowerCamel bool) string {
					return agg.Views[idx].String(prefix, sep, lowerCamel)
				},
				StringPKs: func(prefix, sep string, lowerCamel bool) string {
					return agg.Views[idx].StringPKs(prefix, sep, lowerCamel)
				},
				StringCKs: func(prefix, sep string, lowerCamel bool) string {
					return agg.Views[idx].StringCKs(prefix, sep, lowerCamel)
				},
			},
		)
	}
	return arg
}

func (arg *ModelArg) kind(f protoreflect.FieldDescriptor) string {
	switch f.Cardinality() {
	case protoreflect.Repeated:
		return "repeated"
	default:
		return ""
	}
}
