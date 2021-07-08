package codegen

import (
	"fmt"
	"github.com/ronaksoft/rony"
	parse "github.com/ronaksoft/rony/internal/parser"
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

	// If message is representing a model then following parameters are filled
	IsAggregate   bool
	IsSingleton   bool
	AggregateType string
	Table         *ModelKey
	Views         []*ModelKey
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

	// parse model based on rony proto options
	arg.parseModel(file, gFile, m)
	return arg
}

func (ma *MessageArg) parseModel(file *protogen.File, gFile *protogen.GeneratedFile, m *protogen.Message) {
	// Generate the aggregate description from proto options
	aggregateDesc := strings.Builder{}
	opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)

	ma.IsSingleton = proto.GetExtension(opt, rony.E_RonySingleton).(bool)
	if ma.IsSingleton {
		// if message is going to be singleton then it could not be aggregate
		return
	}

	if entity := proto.GetExtension(opt, rony.E_RonyAggregate).(bool); entity {
		aggrType := proto.GetExtension(opt, rony.E_RonyAggregateType).(string)
		if aggrType == "" {
			panic("define rony_aggregate_type")
		}
		aggregateDesc.WriteString(fmt.Sprintf("{{@model %s}}\n", aggrType))
	}
	if tab := proto.GetExtension(opt, rony.E_RonyAggregateTable).(string); tab != "" {
		aggregateDesc.WriteString(fmt.Sprintf("{{@tab %s}}\n", tab))
	}
	if views := proto.GetExtension(opt, rony.E_RonyAggregateView).([]string); len(views) > 0 {
		for _, view := range views {
			aggregateDesc.WriteString(fmt.Sprintf("{{@view %s}}\n", view))
		}
	}

	// Parse the generated description
	t, err := parse.Parse(string(m.Desc.Name()), aggregateDesc.String())
	if err != nil {
		panic(err)
	}

	// Generate Go and CQL kinds of the fields
	cqlTypes := map[string]string{}
	goTypes := map[string]string{}
	protoTypes := map[string]string{}
	for _, f := range m.Fields {
		protoTypes[f.GoName] = f.Desc.Kind().String()
		cqlTypes[f.GoName] = CqlKind(f.Desc)
		goTypes[f.GoName] = GoKind(file, gFile, f.Desc)
	}
	for _, n := range t.Root.Nodes {
		switch n.Type() {
		case parse.NodeModel:
			ma.AggregateType = n.(*parse.ModelNode).Text
			ma.IsAggregate = true
		case parse.NodeTable:
			ma.Table = &ModelKey{
				Arg: ma,
			}
			nn := n.(*parse.TableNode)
			for _, k := range nn.PartitionKeys {
				ma.Table.pks = append(ma.Table.pks, Prop{
					Name:      k,
					ProtoType: protoTypes[k],
					CqlType:   cqlTypes[k],
					GoType:    goTypes[k],
					Order:     "",
				})
			}
			for _, k := range nn.ClusteringKeys {
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
		case parse.NodeView:
			view := &ModelKey{
				Arg: ma,
			}
			nn := n.(*parse.ViewNode)
			for _, k := range nn.PartitionKeys {
				view.pks = append(view.pks, Prop{
					Name:      k,
					ProtoType: protoTypes[k],
					CqlType:   cqlTypes[k],
					GoType:    goTypes[k],
					Order:     "",
				})
			}

			for _, k := range nn.ClusteringKeys {
				order := ASC
				if strings.HasPrefix(k, "-") {
					order = DESC
				}

				k = strings.TrimLeft(k, "-")
				view.cks = append(view.cks, Prop{
					Name:      k,
					ProtoType: protoTypes[k],
					CqlType:   cqlTypes[k],
					GoType:    goTypes[k],
					Order:     order,
				})
			}
			ma.Views = append(ma.Views, view)
		}
	}
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
	CqlKind     string
	Cardinality string
	HasIndex    bool
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
	arg.CqlKind = CqlKind(f.Desc)
	arg.Cardinality = f.Desc.Cardinality().String()
	arg.ZeroValue = ZeroValue(f.Desc)

	opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
	arg.HasIndex = proto.GetExtension(opt, rony.E_RonyIndex).(bool)
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

type ModelKey struct {
	Arg *MessageArg
	pks []Prop
	cks []Prop
}

func (m *ModelKey) Name() string {
	return m.Arg.Name
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
func (m *ModelKey) Names(filter PropFilter, prefix string, sep string, nameCase TextCase) string {
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
		default:
			sb.WriteString(p.Name)
		}
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
)

type Language string

const (
	LangGo    Language = "GO"
	LangCQL   Language = "CQL"
	LangProto Language = "PROTO"
)
