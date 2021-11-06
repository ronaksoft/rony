package rpc

import (
	"github.com/ronaksoft/rony"
	"google.golang.org/protobuf/proto"

	// "github.com/ronaksoft/rony"
	// "github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
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

func GenFunc(g *protogen.GeneratedFile, _ *codegen.PluginOptions, files ...*protogen.File) error {
	for _, f := range files {
		templateArg := codegen.GenTemplateArg(f)
		if len(templateArg.Services) > 0 {
			g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edge"})
			g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/errors"})
			g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})
			g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "fmt"})
			g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
			g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})

			g.P("var _ = tools.TimeUnix()")
			for _, arg := range templateArg.Services {
				g.P(codegen.ExecTemplate(template.Must(template.New("genServer").Parse(genServer)), arg))
				g.P(codegen.ExecTemplate(template.Must(template.New("genServerWrapper").Parse(genServerWrapper)), arg))
				g.P(codegen.ExecTemplate(template.Must(template.New("genTunnelCommand").Parse(genTunnelCommand)), arg))
				if arg.HasRestProxy {
					g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "net/http"})
					g.P(codegen.ExecTemplate(template.Must(template.New("genServerRestProxy").Parse(genServerRestProxy)), arg))
				}

				if !proto.GetExtension(arg.Options(), rony.E_RonyNoClient).(bool) {
					g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edgec"})
					g.P(codegen.ExecTemplate(template.Must(template.New("genClient").Parse(genClient)), arg))
				}
				cliOpt := proto.GetExtension(arg.Options(), rony.E_RonyCli).(*rony.CliOpt)
				if cliOpt != nil {
					g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/spf13/cobra"})
					g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/config"})
					g.P(codegen.ExecTemplate(template.Must(template.New("genCobraCmd").Parse(genCobraCmd)), arg))
				}
			}
		}
	}

	return nil
}

const genServerRestProxy = `
{{- $service := . }}
{{- range .Methods }}
func (sw *{{$service.NameCC}}Wrapper) {{.NameCC}}RestClient (conn rony.RestConn, ctx *edge.DispatchCtx) error {
	{{- if eq .Input.Pkg "" }}
		req := &{{.Input.Name}}{}
	{{- else }}
		req := &{{.Input.Pkg}}{{.Input.Name}}{}
	{{- end }}
	
	{{- if .Rest.Unmarshal }}
	{{- if .Rest.Json }}
		if len(conn.Body()) > 0 {
			err := req.UnmarshalJSON(conn.Body())
			if err != nil {
				return err 
			}
		}
	{{- else }}
		err := req.Unmarshal(conn.Body())
		if err != nil {
			return err 
		}
	{{- end }}
	{{- end }}

	{{- range .Rest.ExtraCode }}
	{{.}}
	{{- end }}

	ctx.Fill(conn.ConnID(), C_{{.Fullname}}, req)
	return nil
}
func (sw *{{$service.NameCC}}Wrapper) {{.NameCC}}RestServer(conn rony.RestConn, ctx *edge.DispatchCtx) (err error) {
	if !ctx.BufferPop(func(envelope *rony.MessageEnvelope) {
		switch envelope.Constructor {
		case {{.Output.CName}}:
			x := &{{.Output.Name}}{}
			_ = x.Unmarshal(envelope.Message)
			var b []byte
			{{- if .Rest.Json }}
			b, err = x.MarshalJSON()
			{{- else }}
			b, err = x.Marshal()
			{{- end }}
			if err != nil {
				return 
			}
			err = conn.WriteBinary(ctx.StreamID(), b)
			return
		case rony.C_Error:
			x := &rony.Error{}
			_ = x.Unmarshal(envelope.Message)
			err = x
			return
		case rony.C_Redirect:
			x := &rony.Redirect{}
			_ = x.Unmarshal(envelope.Message)
			if len(x.Edges) == 0 || len(x.Edges[0].HostPorts) == 0 {
				break
			}
			switch x.Reason {
			case rony.RedirectReason_ReplicaSetSession:
				conn.Redirect(http.StatusPermanentRedirect, x.Edges[0].HostPorts[0])
			case rony.RedirectReason_ReplicaSetRequest:
				conn.Redirect(http.StatusTemporaryRedirect, x.Edges[0].HostPorts[0])
			}
			return
		}
		err = errors.ErrUnexpectedResponse			
	}) {
		err = errors.ErrInternalServer
	}

	return
}
{{ end }}
`

