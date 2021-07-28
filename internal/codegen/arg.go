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
	file   *protogen.File
	desc   *protogen.Message
	name   string
	pkg    string
	Fields []FieldArg
	C      uint32

	// If message is representing a model then following parameters are filled
	IsAggregate bool
	IsSingleton bool
	RemoteRepo  string
	LocalRepo   string
	Table       *ModelKey
	TableExtra  []Prop
	Views       []*ModelKey
}

func GetMessageArg(m *protogen.Message) MessageArg {
	arg := MessageArg{}
	arg.name = string(m.Desc.Name())
	arg.pkg = string(m.Desc.ParentFile().Package())
	arg.C = crc32.ChecksumIEEE([]byte(m.Desc.Name()))
	for _, f := range m.Fields {
		arg.Fields = append(arg.Fields, GetFieldArg(f))
	}

	// Generate the aggregate description from proto options
	opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)
	arg.RemoteRepo = strings.ToLower(proto.GetExtension(opt, rony.E_RonyRemoteRepo).(string))
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

// FieldArg holds the data needed by the template engine to generate code based on the protogen.Field
type FieldArg struct {
	file        *protogen.File
	desc        protoreflect.FieldDescriptor
	name        string
	pkg         string
	ZeroValue   string
	Kind        string
	GoKind      string
	CqlKind     string
	Cardinality string
	HasIndex    bool
}

func GetFieldArg(f *protogen.Field) FieldArg {
	arg := FieldArg{
		desc: f.Desc,
	}
	arg.name = string(f.Desc.Name())
	if f.Message != nil {
		arg.pkg = string(f.Message.Desc.ParentFile().Package())
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

// ServiceArg holds the data needed by the template engine to generate code based on the protogen.Service
type ServiceArg struct {
	file         *protogen.File
	desc         protoreflect.ServiceDescriptor
	name         string
	C            uint32
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

func GetServiceArg(s *protogen.Service) ServiceArg {
	arg := ServiceArg{
		desc: s.Desc,
	}
	arg.name = s.GoName
	arg.C = crc32.ChecksumIEEE([]byte(arg.name))
	for _, m := range s.Methods {
		ma := GetMethodArg(m)
		if ma.RestEnabled {
			arg.HasRestProxy = true
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

func GetMethodArg(m *protogen.Method) MethodArg {
	arg := MethodArg{
		desc: m.Desc,
	}
	arg.name = m.GoName
	arg.Input = GetMessageArg(m.Input)
	arg.Output = GetMessageArg(m.Output)
	if arg.Input.IsSingleton || arg.Input.IsAggregate {
		panic("method input cannot be aggregate or singleton")
	}
	if arg.Output.IsSingleton || arg.Output.IsAggregate {
		panic("method output cannot be aggregate or singleton")
	}
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

type ModelKey struct {
	Arg   *MessageArg
	pks   []Prop
	cks   []Prop
	index int
	alias string
}

func (m *ModelKey) Name() string {
	return m.Arg.name
}

func (m *ModelKey) Alias() string {
	return m.alias
}

func (m *ModelKey) PartitionKeys() []Prop {
	var all []Prop
	all = append(all, m.pks...)
	return all
}

func (m *ModelKey) PKs() []Prop {
	return m.PartitionKeys()
}

func (m *ModelKey) ClusteringKeys() []Prop {
	var all []Prop
	all = append(all, m.cks...)
	return all
}

func (m *ModelKey) CKs() []Prop {
	return m.ClusteringKeys()
}

func (m *ModelKey) Keys() []Prop {
	var all []Prop
	all = append(all, m.pks...)
	all = append(all, m.cks...)
	return all
}

func (m *ModelKey) Index() int {
	return m.index
}

func (m *ModelKey) IsSubset(n *ModelKey) bool {
	for _, np := range n.Keys() {
		found := false
		for _, mp := range m.Keys() {
			if mp.Name == np.Name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (m *ModelKey) HasProp(name string) bool {
	for _, p := range m.Keys() {
		if p.Name == name {
			return true
		}
	}
	return false
}

// NameTypes is kind of strings.Join function which returns a custom format of combination of model properties.
// This is a helper function used in code generator templates.
func (m *ModelKey) NameTypes(filter PropFilter, namePrefix string, nameCase TextCase, lang Language) string {
	var props []Prop
	switch filter {
	case PropFilterALL:
		props = m.Keys()
	case PropFilterCKs:
		props = m.ClusteringKeys()
	case PropFilterPKs:
		props = m.PartitionKeys()
	}

	sb := strings.Builder{}
	for idx, p := range props {
		if idx != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(namePrefix)
		switch nameCase {
		case LowerCamelCase:
			sb.WriteString(tools.ToLowerCamel(p.Name))
		case CamelCase:
			sb.WriteString(tools.ToCamel(p.Name))
		case KebabCase:
			sb.WriteString(tools.ToKebab(p.Name))
		case SnakeCase:
			sb.WriteString(tools.ToSnake(p.Name))
		default:
			sb.WriteString(p.Name)
		}
		sb.WriteRune(' ')
		switch lang {
		case LangGo:
			sb.WriteString(p.GoType)
		case LangCQL:
			sb.WriteString(p.CqlType)
		case LangProto:
			sb.WriteString(p.ProtoType)
		}
	}
	return sb.String()
}

// Names is kind of strings.Join function which returns a custom format of property names.
// This is a helper function used in code generator templates.
func (m *ModelKey) Names(filter PropFilter, prefix, postfix string, sep string, nameCase TextCase) string {
	var props []Prop
	switch filter {
	case PropFilterALL:
		props = m.Keys()
	case PropFilterCKs:
		props = m.ClusteringKeys()
	case PropFilterPKs:
		props = m.PartitionKeys()
	}

	sb := strings.Builder{}
	for idx, p := range props {
		if idx != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(prefix)
		switch nameCase {
		case LowerCamelCase:
			sb.WriteString(tools.ToLowerCamel(p.Name))
		case CamelCase:
			sb.WriteString(tools.ToCamel(p.Name))
		case KebabCase:
			sb.WriteString(tools.ToKebab(p.Name))
		case SnakeCase:
			sb.WriteString(tools.ToSnake(p.Name))
		default:
			sb.WriteString(p.Name)
		}
		sb.WriteString(postfix)
	}
	return sb.String()
}

type Prop struct {
	Name      string
	ProtoType string
	CqlType   string
	GoType    string
	Order     Order
}

func (p *Prop) NameSC() string {
	return tools.ToSnake(p.Name)
}

type PropFilter string

const (
	PropFilterALL = "ALL"
	PropFilterPKs = "PKs"
	PropFilterCKs = "CKs"
)

type Order string

const (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

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

type ModuleArg struct {
	Service    *ServiceArg
	Aggregates []MessageArg
	Singletons []MessageArg
}

func GetModuleArg(plugin *protogen.Plugin) ModuleArg {
	arg := ModuleArg{}
	for _, f := range plugin.Files {
		if len(f.Services) == 0 {
			continue
		}
		if len(f.Services) != 1 || arg.Service != nil {
			panic("only one service must be defined in each module")
		}
		// arg.Service = GetServiceArg(f, gFile, file.Services[0])
	}

	// arg.Service =
	// for _, m := range file.Messages {
	// 	ma := GetMessageArg(file, gFile, m)
	// 	if ma.IsAggregate {
	// 		arg.Aggregates = append(arg.Aggregates, ma)
	// 	} else if ma.IsSingleton {
	// 		arg.Singletons = append(arg.Singletons, ma)
	// 	}
	// }
	return arg
}
