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

	constTables(file, g)
	initCqlQueries(file, g)
	funcsAndFactories(file, g)
}
func constTables(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// Tables")
	g.P("const (")
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		g.P("Table", m.Desc.Name(), "= \"", tools.ToSnake(string(m.Desc.Name())), "\"")
		for _, v := range mm.ViewParams {
			g.P("View", m.Desc.Name(), "By", v, "= \"", fmt.Sprintf("%s_by_%s", tools.ToSnake(string(m.Desc.Name())), tools.ToSnake(v)), "\"")
		}
	}
	g.P(")")
	g.P()
}

func initCqlQueries(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("func init() {")
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		g.P("cql.AddCqlQuery(`")
		g.P("CREATE TABLE IF NOT EXISTS ", tools.ToSnake(mm.Name), " (")
		for _, fn := range mm.FieldNames {
			g.P(fmt.Sprintf("%s \t %s,", tools.ToSnake(fn), mm.FieldsCql[fn]))
		}
		g.P("data \t blob,")
		pksb := strings.Builder{}
		pksb.WriteRune('(')
		switch {
		case len(mm.Table.PKs)+len(mm.Table.CKs) == 1:
			pksb.WriteString(tools.ToSnake(mm.Table.PKs[0]))
		case len(mm.Table.PKs) == 1:
			pksb.WriteString(tools.ToSnake(mm.Table.PKs[0]))
		default:
			pksb.WriteRune('(')
			for idx, pk := range mm.Table.PKs {
				if idx != 0 {
					pksb.WriteString(", ")
				}
				pksb.WriteString(tools.ToSnake(pk))
			}
			pksb.WriteRune(')')

		}
		for _, ck := range mm.Table.CKs {
			pksb.WriteString(", ")
			pksb.WriteString(tools.ToSnake(ck))
		}
		pksb.WriteRune(')')

		orders := strings.Builder{}
		for idx, k := range mm.Table.Orders {
			kWithoutSign := strings.TrimLeft(k, "-")
			if idx > 0 {
				orders.WriteString(", ")
			}
			if strings.HasPrefix(k, "-") {
				orders.WriteString(fmt.Sprintf("%s DESC", tools.ToSnake(kWithoutSign)))
			} else {
				orders.WriteString(fmt.Sprintf("%s ASC", tools.ToSnake(kWithoutSign)))
			}
		}
		g.P("PRIMARY KEY ", pksb.String())
		if len(mm.Table.Orders) > 0 {
			g.P(") WITH CLUSTERING ORDER BY (", orders.String(), ");")
		} else {
			g.P(");")
		}
		g.P("`)")

		// Create Materialized Views
		for idx, v := range mm.Views {
			g.P("cql.AddCqlQuery(`")
			g.P("CREATE MATERIALIZED VIEW IF NOT EXISTS ",
				fmt.Sprintf("%s_by_%s", tools.ToSnake(string(m.Desc.Name())), tools.ToSnake(mm.ViewParams[idx])),
				" AS ",
			)
			g.P("SELECT *")
			g.P("FROM ", tools.ToSnake(mm.Name))
			pksb := strings.Builder{}
			pksb.WriteRune('(')
			switch {
			case len(v.PKs)+len(v.CKs) == 1:
				g.P("WHERE ", tools.ToSnake(v.PKs[0]), " IS NOT null")
				pksb.WriteString(tools.ToSnake(v.PKs[0]))
			case len(v.PKs) == 1:
				g.P("WHERE ", tools.ToSnake(v.PKs[0]), " IS NOT null")
				pksb.WriteString(tools.ToSnake(v.PKs[0]))
			default:
				pksb.WriteRune('(')
				for idx, pk := range v.PKs {
					if idx != 0 {
						g.P("AND ", tools.ToSnake(v.PKs[idx]), " IS NOT null")
						pksb.WriteString(", ")
					} else {
						g.P("WHERE ", tools.ToSnake(v.PKs[idx]), " IS NOT null")
					}
					pksb.WriteString(tools.ToSnake(pk))
				}
				pksb.WriteRune(')')

			}
			for _, ck := range v.CKs {
				g.P("AND ", tools.ToSnake(ck), " IS NOT null")
				pksb.WriteString(", ")
				pksb.WriteString(tools.ToSnake(ck))
			}
			pksb.WriteRune(')')
			g.P("PRIMARY KEY ", pksb.String())
			g.P("`)")
		}
	}
	g.P("}")

}
func funcsAndFactories(file *protogen.File, g *protogen.GeneratedFile) {
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		funcInsert(mm, g)
		funcGet(mm, g)
		funcListBy(mm, g)
	}
}
func funcInsert(mm *Model, g *protogen.GeneratedFile) {
	g.P("var _", mm.Name, "InsertFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
	g.P("return qb.Insert(Table", mm.Name, ").")
	columns := strings.Builder{}
	binds := strings.Builder{}
	for _, f := range mm.FieldNames {
		columns.WriteString(fmt.Sprintf("%q, ", tools.ToSnake(f)))
		binds.WriteString(fmt.Sprintf("x.%s, ", f))
	}
	columns.WriteString("\"data\"")
	binds.WriteString("b")

	g.P("Columns(", columns.String(), ").")
	g.P("Query(cql.Session())")
	g.P("})")
	g.P("func ", mm.Name, "Insert (x *", mm.Name, ") (err error) {")
	g.P("q := _", mm.Name, "InsertFactory.GetQuery()")
	g.P("defer _", mm.Name, "InsertFactory.Put(q)")
	g.P()
	g.P("mo := proto.MarshalOptions{UseCachedSize: true}")
	g.P("b := pools.Bytes.GetCap(mo.Size(x))")
	g.P("defer pools.Bytes.Put(b)")
	g.P()
	g.P("b, err = mo.MarshalAppend(b, x)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	g.P("q.Bind(", binds.String(), ")")
	g.P("err = cql.Exec(q)")
	g.P("return err")
	g.P("}")
}
func funcGet(mm *Model, g *protogen.GeneratedFile) {
	// Generate Factory
	g.P("var _", mm.Name, "GetFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
	g.P("return qb.Select(Table", mm.Name, ").")
	g.P("Columns(\"data\").")
	where := strings.Builder{}
	args := strings.Builder{}
	bind := strings.Builder{}
	for idx, f := range mm.Table.Keys() {
		if idx != 0 {
			where.WriteString(", ")
			args.WriteString(", ")
			bind.WriteString(", ")
		}
		args.WriteString(fmt.Sprintf("%s %s", tools.ToLowerCamel(f), mm.FieldsGo[f]))
		where.WriteString(fmt.Sprintf("qb.Eq(%q)", tools.ToSnake(f)))
		bind.WriteString(tools.ToLowerCamel(f))
	}
	g.P("Where(", where.String(), ").")
	g.P("Query(cql.Session())")
	g.P("})")
	g.P()

	// Generate Func
	g.P("func ", mm.Name, "Get (", args.String(), ", x *", mm.Name, ") (*", mm.Name, ", error) {")
	g.P("if x == nil {")
	g.P("x = &", mm.Name, "{}")
	g.P("}")
	g.P("q := _", mm.Name, "GetFactory.GetQuery()")
	g.P("defer _", mm.Name, "GetFactory.Put(q)")
	g.P()
	g.P("b := pools.Bytes.GetCap(512)")
	g.P("defer pools.Bytes.Put(b)")
	g.P()
	g.P("q.Bind(", bind.String(), ")")
	g.P("err := cql.Scan(q, &b)")
	g.P("if err != nil {")
	g.P("return x, err")
	g.P("}")
	g.P()
	g.P("err = proto.UnmarshalOptions{Merge:true}.Unmarshal(b, x)")
	g.P("return x, err")
	g.P("}")
	g.P()
}
func funcListBy(mm *Model, g *protogen.GeneratedFile) {
	for idx, v := range mm.ViewParams {
		// Generate Factory
		g.P("var _", mm.Name, "ListBy", v, "Factory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
		g.P("return qb.Select(View", mm.Name, "By", v, ").")
		g.P("Columns(\"data\").")
		where := strings.Builder{}
		args := strings.Builder{}
		bind := strings.Builder{}
		for idx, f := range mm.Views[idx].PKs {
			if idx != 0 {
				bind.WriteString(", ")
				where.WriteString(", ")
				args.WriteString(", ")
			}
			args.WriteString(fmt.Sprintf("%s %s", tools.ToLowerCamel(f), mm.FieldsGo[f]))
			where.WriteString(fmt.Sprintf("qb.Eq(%q)", tools.ToSnake(f)))
			bind.WriteString(tools.ToLowerCamel(f))
		}
		g.P("Where(", where.String(), ").")
		g.P("Query(cql.Session())")
		g.P("})")
		g.P()

		// Generate Function
		g.P("func ", mm.Name, "ListBy", v, " (", args.String(), ", limit int, f func(x *", mm.Name, ") bool) error {")
		g.P("q := _", mm.Name, "ListBy", v, "Factory.GetQuery()")
		g.P("defer _", mm.Name, "ListBy", v, "Factory.Put(q)")
		g.P()
		g.P("b := pools.Bytes.GetCap(512)")
		g.P("defer pools.Bytes.Put(b)")
		g.P()
		g.P("q.Bind(", bind.String(), ")")
		g.P("iter := q.Iter()")
		g.P("for iter.Scan(&b) {")
		g.P("x := Pool", mm.Name, ".Get()")
		g.P("err := proto.UnmarshalOptions{Merge:true}.Unmarshal(b, x)")
		g.P("if err != nil {")
		g.P("Pool", mm.Name, ".Put(x)")
		g.P("return err")
		g.P("}")
		g.P("if limit--; limit <= 0 || !f(x) {")
		g.P("Pool", mm.Name, ".Put(x)")
		g.P("break")
		g.P("}")
		g.P("Pool", mm.Name, ".Put(x)")
		g.P("}")
		g.P("return nil")
		g.P("}")
	}
}
