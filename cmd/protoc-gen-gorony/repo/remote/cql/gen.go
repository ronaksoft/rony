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
	funcs := map[string]interface{}{
		"Singular": func(x string) string {
			return inflection.Singular(x)
		},
		"Plural": func(x string) string {
			return inflection.Plural(x)
		},
		"MVNameSC": func(m codegen.ModelKey) string {
			alias := m.Alias()
			if alias == "" {
				alias = m.Names(codegen.PropFilterALL, "", "", "", codegen.None)
			}
			sb := strings.Builder{}
			sb.WriteString(tools.ToSnake(m.Name()))
			sb.WriteString("_")
			sb.WriteString(tools.ToSnake(alias))
			return sb.String()
		},
		"MVName": func(m codegen.ModelKey) string {
			alias := m.Alias()
			if alias == "" {
				alias = m.Names(codegen.PropFilterALL, "", "", "", codegen.None)
			}
			sb := strings.Builder{}
			sb.WriteString(m.Name())
			sb.WriteString(alias)
			return sb.String()
		},
		"MVAlias": func(m codegen.ModelKey, prefix string) string {
			if m.Alias() != "" {
				return m.Alias()
			}
			sb := strings.Builder{}
			sb.WriteString(prefix)
			sb.WriteString(m.Names(codegen.PropFilterALL, "", "", "", codegen.None))

			return sb.String()
		},
		"Columns": func(m codegen.ModelKey) string {
			sb := strings.Builder{}
			sb.WriteString(m.Names(codegen.PropFilterALL, "\"", "\"", ", ", codegen.SnakeCase))
			sb.WriteString(", \"sdata\"")
			return sb.String()
		},
		"ColumnsValue": func(m codegen.ModelKey, prefix, postfix string) string {
			textCase := codegen.LowerCamelCase
			if prefix != "" {
				textCase = codegen.None
			}
			sb := strings.Builder{}
			sb.WriteString(m.Names(codegen.PropFilterALL, prefix, postfix, ", ", textCase))
			return sb.String()
		},
		"ColumnsValuePKs": func(m codegen.ModelKey, prefix, postfix string) string {
			textCase := codegen.LowerCamelCase
			if prefix != "" {
				textCase = codegen.None
			}
			sb := strings.Builder{}
			sb.WriteString(m.Names(codegen.PropFilterPKs, prefix, postfix, ", ", textCase))
			return sb.String()
		},
		"Where": func(m codegen.ModelKey) string {
			sb := strings.Builder{}
			sb.WriteString(m.Names(codegen.PropFilterALL, "qb.Eq(\"", "\")", ", ", codegen.SnakeCase))
			return sb.String()
		},
		"PartKeys": func(m codegen.ModelKey) string {
			return m.Names(codegen.PropFilterPKs, "\"", "\"", ", ", codegen.SnakeCase)
		},
		"SortKeys": func(m codegen.ModelKey) string {
			return m.Names(codegen.PropFilterCKs, "\"", "\"", ", ", codegen.SnakeCase)
		},
		"FuncArgs": func(m codegen.ModelKey, prefix string) string {
			textCase := codegen.LowerCamelCase
			if prefix != "" {
				textCase = codegen.None
			}
			return m.NameTypes(codegen.PropFilterALL, prefix, textCase, codegen.LangGo)
		},
	}
	for _, m := range g.f.Messages {
		arg := codegen.GetMessageArg(g.f, g.g, m)
		if arg.IsAggregate {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/pools"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx/v2/table"})

			g.g.P(g.Exec(template.Must(template.New("genRemoteRepo").Funcs(funcs).Parse(genRemoteRepo)), arg))
			g.g.P(g.Exec(template.Must(template.New("genCRUD").Funcs(funcs).Parse(genCRUD)), arg))
			g.g.P(g.Exec(template.Must(template.New("genListByPK").Funcs(funcs).Parse(genListByPK)), arg))
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
			alias := m.Alias()
			if alias == "" {
				alias = m.Names(codegen.PropFilterALL, "", "", "", codegen.None)
			}
			sb := strings.Builder{}
			sb.WriteString(tools.ToSnake(m.Name()))
			sb.WriteString("_")
			sb.WriteString(tools.ToSnake(alias))
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

