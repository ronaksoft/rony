package cqlmodel

import (
	"fmt"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/aggregate"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
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

// Generate generates the repo functions for messages which are identified as model with {{@entity cql}}
func Generate(file *protogen.File, g *protogen.GeneratedFile) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/pools"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/repo/cql"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx/v2/qb"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})

	constTables(file, g)
	initCqlQueries(file, g)
	funcsAndFactories(file, g)
}
func constTables(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// Tables")
	g.P("const (")
	for _, m := range file.Messages {
		mm := aggregate.GetAggregates()[string(m.Desc.Name())]
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
		mm := aggregate.GetAggregates()[string(m.Desc.Name())]
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
		mm := aggregate.GetAggregates()[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		funcInsert(mm, g)
		funcGet(mm, g)
		funcListBy(mm, g)
	}
}
func funcInsert(mm *aggregate.Aggregate, g *protogen.GeneratedFile) {
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
func funcGet(mm *aggregate.Aggregate, g *protogen.GeneratedFile) {
	// Generate Factory
	g.P("var _", mm.Name, "GetFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
	g.P("return qb.Select(Table", mm.Name, ").")
	g.P("Columns(\"data\").")
	where := strings.Builder{}
	args := strings.Builder{}
	bind := strings.Builder{}
	for idx, f := range mm.Table.PKs {
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
func funcListBy(mm *aggregate.Aggregate, g *protogen.GeneratedFile) {
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
		g.P("func ", mm.Name, "ListBy", v, " (", args.String(), ", limit int32, f func(x *", mm.Name, ") bool) error {")
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
