package codegen

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

// MessageArg holds the data needed by the template engine to generate code based on the protogen.Message
type MessageArg struct {
	desc     protoreflect.MessageDescriptor
	Fullname string
	Name     string
	NameCC   string // LowerCamelCase(Name)
	CName    string // [Pkg.]C_Name
	Pkg      string
	C        uint32
	Fields   []FieldArg
}

func GetMessageArg(file *protogen.File, gFile *protogen.GeneratedFile, m *protogen.Message) MessageArg {
	arg := MessageArg{
		desc: m.Desc,
	}
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

// FieldArg holds the data needed by the template engine to generate code based on the protogen.Field
type FieldArg struct {
	desc        protoreflect.FieldDescriptor
	Name        string // Name of the field
	NameCC      string // LowerCamelCase(Name)
	Pkg         string
	Type        string
	ZeroValue   string
	Kind        string
	GoKind      string
	Cardinality string
}

func GetFieldArg(file *protogen.File, gFile *protogen.GeneratedFile, f *protogen.Field) FieldArg {
	arg := FieldArg{
		desc: f.Desc,
	}
	arg.Name = f.GoName
	arg.NameCC = tools.ToLowerCamel(arg.Name)
	arg.Pkg, arg.Type = DescParts(file, gFile, f.Desc.Message())
	arg.Kind = f.Desc.Kind().String()
	arg.GoKind = GoKind(file, gFile, f.Desc)
	arg.Cardinality = f.Desc.Cardinality().String()
	arg.ZeroValue = ZeroValue(f.Desc)
	return arg
}

// ServiceArg holds the data needed by the template engine to generate code based on the protogen.Service
type ServiceArg struct {
	desc    protoreflect.ServiceDescriptor
	Name    string
	NameCC  string // LowerCamelCase(Name)
	NameKC  string // KebabCase(Name)
	C       uint32
	Methods []MethodArg
}

func GetServiceArg(file *protogen.File, gFile *protogen.GeneratedFile, s *protogen.Service) ServiceArg {
	arg := ServiceArg{
		desc: s.Desc,
	}
	arg.Name = s.GoName
	arg.NameCC = tools.ToLowerCamel(arg.Name)
	arg.NameKC = tools.ToKebab(arg.Name)
	arg.C = crc32.ChecksumIEEE([]byte(arg.Name))
	for _, m := range s.Methods {
		arg.Methods = append(arg.Methods, GetMethodArg(file, gFile, m))
	}
	return arg
}

// MethodArg holds the data needed by the template engine to generate code based on the protogen.Method
type MethodArg struct {
	desc        protoreflect.MethodDescriptor
	Name        string
	NameCC      string // LowerCamelCase(Name)
	NameKC      string // KebabCase(Name)
	Input       MessageArg
	Output      MessageArg
	RestEnabled bool
	TunnelOnly  bool
	Rest        RestArg
}

type RestArg struct {
	Method    string
	Path      string
	Json      bool
	Unmarshal bool
	ExtraCode []string
}

func GetMethodArg(file *protogen.File, gFile *protogen.GeneratedFile, m *protogen.Method) MethodArg {
	arg := MethodArg{
		desc: m.Desc,
	}
	arg.Name = m.GoName
	arg.NameCC = tools.ToLowerCamel(arg.Name)
	arg.NameKC = tools.ToKebab(arg.Name)
	arg.Input = GetMessageArg(file, gFile, m.Input)
	arg.Output = GetMessageArg(file, gFile, m.Output)
	opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
	restOpt := proto.GetExtension(opt, rony.E_RonyRest).(*rony.RestOpt)
	rest := RestArg{Method: "", Path: "", Json: false}
	if restOpt != nil {
		rest.Method = restOpt.GetMethod()
		rest.Path = fmt.Sprintf("/%s", strings.Trim(restOpt.GetPath(), "/"))
		rest.Json = restOpt.GetJsonEncode()

		var pathVars []string
		bindVars := map[string]string{}
		for _, pv := range strings.Split(rest.Path, "/") {
			if !strings.HasPrefix(pv, ":") {
				continue
			}
			pathVars = append(pathVars, strings.TrimLeft(pv, ":"))
		}
		for _, bv := range restOpt.GetBindVariables() {
			parts := strings.SplitN(strings.TrimSpace(bv), "=", 2)
			if len(parts) == 2 {
				bindVars[parts[0]] = parts[1]
			}
		}

		if len(arg.Input.Fields) > len(pathVars) {
			rest.Unmarshal = true
		}

		for _, pathVar := range pathVars {
			varName := pathVar
			if bindVars[pathVar] != "" {
				varName = bindVars[pathVar]
			}
			for _, f := range m.Input.Fields {
				if f.Desc.JSONName() == varName {
					var ec string

					switch f.Desc.Kind() {
					case protoreflect.Int64Kind, protoreflect.Sfixed64Kind:
						ec = fmt.Sprint("req.", f.Desc.Name(), "= tools.StrToInt64(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
					case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
						ec = fmt.Sprint("req.", f.Desc.Name(), "= tools.StrToUInt64(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
					case protoreflect.Int32Kind, protoreflect.Sfixed32Kind:
						ec = fmt.Sprint("req.", f.Desc.Name(), "= tools.StrToInt32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
					case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
						ec = fmt.Sprint("req.", f.Desc.Name(), "= tools.StrToUInt32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
					case protoreflect.StringKind:
						ec = fmt.Sprint("req.", f.Desc.Name(), "= tools.GetString(conn.Get(\"", pathVar, "\"), \"\")")
					case protoreflect.BytesKind:
						ec = fmt.Sprint("req.", f.Desc.Name(), "= tools.S2B(tools.GetString(conn.Get(\"", pathVar, "\"), \"\"))")
					case protoreflect.DoubleKind:
						ec = fmt.Sprint("req.", f.Desc.Name(), "= tools.StrToFloat32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
					default:
						ec = ""
					}

					if ec != "" {
						rest.ExtraCode = append(rest.ExtraCode, ec)
					}
				}
			}
		}

	}

	arg.RestEnabled = restOpt != nil
	arg.TunnelOnly = proto.GetExtension(opt, rony.E_RonyInternal).(bool)
	arg.Rest = rest

	return arg
}

type ModelArg struct {
	Name    string
	Message MessageArg
	Table   ModelKey
	Views   []ModelKey
}

type ModelFieldArg struct {
	FieldArg
	HasIndex bool
}

type ModelKey struct {
	pks []Prop
	cks []Prop
}

func (m *ModelKey) PartitionKeys() []Prop {
	return m.pks
}

func (m *ModelKey) ClusteringKeys() []Prop {
	return m.cks
}

func (m *ModelKey) Keys() []Prop {
	var all []Prop
	all = append(all, m.pks...)
	all = append(all, m.cks...)
	return all
}

func (m *ModelKey) StringNameTypes(namePrefix string, filter PropFilter) string {
	panic("implement me")
}

func (m *ModelKey) StringNames(prefix string, sep string, nameCase TextCase) string {
	panic("implement me")
}

type Prop struct {
	Name      string
	ProtoType string
	CqlType   string
	GoType    string
	Order     Order
}

type PropFilter string

const (
	PropFilterALL = "ALL"
	PropFilterPKs = "PKs"
	PropFilterCKs = "CKs"
)

type Order string

const (
	ASC  Order = "asc"
	DESC Order = "desc"
)

type TextCase string

const (
	CamelCase      TextCase = "CC"
	LowerCamelCase TextCase = "LCC"
	KebabCase      TextCase = "KC"
)
