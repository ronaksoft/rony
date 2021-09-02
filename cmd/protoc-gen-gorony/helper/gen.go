package helper

import (
	"fmt"
	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
	"text/template"
)

/*
   Creation Time: 2021 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Generator struct {
	f   *protogen.File
	g   *protogen.GeneratedFile
	opt *codegen.PluginOptions
	initFuncBlock *strings.Builder
}

func New(f *protogen.File, g *protogen.GeneratedFile, options *codegen.PluginOptions) *Generator {
	return &Generator{
		f:   f,
		g:   g,
		opt: options,
		initFuncBlock: &strings.Builder{},
	}
}

func (g *Generator) Generate() {
	g.g.P("package ", g.f.GoPackageName)
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "sync"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/registry"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/encoding/protojson"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/pools"})
	if !g.opt.NoEdgeDependency {
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edge"})
	}

	g.g.P("var _ = pools.Imported")


	for _, m := range g.f.Messages {
		arg := codegen.GetMessageArg(m).With(g.f)
		for _, f := range arg.Fields {
			if f.Pkg() != "" {
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: f.ImportPath})
			}
		}
		g.appendToInit(fmt.Sprintf("registry.RegisterConstructor(%d, %q)", arg.C, arg.Name()))
		g.g.P(g.Exec(template.Must(template.New("genPool").Parse(genPool)), arg))
		g.g.P(g.Exec(template.Must(template.New("genDeepCopy").Parse(genDeepCopy)), arg))
		g.g.P(g.Exec(template.Must(template.New("genClone").Parse(genClone)), arg))
		g.g.P(g.Exec(template.Must(template.New("genSerializers").Parse(genSerializers)), arg))

		if !g.opt.NoEdgeDependency {
			g.g.P(g.Exec(template.Must(template.New("genPushToContext").Parse(genPushToContext)), arg))
		}
	}
	for _, st := range g.f.Services {
		arg := codegen.GetServiceArg(st).With(g.f)
		for _, m := range arg.Methods {
			if m.Output.Pkg() != "" {
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: m.Output.ImportPath})
			}
			if m.Input.Pkg() != "" {
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: m.Input.ImportPath})
			}
			g.appendToInit(fmt.Sprintf("registry.RegisterConstructor(%d, %q)", m.C, m.Fullname()))
			g.g.P("const C_", m.Fullname(), " int64 = ", fmt.Sprintf("%d", m.C))
		}
	}

	if g.initFuncBlock.Len() > 0 {
		g.g.P("// register constructors of the messages to the registry package")
		g.g.P("func init() {")
		g.g.P(g.initFuncBlock.String())
		g.g.P("}")
	}
}

func (g *Generator) Exec(t *template.Template, v interface{}) string {
	sb := &strings.Builder{}
	if err := t.Execute(sb, v); err != nil {
		panic(err)
	}

	return sb.String()
}

func (g *Generator) appendToInit(x string) {
	g.initFuncBlock.WriteString(x)
	g.initFuncBlock.WriteRune('\n')
}

const genPool = `
const C_{{.Name}} int64 = {{.C}}
type pool{{.Name}} struct {
	pool sync.Pool
}

func (p *pool{{.Name}}) Get() *{{.Name}} {
	x, ok := p.pool.Get().(*{{.Name}})
	if !ok {
		x = &{{.Name}}{}
	}
	{{ range .Fields }}
	{{- if and (eq .Cardinality "optional") (eq .Kind "message") }}
		{{- if ne .Pkg ""}}
			x.{{.Name}} = {{.Pkg}}.Pool{{.Type}}.Get()
		{{- else}}
			x.{{.Name}} = Pool{{.Type}}.Get()
		{{- end}}
	{{- end }}
	{{ end }}
	return x
}

func (p *pool{{.Name}}) Put(x *{{.Name}}) {
	if x == nil {
		return
	}
	
	{{ range .Fields }}
		{{- if eq .Cardinality "repeated" }}
			{{- if eq .Kind "message" }}
				for _, z := range x.{{.Name}} {
					{{- if ne .Pkg ""}}
						{{.Pkg}}.Pool{{.Type}}.Put(z)
					{{- else}}
						Pool{{.Type}}.Put(z)
					{{- end}}
				}
				x.{{.Name}} = x.{{.Name}}[:0]
			{{- else if eq .Kind "bytes" }}
				for _, z := range x.{{.Name}} {
					pools.Bytes.Put(z)
				}
				x.{{.Name}} = x.{{.Name}}[:0]
			{{- else }}
				x.{{- .Name}} = x.{{.Name}}[:0]
			{{- end }}
		{{- else }}
			{{- if eq .Kind "bytes" }}
				x.{{.Name}} = x.{{.Name}}[:0]
			{{- else if eq .Kind "message" }}
				{{- if ne .Pkg ""}}
					{{.Pkg}}.Pool{{.Type}}.Put(x.{{.Name}})
				{{- else}}
					Pool{{.Type}}.Put(x.{{.Name}})
				{{- end}}
			{{- else }}
				x.{{.Name}} = {{.ZeroValue}}
			{{- end }}
		{{- end }}
	{{- end }}
	
	p.pool.Put(x)
}

var Pool{{.Name}} = pool{{.Name}}{}
`

const genDeepCopy = `
func (x *{{.Name}}) DeepCopy(z *{{.Name}}) {
	{{- range .Fields -}}
		{{- if eq .Cardinality "repeated" }}
			{{- if eq .Kind "message" }}
				for idx := range x.{{.Name}} {
					if x.{{.Name}}[idx] == nil {
						continue
					}
					{{- if eq .Pkg "" }}
						xx := Pool{{.Type}}.Get()
					{{- else }}
						xx := {{.Pkg}}.Pool{{.Type}}.Get()
					{{- end }}
					x.{{.Name}}[idx].DeepCopy(xx)
					z.{{.Name}} = append(z.{{.Name}}, xx)
				}
			{{- else if eq .Kind "bytes" }}
				z.{{.Name}} = z.{{.Name}}[:0]
				zl := len(z.{{.Name}})
				for idx := range x.{{.Name}} {
					if idx <  zl {
						z.{{.Name}} = append(z.{{.Name}}, append(z.{{.Name}}[idx][:0], x.{{.Name}}[idx]...))
					} else {
						zb := pools.Bytes.GetCap(len(x.{{.Name}}[idx]))
						z.{{.Name}} = append(z.{{.Name}}, append(zb, x.{{.Name}}[idx]...))
					}
				}
			{{- else }}
				z.{{.Name}} = append(z.{{.Name}}[:0], x.{{.Name}}...)
			{{- end }}
		{{- else }}
			{{- if eq .Kind "message" }}
				if x.{{.Name}} != nil {
					if z.{{.Name}} == nil {
						{{- if eq .Pkg "" }}
							z.{{.Name}} = Pool{{.Type}}.Get()
						{{- else }}
							z.{{.Name}} = {{.Pkg}}.Pool{{.Type}}.Get()
						{{- end }}
					}
					x.{{.Name}}.DeepCopy(z.{{.Name}})
				} else {
					{{- if eq .Pkg "" }}
						Pool{{.Type}}.Put(z.{{.Name}})
					{{- else }}
						{{.Pkg}}.Pool{{.Type}}.Put(z.{{.Name}})
					{{- end }}
					z.{{.Name}} = nil 
				}
			{{- else if eq .Kind "bytes" }}
				z.{{.Name}} = append(z.{{.Name}}[:0], x.{{.Name}}...)
			{{- else }}
				z.{{.Name}} = x.{{.Name}}
			{{- end }}
		{{- end }}
	{{- end }}
}

`

const genClone = `
func (x *{{.Name}}) Clone() *{{.Name}} {
	z := &{{.Name}}{}
	x.DeepCopy(z)
	return z
}

`

const genPushToContext = `
	func (x *{{.Name}}) PushToContext(ctx *edge.RequestCtx) {
		ctx.PushMessage({{.CName}}, x)
	}
`

const genSerializers = `
	func (x *{{.Name}}) Unmarshal(b []byte) error {
		return proto.UnmarshalOptions{Merge: true}.Unmarshal(b, x)
	}

	func (x *{{.Name}}) Marshal() ([]byte, error)  {
		return proto.Marshal(x)
	}
	
	func (x *{{.Name}}) UnmarshalJSON(b []byte) error {
		return protojson.Unmarshal(b, x)
	}

	func (x *{{.Name}}) MarshalJSON() ([]byte, error)  {
		return protojson.Marshal(x)
	}
	
`