const genServer = `
type I{{.Name}} interface {
{{- range .Methods }}
	{{.Name}} (ctx *edge.RequestCtx, req *{{.Input.Fullname}}, res *{{.Output.Fullname}}) *rony.Error
{{- end }}
}

func Register{{.Name}} (h I{{.Name}}, e *edge.Server, preHandlers ...edge.Handler) {
	w := {{.NameCC}}Wrapper {
		h: h,
	}
	w.Register(e, func(c uint64) []edge.Handler {
		return preHandlers
	})
}

func Register{{.Name}}WithFunc(h I{{.Name}}, e *edge.Server, handlerFunc func(c uint64) []edge.Handler) {
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
		req := &{{.Input.Name}}{}
	{{- else }}
		req := &{{.Input.Pkg}}.{{.Input.Name}}{}
	{{- end }}
	{{- if eq .Output.Pkg "" }}
		res := &{{.Output.Name}}{}
	{{- else }}
		res := &{{.Output.Pkg}}.{{.Output.Name}}{}
	{{- end }}

	err := proto.UnmarshalOptions{Merge:true}.Unmarshal(in.Message, req)
	if err != nil {
		ctx.PushError(errors.ErrInvalidRequest)
		return
	}

	rErr := sw.h.{{.Name}} (ctx, req, res)
	if rErr != nil {
		ctx.PushError(rErr)
		return
	}
	if !ctx.Stopped() {
	{{- if eq .Output.Pkg "" }}
		ctx.PushMessage(C_{{.Output.Name}}, res)
	{{- else }}
		ctx.PushMessage({{.Output.Pkg}}.C_{{.Output.Name}}, res)
	{{- end }}
	}
}
{{- end }}

func (sw *{{.NameCC}}Wrapper) Register (e *edge.Server, handlerFunc func(c uint64) []edge.Handler) {
	if handlerFunc == nil {
		handlerFunc = func(c uint64) []edge.Handler {
			return nil
		}
	}
	
	
	{{- range .Methods }}
	e.SetHandler(
		edge.NewHandlerOptions().SetConstructor(C_{{.Fullname}}).
		SetHandler(handlerFunc(C_{{.Fullname}})...).
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
func (c *{{$serviceName}}Client) {{.Name}} (
	req *{{.Input.Fullname}}, kvs ...*rony.KeyValue,
) (*{{.Output.Fullname}}, error) {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(c.c.GetRequestID(), C_{{.Fullname}}, req, kvs ...)
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
func TunnelRequest{{$serviceName}}{{.Name}} (
	ctx *edge.RequestCtx, replicaSet uint64, 
	req *{{.Input.Fullname}}, res *{{.Output.Fullname}}, 
	kvs ...*rony.KeyValue,
) error {
	out := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(out)
	in := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(in)
	out.Fill(ctx.ReqID(), C_{{.Fullname}}, req, kvs...)
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
					config.StringFlag("{{.NameCC}}", "{{.DefaultValue}}", "{{.HelpText}}"),
				{{- else if eq .GoKind "int64" }}
					config.Int64Flag("{{.NameCC}}", tools.StrToInt64("{{.DefaultValue}}"), "{{.HelpText}}"),
				{{- else if eq .GoKind "uint64" }}
					config.Uint64Flag("{{.NameCC}}", tools.StrToUInt64("{{.DefaultValue}}"), "{{.HelpText}}"),
				{{- else if eq .GoKind "int32" }}
					config.Int32Flag("{{.NameCC}}", tools.StrToInt32("{{.DefaultValue}}"), "{{.HelpText}}"),
				{{- else if eq .GoKind "uint32" }}
					config.Uint32Flag("{{.NameCC}}", tools.StrToUInt32("{{.DefaultValue}}"), "{{.HelpText}}"),
				{{- else if eq .GoKind "bool" }}
					config.BoolFlag("{{.NameCC}}", false, "{{.HelpText}}"),
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
