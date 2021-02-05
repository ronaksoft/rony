package aggregate

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	parse "github.com/ronaksoft/rony/internal/parser"
	"github.com/ronaksoft/rony/tools"
	"github.com/scylladb/gocqlx/v2/qb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"hash/crc32"
	"strings"
)

/*
   Creation Time: 2021 - Jan - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var _ = qb.ASC

var loadedAggregates = map[string]*Aggregate{}
var loadedFields = map[string]struct{}{}

type Aggregate struct {
	Type       string
	Name       string
	Table      Key
	Views      []Key
	ViewParams []string
	FieldNames []string
	FieldsCql  map[string]string
	FieldsGo   map[string]string
}

func (m *Aggregate) FuncArgs(keyPrefix string, key Key) string {
	sb := strings.Builder{}
	for idx, k := range key.Keys() {
		if idx != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", keyPrefix, k)))
		sb.WriteRune(' ')
		sb.WriteString(m.FieldsGo[k])
	}
	return sb.String()
}

func (m *Aggregate) FuncArgsPKs(keyPrefix string, key Key) string {
	sb := strings.Builder{}
	for idx, k := range key.PKs {
		if idx != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", keyPrefix, k)))
		sb.WriteRune(' ')
		sb.WriteString(m.FieldsGo[k])
	}
	return sb.String()
}

func (m *Aggregate) FuncArgsCKs(keyPrefix string, key Key) string {
	sb := strings.Builder{}
	for idx, k := range key.CKs {
		if idx != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", keyPrefix, k)))
		sb.WriteRune(' ')
		sb.WriteString(m.FieldsGo[k])
	}
	return sb.String()
}

type Key struct {
	PKs    []string
	CKs    []string
	Orders []string
}

func (k *Key) Keys() []string {
	keys := make([]string, 0, len(k.PKs)+len(k.CKs))
	keys = append(keys, k.PKs...)
	keys = append(keys, k.CKs...)
	return keys
}

func (k *Key) StringPKs(prefix string, sep string, lowerCamel bool) string {
	sb := strings.Builder{}
	for idx, k := range k.PKs {
		if idx != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(prefix)
		if lowerCamel {
			sb.WriteString(tools.ToLowerCamel(k))
		} else {
			sb.WriteString(k)
		}

	}
	return sb.String()
}

func (k *Key) StringCKs(prefix string, sep string, lowerCamel bool) string {
	sb := strings.Builder{}
	for idx, k := range k.CKs {
		if idx != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(prefix)
		if lowerCamel {
			sb.WriteString(tools.ToLowerCamel(k))
		} else {
			sb.WriteString(k)
		}

	}
	return sb.String()
}

func (k *Key) ChecksumPKs(keyPrefix string, sep string, lowerCamel bool) uint32 {
	return crc32.ChecksumIEEE(tools.StrToByte(k.StringPKs(keyPrefix, sep, lowerCamel)))
}

func (k *Key) String(prefix string, sep string, lowerCamel bool) string {
	sb := strings.Builder{}
	for idx, k := range k.Keys() {
		if idx != 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(prefix)
		if lowerCamel {
			sb.WriteString(tools.ToLowerCamel(k))
		} else {
			sb.WriteString(k)
		}

	}
	return sb.String()
}

func (k *Key) Checksum() uint32 {
	return crc32.ChecksumIEEE(tools.StrToByte(k.String("", ",", false)))
}

// ResetAggregates reset the internal data
func ResetAggregates() {
	loadedAggregates = map[string]*Aggregate{}
}

// FillAggregate fills the in global loadedAggregates with parsed data
func FillAggregate(file *protogen.File, g *protogen.GeneratedFile, m *protogen.Message) {
	var (
		isAggregate = false
		mm          = Aggregate{
			FieldsCql: make(map[string]string),
			FieldsGo:  make(map[string]string),
		}
	)

	// Generate the aggregate description from proto options
	aggregateDesc := strings.Builder{}
	opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)

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
	if view := proto.GetExtension(opt, rony.E_RonyAggregateView).(string); view != "" {
		aggregateDesc.WriteString(fmt.Sprintf("{{@view %s}}\n", view))
	}

	// Parse the generated description
	t, err := parse.Parse(string(m.Desc.Name()), aggregateDesc.String())
	if err != nil {
		panic(err)
	}
	fields := make(map[string]struct{})
	for _, n := range t.Root.Nodes {
		switch n.Type() {
		case parse.NodeModel:
			nn := n.(*parse.ModelNode)
			mm.Type = nn.Text
			isAggregate = true
		case parse.NodeTable:
			pk := Key{}
			nn := n.(*parse.TableNode)
			for _, k := range nn.PartitionKeys {
				fields[k] = struct{}{}
				pk.PKs = append(pk.PKs, k)
			}
			for _, k := range nn.ClusteringKeys {
				kWithoutSign := strings.TrimLeft(k, "-")
				fields[kWithoutSign] = struct{}{}
				pk.Orders = append(pk.Orders, k)
				pk.CKs = append(pk.CKs, kWithoutSign)
			}
			mm.Table = pk
		case parse.NodeView:
			pk := Key{}
			nn := n.(*parse.ViewNode)
			sb := strings.Builder{}
			for _, k := range nn.PartitionKeys {
				fields[k] = struct{}{}
				pk.PKs = append(pk.PKs, k)
				sb.WriteString(k)
			}
			mm.ViewParams = append(mm.ViewParams, sb.String())
			for _, k := range nn.ClusteringKeys {
				kWithoutSign := strings.TrimLeft(k, "-")
				fields[kWithoutSign] = struct{}{}
				pk.Orders = append(pk.Orders, k)
				pk.CKs = append(pk.CKs, kWithoutSign)
			}
			mm.Views = append(mm.Views, pk)
		}
	}
	for f := range fields {
		loadedFields[f] = struct{}{}
		mm.FieldNames = append(mm.FieldNames, f)
	}
	if isAggregate {
		for _, f := range m.Fields {
			mm.FieldsCql[f.GoName] = z.CqlKind(f.Desc)
			mm.FieldsGo[f.GoName] = z.GoKind(file, g, f.Desc)
		}
		mm.Name = string(m.Desc.Name())
		loadedAggregates[mm.Name] = &mm
	}
}

// GetAggregates
func GetAggregates() map[string]*Aggregate {
	return loadedAggregates
}
