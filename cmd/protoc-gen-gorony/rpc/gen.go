package rpc

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
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
	f *protogen.File
	g *protogen.GeneratedFile
}

func New(f *protogen.File, g *protogen.GeneratedFile) *Generator {
	return &Generator{
		f: f,
		g: g,
	}
}

func (g *Generator) Generate() {
	if len(g.f.Services) > 0 {
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edge"})
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/errors"})
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "fmt"})
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})

		g.g.P("var _ = tools.TimeUnix()")
		for _, s := range g.f.Services {
			arg := z.GetServiceArg(g.f, g.g, s)
			g.g.P(g.Exec(template.Must(template.New("genServer").Parse(genServer)), arg))
			g.g.P(g.Exec(template.Must(template.New("genServerWrapper").Parse(genServerWrapper)), arg))
			g.g.P(g.Exec(template.Must(template.New("genTunnelCommand").Parse(genTunnelCommand)), arg))

			g.genServerRestProxy(s)

			opt, _ := s.Desc.Options().(*descriptorpb.ServiceOptions)
			if !proto.GetExtension(opt, rony.E_RonyNoClient).(bool) {
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edgec"})
				g.g.P(g.Exec(template.Must(template.New("genClient").Parse(genClient)), arg))
			}
			if proto.GetExtension(opt, rony.E_RonyCobraCmd).(bool) {
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/spf13/cobra"})
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/config"})
				g.g.P(g.Exec(template.Must(template.New("genCobraCmd").Parse(genCobraCmd)), arg))
			}
		}
	}
}

func (g *Generator) Exec(t *template.Template, v interface{}) string {
	sb := &strings.Builder{}
	err := t.Execute(sb, v)
	if err != nil {
		panic(err)
	}
	return sb.String()
}

