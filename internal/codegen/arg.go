package codegen

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
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

type TextCase string

const (
	None           TextCase = ""
	CamelCase      TextCase = "CC"
	LowerCamelCase TextCase = "LCC"
	KebabCase      TextCase = "KC"
	SnakeCase      TextCase = "SC"
)

type Language string

const (
	LangGo    Language = "GO"
	LangCQL   Language = "CQL"
	LangProto Language = "PROTO"
)

type TemplateArg struct {
	Messages []MessageArg
	Services []ServiceArg
}

func GenTemplateArg(f *protogen.File) *TemplateArg {
	arg := &TemplateArg{}
	for _, m := range f.Messages {
		arg.Messages = append(arg.Messages, GetMessageArg(m).With(f))
	}
	for _, s := range f.Services {
		arg.Services = append(arg.Services, getServiceArg(s).With(f))
	}

	return arg
}

// MessageArg holds the data needed by the template engine to generate code based on the protogen.Message
type MessageArg struct {
	file       *protogen.File
	desc       *protogen.Message
	name       string
	pkg        string
	Fields     []FieldArg
	C          uint64
	ImportPath protogen.GoImportPath

	// If message is representing a model then following parameters are filled
	IsAggregate bool
	IsSingleton bool
	GlobalRepo  string
	LocalRepo   string
	Table       *ModelKey
	TableExtra  []Prop
	Views       []*ModelKey
}

func GetMessageArg(m *protogen.Message) MessageArg {
	arg := MessageArg{}
	arg.name = string(m.Desc.Name())
	arg.pkg = string(m.Desc.ParentFile().Package())
	arg.C = CrcHash([]byte(m.Desc.Name()))
	arg.ImportPath = m.GoIdent.GoImportPath
	for _, f := range m.Fields {
		arg.Fields = append(arg.Fields, getFieldArg(f))
	}

	// Generate the aggregate description from proto options
	opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)
	arg.GlobalRepo = strings.ToLower(proto.GetExtension(opt, rony.E_RonyGlobalRepo).(string))
	arg.LocalRepo = strings.ToLower(proto.GetExtension(opt, rony.E_RonyLocalRepo).(string))
	arg.IsSingleton = proto.GetExtension(opt, rony.E_RonySingleton).(bool)
	if arg.IsSingleton {
		// if message is going to be singleton then it could not be aggregate
		return arg
	}

	// If there is no table defined then it is not an aggregate
	if proto.GetExtension(opt, rony.E_RonyTable).(*rony.PrimaryKeyOpt) == nil {
		return arg
	}

	arg.IsAggregate = true
	arg.parsePrimaryKeyOpt(m)

	return arg
}

