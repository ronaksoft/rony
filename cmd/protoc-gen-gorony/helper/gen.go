package helper

import (
	"fmt"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	"google.golang.org/protobuf/compiler/protogen"
	"hash/crc32"
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
	f       *protogen.File
	g       *protogen.GeneratedFile
	plugins map[string]struct{}
}

func New(f *protogen.File, g *protogen.GeneratedFile, plugins map[string]struct{}) *Generator {
	return &Generator{
		f:       f,
		g:       g,
		plugins: plugins,
	}
}

func (g *Generator) Generate() {
	g.g.P("package ", g.f.GoPackageName)
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "sync"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/registry"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/encoding/protojson"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/pools"})
	if _, ok := g.plugins["no_edge_dep"]; !ok {
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edge"})
	}

	g.g.P("var _ = pools.Imported")
	initFunc := &strings.Builder{}
	initFunc.WriteString("func init() {\n")
	for _, m := range g.f.Messages {
		g.genPool(m, initFunc)
		g.genDeepCopy(m)
		g.genMarshal(m)
		g.genUnmarshal(m)
		g.genMarshalJson(m)
		g.genUnmarshalJson(m)

		if _, ok := g.plugins["no_edge_dep"]; !ok {
			g.genPushToContext(m)
		}
	}
	for _, st := range g.f.Services {
		for _, m := range st.Methods {
			methodName := fmt.Sprintf("%s%s", st.Desc.Name(), m.Desc.Name())
			constructor := crc32.ChecksumIEEE([]byte(methodName))
			initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, methodName))
			g.g.P("const C_", methodName, " int64 = ", fmt.Sprintf("%d", constructor))
		}
	}
	initFunc.WriteString("}")
	g.g.P("")
	g.g.P(initFunc.String())
	g.g.P()
}

func (g *Generator) Exec(t *template.Template, v interface{}) string {
	sb := &strings.Builder{}
	err := t.Execute(sb, v)
	if err != nil {
		panic(err)
	}
	return sb.String()
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
	return x
}

func (p *pool{{.Name}}) Put(x *{{.Name}}) {
	if x == nil {
		return
	}
	
	{{ range .Fields}}
	{{- if eq .Kind "repeated"}}
		x.{{- .Name}} = x.{{.Name}}[:0]
	{{- else if eq .Kind "repeated/bytes"}}
		for _, z := range x.{{.Name}} {
			pools.Bytes.Put(z)
		}
		x.{{- .Name}} = x.{{.Name}}[:0]
	{{- else if eq .Kind "repeated/msg"}}
		for _, z := range x.{{.Name}} {
			{{- if ne .Pkg ""}}
				{{.Pkg}}.Pool{{.Type}}.Put(z)
			{{- else}}
				Pool{{.Type}}.Put(z)
			{{- end}}
		}
		x.{{.Name}} = x.{{.Name}}[:0]
	{{- else if eq .Kind "bytes"}}
		x.{{.Name}} = x.{{.Name}}[:0]
	{{- else if eq .Kind "msg"}}
		{{- if ne .Pkg ""}}
			{{.Pkg}}.Pool{{.Type}}.Put(x.{{.Name}})
		{{- else}}
			Pool{{.Type}}.Put(x.{{.Name}})
		{{- end}}
	{{- else}}
		x.{{.Name}} = {{.ZeroValue}}
	{{- end}}
	{{- end}}

	p.pool.Put(x)
}

var Pool{{.Name}} = pool{{.Name}}{}
`

func (g *Generator) genPool(m *protogen.Message, initFunc *strings.Builder) {
	arg := z.TemplateArg{
		Name: string(m.Desc.Name()),
		C:    crc32.ChecksumIEEE([]byte(m.Desc.Name())),
	}
	initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", arg.C, arg.Name))
	for _, ft := range m.Fields {
		arg.AddField(g.f, g.g, ft.Desc)
	}
	g.g.P(g.Exec(template.Must(template.New("genPool").Parse(genPool)), arg))
}

const genDeepCopy = `
func (x *{{.Name}}) DeepCopy(z *{{.Name}}) {
	{{- range .Fields -}}
	{{- if eq .Kind "repeated" }}
		z.{{.Name}} = append(z.{{.Name}}[:0], x.{{.Name}}...) 
	{{- else if eq .Kind "repeated/msg" }}
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
	{{- else if eq .Kind "repeated/bytes" }}
		z.{{.Name}} = z.{{.Name}}[:0]
		zl := len(z.{{.Name}})
		for idx := range x.{{.Name}} {
			if idx <  zl {
				z.{{.Name}} = append(z.{{.Name}}, append(z.{{.Name}}[idx][:0], x.{{.Name}}[idx]...))
			} else {
				z.{{.Name}} = append(z.{{.Name}}, append(nil, x.{{.Name}}[idx]...))
			}
		}
	{{- else if eq .Kind "msg" }}
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
			// TODO:: release to pool
			z.{{.Name}} = nil 
		}
	{{- else if eq .Kind "bytes" }}
		z.{{.Name}} = append(z.{{.Name}}[:0], x.{{.Name}}...)
	{{- else }}
		z.{{.Name}} = x.{{.Name}}
	{{- end }}
	{{- end }}
}

`

func (g *Generator) genDeepCopy(m *protogen.Message) {
	arg := z.TemplateArg{
		Name: string(m.Desc.Name()),
		C:    crc32.ChecksumIEEE([]byte(m.Desc.Name())),
	}
	for _, ft := range m.Fields {
		arg.AddField(g.f, g.g, ft.Desc)
	}
	g.g.P(g.Exec(template.Must(template.New("genDeepCopy").Parse(genDeepCopy)), arg))
}

func (g *Generator) genPushToContext(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") PushToContext(ctx *edge.RequestCtx) {")
	g.g.P("ctx.PushMessage(C_", mtName, ", x)")
	g.g.P("}")
	g.g.P()
}
func (g *Generator) genUnmarshal(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") Unmarshal(b []byte) error {")
	g.g.P("return proto.UnmarshalOptions{}.Unmarshal(b, x)")
	g.g.P("}")
	g.g.P()
}
func (g *Generator) genMarshal(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") Marshal() ([]byte, error) {")
	g.g.P("return proto.Marshal(x)")
	g.g.P("}")
	g.g.P()
}
func (g *Generator) genUnmarshalJson(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") UnmarshalJSON(b []byte) error {")
	g.g.P("return protojson.Unmarshal(b, x)")
	g.g.P("}")
	g.g.P()
}
func (g *Generator) genMarshalJson(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") MarshalJSON() ([]byte, error) {")
	g.g.P("return protojson.Marshal(x)")
	g.g.P("}")
	g.g.P()
}
