package z

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
	p, t := DescName(file, g, desc)
	if p == "" {
		return fmt.Sprintf("C_%s", t)
	} else {
		return fmt.Sprintf("%s.C_%s", p, t)
	}
}

func Name(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) string {
	p, t := DescName(file, g, desc)
	if p == "" {
		return fmt.Sprintf("%s", t)
	} else {
		return fmt.Sprintf("%s.%s", p, t)
	}
}

func PackageName(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) string {
	pkg, _ := DescName(file, g, desc)
	return pkg
}

// descName returns the package and ident name
func DescName(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) (string, string) {
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

// zeroValue returns the equal zero value based on the input type
func ZeroValue(t protoreflect.Kind) string {
	switch t {
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