func (ma *MessageArg) parsePrimaryKeyOpt(m *protogen.Message) {
	opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)
	tablePK := proto.GetExtension(opt, rony.E_RonyTable).(*rony.PrimaryKeyOpt)
	viewPKs := proto.GetExtension(opt, rony.E_RonyView).([]*rony.PrimaryKeyOpt)

	// Generate Go and CQL kinds of the fields
	cqlTypes := map[string]string{}
	goTypes := map[string]string{}
	protoTypes := map[string]string{}
	uniqueView := map[string]*ModelKey{}
	extraProp := map[string]Prop{}
	for _, f := range m.Fields {
		protoTypes[f.GoName] = f.Desc.Kind().String()
		cqlTypes[f.GoName] = CqlKind(f.Desc)
		goTypes[f.GoName] = GoKind(f.Desc)
	}

	// Fill Table's ModelKey
	ma.Table = &ModelKey{
		Arg:   ma,
		alias: tablePK.GetAlias(),
	}
	for _, k := range tablePK.PartKey {
		ma.Table.pks = append(ma.Table.pks, Prop{
			Name:      k,
			ProtoType: protoTypes[k],
			CqlType:   cqlTypes[k],
			GoType:    goTypes[k],
			Order:     "",
		})
	}
	for _, k := range tablePK.SortKey {
		order := ASC
		if strings.HasPrefix(k, "-") {
			order = DESC
		}

		k = strings.TrimLeft(k, "-")
		ma.Table.cks = append(ma.Table.cks, Prop{
			Name:      k,
			ProtoType: protoTypes[k],
			CqlType:   cqlTypes[k],
			GoType:    goTypes[k],
			Order:     order,
		})
	}

	// Fill Views' ModelKey
	for _, v := range viewPKs {
		view := &ModelKey{
			Arg:   ma,
			alias: v.GetAlias(),
		}
		for _, k := range v.PartKey {
			p := Prop{
				Name:      k,
				ProtoType: protoTypes[k],
				CqlType:   cqlTypes[k],
				GoType:    goTypes[k],
				Order:     "",
			}
			if !ma.Table.HasProp(k) {
				extraProp[k] = p
			}
			view.pks = append(view.pks, p)
		}
		for _, k := range v.SortKey {
			order := ASC
			if strings.HasPrefix(k, "-") {
				order = DESC
			}
			k = strings.TrimLeft(k, "-")
			p := Prop{
				Name:      k,
				ProtoType: protoTypes[k],
				CqlType:   cqlTypes[k],
				GoType:    goTypes[k],
				Order:     order,
			}
			if !ma.Table.HasProp(k) {
				extraProp[k] = p
			}
			view.cks = append(view.cks, p)
		}
		if oldView, ok := uniqueView[view.Names(PropFilterPKs, "", "", "", None)]; ok {
			view.index = oldView.index + 1
		}
		uniqueView[view.Names(PropFilterPKs, "", "", "", None)] = view
		ma.Views = append(ma.Views, view)
	}
	for _, p := range extraProp {
		ma.TableExtra = append(ma.TableExtra, p)
	}
}

func (ma *MessageArg) currentPkg() string {
	var pkg string
	if ma.file != nil {
		pkg = string(ma.file.GoPackageName)
	}

	return pkg
}

func (ma MessageArg) Name() string {
	return ma.name
}

func (ma MessageArg) NameCC() string {
	return tools.ToLowerCamel(ma.name)
}

func (ma MessageArg) NameKC() string {
	return tools.ToKebab(ma.name)
}

func (ma MessageArg) NameSC() string {
	return tools.ToSnake(ma.name)
}

func (ma MessageArg) Pkg() string {
	if ma.currentPkg() == ma.pkg {
		return ""
	}

	return ma.pkg
}

func (ma MessageArg) ViewsByPK() map[string][]*ModelKey {
	var res = map[string][]*ModelKey{}
	for _, v := range ma.Views {
		k := v.Names(PropFilterPKs, "", "", "", None)
		res[k] = append(res[k], v)
	}

	return res
}

func (ma MessageArg) Fullname() string {
	if ma.pkg == ma.currentPkg() {
		return ma.name
	} else {
		return fmt.Sprintf("%s.%s", ma.pkg, ma.name)
	}
}

func (ma MessageArg) CName() string {
	return fmt.Sprintf("C_%s", ma.name)
}

func (ma MessageArg) With(f *protogen.File) MessageArg {
	ma.file = f
	for idx := range ma.Fields {
		ma.Fields[idx] = ma.Fields[idx].With(f)
	}

	return ma
}

func (ma MessageArg) Options() *descriptorpb.MessageOptions {
	opt, _ := ma.desc.Desc.Options().(*descriptorpb.MessageOptions)

	return opt
}

// FieldArg holds the data needed by the template engine to generate code based on the protogen.Field
type FieldArg struct {
	file         *protogen.File
	desc         protoreflect.FieldDescriptor
	name         string
	pkg          string
	ImportPath   protogen.GoImportPath
	ZeroValue    string
	Kind         string
	GoKind       string
	CqlKind      string
	Cardinality  string
	HasIndex     bool
	HelpText     string
	DefaultValue string
}

