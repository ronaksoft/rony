package z

import (
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

type TemplateArg struct {
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

func GetTemplateArg(f *protogen.File, gf *protogen.GeneratedFile, m *protogen.Message) TemplateArg {
	arg := TemplateArg{
		Name: string(m.Desc.Name()),
		C:    crc32.ChecksumIEEE([]byte(m.Desc.Name())),
	}
	for _, ft := range m.Fields {
		arg.AddField(f, gf, ft.Desc)
	}
	return arg
}
func (arg *TemplateArg) kind(f protoreflect.FieldDescriptor) string {
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

func (arg *TemplateArg) AddField(f *protogen.File, gf *protogen.GeneratedFile, fd protoreflect.FieldDescriptor) {
	name := string(fd.Name())
	pkg, typ := DescParts(f, gf, fd.Message())
	arg.Fields = append(arg.Fields, struct {
		Kind      string
		Pkg       string
		Type      string
		Name      string
		ZeroValue string
	}{Kind: arg.kind(fd), Pkg: pkg, Type: typ, Name: name, ZeroValue: ZeroValue(fd)})
}
