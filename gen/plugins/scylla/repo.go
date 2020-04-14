package scylla

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gen"
	"github.com/iancoleman/strcase"
	"strings"
)

/*
   Creation Time: 2020 - Apr - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type RepoPlugin struct {
	g *gen.Generator
}

func (s *RepoPlugin) Name() string {
	return "Scylla"
}

func (s *RepoPlugin) Init(g *gen.Generator) {
	s.g = g
}

func (s *RepoPlugin) Generate(desc *gen.Descriptor) {

	s.generateVars(desc)

	// Generate Init Function
	s.g.P("func Init(s *gocql.Session) {")
	s.g.In()
	s.g.P("dbs = s")
	s.g.Out()
	s.g.P("}")
	s.g.Nl() // End of Init

	s.generateCreateTables(desc)
	s.generateSaveModel(desc)
	s.generateGet(desc)
	s.generateList(desc)
	s.generateDelete(desc)

}
func (s *RepoPlugin) generateVars(desc *gen.Descriptor) {
	s.g.P("var (")
	s.g.In()
	s.g.P("dbs	*gocql.Session")
	for _, m := range desc.Models {
		// Validate the model and panic if it is not validated
		validateModel(&m)
		tbVarName := fmt.Sprintf("tb%s", strcase.ToCamel(m.Name))
		tbValName := fmt.Sprintf("\"%s\"", strcase.ToSnake(m.Name))
		s.g.P(tbVarName, "=", tbValName)

		// Define Get
		qpVarName := fmt.Sprintf("qpGet%s", strcase.ToCamel(m.Name))
		s.g.P(qpVarName, "=", "sync.Pool {")
		s.g.In()
		s.g.P("New: func() interface{} {")
		s.g.In()
		s.g.Pns("stmt, name := qb.Select(", tbVarName, ").")
		s.g.P("Columns(\"data\").")
		s.g.P("Where(")
		s.g.In()
		for _, k := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			s.g.Pns("qb.Eq(\"", strcase.ToSnake(k), "\"),")
		}
		s.g.Out()
		s.g.P(").ToCql()")
		s.g.P("return gocqlx.Query(dbs.Query(stmt), name)")
		s.g.Out()
		s.g.P("},")
		s.g.Out()
		s.g.P("}")

		// Define Save
		qpVarName = fmt.Sprintf("qpSave%s", strcase.ToCamel(m.Name))
		s.g.P(qpVarName, "=", "sync.Pool {")
		s.g.In()
		s.g.P("New: func() interface{} {")
		s.g.In()
		s.g.Pns("stmt, name := qb.Insert(", tbVarName, ").")
		colNames := strings.Builder{}
		for _, k := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			colNames.WriteRune('"')
			colNames.WriteString(strcase.ToSnake(k))
			colNames.WriteRune('"')
			colNames.WriteRune(',')
		}
		colNames.WriteString("\"data\"")
		s.g.Pns("Columns(", colNames.String(), ").")
		s.g.P("ToCql()")
		s.g.P("return gocqlx.Query(dbs.Query(stmt), name)")
		s.g.Out()
		s.g.P("},")
		s.g.Out()
		s.g.P("}")

		// Define Delete
		qpVarName = fmt.Sprintf("qpDelete%s", strcase.ToCamel(m.Name))
		s.g.P(qpVarName, "=", "sync.Pool {")
		s.g.In()
		s.g.P("New: func() interface{} {")
		s.g.In()
		s.g.Pns("stmt, name := qb.Delete(", tbVarName, ").")
		s.g.P("Where(")
		s.g.In()
		for _, k := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			s.g.Pns("qb.Eq(\"", strcase.ToSnake(k), "\"),")
		}
		s.g.Out()
		s.g.P(").ToCql()")
		s.g.P("return gocqlx.Query(dbs.Query(stmt), name)")
		s.g.Out()
		s.g.P("},")
		s.g.Out()
		s.g.P("}")

		for _, fk := range append(m.FilterKeys) {
			tbVarName := fmt.Sprintf("tb%sBy%s", strcase.ToCamel(m.Name), strcase.ToCamel(fk.Name))
			tbValName := fmt.Sprintf("\"%s_by_%s\"", strcase.ToSnake(m.Name), strcase.ToSnake(fk.Name))
			s.g.P(tbVarName, "=", tbValName)

			// Define GetBy
			qpVarName := fmt.Sprintf("qpGet%sBy%s", strcase.ToCamel(m.Name), strcase.ToCamel(fk.Name))
			s.g.P(qpVarName, "=", "sync.Pool {")
			s.g.In()
			s.g.P("New: func() interface{} {")
			s.g.In()
			s.g.Pns("stmt, name := qb.Select(", tbVarName, ").")
			s.g.P("Columns(\"data\").")
			s.g.P("Where(")
			s.g.In()
			for _, k := range append(fk.PartitionKeys, fk.ClusteringKeys...) {
				s.g.Pns("qb.Eq(\"", strcase.ToSnake(k), "\"),")
			}
			s.g.Out()
			s.g.P(").ToCql()")
			s.g.P("return gocqlx.Query(dbs.Query(stmt), name)")
			s.g.Out()
			s.g.P("},")
			s.g.Out()
			s.g.P("}")

		}
	}
	s.g.Out()
	s.g.P(")")
	s.g.Nl(2) // End of Vars
}
func (s *RepoPlugin) generateCreateTables(desc *gen.Descriptor) {
	// Generate CreateTables Function
	s.g.P("func CreateTables(s *gocql.Session, db string) error {")
	s.g.In()

	// Define Vars
	s.g.P("var (")
	s.g.In()
	s.g.P("q     *gocql.Query")
	s.g.P("err   error")
	s.g.P("stmtB strings.Builder")
	s.g.Out()
	s.g.P(")")
	s.g.Nl() // End of Vars

	s.g.P("createKeySpaceQuery := fmt.Sprintf(\"CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}\", db)")
	s.g.P("err = s.Query(createKeySpaceQuery).Exec()")
	s.g.P("if err != nil {")
	s.g.In()
	s.g.P("return err ")
	s.g.Out()
	s.g.P("}")

	for _, m := range desc.Models {
		tbValName := fmt.Sprintf("%s", strcase.ToSnake(m.Name))
		primaryKey := fmt.Sprintf("(%s)", strings.Join(snakeCaseAll(m.PrimaryKey.PartitionKeys), ","))
		if len(m.PrimaryKey.ClusteringKeys) > 0 {
			primaryKey = fmt.Sprintf("(%s, %s)", primaryKey, strings.Join(snakeCaseAll(m.PrimaryKey.ClusteringKeys), ","))
		}
		s.g.P("stmtB.Reset()")
		s.g.P("stmtB.WriteString(\"CREATE TABLE IF NOT EXISTS \")")
		s.g.Pns("stmtB.WriteString(fmt.Sprintf(\"%s.", tbValName, " (\\n\", db))")
		for _, pn := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			p, err := m.GetProperty(pn)
			if err != nil {
				panic(err)
			}
			s.g.Pns(
				"stmtB.WriteString(\"",
				fmt.Sprintf("%s %s,", strcase.ToSnake(p.Def.Name()), p.ToScyllaType()),
				"\\n\")",
			)
		}
		s.g.P("stmtB.WriteString(\"data blob,\\n\")")
		s.g.Pns("stmtB.WriteString(\"PRIMARY KEY ", primaryKey, "\")")
		s.g.P("stmtB.WriteString(\");\")")
		s.g.P("q = s.Query(stmtB.String()).RetryPolicy(nil)")
		s.g.P("err = q.Exec()")
		s.g.P("q.Release()")
		s.g.P("if err != nil {")
		s.g.In()
		s.g.P("return err")
		s.g.Out()
		s.g.P("}")
		s.g.Nl()

		for _, fk := range m.FilterKeys {
			mvValName := fmt.Sprintf("%s_by_%s", strcase.ToSnake(m.Name), strcase.ToSnake(fk.Name))
			primaryKey := fmt.Sprintf("(%s)", strings.Join(snakeCaseAll(fk.PartitionKeys), ","))
			if len(fk.ClusteringKeys) > 0 {
				primaryKey = fmt.Sprintf("(%s, %s)", primaryKey, strings.Join(snakeCaseAll(fk.ClusteringKeys), ","))
			}
			s.g.P("stmtB.Reset()")
			s.g.P("stmtB.WriteString(\"CREATE MATERIALIZED VIEW IF NOT EXISTS \")")
			s.g.Pns("stmtB.WriteString(fmt.Sprintf(\"%s.", mvValName, " AS\\n\", db))")
			s.g.P("stmtB.WriteString(\"SELECT *\\n\")")
			s.g.Pns("stmtB.WriteString(fmt.Sprintf(\"FROM %s.", tbValName, "\\n\", db))")
			for idx, pn := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
				p, err := m.GetProperty(pn)
				if err != nil {
					panic(err)
				}
				if idx == 0 {
					s.g.P("stmtB.WriteString(\"WHERE", strcase.ToSnake(p.Def.Name()), "IS NOT NULL\\n\")")
				} else {
					s.g.P("stmtB.WriteString(\"AND", strcase.ToSnake(p.Def.Name()), "IS NOT NULL\\n\")")
				}
			}
			s.g.Pns("stmtB.WriteString(\"PRIMARY KEY ", primaryKey, "\")")
			s.g.P("q = s.Query(stmtB.String()).RetryPolicy(nil)")
			s.g.P("err = q.Exec()")
			s.g.P("if err != nil {")
			s.g.In()
			s.g.P("return err")
			s.g.Out()
			s.g.P("}")
			s.g.Nl()
		}

	}

	s.g.P("return nil")
	s.g.Out()
	s.g.P("}") // End of CreateTables
	s.g.Nl()
}
func (s *RepoPlugin) generateSaveModel(desc *gen.Descriptor) {
	for _, m := range desc.Models {
		s.g.Pns("func Save", strcase.ToCamel(m.Name), "(x *", m.Name, ") error {")
		s.g.In()
		s.g.P("b := pbytes.GetLen(x.Size())")
		s.g.P("defer pbytes.Put(b)")
		s.g.P("n, err := x.MarshalTo(b)")
		s.g.P("if err != nil {")
		s.g.In()
		s.g.P("return err")
		s.g.Out()
		s.g.P("}")
		s.g.Nl()
		qpVarName := fmt.Sprintf("qpSave%s", strcase.ToCamel(m.Name))
		s.g.P(fmt.Sprintf("q := %s.Get().(*gocqlx.Queryx)", qpVarName))
		s.g.P(fmt.Sprintf("defer %s.Put(q)", qpVarName))
		s.g.Nl()
		binds := strings.Builder{}
		for _, pn := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			binds.WriteString("x.")
			binds.WriteString(strcase.ToCamel(pn))
			binds.WriteString(", ")
		}
		binds.WriteString("b[:n]")
		s.g.P(fmt.Sprintf("q.Bind(%s)", binds.String()))
		s.g.P("err = q.Exec()")
		s.g.P("return err")
		s.g.Out()
		s.g.P("}")
		s.g.Nl()
	}
}
func (s *RepoPlugin) generateGet(desc *gen.Descriptor) {
	for _, m := range desc.Models {
		args := strings.Builder{}
		for idx, pn := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			p, _ := m.GetProperty(pn)
			if idx != 0 {
				args.WriteRune(',')
			}
			args.WriteString(strcase.ToLowerCamel(p.Def.Name()))
			args.WriteRune(' ')
			if p.CheckOption(gen.Slice) {
				args.WriteString("[]")
			}
			args.WriteString(p.Def.Type())
		}
		s.g.Pns("func Get", strcase.ToCamel(m.Name), "(", args.String(), ") (*", m.Name, ", error) {")
		s.g.In()
		s.g.P("var (")
		s.g.In()
		s.g.P("b []byte")
		s.g.Out()
		s.g.P(")")
		s.g.Nl()
		qpVarName := fmt.Sprintf("qpGet%s", strcase.ToCamel(m.Name))
		s.g.P(fmt.Sprintf("q := %s.Get().(*gocqlx.Queryx)", qpVarName))
		s.g.P(fmt.Sprintf("defer %s.Put(q)", qpVarName))
		s.g.Nl()
		binds := strings.Builder{}
		for idx, pn := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			if idx != 0 {
				binds.WriteString(", ")
			}
			binds.WriteString(strcase.ToLowerCamel(pn))
		}
		s.g.P(fmt.Sprintf("q.Bind(%s)", binds.String()))
		s.g.P("err := q.Scan(&b)")
		s.g.P("if err != nil {")
		s.g.In()
		s.g.P("return nil, err")
		s.g.Out()
		s.g.P("}")
		s.g.Nl()

		s.g.Pns("v := &", strcase.ToCamel(m.Name), "{}")
		s.g.P("err = v.Unmarshal(b)")
		s.g.P("if err != nil {")
		s.g.In()
		s.g.P("return nil, err")
		s.g.Out()
		s.g.P("}")
		s.g.Nl()
		s.g.P("return v, nil")
		s.g.Out()
		s.g.P("}")
		s.g.Nl()

		for _, fk := range m.FilterKeys {
			args := strings.Builder{}
			for idx, pn := range append(fk.PartitionKeys, fk.ClusteringKeys...) {
				p, _ := m.GetProperty(pn)
				if idx != 0 {
					args.WriteRune(',')
				}
				args.WriteString(strcase.ToLowerCamel(p.Def.Name()))
				args.WriteRune(' ')
				if p.CheckOption(gen.Slice) {
					args.WriteString("[]")
				}
				args.WriteString(p.Def.Type())
			}
			s.g.Pns("func Get", strcase.ToCamel(m.Name), "By", strcase.ToCamel(fk.Name), "(", args.String(), ") (*", m.Name, ", error) {")
			s.g.In()
			s.g.P("var (")
			s.g.In()
			s.g.P("b []byte")
			s.g.Out()
			s.g.P(")")
			s.g.Nl()
			qpVarName := fmt.Sprintf("qpGet%sBy%s", strcase.ToCamel(m.Name), strcase.ToCamel(fk.Name))
			s.g.P(fmt.Sprintf("q := %s.Get().(*gocqlx.Queryx)", qpVarName))
			s.g.P(fmt.Sprintf("defer %s.Put(q)", qpVarName))
			s.g.Nl()
			binds := strings.Builder{}
			for idx, pn := range append(fk.PartitionKeys, fk.ClusteringKeys...) {
				if idx != 0 {
					binds.WriteString(", ")
				}
				binds.WriteString(strcase.ToLowerCamel(pn))
			}
			s.g.P(fmt.Sprintf("q.Bind(%s)", binds.String()))
			s.g.P("err := q.Scan(&b)")
			s.g.P("if err != nil {")
			s.g.In()
			s.g.P("return nil, err")
			s.g.Out()
			s.g.P("}")
			s.g.Nl()

			s.g.Pns("v := &", strcase.ToCamel(m.Name), "{}")
			s.g.P("err = v.Unmarshal(b)")
			s.g.P("if err != nil {")
			s.g.In()
			s.g.P("return nil, err")
			s.g.Out()
			s.g.P("}")
			s.g.Nl()
			s.g.P("return v, nil")
			s.g.Out()
			s.g.P("}")
			s.g.Nl()
		}
	}
}
func (s *RepoPlugin) generateDelete(desc *gen.Descriptor) {
	for _, m := range desc.Models {
		args := strings.Builder{}
		for idx, pn := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			p, _ := m.GetProperty(pn)
			if idx != 0 {
				args.WriteRune(',')
			}
			args.WriteString(strcase.ToLowerCamel(p.Def.Name()))
			args.WriteRune(' ')
			if p.CheckOption(gen.Slice) {
				args.WriteString("[]")
			}
			args.WriteString(p.Def.Type())
		}
		s.g.Pns("func Delete", strcase.ToCamel(m.Name), "(", args.String(), ") error {")
		s.g.In()
		qpVarName := fmt.Sprintf("qpDelete%s", strcase.ToCamel(m.Name))
		s.g.P(fmt.Sprintf("q := %s.Get().(*gocqlx.Queryx)", qpVarName))
		binds := strings.Builder{}
		for idx, pn := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			if idx != 0 {
				binds.WriteString(", ")
			}
			binds.WriteString(strcase.ToLowerCamel(pn))
		}
		s.g.P(fmt.Sprintf("q.Bind(%s)", binds.String()))
		s.g.P("err := q.Exec()")
		s.g.P(fmt.Sprintf("%s.Put(q)", qpVarName))
		s.g.P("return err")
		s.g.Out()
		s.g.P("}")
		s.g.Nl()
	}
}
func (s *RepoPlugin) generateList(desc *gen.Descriptor) {

}

func (s *RepoPlugin) GeneratePrepend(desc *gen.Descriptor) {
	s.g.P("package", strcase.ToLowerCamel(desc.Name))
	s.g.Nl()
	s.g.P("import (")
	s.g.In()
	s.g.P("\"fmt\"")
	s.g.P("\"github.com/gobwas/pool/pbytes\"")
	s.g.P("\"github.com/gocql/gocql\"")
	s.g.P("\"github.com/scylladb/gocqlx\"")
	s.g.P("\"github.com/scylladb/gocqlx/qb\"")
	s.g.P("\"strings\"")
	s.g.P("\"sync\"")
	s.g.Out()
	s.g.P(")")
	s.g.Nl(2)
	s.g.P("/*")
	s.g.In()
	s.g.P("   Creation Time: YYYY - MMM - DD")
	s.g.P("   Auto Generated by Rony's Code Generator")
	s.g.Out()
	s.g.P("*/")

}