func getFieldArg(f *protogen.Field) FieldArg {
	arg := FieldArg{
		desc: f.Desc,
	}

	arg.name = f.GoName
	if f.Message != nil {
		arg.pkg = string(f.Message.Desc.ParentFile().Package())
		arg.ImportPath = f.Message.GoIdent.GoImportPath
	} else {
		arg.pkg = string(f.Desc.ParentFile().Package())
	}

	arg.Kind = f.Desc.Kind().String()
	arg.GoKind = GoKind(f.Desc)
	arg.CqlKind = CqlKind(f.Desc)
	arg.Cardinality = f.Desc.Cardinality().String()
	arg.ZeroValue = ZeroValue(f.Desc)

	opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
	arg.HasIndex = proto.GetExtension(opt, rony.E_RonyIndex).(bool)
	arg.HelpText = proto.GetExtension(opt, rony.E_RonyHelp).(string)
	arg.DefaultValue = proto.GetExtension(opt, rony.E_RonyDefault).(string)

	return arg
}

func (fa FieldArg) currentPkg() string {
	var pkg string
	if fa.file != nil {
		pkg = string(fa.file.GoPackageName)
	}

	return pkg
}

func (fa FieldArg) Name() string {
	return fa.name
}

func (fa FieldArg) NameCC() string {
	return tools.ToLowerCamel(fa.name)
}

func (fa FieldArg) NameKC() string {
	return tools.ToKebab(fa.name)
}

func (fa FieldArg) NameSC() string {
	return tools.ToSnake(fa.name)
}

func (fa FieldArg) Type() string {
	if fa.desc.Message() != nil {
		return string(fa.desc.Message().Name())
	}

	return fa.GoKind
}

func (fa FieldArg) Pkg() string {
	if fa.currentPkg() == fa.pkg {
		return ""
	}

	return fa.pkg
}

func (fa FieldArg) With(f *protogen.File) FieldArg {
	fa.file = f

	return fa
}

func (fa FieldArg) Options() *descriptorpb.FieldOptions {
	opt, _ := fa.desc.Options().(*descriptorpb.FieldOptions)

	return opt
}

// ServiceArg holds the data needed by the template engine to generate code based on the protogen.Service
type ServiceArg struct {
	file         *protogen.File
	desc         protoreflect.ServiceDescriptor
	name         string
	C            uint64
	Methods      []MethodArg
	HasRestProxy bool
}

func (sa ServiceArg) currentPkg() string {
	var pkg string
	if sa.file != nil {
		pkg = string(sa.file.GoPackageName)
	}

	return pkg
}

func (sa ServiceArg) Name() string {
	return sa.name
}

func (sa ServiceArg) NameCC() string {
	return tools.ToLowerCamel(sa.name)
}

func (sa ServiceArg) NameKC() string {
	return tools.ToKebab(sa.name)
}

func (sa ServiceArg) NameSC() string {
	return tools.ToSnake(sa.name)
}

func (sa ServiceArg) With(f *protogen.File) ServiceArg {
	sa.file = f
	for idx := range sa.Methods {
		sa.Methods[idx] = sa.Methods[idx].With(sa.file)
	}

	return sa
}

func (sa ServiceArg) Options() *descriptorpb.ServiceOptions {
	opt, _ := sa.desc.Options().(*descriptorpb.ServiceOptions)

	return opt
}

func getServiceArg(s *protogen.Service) ServiceArg {
	arg := ServiceArg{
		desc: s.Desc,
	}
	arg.name = string(s.Desc.Name())
	arg.C = CrcHash([]byte(arg.name))
	for _, m := range s.Methods {
		ma := getMethodArg(s, m)
		if ma.RestEnabled {
			arg.HasRestProxy = true
		}
		if arg.currentPkg() != "" && arg.currentPkg() != ma.Input.pkg {
			panic("input must be with in the same package of its service")
		}
		arg.Methods = append(arg.Methods, ma)
	}

	return arg
}

// MethodArg holds the data needed by the template engine to generate code based on the protogen.Method
type MethodArg struct {
	desc        protoreflect.MethodDescriptor
	file        *protogen.File
	name        string
	fullname    string
	C           uint64
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
	PathVars  map[string]protoreflect.Kind
}

