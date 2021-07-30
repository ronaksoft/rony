package module

import (
	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
	"text/template"
)

/*
   Creation Time: 2021 - Jul - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Generator struct {
	p *protogen.Plugin
	g *protogen.GeneratedFile
}

func New(p *protogen.Plugin) *Generator {
	return &Generator{
		p: p,
	}
}

func (g *Generator) Generate() error {
	arg := codegen.GetModuleArg(g.p)
	g.g = g.p.NewGeneratedFile("module.rony.go", arg.ImportPath)
	g.g.P("package ", arg.PackageName)

	if len(arg.LocalRepos) > 0 {
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
	}
	if len(arg.RemoteRepos) > 0 {
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx"})
	}
	g.g.P(g.Exec(template.Must(template.New("genModuleBase").Funcs(funcs).Parse(genModuleBase)), arg))
	g.g.P(g.Exec(template.Must(template.New("genLocalRepos").Funcs(funcs).Parse(genLocalRepos)), arg))
	g.g.P(g.Exec(template.Must(template.New("genRemoteRepos").Funcs(funcs).Parse(genRemoteRepos)), arg))

	return nil
}

func (g *Generator) Exec(t *template.Template, v interface{}) string {
	sb := &strings.Builder{}
	if err := t.Execute(sb, v); err != nil {
		panic(err)
	}

	return sb.String()
}

var funcs = map[string]interface{}{
	"InitArgs": func(m codegen.ModuleArg) string {
		sb := strings.Builder{}
		cnt := 0
		for _, r := range m.LocalRepos {
			if cnt > 0 {
				sb.WriteString(", ")
			}
			switch r {
			case "store":
				sb.WriteString("store rony.Store")
			}
			cnt++
		}
		for _, r := range m.RemoteRepos {
			if cnt > 0 {
				sb.WriteString(", ")
			}
			switch r {
			case "cql":
				sb.WriteString(", session gocqlx.Session")
			}
			cnt++
		}
		return sb.String()

	},
}

const genLocalRepos = `
{{- if gt (len .LocalRepos) 0 }}
	type LocalRepos struct {
	{{- range .Aggregates -}}
	{{ if ne .LocalRepo "" }}
		{{.Name}} *{{.Name}}LocalRepo 
	{{- end -}}
	{{- end -}}
	{{- range .Singletons -}}
	{{ if ne .LocalRepo "" -}}
		{{.Name}} *{{.Name}}LocalRepo 
	{{- end -}}
	{{- end -}}
	}
	
	func newLocalRepos(s rony.Store) LocalRepos {
		return LocalRepos {
	{{ range .Aggregates -}}
	{{ if ne .LocalRepo "" -}}
		{{.Name}}: New{{.Name}}LocalRepo(s), 
	{{ end -}}
	{{- end -}}
	{{- range .Singletons -}}
	{{ if ne .LocalRepo "" -}}
		{{.Name}}: New{{.Name}}LocalRepo(s),
	{{ end -}}
	{{- end -}}
		}
	}	
{{- end }}

`

const genRemoteRepos = `
{{- if gt (len .RemoteRepos) 0 }}
	type RemoteRepos struct {
	{{ range .Aggregates -}}
	{{ if ne .RemoteRepo "" -}}
		{{.Name}} {{.Name}}RemoteRepo 
	{{- end -}}
	{{- end -}}
	{{- range .Singletons -}}
	{{ if ne .RemoteRepo "" -}}
		{{.Name}} {{.Name}}RemoteRepo 
	{{ end -}}
	{{- end -}}
	}
	
	func newRemoteRepos(s gocqlx.Session) RemoteRepos {
		return RemoteRepos {
	{{- range .Aggregates -}}
	{{- if ne .RemoteRepo "" -}}
		{{.Name}}: New{{.Name}}RemoteRepo(s), 
	{{- end -}}
	{{- end -}}
	{{- range .Singletons -}}
	{{- if ne .RemoteRepo "" -}}
		{{.Name}}: New{{.Name}}RemoteRepo(s),
	{{- end -}}
	{{- end -}}
		}
	}
{{- end }}
`

const genModuleBase = `
type ModuleBase struct {
{{- if gt (len .LocalRepos) 0 }}
	local LocalRepos
{{- end }}
{{- if gt (len .RemoteRepos) 0 }}
	remote RemoteRepos
{{- end }}
}

func New({{InitArgs .}}) ModuleBase {
	m := ModuleBase{
{{- if gt (len .LocalRepos) 0 }}
		local: newLocalRepos(store),
{{- end }}
{{- if gt (len .RemoteRepos) 0 }}
		remote: newRemoteRepos(session)
{{- end }}
	}
	return m
}
{{ if gt (len .LocalRepos) 0 }}
	func (m ModuleBase) Local() LocalRepos {
		return m.local
	}
	
	func (m ModuleBase) L() LocalRepos {
		return m.local
	}
{{- end }}

{{ if gt (len .RemoteRepos) 0 }}
	func (m ModuleBase) Remote() RemoteRepos {
		return m.remote
	}
	
	func (m ModuleBase) R() RemoteRepos {
		return m.remote
	}
{{- end }}



`