func (g *Generator) genServerRestProxy(s *protogen.Service) {
	for _, m := range s.Methods {
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		restOpt, _ := proto.GetExtension(opt, rony.E_RonyRest).(*rony.RestOpt)
		if restOpt == nil {
			continue
		}
		g.g.P("// ", restOpt.String())

		g.createRestClient(s, m, restOpt)
		g.createRestServer(s, m, restOpt)
	}
}
func (g *Generator) createRestClient(s *protogen.Service, m *protogen.Method, opt *rony.RestOpt) {
	serviceName := string(s.Desc.Name())
	methodName := string(m.Desc.Name())
	methodConstructor := fmt.Sprintf("C_%s%s", serviceName, methodName)
	inputPkg, inputType := z.DescParts(g.f, g.g, m.Desc.Input())

	g.g.P(
		"func (sw *", tools.ToLowerCamel(serviceName), "Wrapper) ",
		tools.ToLowerCamel(methodName), "RestClient (conn rony.RestConn, ctx *edge.DispatchCtx) error {",
	)
	if inputPkg == "" {
		g.g.P("req := Pool", inputType, ".Get()")
		g.g.P("defer Pool", inputType, ".Put(req)")
	} else {
		g.g.P("req := ", inputPkg, ".Pool", inputType, ".Get()")
		g.g.P("defer ", inputPkg, ".Pool", inputType, ".Put(req)")
	}

	var pathVars []string
	path := fmt.Sprintf("/%s", strings.Trim(opt.GetPath(), "/"))
	for _, pv := range strings.Split(path, "/") {
		if !strings.HasPrefix(pv, ":") {
			continue
		}
		pathVars = append(pathVars, strings.TrimLeft(pv, ":"))
	}

	bindVars := map[string]string{}
	for _, bv := range strings.Split(opt.GetBindVariables(), ",") {
		parts := strings.SplitN(strings.TrimSpace(bv), "=", 2)
		if len(parts) == 2 {
			bindVars[parts[0]] = parts[1]
		}
	}

	if len(m.Input.Fields) > len(pathVars) {
		if opt.GetJsonEncode() {
			g.g.P("err := req.UnmarshalJSON(conn.Body())")
			g.g.P("if err != nil {")
			g.g.P("return err")
			g.g.P("}")
		} else {
			g.g.P("err := req.Unmarshal(conn.Body())")
			g.g.P("if err != nil {")
			g.g.P("return err")
			g.g.P("}")
		}
	}

	// Try to bind path variables to the input message
	for _, pathVar := range pathVars {
		varName := pathVar
		if bindVars[pathVar] != "" {
			varName = bindVars[pathVar]
		}
		for _, f := range m.Input.Fields {
			if f.Desc.JSONName() == varName {
				switch f.Desc.Kind() {
				case protoreflect.Int64Kind, protoreflect.Sfixed64Kind:
					g.g.P("req.", f.Desc.Name(), "= tools.StrToInt64(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
				case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					g.g.P("req.", f.Desc.Name(), "= tools.StrToUInt64(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
				case protoreflect.Int32Kind, protoreflect.Sfixed32Kind:
					g.g.P("req.", f.Desc.Name(), "= tools.StrToInt32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
					g.g.P("req.", f.Desc.Name(), "= tools.StrToUInt32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
				case protoreflect.StringKind:
					g.g.P("req.", f.Desc.Name(), "= tools.GetString(conn.Get(\"", pathVar, "\"), \"\")")
				case protoreflect.BytesKind:
					g.g.P("req.", f.Desc.Name(), "= tools.S2B(tools.GetString(conn.Get(\"", pathVar, "\"), \"\"))")
				case protoreflect.DoubleKind:
					g.g.P("req.", f.Desc.Name(), "= tools.StrToFloat32(tools.GetString(conn.Get(\"", pathVar, "\"), \"0\"))")
				}
			}
		}
	}

	g.g.P("ctx.FillEnvelope(conn.ConnID(), ", methodConstructor, ", req)")
	g.g.P("return nil")
	g.g.P("}") // end of client side func block
	g.g.P()

}
func (g *Generator) createRestServer(s *protogen.Service, m *protogen.Method, opt *rony.RestOpt) {
	serviceName := string(s.Desc.Name())
	methodName := string(m.Desc.Name())
	outputName := z.DescName(g.f, g.g, m.Desc.Output())
	outputConstructor := fmt.Sprintf("C_%s", outputName)
	g.g.P(
		"func (sw *", tools.ToLowerCamel(serviceName), "Wrapper) ",
		tools.ToLowerCamel(methodName), "RestServer (conn rony.RestConn, ctx *edge.DispatchCtx) error {",
	)

	g.g.P("envelope := ctx.BufferPop()")
	g.g.P("if envelope == nil {")
	g.g.P("return errors.ErrInternalServer")
	g.g.P("}")

	g.g.P("switch envelope.Constructor {")
	g.g.P("case ", outputConstructor, ":")
	g.g.P("x := &", outputName, "{}")
	g.g.P("_ = x.Unmarshal(envelope.Message)")
	if opt.GetJsonEncode() {
		g.g.P("b, err := x.MarshalJSON()")
	} else {
		g.g.P("b, err := x.Marshal()")
	}
	g.g.P("if err != nil {")
	g.g.P("return err")
	g.g.P("}")
	g.g.P("return conn.WriteBinary(ctx.StreamID(), b)")
	g.g.P()
	g.g.P("case rony.C_Error:")
	g.g.P("x := &rony.Error{}")
	g.g.P("_ = x.Unmarshal(envelope.Message)")
	g.g.P()
	g.g.P("default:")
	g.g.P("return errors.ErrUnexpectedResponse")
	g.g.P("}")
	g.g.P()
	g.g.P("return errors.ErrInternalServer")
	g.g.P("}") // end of server side func block
	g.g.P()
}

const genServer = `
type I{{.Name}} interface {
{{- range .Methods }}
	{{.Name}} (ctx *edge.RequestCtx, req *{{.Input.Name}}, res *{{.Output.Name}})
{{- end }}
}

func Register{{.Name}} (h I{{.Name}}, e *edge.Server, preHandlers ...edge.Handler) {
	w := {{.NameCC}}Wrapper {
		h: h,
	}
	w.Register(e, func(c int64) []edge.Handler {
		return preHandlers
	})
}

func Register{{.Name}}WithFunc(h I{{.Name}}, e *edge.Server, handlerFunc func(c int64) []edge.Handler) {
	w := {{.NameCC}}Wrapper {
		h: h,
	}
	w.Register(e, handlerFunc)
}
`

const genServerWrapper = `
type {{.NameCC}}Wrapper struct {
	h I{{.Name}}
}

{{- $serviceNameCC := .NameCC -}}
{{- $serviceName := .Name -}}
{{- range .Methods }}
func (sw *{{$serviceNameCC}}Wrapper) {{.NameCC}}Wrapper(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	{{- if eq .Input.Pkg "" }}
		req := Pool{{.Input.Name}}.Get()
		defer Pool{{.Input.Name}}.Put(req)
	{{- else }}
		req := {{.Input.Pkg}}Pool{{.Input.Name}}.Get()
		defer {{.Input.Pkg}}Pool{{.Input.Name}}.Put(req)
	{{- end }}
	{{- if eq .Output.Pkg "" }}
		res := Pool{{.Output.Name}}.Get()
		defer Pool{{.Output.Name}}.Put(res)
	{{- else }}
		res := {{.Output.Pkg}}Pool{{.Output.Name}}.Get()
		defer {{.Output.Pkg}}Pool{{.Output.Name}}.Put(res)
	{{- end }}

	err := proto.UnmarshalOptions{Merge:true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	sw.h.{{.Name}} (ctx, req, res)
	if !ctx.Stopped() {
	{{- if eq .Output.Pkg "" }}
		ctx.PushMessage(C_{{.Output.Name}}, res)
	{{- else }}
		ctx.PushMessage({{.Output.Pkg}}.C_{{.Output.Name}}, res)
	{{- end }}
	}
}
{{- end }}

func (sw *{{.NameCC}}Wrapper) Register (e *edge.Server, handlerFunc func(c int64) []edge.Handler) {
	if handlerFunc == nil {
		handlerFunc = func(c int64) []edge.Handler {
			return nil
		}
	}
	
	
	{{- range .Methods }}
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_{{$serviceName}}{{.Name}}).
		SetHandler(handlerFunc(C_{{$serviceName}}{{.Name}})...).
        Append(sw.{{.NameCC}}Wrapper)
		{{- if .TunnelOnly }}.TunnelOnly(){{- end }},
	)
	{{- if .RestEnabled }}
	e.SetRestProxy(
		"{{.Rest.Method}}", "{{.Rest.Path}}",
		edge.NewRestProxy(sw.{{.NameCC}}RestClient, sw.{{.NameCC}}RestServer),
	)
	{{- end }}
	{{- end }}
}
`

const genClient = `
type {{.Name}}Client struct {
	c edgec.Client
}

func New{{.Name}}Client(ec edgec.Client) *{{.Name}}Client {
	return &{{.Name}}Client{
		c: ec,
	}
}

{{- $serviceName := .Name -}}
{{- range .Methods }}
{{- if not .TunnelOnly }}
func (c *{{$serviceName}}Client) {{.Name}} (req *{{.Input.Name}}, kvs ...*rony.KeyValue) (*{{.Output.Name}}, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_{{$serviceName}}{{.Name}}, req, kvs ...)
	err := c.c.Send(out, in)
	if err != nil {
		return nil, err
	}
	switch in.GetConstructor() {
	{{- if eq .Output.Pkg "" }}
	case C_{{.Output.Name}}:
	{{- else }}
	case {{.Output.Pkg}}.C_{{.Output.Name}}:
	{{- end }}
		x := &{{.Output.Name}}{}
		_ = proto.Unmarshal(in.Message, x)
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.Message)
		return nil, x
	default:
		return nil, fmt.Errorf("unkown message :%d", in.GetConstructor())
	}
}
{{- end }}
{{- end }}
`

const genTunnelCommand = `
{{- $serviceName := .Name -}}
{{- range .Methods }}
func TunnelRequest{{$serviceName}}{{.Name}} (ctx *edge.RequestCtx, replicaSet uint64, req *{{.Input.Name}}, res *{{.Output.Name}}, kvs ...*rony.KeyValue) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_{{$serviceName}}{{.Name}}, req, kvs...)
	err := ctx.TunnelRequest(replicaSet, out, in)
	if err != nil {
		return err
	}

	switch in.GetConstructor() {
	case {{.Output.CName}}:
		_ = res.Unmarshal(in.GetMessage())
		return nil 
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(in.GetMessage())
		return x
	default:
		return errors.ErrUnexpectedTunnelResponse
	}
}
{{- end }}
`

const genCobraCmd = `
func prepare{{.Name}}Command(cmd *cobra.Command, c edgec.Client) (*{{.Name}}Client, error) {
	// Bind current flags to the registered flags in config package
	err := config.BindCmdFlags(cmd)
	if err != nil {
		return nil, err
	}
	
	return New{{.Name}}Client(c), nil
}

{{- $serviceName := .Name -}}
{{- range .Methods }}
	{{- if not .TunnelOnly }}
		var gen{{$serviceName}}{{.Name}}Cmd = func(h I{{$serviceName}}Cli, c edgec.Client) *cobra.Command {
			cmd := &cobra.Command {
				Use: "{{.NameKC}}",
				RunE: func(cmd *cobra.Command, args []string) error {
					cli, err := prepare{{$serviceName}}Command(cmd, c)
					if err != nil {
						return err 
					}
					return h.{{.Name}}(cli, cmd, args)
				},
			}
			config.SetFlags(cmd,
			{{- range .Input.Fields }}
				{{- if or (eq .GoKind "string") (eq .GoKind "[]byte") }}
					config.StringFlag("{{.NameCC}}", "", ""),
				{{- else if eq .GoKind "int64" }}
					config.Int64Flag("{{.NameCC}}", 0, ""),
				{{- else if eq .GoKind "uint64" }}
					config.Uint64Flag("{{.NameCC}}", 0, ""),
				{{- else if eq .GoKind "int32" }}
					config.Int32Flag("{{.NameCC}}", 0, ""),
				{{- else if eq .GoKind "uint32" }}
					config.Uint32Flag("{{.NameCC}}", 0, ""),
				{{- else if eq .GoKind "bool" }}
					config.BoolFlag("{{.NameCC}}", false, ""),
				{{- end }}
			{{- end }}
			)
			return cmd
		}
	{{- end }}
{{- end }}

type I{{.Name}}Cli interface {
{{- range .Methods }}
	{{- if not .TunnelOnly }}
		{{.Name}} (cli *{{$serviceName}}Client, cmd *cobra.Command, args []string) error
	{{- end }}
{{- end }}
}

func Register{{$serviceName}}Cli (h I{{$serviceName}}Cli, c edgec.Client, rootCmd *cobra.Command) {
	rootCmd.AddCommand(
	{{- range .Methods }}
	{{- if not .TunnelOnly }}
		gen{{$serviceName}}{{.Name}}Cmd(h ,c),
	{{- end }}
	{{- end }}
	)
}
`