func getRestArg(m *protogen.Method, arg *MethodArg) {
	opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
	restOpt := proto.GetExtension(opt, rony.E_RonyRest).(*rony.RestOpt)
	if restOpt == nil {
		return
	}

	//rest := RestArg{Method: "", Path: "", Json: false}
	arg.Rest.Method = restOpt.GetMethod()
	arg.Rest.Path = fmt.Sprintf("/%s", strings.Trim(restOpt.GetPath(), "/"))
	arg.Rest.Json = restOpt.GetJsonEncode()

	var pathVars []string
	bindVars := map[string]string{}
	for _, pv := range strings.Split(arg.Rest.Path, "/") {
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
		arg.Rest.Unmarshal = true
	}

	for _, pathVar := range pathVars {
		varName := pathVar
		if _, ok := bindVars[pathVar]; ok {
			varName = bindVars[pathVar]
		}
		for _, f := range m.Input.Fields {
			if string(f.Desc.Name()) == varName {
				var ec string

				switch f.Desc.Kind() {
				case protoreflect.Int64Kind, protoreflect.Sfixed64Kind:
					ec = fmt.Sprint(
						"req.",
						f.GoName, "= tools.StrToInt64(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))",
					)
				case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					ec = fmt.Sprint(
						"req.",
						f.GoName, "= tools.StrToUInt64(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))",
					)
				case protoreflect.Int32Kind, protoreflect.Sfixed32Kind:
					ec = fmt.Sprint(
						"req.",
						f.GoName, "= tools.StrToInt32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))",
					)
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
					ec = fmt.Sprint(
						"req.",
						f.GoName, "= tools.StrToUInt32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))",
					)
				case protoreflect.StringKind:
					ec = fmt.Sprint(
						"req.",
						f.GoName, "= tools.GetString(conn.Get(\"", pathVar, "\"), \"\")",
					)
				case protoreflect.BytesKind:
					ec = fmt.Sprint(
						"req.",
						f.GoName, "= tools.S2B(tools.GetString(conn.Get(\"", pathVar, "\"), \"\"))",
					)
				case protoreflect.DoubleKind:
					ec = fmt.Sprint(
						"req.",
						f.GoName, "= tools.StrToFloat32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))",
					)
				default:
					ec = ""
				}

				if ec != "" {
					arg.Rest.ExtraCode = append(arg.Rest.ExtraCode, ec)
					arg.Rest.PathVars[pathVar] = f.Desc.Kind()
				}
			}
		}
	}

	return
}

func (ma MethodArg) Fullname() string {
	return ma.fullname
}

func (ma MethodArg) Name() string {
	return ma.name
}

func (ma MethodArg) NameCC() string {
	return tools.ToLowerCamel(ma.name)
}

func (ma MethodArg) NameKC() string {
	return tools.ToKebab(ma.name)
}

func (ma MethodArg) NameSC() string {
	return tools.ToSnake(ma.name)
}

func (ma MethodArg) With(f *protogen.File) MethodArg {
	ma.file = f
	ma.Input = ma.Input.With(f)
	ma.Output = ma.Output.With(f)

	return ma
}

func (ma MethodArg) Options() *descriptorpb.MethodOptions {
	opt, _ := ma.desc.Options().(*descriptorpb.MethodOptions)

	return opt
}

func getMethodArg(s *protogen.Service, m *protogen.Method) MethodArg {
	arg := MethodArg{
		desc: m.Desc,
		Rest: RestArg{
			Method:   "",
			Path:     "",
			Json:     false,
			PathVars: map[string]protoreflect.Kind{},
		},
	}
	arg.fullname = fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())
	arg.name = string(m.Desc.Name())
	arg.C = CrcHash([]byte(fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())))
	arg.Input = GetMessageArg(m.Input)
	arg.Output = GetMessageArg(m.Output)
	if arg.Input.IsSingleton || arg.Input.IsAggregate {
		panic("method input cannot be aggregate or singleton")
	}
	if arg.Output.IsSingleton || arg.Output.IsAggregate {
		panic("method output cannot be aggregate or singleton")
	}

	opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
	getRestArg(m, &arg)

	arg.TunnelOnly = proto.GetExtension(opt, rony.E_RonyInternal).(bool)

	return arg
}
