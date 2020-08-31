package main

import (
	parse "git.ronaksoft.com/ronak/rony/internal/parser"
	"github.com/scylladb/gocqlx/v2/qb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"strings"
)

/*
   Creation Time: 2020 - Aug - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var _ = qb.ASC

var _Models = map[string]*Model{}
var _Fields = map[string]struct{}{}

type Model struct {
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

func resetModels() {
	_Models = map[string]*Model{}
}

// fillModel fills the in global _Models with parsed data
func fillModel(m *protogen.Message) {
	var (
		isModel = false
		mm      = Model{
			FieldsCql: make(map[string]string),
			FieldsGo:  make(map[string]string),
		}
	)

	t, err := parse.Parse(string(m.Desc.Name()), string(m.Comments.Leading))
	if err != nil {
		panic(err)
	}
	fields := make(map[string]struct{})
	for _, n := range t.Root.Nodes {
		switch n.Type() {
		case parse.NodeModel:
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
		_Fields[f] = struct{}{}
		mm.FieldNames = append(mm.FieldNames, f)
	}
	if isModel {
		for _, f := range m.Fields {
			mm.FieldsCql[f.GoName] = kindCql(f.Desc.Kind())
			mm.FieldsGo[f.GoName] = kindGo(f.Desc.Kind())
		}
		mm.Name = string(m.Desc.Name())
		_Models[mm.Name] = &mm
	}
}

func getModels() map[string]*Model {
	return _Models
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

// deskName returns the package and ident name
func descName(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) (string, string) {
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
