package z

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"hash/crc32"
	"strings"
)

/*
   Creation Time: 2021 - Jul - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type MessageArg struct {
	Fullname string
	Name     string
	NameCC   string // LowerCamelCase(Name)
	CName    string // ConstructorName([Pkg.]C_Name)
	Pkg      string
	C        uint32
	Fields   []FieldArg
}

func GetMessageArg(file *protogen.File, gFile *protogen.GeneratedFile, m *protogen.Message) MessageArg {
	arg := MessageArg{}
	arg.Pkg, arg.Name = DescParts(file, gFile, m.Desc)
	arg.NameCC = tools.ToLowerCamel(arg.Name)
	arg.C = crc32.ChecksumIEEE([]byte(m.Desc.Name()))
	if arg.Pkg == "" {
		arg.CName = fmt.Sprintf("C_%s", arg.Name)
		arg.Fullname = arg.Name
	} else {
		arg.CName = fmt.Sprintf("%s.C_%s", arg.Pkg, arg.Name)
		arg.Fullname = fmt.Sprintf("%s.%s", arg.Pkg, arg.Name)
	}
	for _, f := range m.Fields {
		arg.Fields = append(arg.Fields, GetFieldArg(file, gFile, f))
	}
	return arg
}

type FieldArg struct {
	Name        string // Name of the field
	NameCC      string // LowerCamelCase(Name)
	Pkg         string
	Type        string
	ZeroValue   string
	Kind        string
	GoKind      string
	Cardinality string
	desc        protoreflect.FieldDescriptor
}

func GetFieldArg(file *protogen.File, gFile *protogen.GeneratedFile, f *protogen.Field) FieldArg {
	arg := FieldArg{}
	arg.Name = f.GoName
	arg.NameCC = tools.ToLowerCamel(arg.Name)
	arg.Pkg, arg.Type = DescParts(file, gFile, f.Desc.Message())
	arg.Kind = f.Desc.Kind().String()
	arg.GoKind = GoKind(file, gFile, f.Desc)
	arg.Cardinality = f.Desc.Cardinality().String()
	arg.ZeroValue = ZeroValue(f.Desc)
	return arg
}

type ModelArg struct {
	Message  MessageArg
	HasIndex bool
	Table    DBArg
	Views    []DBArg
	Fields   []FieldArg
}

type ModelFieldArg struct {
	FieldArg
	HasIndex   bool
	DBIndexKey func(pre, post string) string
}

type DBArg struct {
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

type ServiceArg struct {
	Name    string
	NameCC  string // LowerCamelCase(Name)
	NameKC  string // KebabCase(Name)
	C       uint32
	Methods []MethodArg
}

func GetServiceArg(file *protogen.File, gFile *protogen.GeneratedFile, s *protogen.Service) ServiceArg {
	arg := ServiceArg{}
	arg.Name = s.GoName
	arg.NameCC = tools.ToLowerCamel(arg.Name)
	arg.NameKC = tools.ToKebab(arg.Name)
	arg.C = crc32.ChecksumIEEE([]byte(arg.Name))
	for _, m := range s.Methods {
		arg.Methods = append(arg.Methods, GetMethodArg(file, gFile, m))
	}
	return arg
}

type MethodArg struct {
	Name        string
	NameCC      string // LowerCamelCase(Name)
	NameKC      string // KebabCase(Name)
	Input       MessageArg
	Output      MessageArg
	RestEnabled bool
	TunnelOnly  bool
	Rest        struct {
		Method string
		Path   string
		Json   bool
	}
}

func GetMethodArg(file *protogen.File, gFile *protogen.GeneratedFile, m *protogen.Method) MethodArg {
	arg := MethodArg{}
	arg.Name = m.GoName
	arg.NameCC = tools.ToLowerCamel(arg.Name)
	arg.NameKC = tools.ToKebab(arg.Name)
	arg.Input = GetMessageArg(file, gFile, m.Input)
	arg.Output = GetMessageArg(file, gFile, m.Output)
	opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
	restOpt := proto.GetExtension(opt, rony.E_RonyRest).(*rony.RestOpt)
	rest := struct {
		Method string
		Path   string
		Json   bool
	}{Method: "", Path: "", Json: false}
	if restOpt != nil {
		rest.Method = restOpt.GetMethod()
		rest.Path = fmt.Sprintf("/%s", strings.Trim(restOpt.GetPath(), "/"))
		rest.Json = restOpt.GetJsonEncode()
	}
	arg.RestEnabled = restOpt != nil
	arg.TunnelOnly = proto.GetExtension(opt, rony.E_RonyInternal).(bool)
	arg.Rest = rest
	return arg
}
