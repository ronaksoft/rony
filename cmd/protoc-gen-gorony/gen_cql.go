package main

import (
	"fmt"
	parse "git.ronaksoft.com/ronak/rony/internal/parser"
	"git.ronaksoft.com/ronak/rony/tools"
	"github.com/scylladb/gocqlx/v2/qb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

/*
   Creation Time: 2020 - Aug - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var _ = qb.ASC

var _Models = map[string]*model{}
var _Fields = map[string]struct{}{}

type model struct {
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

// fillModel fills the in global _Models with parsed data
func fillModel(m *protogen.Message) {
	var (
		isModel = false
		mm      = model{
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

func GenCql(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("package ", file.GoPackageName)
	g.P("import (")
	// g.P("\"git.ronaksoft.com/ronak/rony\"")
	g.P("\"git.ronaksoft.com/ronak/rony/pools\"")
	g.P("\"git.ronaksoft.com/ronak/rony/repo/cql\"")
	g.P("\"github.com/scylladb/gocqlx/v2\"")
	g.P("\"github.com/scylladb/gocqlx/v2/qb\"")
	g.P("\"google.golang.org/protobuf/proto\"")
	g.P(")")
	g.P()

	g.P("const (")
	g.P("colData = \"data\"")
	g.P("colVersion = \"ver\"")
	g.P(")")
	g.P()

	genTableDefs(file, g)
	genColumnDefs(file, g)
	genCreateTableCql(file, g)
	genFuncsFactories(file, g)
}
func genTableDefs(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// Tables")
	g.P("const (")
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		g.P("Table", m.Desc.Name(), "= \"", tools.CamelToSnakeASCII(string(m.Desc.Name())), "\"")
		for _, v := range mm.ViewParams {
			g.P("View", m.Desc.Name(), "By", v, "= \"", fmt.Sprintf("%s_by_%s", tools.CamelToSnakeASCII(string(m.Desc.Name())), tools.CamelToSnakeASCII(v)), "\"")
		}
	}
	g.P(")")
	g.P()
}
func genColumnDefs(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// Columns")
	g.P("const (")
	for f := range _Fields {
		g.P("Col", f, " = \"", tools.CamelToSnakeASCII(f), "\"")
	}
	g.P(")")
	g.P()
}
func genCreateTableCql(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// CQLs Create")
	g.P("const (")
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		g.P("_", mm.Name, "CqlCreateTable = `")
		g.P("CREATE TABLE IF NOT EXISTS ", tools.CamelToSnakeASCII(mm.Name), " (")
		for _, fn := range mm.FieldNames {
			g.P(fmt.Sprintf("%s \t %s,", tools.CamelToSnakeASCII(fn), mm.FieldsCql[fn]))
		}
		g.P("data \t blob,")
		pksb := strings.Builder{}
		pksb.WriteRune('(')
		switch {
		case len(mm.Table.PKs)+len(mm.Table.CKs) == 1:
			pksb.WriteString(tools.CamelToSnakeASCII(mm.Table.PKs[0]))
		case len(mm.Table.PKs) == 1:
			pksb.WriteString(tools.CamelToSnakeASCII(mm.Table.PKs[0]))
		default:
			pksb.WriteRune('(')
			for idx, pk := range mm.Table.PKs {
				if idx != 0 {
					pksb.WriteString(", ")
				}
				pksb.WriteString(tools.CamelToSnakeASCII(pk))
			}
			pksb.WriteRune(')')

		}
		for _, ck := range mm.Table.CKs {
			pksb.WriteString(", ")
			pksb.WriteString(tools.CamelToSnakeASCII(ck))
		}
		pksb.WriteRune(')')

		orders := strings.Builder{}
		for idx, k := range mm.Table.Orders {
			kWithoutSign := strings.TrimLeft(k, "-")
			if idx > 0 {
				orders.WriteString(", ")
			}
			if strings.HasPrefix(k, "-") {
				orders.WriteString(fmt.Sprintf("%s DESC", tools.CamelToSnakeASCII(kWithoutSign)))
			} else {
				orders.WriteString(fmt.Sprintf("%s ASC", tools.CamelToSnakeASCII(kWithoutSign)))
			}
		}
		g.P("PRIMARY KEY ", pksb.String())
		g.P(") WITH CLUSTERING ORDER BY (", orders.String(), ");")
		g.P("`")

		// Create Materialized Views
		for idx, v := range mm.Views {
			g.P("_", mm.Name, mm.ViewParams[idx], "CqlCreateMaterializedView = `")
			g.P("CREATE MATERIALIZED VIEW ",
				fmt.Sprintf("%s_by_%s", tools.CamelToSnakeASCII(string(m.Desc.Name())), tools.CamelToSnakeASCII(mm.ViewParams[idx])),
				" AS ",
			)
			g.P("SELECT *")
			g.P("FROM ", tools.CamelToSnakeASCII(mm.Name))
			pksb := strings.Builder{}
			pksb.WriteRune('(')
			switch {
			case len(v.PKs)+len(v.CKs) == 1:
				g.P("WHERE ", tools.CamelToSnakeASCII(v.PKs[0]), " IS NOT null")
				pksb.WriteString(tools.CamelToSnakeASCII(v.PKs[0]))
			case len(v.PKs) == 1:
				g.P("WHERE ", tools.CamelToSnakeASCII(v.PKs[0]), " IS NOT null")
				pksb.WriteString(tools.CamelToSnakeASCII(v.PKs[0]))
			default:
				pksb.WriteRune('(')
				for idx, pk := range v.PKs {
					if idx != 0 {
						g.P("AND ", tools.CamelToSnakeASCII(v.PKs[idx]), " IS NOT null")
						pksb.WriteString(", ")
					} else {
						g.P("WHERE ", tools.CamelToSnakeASCII(v.PKs[idx]), " IS NOT null")
					}
					pksb.WriteString(tools.CamelToSnakeASCII(pk))
				}
				pksb.WriteRune(')')

			}
			for _, ck := range v.CKs {
				g.P("AND ", tools.CamelToSnakeASCII(ck), " IS NOT null")
				pksb.WriteString(", ")
				pksb.WriteString(tools.CamelToSnakeASCII(ck))
			}
			pksb.WriteRune(')')
			g.P("PRIMARY KEY ", pksb.String())
			g.P("`")
		}
	}
	g.P(")")
}
func genFuncsFactories(file *protogen.File, g *protogen.GeneratedFile) {
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		insert(mm, g)
		get(mm, g)
	}

}
func insert(mm *model, g *protogen.GeneratedFile) {
	g.P("var _", mm.Name, "InsertFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
	g.P("return qb.Insert(Table", mm.Name, ").")
	sb := strings.Builder{}
	for _, f := range mm.FieldNames {
		sb.WriteString("Col")
		sb.WriteString(f)
		sb.WriteString(", ")
	}
	sb.WriteString("colData")
	g.P("Columns(", sb.String(), ").")
	g.P("Query(cql.Session())")
	g.P("})")

	g.P("func ", mm.Name, "Insert (x *", mm.Name, ") (err error) {")
	g.P("q := _", mm.Name, "InsertFactory.GetQuery()")
	g.P("defer _", mm.Name, "InsertFactory.Put(q)")
	g.P()
	g.P("mo := proto.MarshalOptions{UseCachedSize: true}")
	g.P("data := pools.Bytes.GetCap(mo.Size(x))")
	g.P("defer pools.Bytes.Put(data)")
	g.P()
	g.P("data, err = mo.MarshalAppend(data, x)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()

	sb = strings.Builder{}
	for _, f := range mm.FieldNames {
		sb.WriteString("x.")
		sb.WriteString(f)
		sb.WriteString(", ")
	}
	sb.WriteString("data")
	g.P("q.Bind(", sb.String(), ")")
	g.P("err = cql.Exec(q)")
	g.P("return err")
	g.P("}")
}
func get(mm *model, g *protogen.GeneratedFile) {
	g.P("var _", mm.Name, "GetFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
	g.P("return qb.Select(Table", mm.Name, ").")
	g.P("Columns(colData).")
	sb := strings.Builder{}
	for idx, f := range mm.FieldNames {
		if idx != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("qb.Eq(Col")
		sb.WriteString(f)
		sb.WriteString(")")
	}
	g.P("Where(", sb.String(), ").")
	g.P("Query(cql.Session())")
	g.P("})")

	sb = strings.Builder{}
	sb2 := strings.Builder{}
	for idx, f := range mm.FieldNames {
		if idx != 0 {
			sb.WriteString(", ")
			sb2.WriteString(", ")
		}
		sb.WriteString(tools.CamelToSnakeASCII(f))
		sb2.WriteString(tools.CamelToSnakeASCII(f))
		sb.WriteRune(' ')
		sb.WriteString(mm.FieldsGo[f])
	}
	g.P("func ", mm.Name, "Get (", sb.String(), ") (m *", mm.Name, ", err error) {")
	g.P("q := _", mm.Name, "GetFactory.GetQuery()")
	g.P("defer _", mm.Name, "GetFactory.Put(q)")
	g.P()
	g.P("data := pools.Bytes.GetCap(512)")
	g.P("defer pools.Bytes.Put(data)")
	g.P()
	g.P("q.Bind(", sb2.String(), ")")
	g.P("err = cql.Scan(q, data)")
	g.P("if err != nil {")
	g.P("return")
	g.P("}")
	g.P()
	g.P("err = proto.UnmarshalOptions{Merge: true}.Unmarshal(data, m)")
	g.P("return m, err")

	g.P("}")
	g.P()
}
func listBy(mm *model, g *protogen.GeneratedFile) {

}
