package codegen

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"hash/crc32"
	"hash/crc64"
	"strings"
	"text/template"
)

/*
   Creation Time: 2021 - Jan - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	CrcBits = 64
	CrcTab  = crc64.MakeTable(crc64.ISO)
)

// ZeroValue returns the equal zero value based on the input type
func ZeroValue(f protoreflect.FieldDescriptor) string {
	switch f.Kind() {
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.StringKind:
		return "\"\""
	case protoreflect.MessageKind:
		return "nil"
	case protoreflect.BytesKind:
		return "nil"
	default:
		return "0"
	}
}

func GoKind(d protoreflect.FieldDescriptor) string {
	switch d.Kind() {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "int64"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "uint64"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "[]byte"
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.EnumKind:
		return string(d.Enum().Name())
	}

	return "unsupported"
}

func CqlKind(d protoreflect.FieldDescriptor) string {
	switch d.Kind() {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "int"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "bigint"
	case protoreflect.DoubleKind:
		return "double"
	case protoreflect.FloatKind:
		return "float"
	case protoreflect.BytesKind, protoreflect.StringKind:
		return "blob"
	case protoreflect.BoolKind:
		return "boolean"
	case protoreflect.EnumKind:
		return "int"
	}

	return "unsupported"
}

func CrcHash(data []byte) uint64 {
	if CrcBits == 64 {
		return crc64.Checksum(data, CrcTab)
	} else {
		return uint64(crc32.Checksum(data, crc32.IEEETable))
	}
}

func ExecTemplate(t *template.Template, v interface{}) string {
	sb := &strings.Builder{}
	if err := t.Execute(sb, v); err != nil {
		panic(err)
	}
	return sb.String()
}
