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

type Arg struct {
	Prefix   string
	Postfix  string
	Name     string
	Pkg      string
	Type     string
	C        uint32
	DBKey    string
	HasIndex bool
	Fields   []struct {
		Kind             string
		Pkg              string
		Type             string
		Name             string
		ZeroValue        string
		DBKey            string
		DBKeyWithPostfix string
		DBPrefix         string
		HasIndex         bool
	}
	Views []struct {
		Keys  []string
		DBKey string
	}
}

func GetArg(g *Generator, m *protogen.Message, agg *Aggregate, prefix string) Arg {
	arg := Arg{
		Name:     string(m.Desc.Name()),
		C:        crc32.ChecksumIEEE([]byte(m.Desc.Name())),
		HasIndex: agg.HasIndex,
		Prefix:   prefix,
		Postfix:  "[idx]",
	}
	arg.DBKey = agg.Table.DBKey(arg.Prefix)
	arg.Pkg, arg.Type = z.DescParts(g.f, g.g, m.Desc)
	for _, ft := range m.Fields {
		name := string(ft.Desc.Name())
		pkg, typ := z.DescParts(g.f, g.g, ft.Desc.Message())
		opt, _ := ft.Desc.Options().(*descriptorpb.FieldOptions)
		arg.Fields = append(
			arg.Fields,
			struct {
				Kind             string
				Pkg              string
				Type             string
				Name             string
				ZeroValue        string
				DBKey            string
				DBKeyWithPostfix string
				DBPrefix         string
				HasIndex         bool
			}{
				Kind: arg.kind(ft.Desc),
				Pkg:  pkg, Type: typ,
				Name:             name,
				ZeroValue:        z.ZeroValue(ft.Desc),
				DBKey:            indexKey(agg, name, arg.Prefix, ""),
				DBKeyWithPostfix: indexKey(agg, name, arg.Prefix, arg.Postfix),
				DBPrefix:         indexPrefix(agg, name, tools.ToLowerCamel(inflection.Singular(name))),
				HasIndex:         proto.GetExtension(opt, rony.E_RonyIndex).(bool),
			},
		)
	}
	for idx, v := range agg.Views {
		arg.Views = append(arg.Views,
			struct {
				Keys  []string
				DBKey string
			}{
				Keys:  v.Keys(),
				DBKey: agg.Views[idx].DBKey(arg.Prefix),
			},
		)
	}
	return arg
}

func (arg *Arg) kind(f protoreflect.FieldDescriptor) string {
	switch f.Cardinality() {
	case protoreflect.Repeated:
		return "repeated"
	default:
		return ""
	}
}

func (arg *Arg) With(prefix, postfix string) *Arg {
	arg.Prefix = prefix
	arg.Postfix = postfix
	return arg
}
