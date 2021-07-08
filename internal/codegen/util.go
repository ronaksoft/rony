package codegen

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

/*
   Creation Time: 2021 - Jan - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func Constructor(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) string {
	p, t := DescParts(file, g, desc)
	if p == "" {
		return fmt.Sprintf("C_%s", t)
	} else {
		return fmt.Sprintf("%s.C_%s", p, t)
	}
}

func Name(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) string {
	p, t := DescParts(file, g, desc)
	if p == "" {
		return t
	} else {
		return fmt.Sprintf("%s.%s", p, t)
	}
}

func PackageName(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) string {
	pkg, _ := DescParts(file, g, desc)
	return pkg
}

// DescParts returns the package and identifier name
func DescParts(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) (string, string) {
	if desc == nil {
		return "", ""
	}
	if string(desc.FullName().Parent()) == string(file.GoPackageName) {
		return "", string(desc.Name())
	} else {
		fd, ok := desc.ParentFile().Options().(*descriptorpb.FileOptions)
		if ok {
			g.QualifiedGoIdent(protogen.GoImportPath(fd.GetGoPackage()).Ident(string(desc.Name())))
		}
		return string(desc.ParentFile().Package()), string(desc.Name())
	}
}

// DescName return the fullname of the descriptor. i.e. package.identifier
func DescName(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) string {
	pkg, name := DescParts(file, g, desc)
	if pkg == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", pkg, name)
}

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

func GoKind(file *protogen.File, g *protogen.GeneratedFile, d protoreflect.FieldDescriptor) string {
	pkg := PackageName(file, g, d.Message())
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
		if pkg == "" {
			return string(d.Name())
		} else {
			return fmt.Sprintf("%s.%s", pkg, d.Name())
		}
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
	}
	return "unsupported"
	// panic(fmt.Sprintf("unsupported kindCql: %v", k.String()))
}