const genRemoteRepo = `
{{$repoName := print .Name "RemoteRepo"}}
type {{$repoName}} struct {
	qp map[string]*pools.QueryPool
	t *table.Table
	v map[string]*table.Table
	s gocqlx.Session
}

func New{{$repoName}}(s gocqlx.Session) *{{$repoName}} {
	r := &{{$repoName}}{
		s: s,
		t: table.New(table.Metadata{
			Name: "tab_{{Singular .NameSC}}",
			Columns: []string{ {{- Columns .Table -}} },
			PartKey: []string{ {{- PartKeys .Table -}} },
			SortKey: []string{ {{- SortKeys .Table -}} },
		}),
		v: map[string]*table.Table{
		{{- range .Views }}
		"{{MVAlias . ""}}": table.New(table.Metadata{
			Name: "view_{{MVNameSC .}}",
			Columns: []string{ {{- Columns . -}} },
			PartKey: []string{ {{- PartKeys . -}} },
			SortKey: []string{ {{- SortKeys . -}} },
		}),
		{{- end }}
		},
	}
    
	r.qp = map[string]*pools.QueryPool{
		"insertIF": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Unique().Query(s)
		}),
		"insert": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.InsertBuilder().Query(s)
		}),
		"update": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.UpdateBuilder().Set("sdata").Query(s)
		}),
		"delete": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.DeleteBuilder().Query(s)
		}),
		"get": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.t.GetQuery(s)
		}),
		{{- range .Views }}
		"getBy{{MVAlias . ""}}": pools.NewQueryPool(func() *gocqlx.Queryx {
			return r.v["{{MVAlias . ""}}"].GetQuery(s)
		}),
		{{- end }}

	}
	return r
}

func (r *{{$repoName}}) Table() *table.Table {
	return r.t
}

func (r *{{$repoName}}) T() *table.Table {
	return r.t
}

{{ range .Views }}
func (r *{{$repoName}}) {{MVAlias . "MV"}}() *table.Table {
	return r.v["{{MVAlias . ""}}"]
}
{{ end }}
`

const genCRUD = `
{{$repoName := print .Name "RemoteRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) Insert(m *{{$modelName}}, replace bool) error {
	buf := pools.Buffer.FromProto(m)
	defer pools.Buffer.Put(buf)
	
	var q *gocqlx.Queryx
	if replace {
		q = r.qp["insertIF"].GetQuery()
		defer r.qp["insertIF"].Put(q)
	} else {
		q = r.qp["insert"].GetQuery()
		defer r.qp["insert"].Put(q)
	}
	

	q.Bind({{ColumnsValue .Table "m." ""}}, *buf.Bytes())
	return q.Exec()
}

func (r *{{$repoName}}) Update(m *{{$modelName}}) error {
	buf := pools.Buffer.FromProto(m)
	defer pools.Buffer.Put(buf)
	
	q := r.qp["update"].GetQuery()
	defer r.qp["update"].Put(q)

	
	q.Bind(*buf.Bytes(), {{ColumnsValue .Table "m." ""}})
	return q.Exec()
}

func (r *{{$repoName}}) Delete({{FuncArgs .Table ""}}) error {
	q := r.qp["delete"].GetQuery()
	defer r.qp["delete"].Put(q)

	
	q.Bind({{ColumnsValue .Table "" ""}})
	return q.Exec()
}

func (r *{{$repoName}}) Get({{FuncArgs .Table ""}}, m *{{$modelName}}) (*{{$modelName}}, error) {
	q := r.qp["get"].GetQuery()
	defer r.qp["get"].Put(q)

	if m == nil {
		m = &{{$modelName}}{}
	}

	q.Bind({{ColumnsValue .Table "" ""}})

	var b []byte
	err := q.Scan({{ColumnsValue .Table "&m." ""}}, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}
{{ range .Views }}
func (r *{{$repoName}}) GetBy{{MVAlias . ""}} ({{FuncArgs . ""}}, m *{{$modelName}}) (*{{$modelName}}, error) {
	q := r.qp["getBy{{MVAlias . ""}}"].GetQuery()
	defer r.qp["getBy{{MVAlias . ""}}"].Put(q)

	if m == nil {
		m = &{{$modelName}}{}
	}

	q.Bind({{ColumnsValue . "" ""}})

	var b []byte
	err := q.Scan({{ColumnsValue . "&m." ""}}, &b)
	if err != nil {
		return m, err
	}
	err = m.Unmarshal(b)
	return m, err
}

{{ end }}
`

const genListByPK = `
{{$repoName := print .Name "RemoteRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) List(pk {{$modelName}}PrimaryKey, limit uint) ([]*{{$modelName}}, error) {
	var (
		q *gocqlx.Queryx
		res []*{{$modelName}}
		err error
	)

	switch pk := pk.(type) {
	case {{$modelName}}PK:
		q = r.t.SelectBuilder().Limit(limit).Query(r.s)
		q.Bind({{ColumnsValuePKs .Table "pk." ""}})
{{ range .Views }}
	case {{MVName .}}PK:
		q = r.v["{{MVAlias . ""}}"].SelectBuilder().Limit(limit).Query(r.s)
		q.Bind({{ColumnsValuePKs . "pk." ""}})
{{ end }}
	default:
		panic("BUG!! incorrect mount key")
	}
	err = q.SelectRelease(&res)

	return res, err
}
`
