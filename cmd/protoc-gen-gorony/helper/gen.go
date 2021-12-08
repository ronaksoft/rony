package helper

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
)

/*
   Creation Time: 2021 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func GenFunc(g *protogen.GeneratedFile, opt *codegen.PluginOptions, files ...*protogen.File) error {
	initBlock := &strings.Builder{}
	appendToInit := func(x string) {
		initBlock.WriteString(x)
		initBlock.WriteRune('\n')
	}
	g.P("package ", files[0].GoPackageName)
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "sync"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/registry"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/encoding/protojson"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/pools"})
	if !opt.NoEdgeDependency {
		g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edge"})
	}

	g.P("var _ = pools.Imported")
	for _, f := range files {
		templateArg := codegen.GenTemplateArg(f)
		for _, arg := range templateArg.Messages {
			for _, f := range arg.Fields {
				if f.Pkg() != "" {
					g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: f.ImportPath})
				}
			}
			appendToInit(fmt.Sprintf("registry.Register(%d, %q, factory%s)", arg.C, arg.Name(), arg.Name()))
			g.P(codegen.ExecTemplate(template.Must(template.New("genPool").Parse(genPool)), arg))
			g.P(codegen.ExecTemplate(template.Must(template.New("genDeepCopy").Parse(genDeepCopy)), arg))
			g.P(codegen.ExecTemplate(template.Must(template.New("genClone").Parse(genClone)), arg))
			g.P(codegen.ExecTemplate(template.Must(template.New("genSerializers").Parse(genSerializers)), arg))
			if arg.IsEnvelope() {
				g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/goccy/go-json"})
			}
			g.P(codegen.ExecTemplate(template.Must(template.New("genFactory").Parse(genFactory)), arg))

			if !opt.NoEdgeDependency {
				g.P(codegen.ExecTemplate(template.Must(template.New("genPushToContext").Parse(genPushToContext)), arg))
			}
		}
		for _, arg := range templateArg.Services {
			for _, m := range arg.Methods {
				if m.Output.Pkg() != "" {
					g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: m.Output.ImportPath})
				}
				if m.Input.Pkg() != "" {
					g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: m.Input.ImportPath})
				}
				appendToInit(fmt.Sprintf("registry.Register(%d, %q, factory%s)", m.C, m.Fullname(), m.Input.Name()))
				g.P("const C_", m.Fullname(), " uint64 = ", fmt.Sprintf("%d", m.C))
			}
		}

		if initBlock.Len() > 0 {
			g.P("// register constructors of the messages to the registry package")
			g.P("func init() {")
			g.P(initBlock.String())
			g.P("}")
		}
	}

	return nil
}

const genPool = `
const C_{{.Name}} uint64 = {{.C}}
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

	{{ if not .SkipJson }}
		func (x *{{.Name}}) UnmarshalJSON(b []byte) error {
		{{- if and .IsEnvelope }}
			je := registry.JSONEnvelope{}
			err := go_json.Unmarshal(b, &je)
			if err != nil {
				return err
			}
		
			x.Constructor = registry.N(je.Constructor)
		
			m, err := registry.Get(x.Constructor)
			if err != nil {
				return err
			}
		
			err = m.UnmarshalJSON(je.Message)
			if err != nil {
				return err
			}
		
			x.Message, err = proto.Marshal(m)
			if err != nil {
				return err
			}
		
			return nil
		{{- else }}
			return protojson.Unmarshal(b, x)
		{{- end }}
		}
	
		func (x *{{.Name}}) MarshalJSON() ([]byte, error)  {
		{{- if .IsEnvelope }}
			m, err := registry.UnwrapJSON(x)
			if err != nil {
				return nil, err
			}
		
			je := registry.JSONEnvelope{
				Constructor: registry.C(x.Constructor),
			}
		
			je.Message, err = m.MarshalJSON()
			if err != nil {
				return nil, err
			}
		
			return go_json.Marshal(je)
		{{- else }}
			return protojson.Marshal(x)
		{{- end }}
		}
	{{- end }}
`

const genFactory = `
func factory{{.Name}} () registry.Message {
	return &{{.Name}}{}
}
`
