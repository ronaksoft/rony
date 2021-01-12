package model

import (
	"fmt"
	"github.com/ronaksoft/rony"
	parse "github.com/ronaksoft/rony/internal/parser"
	"github.com/scylladb/gocqlx/v2/qb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
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
	Table      PrimaryKey
	Views      []PrimaryKey
	ViewParams []string
	FieldNames []string
	FieldsCql  map[string]string
	FieldsGo   map[string]string
}

type PrimaryKey struct {
	PKs    []string
	CKs    []string
	Orders []string
}

func (pk *PrimaryKey) Keys() []string {
	keys := make([]string, 0, len(pk.PKs)+len(pk.CKs))
	keys = append(keys, pk.PKs...)
	keys = append(keys, pk.CKs...)
	return keys
}

// ResetModels reset the internal data
func ResetModels() {
	loadedModels = map[string]*Model{}
}

// FillModel fills the in global loadedModels with parsed data
func FillModel(m *protogen.Message) {
	var (
		isModel = false
		mm      = Model{
			FieldsCql: make(map[string]string),
			FieldsGo:  make(map[string]string),
		}
	)

	modelDesc := strings.Builder{}
	opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)
	if tab := proto.GetExtension(opt, rony.E_Table).(string); tab != "" {
		modelDesc.WriteString(fmt.Sprintf("{{@tab %s}}\n", tab))
	}
	if view := proto.GetExtension(opt, rony.E_View).(string); view != "" {
		modelDesc.WriteString(fmt.Sprintf("{{@view %s}}\n", view))
	}

	t, err := parse.Parse(string(m.Desc.Name()), modelDesc.String())
	if err != nil {
		panic(err)
	}
	fields := make(map[string]struct{})
	for _, n := range t.Root.Nodes {
		switch n.Type() {
		case parse.NodeEntity:
			nn := n.(*parse.EntityNode)
			mm.Type = nn.Text
			isModel = true
		case parse.NodeTable:
			pk := PrimaryKey{}
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
			pk := PrimaryKey{}
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
			mm.FieldsCql[f.GoName] = kindCql(f.Desc.Kind())
			mm.FieldsGo[f.GoName] = kindGo(f.Desc.Kind())
		}
		mm.Name = string(m.Desc.Name())
		loadedModels[mm.Name] = &mm
	}
}

// GetModels
func GetModels() map[string]*Model {
	return loadedModels
}

// kindCql converts the proto buffer type to cql types
func kindCql(k protoreflect.Kind) string {
	switch k {
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

// kindGo converts proto buffer types to golang types
func kindGo(k protoreflect.Kind) string {
	switch k {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind:
		return "int32"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "uint32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind:
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
	}
	return "unsupported"
	// panic(fmt.Sprintf("unsupported kindGo: %v", k.String()))
}
