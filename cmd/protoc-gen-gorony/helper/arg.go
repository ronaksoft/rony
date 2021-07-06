package helper

import (
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	Name   string
	Pkg    string
	Type   string
	C      uint32
	Fields []struct {
		Kind      string
		Pkg       string
		Type      string
		Name      string
		ZeroValue string
	}
}

func GetArg(g *Generator, m *protogen.Message) Arg {
	arg := Arg{
		Name: string(m.Desc.Name()),
		C:    crc32.ChecksumIEEE([]byte(m.Desc.Name())),
	}
	arg.Pkg, arg.Type = z.DescParts(g.f, g.g, m.Desc)
	for _, ft := range m.Fields {
		name := string(ft.Desc.Name())
		pkg, typ := z.DescParts(g.f, g.g, ft.Desc.Message())
		arg.Fields = append(arg.Fields, struct {
			Kind      string
			Pkg       string
			Type      string
			Name      string
			ZeroValue string
		}{Kind: arg.kind(ft.Desc), Pkg: pkg, Type: typ, Name: name, ZeroValue: z.ZeroValue(ft.Desc)})
	}
	return arg
}

func (arg *Arg) kind(f protoreflect.FieldDescriptor) string {
	switch f.Cardinality() {
	case protoreflect.Repeated:
		switch f.Kind() {
		case protoreflect.MessageKind:
			return "repeated/msg"
		case protoreflect.BytesKind:
			return "repeated/bytes"
		default:
			return "repeated"
		}
	default:
		switch f.Kind() {
		case protoreflect.MessageKind:
			return "msg"
		case protoreflect.BytesKind:
			return "bytes"
		default:
			return ""
		}
	}
}
