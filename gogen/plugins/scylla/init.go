package scylla

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gogen"
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
	g *gogen.Generator
}

func (s *RepoPlugin) Name() string {
	return "Scylla"
}

func (s *RepoPlugin) Init(g *gogen.Generator) {
	s.g = g
}

func (s *RepoPlugin) Generate(desc *gogen.Descriptor) {
	s.g.P("var (")
	s.g.In()
	s.g.P("dbs	*gocql.Session")
	for _, m := range desc.Models {
		// Validate the model and panic if it is not validated
		validateModel(&m)

		for _, fk := range append(m.FilterKeys, m.PrimaryKey) {
			tbVarName := fmt.Sprintf("tb%sBy%s", strcase.ToCamel(m.Name), strcase.ToCamel(fk.Name))
			tbValName := fmt.Sprintf("\"%s_by_%s\"", strcase.ToSnake(m.Name), strcase.ToSnake(fk.Name))
			qpVarName := fmt.Sprintf("qpGet%sBy%s", strcase.ToCamel(m.Name), strcase.ToCamel(fk.Name))

			s.g.P(tbVarName, "=", tbValName)
			s.g.P(qpVarName, "=", "sync.Pool {")
			s.g.In()
			s.g.P("New: func() interface{} {")
			s.g.In()
			s.g.Pns("stmt, name := qb.Select(", tbVarName, ").")
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

	// Define Init Function
	s.g.P("func Init(s *gocql.Session) {")
	s.g.In()
	s.g.P("dbs = s")
	s.g.Out()
	s.g.P("}")
	s.g.Nl(2) // End of Init

	s.generateCreateTables(desc)
	s.generateInteractions(desc)
}
func (s *RepoPlugin) generateCreateTables(desc *gogen.Descriptor) {
	// Define CreateTables Function
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

	for _, m := range desc.Models {
		tbValName := fmt.Sprintf("%s_by_%s", strcase.ToSnake(m.Name), strcase.ToSnake(m.PrimaryKey.Name))
		primaryKey := fmt.Sprintf("(%s)", strings.Join(snakeCaseAll(m.PrimaryKey.PartitionKeys), ","))
		if len(m.PrimaryKey.ClusteringKeys) > 0 {
			primaryKey = fmt.Sprintf("(%s, %s)", primaryKey, strings.Join(snakeCaseAll(m.PrimaryKey.ClusteringKeys), ","))
		}
		s.g.P("stmtB.Reset()")
		s.g.P("stmtB.WriteString(\"CREATE TABLE IF NOT EXISTS \")")
		s.g.Pns("stmtB.WriteString(fmt.Sprintf(\"%s.", tbValName, "\\n\", db))")
		for _, pn := range append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...) {
			p, err := m.GetProperty(pn)
			if err != nil {
				panic(err)
			}
			s.g.Pns(
				"stmtB.WriteString(\"",
				fmt.Sprintf("%s %s,", strcase.ToSnake(p.Name), p.ToScyllaType()),
				"\\n\")",
			)
		}
		s.g.P("stmtB.WriteString(\"data blob,\\n\")")
		s.g.Pns("stmtB.WriteString(\"PRIMARY KEY ", primaryKey, "\")")
		s.g.P("q = s.Query(stmtB.String()).RetryPolicy(nil)")
		s.g.P("err = q.Exec()")
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
					s.g.P("stmtB.WriteString(\"WHERE", strcase.ToSnake(p.Name), "IS NOT NULL\\n\")")
				} else {
					s.g.P("stmtB.WriteString(\"AND", strcase.ToSnake(p.Name), "IS NOT NULL\\n\")")
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
}

func (s *RepoPlugin) generateInteractions(desc *gogen.Descriptor) {

}

func (s *RepoPlugin) GeneratePrepend(desc *gogen.Descriptor) {
	s.g.P("package", strcase.ToLowerCamel(desc.Name))
	s.g.Nl()
	s.g.P("import (")
	s.g.In()
	s.g.P("\"fmt\"")
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

/*
	Helper Functions
*/

func snakeCaseAll(s []string) []string {
	o := make([]string, len(s))
	for idx := range s {
		o[idx] = strcase.ToSnake(s[idx])
	}
	return o
}

func validateModel(m *gogen.Model) {
	keys := append(m.PrimaryKey.PartitionKeys, m.PrimaryKey.ClusteringKeys...)
	for _, fk := range m.FilterKeys {
		for _, k := range append(fk.PartitionKeys, fk.ClusteringKeys...) {
			if !inArray(keys, k) {
				panic(fmt.Sprintf("key %s(%s) does not exists in the PrimaryKey for model (%s)", fk.Name, k, m))
			}
		}
	}
}

func inArray(arr []string, v string) bool {
	for i := range arr {
		if arr[i] == v {
			return true
		}
	}
	return false
}
