package cql

import (
	"github.com/jinzhu/inflection"
	"github.com/ronaksoft/rony/internal/codegen"
	parse "github.com/ronaksoft/rony/internal/parser"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
	"text/template"
)

/*
   Creation Time: 2021 - Jul - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type genType int

const (
	GO genType = iota
	CQL
)

type Generator struct {
	f *protogen.File
	g *protogen.GeneratedFile
	t genType
}

func NewGO(f *protogen.File, g *protogen.GeneratedFile) *Generator {
	return &Generator{
		f: f,
		g: g,
		t: GO,
	}
}

func NewCQL(f *protogen.File, g *protogen.GeneratedFile) *Generator {
	return &Generator{
		f: f,
		g: g,
		t: CQL,
	}
}

func Check(f *protogen.File) bool {
	for _, m := range f.Messages {
		t, _ := codegen.Parse(m)
		if t == nil {
			continue
		}
		for _, n := range t.Root.Nodes {
			switch n.Type() {
			case parse.NodeModel:
				return true
			}
		}
	}
	return false
}

func (g *Generator) Generate() {
	switch g.t {
	case GO:
		g.generateGo()
	case CQL:
		g.generateCql()

	}
}

func (g *Generator) generateGo() {
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "bytes"})
	g.g.P("var _ = bytes.MinRead")

	for _, m := range g.f.Messages {
		arg := codegen.GetMessageArg(g.f, g.g, m)
		if arg.IsSingleton {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})

			// g.g.P(g.Exec(template.Must(template.New("genSingletonSave").Funcs(singletonFuncs).Parse(genSingletonSave)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genSingletonRead").Funcs(singletonFuncs).Parse(genSingletonRead)), arg))
		} else if arg.IsAggregate {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})

			// g.g.P(g.Exec(template.Must(template.New("genAggregateCreate").Funcs(aggregateFuncs).Parse(genAggregateCreate)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genAggregateUpdate").Funcs(aggregateFuncs).Parse(genAggregateUpdate)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genAggregateSave").Funcs(aggregateFuncs).Parse(genAggregateSave)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genAggregateRead").Funcs(aggregateFuncs).Parse(genAggregateRead)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genAggregateDelete").Funcs(aggregateFuncs).Parse(genAggregateDelete)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genAggregateHelpers").Funcs(aggregateFuncs).Parse(genAggregateHelpers)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genIter").Funcs(aggregateFuncs).Parse(genIter)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genList").Funcs(aggregateFuncs).Parse(genList)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genIterByPK").Funcs(aggregateFuncs).Parse(genIterByPK)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genIListByPK").Funcs(aggregateFuncs).Parse(genListByPK)), arg))
			// g.g.P(g.Exec(template.Must(template.New("genListByIndex").Funcs(aggregateFuncs).Parse(genListByIndex)), arg))
		}
	}
}

func (g *Generator) generateCql() {
	funcs := map[string]interface{}{
		"Singular": func(x string) string {
			return inflection.Singular(x)
		},
		"Plural": func(x string) string {
			return inflection.Plural(x)
		},
		"PrimaryKey": func(m codegen.ModelKey) string {
			sb := strings.Builder{}
			sb.WriteString("((")
			sb.WriteString(m.Names(codegen.PropFilterPKs, "", ", ", codegen.SnakeCase))
			sb.WriteString(")")
			if len(m.CKs()) > 0 {
				sb.WriteString(", ")
				sb.WriteString(m.Names(codegen.PropFilterCKs, "", ", ", codegen.SnakeCase))
			}
			sb.WriteString(")")
			return sb.String()
		},
		"WithClusteringKey": func(m codegen.ModelKey) string {
			if len(m.CKs()) == 0 {
				return ""
			}
			sb := strings.Builder{}
			sb.WriteString(" WITH CLUSTERING ORDER BY (")
			for idx, k := range m.CKs() {
				if idx > 0 {
					sb.WriteString(", ")
				}
				sb.WriteString(tools.ToSnake(k.Name))
				sb.WriteRune(' ')
				if k.Order == codegen.ASC {
					sb.WriteString("ASC")
				} else {
					sb.WriteString("DESC")
				}
			}
			sb.WriteString(")")
			return sb.String()
		},
		"MVName": func(m codegen.ModelKey) string {
			sb := strings.Builder{}
			sb.WriteString(tools.ToSnake(m.Name()))
			sb.WriteString("_by_")
			sb.WriteString(m.Names(codegen.PropFilterPKs, "", "_", codegen.SnakeCase))
			return sb.String()
		},
		"MVWhere": func(m codegen.ModelKey) string {
			sb := strings.Builder{}
			for idx, k := range m.Keys() {
				if idx == 0 {
					sb.WriteString("WHERE ")
				} else {
					sb.WriteString("\r\n")
					sb.WriteString("AND ")
				}
				sb.WriteString(tools.ToSnake(k.Name))
				sb.WriteString(" IS NOT null")
			}
			return sb.String()
		},
	}
	for _, m := range g.f.Messages {
		arg := codegen.GetMessageArg(g.f, g.g, m)
		if arg.IsAggregate {
			// g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx/v2/table"})
			g.g.P(g.Exec(template.Must(template.New("genCQL").Funcs(funcs).Parse(genCQL)), arg))
		}
	}
}

func (g *Generator) Exec(t *template.Template, v interface{}) string {
	sb := &strings.Builder{}
	if err := t.Execute(sb, v); err != nil {
		panic(err)
	}

	return sb.String()
}

func x() {
	// codegen.MessageArg{}.Table.
}

const genCQL = `
{{$model := .}}
CREATE TABLE tab_{{Singular .NameSC}} 
(
{{- range .Fields }}
	{{.NameSC}}  {{.CqlKind}},	
{{- end }}
	sdata  blob,
	PRIMARY KEY {{PrimaryKey .Table}}
){{WithClusteringKey .Table}};
{{ range .Views }}
CREATE MATERIALIZED VIEW view_{{MVName .}} AS
SELECT *
FROM tab_{{Singular $model.NameSC}}
{{MVWhere .}}
PRIMARY KEY {{PrimaryKey .}}
{{- WithClusteringKey . -}};
{{ end }}
`

const genCreate = ``
const genCreateIF = ``
const genUpdate = ``
const genDelete = ``
const genRead = ``
const genListByPK = ``
const genListByIndex = ``
