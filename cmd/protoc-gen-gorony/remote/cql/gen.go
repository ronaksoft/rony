package cql

import (
	"github.com/jinzhu/inflection"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/codegen"
	parse "github.com/ronaksoft/rony/internal/parser"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
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
		opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)
		if proto.GetExtension(opt, rony.E_RonyTable).(*rony.PrimaryKeyOpt) != nil {
			return true
		}
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
	for _, m := range g.f.Messages {
		arg := codegen.GetMessageArg(g.f, g.g, m)
		funcs := map[string]interface{}{
			"Singular": func(x string) string {
				return inflection.Singular(x)
			},
			"Plural": func(x string) string {
				return inflection.Plural(x)
			},
			"MVNameSC": func(m codegen.ModelKey) string {
				sb := strings.Builder{}
				sb.WriteString(tools.ToSnake(m.Name()))
				sb.WriteString("_by_")
				sb.WriteString(m.Names(codegen.PropFilterPKs, "", "", "_", codegen.SnakeCase))
				if m.OrderByAlias() != "" {
					sb.WriteString("_")
					sb.WriteString(tools.ToSnake(m.OrderByAlias()))
				} else if m.Index() > 0 {
					sb.WriteString("_")
					sb.WriteString(tools.IntToStr(m.Index()))
				}
				return sb.String()
			},
			"MVName": func(m codegen.ModelKey) string {
				sb := strings.Builder{}
				sb.WriteString(m.Name())
				sb.WriteString("By")
				sb.WriteString(m.Names(codegen.PropFilterPKs, "", "", "", codegen.None))
				if m.OrderByAlias() != "" {
					sb.WriteString(m.OrderByAlias())
				} else if m.Index() > 0 {
					sb.WriteString(tools.IntToStr(m.Index()))
				}
				return sb.String()
			},
			"Columns": func(m codegen.ModelKey) string {
				sb := strings.Builder{}
				sb.WriteString(m.Names(codegen.PropFilterALL, "\"", "\"", ", ", codegen.SnakeCase))
				sb.WriteString(", \"sdata\"")
				return sb.String()
			},
			"PartKeys": func(m codegen.ModelKey) string {
				return m.Names(codegen.PropFilterPKs, "\"", "\"", ", ", codegen.SnakeCase)
			},
			"SortKeys": func(m codegen.ModelKey) string {
				return m.Names(codegen.PropFilterCKs, "\"", "\"", ", ", codegen.SnakeCase)
			},
		}
		if arg.IsAggregate {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx/v2/table"})

			g.g.P(g.Exec(template.Must(template.New("genTable").Funcs(funcs).Parse(genTable)), arg))
			g.g.P(g.Exec(template.Must(template.New("genRemoteRepo").Funcs(funcs).Parse(genRemoteRepo)), arg))
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
			sb.WriteString(m.Names(codegen.PropFilterPKs, "", "", ", ", codegen.SnakeCase))
			sb.WriteString(")")
			if len(m.CKs()) > 0 {
				sb.WriteString(", ")
				sb.WriteString(m.Names(codegen.PropFilterCKs, "", "", ", ", codegen.SnakeCase))
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
		"MVNameSC": func(m codegen.ModelKey) string {
			sb := strings.Builder{}
			sb.WriteString(tools.ToSnake(m.Name()))
			sb.WriteString("_by_")
			sb.WriteString(m.Names(codegen.PropFilterPKs, "", "", "_", codegen.SnakeCase))
			if m.OrderByAlias() != "" {
				sb.WriteString("_")
				sb.WriteString(tools.ToSnake(m.OrderByAlias()))
			} else if m.Index() > 0 {
				sb.WriteString("_")
				sb.WriteString(tools.IntToStr(m.Index()))
			}
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
CREATE MATERIALIZED VIEW view_{{MVNameSC .}} AS
SELECT *
FROM tab_{{Singular $model.NameSC}}
{{MVWhere .}}
PRIMARY KEY {{PrimaryKey .}}
{{- WithClusteringKey . -}};
{{ end }}
`
const genTable = `
var tab{{.Name}} = table.New(table.Metadata{
	Name: "tab_{{Singular .NameSC}}",
	Columns: []string{ {{- Columns .Table -}} },
	PartKey: []string{ {{- PartKeys .Table -}} },
	SortKey: []string{ {{- SortKeys .Table -}} },
})
{{ range .Views }}
var view{{MVName .}} = table.New(table.Metadata{
	Name: "view_{{MVNameSC .}}",
	Columns: []string{ {{- Columns . -}} },
	PartKey: []string{ {{- PartKeys . -}} },
	SortKey: []string{ {{- SortKeys . -}} },
})
{{ end }}
`
const genRemoteRepo = `
type {{.Name}}RemoteRepo struct {
	s gocqlx.Session
}

func New{{.Name}}RemoteRepo(s gocqlx.Session) *{{.Name}}RemoteRepo {
	return &{{.Name}}RemoteRepo{
		s: s,
	}
}
`
const genCreate = ``
const genCreateIF = ``
const genUpdate = ``
const genDelete = ``
const genRead = ``
const genListByPK = ``
const genListByIndex = ``
