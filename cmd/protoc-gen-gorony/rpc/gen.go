package rpc

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
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
			arg := codegen.GetServiceArg(g.f, g.g, s)
			g.g.P(g.Exec(template.Must(template.New("genServer").Parse(genServer)), arg))
			g.g.P(g.Exec(template.Must(template.New("genServerWrapper").Parse(genServerWrapper)), arg))
			g.g.P(g.Exec(template.Must(template.New("genTunnelCommand").Parse(genTunnelCommand)), arg))
			g.g.P(g.Exec(template.Must(template.New("genServerRestProxy").Parse(genServerRestProxy)), arg))

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
	if err := t.Execute(sb, v); err != nil {
		panic(err)
	}

	return sb.String()
}

const genServerRestProxy = `
{{- $service := . }}
{{- range .Methods }}
func (sw *{{$service.NameCC}}Wrapper) {{.NameCC}}RestClient (conn rony.RestConn, ctx *edge.DispatchCtx) error {
	{{- if eq .Input.Pkg "" }}
		req := Pool{{.Input.Name}}.Get()
		defer Pool{{.Input.Name}}.Put(req)
	{{- else }}
		req := {{.Input.Pkg}}Pool{{.Input.Name}}.Get()
		defer {{.Input.Pkg}}Pool{{.Input.Name}}.Put(req)
	{{- end }}
	
	{{- if .Rest.Unmarshal }}
	{{- if .Rest.Json }}
		err := req.UnmarshalJSON(conn.Body())
	{{- else }}
		err := req.Unmarshal(conn.Body())
	{{- end }}
		if err != nil {
			return err 
		}
	{{- end }}

	{{- range .Rest.ExtraCode }}
	{{.}}
	{{- end }}

	ctx.FillEnvelope(conn.ConnID(), C_{{$service.Name}}{{.Name}}, req)
	return nil
}
func (sw *{{$service.NameCC}}Wrapper) {{.NameCC}}RestServer(conn rony.RestConn, ctx *edge.DispatchCtx) error {
	envelope := ctx.BufferPop()
	if envelope == nil {
		return errors.ErrInternalServer
	}
	
	switch envelope.Constructor {
	case {{.Output.CName}}:
		x := &{{.Output.Name}}{}
		_ = x.Unmarshal(envelope.Message)
		{{- if .Rest.Json }}
		b, err := x.MarshalJSON()
		{{- else }}
		b, err := x.Marshal()
		{{- end }}
		if err != nil {
			return err
		}
		return conn.WriteBinary(ctx.StreamID(), b)
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(envelope.Message)
		return errors.ErrInternalServer
	}
	return errors.ErrUnexpectedResponse
}
{{ end }}
`

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
		return nil, fmt.Errorf("unknown message :%d", in.GetConstructor())
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
	subCommand := &cobra.Command{
		Use: "{{$serviceName}}",
	}
	subCommand.AddCommand(
	{{- range .Methods }}
	{{- if not .TunnelOnly }}
		gen{{$serviceName}}{{.Name}}Cmd(h ,c),
	{{- end }}
	{{- end }}
	)

	rootCmd.AddCommand(subCommand)
}
`
