package codegen

import (
	"github.com/ronaksoft/rony/tools"
	"strings"
)

/*
   Creation Time: 2021 - Jul - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

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
