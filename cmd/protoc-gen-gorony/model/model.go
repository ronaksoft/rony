package model

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

var loadedModels = map[string]*Model{}
var loadedFields = map[string]struct{}{}

type Model struct {
	Type       string
	Name       string
	Table      Key
	Views      []Key
	ViewParams []string
	FieldNames []string
	FieldsCql  map[string]string
	FieldsGo   map[string]string
}

func (m *Model) FuncArgs(key Key, onlyPKs bool) string {
	sb := strings.Builder{}
	keys := key.PKs
	if !onlyPKs {
		keys = key.Keys()
	}
	for idx, k := range keys {
		if idx != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(k))
		sb.WriteRune(' ')
		sb.WriteString(m.FieldsGo[k])
	}
	return sb.String()
}

func (m *Model) FuncArgsWithPrefix(prefix string, pk Key, onlyPKs bool) string {
	sb := strings.Builder{}
	keys := pk.PKs
	if !onlyPKs {
		keys = pk.Keys()
	}
	for idx, k := range keys {
		if idx != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(tools.ToLowerCamel(fmt.Sprintf("%s%s", prefix, k)))
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
			sb.WriteRune(',')
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

// ResetModels reset the internal data
func ResetModels() {
	loadedModels = map[string]*Model{}
}

// FillModel fills the in global loadedModels with parsed data
func FillModel(file *protogen.File, g *protogen.GeneratedFile, m *protogen.Message) {
	var (
		isModel = false
		mm      = Model{
			FieldsCql: make(map[string]string),
			FieldsGo:  make(map[string]string),
		}
	)

	modelDesc := strings.Builder{}
	opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)
	if entity := proto.GetExtension(opt, rony.E_RonyModel).(bool); entity {
		modelDesc.WriteString(fmt.Sprintf("{{@model %s}}\n", m.Desc.Name()))
	}
	if tab := proto.GetExtension(opt, rony.E_RonyTable).(string); tab != "" {
		modelDesc.WriteString(fmt.Sprintf("{{@tab %s}}\n", tab))
	}
	if view := proto.GetExtension(opt, rony.E_RonyView).(string); view != "" {
		modelDesc.WriteString(fmt.Sprintf("{{@view %s}}\n", view))
	}

	t, err := parse.Parse(string(m.Desc.Name()), modelDesc.String())
	if err != nil {
		panic(err)
	}
	fields := make(map[string]struct{})
	for _, n := range t.Root.Nodes {
		switch n.Type() {
		case parse.NodeModel:
			nn := n.(*parse.ModelNode)
			mm.Type = nn.Text
			isModel = true
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
	if isModel {
		for _, f := range m.Fields {
			mm.FieldsCql[f.GoName] = z.CqlKind(f.Desc)
			mm.FieldsGo[f.GoName] = z.GoKind(file, g, f.Desc)
		}
		mm.Name = string(m.Desc.Name())
		loadedModels[mm.Name] = &mm
	}
}

// GetModels
func GetModels() map[string]*Model {
	return loadedModels
}
