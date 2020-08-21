package main

import (
	"fmt"
	parse "git.ronaksoftware.com/ronak/rony/internal/parser"
	"github.com/scylladb/go-reflectx"
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
	TableName  string
	Table      PrimaryKey
	ViewNames  []string
	Views      []PrimaryKey
	FieldNames []string
	Fields     map[string]string
}

type PrimaryKey struct {
	PKs []string
	CKs []string
}

func kind(k protoreflect.Kind) string {
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
	}
	panic(fmt.Sprintf("unsupported kind: %v", k.String()))
}

func fillModel(m *protogen.Message) {
	var (
		isModel = false
		mm      = model{
			Fields: make(map[string]string),
		}
	)
	for _, f := range m.Fields {
		mm.Fields[f.GoName] = kind(f.Desc.Kind())
	}

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
			mm.TableName = fmt.Sprintf("t_%s", reflectx.CamelToSnakeASCII(string(m.Desc.Name())))
			nn := n.(*parse.TableNode)
			for _, k := range nn.PartitionKeys {
				fields[k] = struct{}{}
				mm.Table.PKs = append(mm.Table.PKs, k)

			}
			for _, k := range nn.ClusteringKeys {
				fields[k] = struct{}{}
				mm.Table.CKs = append(mm.Table.CKs, k)
			}
		case parse.NodeView:
			pk := PrimaryKey{}
			nn := n.(*parse.ViewNode)
			sb := strings.Builder{}
			for idx, k := range nn.PartitionKeys {
				fields[k] = struct{}{}
				if idx != 0 {
					sb.WriteString("_")
				}
				sb.WriteString(reflectx.CamelToSnakeASCII(k))
			}
			mm.ViewNames = append(mm.ViewNames, fmt.Sprintf("mv_%s_by_%s", reflectx.CamelToSnakeASCII(string(m.Desc.Name())), sb.String()))
			for _, k := range nn.ClusteringKeys {
				fields[k] = struct{}{}
				pk.CKs = append(pk.CKs, k)
			}
			mm.Views = append(mm.Views, pk)
		}
	}
	for f := range fields {
		_Fields[f] = struct{}{}
		mm.FieldNames = append(mm.FieldNames, f)
	}
	if isModel {
		mm.Name = string(m.Desc.Name())
		_Models[mm.Name] = &mm
	}
}

func GenCql(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("package ", file.GoPackageName)
	g.P("import (")
	// g.P("\"git.ronaksoftware.com/ronak/rony\"")
	g.P("\"git.ronaksoftware.com/ronak/rony/pools\"")
	g.P("\"git.ronaksoftware.com/ronak/rony/repo/cql\"")
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

	for _, m := range file.Messages {
		fillModel(m)
	}

	genTableDefs(file, g)
	genColumnDefs(file, g)
	genCqls(file, g)
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
		g.P("Table", m.Desc.Name(), "= \"t_", reflectx.CamelToSnakeASCII(string(m.Desc.Name())), "\"")
		for _, v := range mm.ViewNames {
			g.P("View", m.Desc.Name(), "= \"", v, "\"")
		}
	}
	g.P(")")
	g.P()
}

func genColumnDefs(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// Columns")
	g.P("const (")
	for f := range _Fields {
		g.P("Col", f, " = \"", reflectx.CamelToSnakeASCII(f), "\"")
	}
	g.P(")")
	g.P()
}

func genCqls(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// CQLs Create")
	g.P("const (")
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		g.P("_", mm.Name, "CqlCreateTable = `")
		g.P("CREATE TABLE IF NOT EXISTS ", mm.TableName, " (")
		for _, fn := range mm.FieldNames {
			g.P(fmt.Sprintf("%s \t %s,", reflectx.CamelToSnakeASCII(fn), mm.Fields[fn]))
		}
		g.P("data \t blob,")
		pksb := strings.Builder{}
		pksb.WriteRune('(')
		switch {
		case len(mm.Table.PKs)+len(mm.Table.CKs) == 1:
			pksb.WriteString(reflectx.CamelToSnakeASCII(mm.Table.PKs[0]))
		case len(mm.Table.PKs) == 1:
			pksb.WriteString(reflectx.CamelToSnakeASCII(mm.Table.PKs[0]))
		default:
			pksb.WriteRune('(')
			for idx, pk := range mm.Table.PKs {
				if idx != 0 {
					pksb.WriteString(", ")
				}
				pksb.WriteString(reflectx.CamelToSnakeASCII(pk))
			}
			pksb.WriteRune(')')

		}
		for _, ck := range mm.Table.CKs {
			pksb.WriteString(", ")
			pksb.WriteString(reflectx.CamelToSnakeASCII(ck))
		}
		pksb.WriteRune(')')

		g.P("PRIMARY KEY ", pksb.String())
		g.P(")")
		g.P("`")
	}
	g.P(")")
}

func genFuncsFactories(file *protogen.File, g *protogen.GeneratedFile) {
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
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

}
