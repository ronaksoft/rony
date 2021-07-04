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
			g.genServer(s)
			g.genServerWrapper(s)
			g.genServerRestProxy(s)
			g.genTunnelCommand(s)

			opt, _ := s.Desc.Options().(*descriptorpb.ServiceOptions)
			if !proto.GetExtension(opt, rony.E_RonyNoClient).(bool) {
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edgec"})
				g.genClient(s)
			}
			if proto.GetExtension(opt, rony.E_RonyCobraCmd).(bool) {
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/spf13/cobra"})
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/config"})
				g.genCobraCmd(s)
			}
		}
	}
}

func (g *Generator) genServer(s *protogen.Service) {
	serviceName := string(s.Desc.Name())
	// Define the service interface
	g.g.P("type I", s.Desc.Name(), " interface {")
	for _, m := range s.Methods {
		inputName := z.Name(g.f, g.g, m.Desc.Input())
		outputName := z.Name(g.f, g.g, m.Desc.Output())
		g.g.P(m.Desc.Name(), "(ctx *edge.RequestCtx, req *", inputName, ", res *", outputName, ")")
	}
	g.g.P("}")
	g.g.P()
	g.g.P()
	// - end of service interface definition

	g.g.P("func Register", s.Desc.Name(), "(h I", s.Desc.Name(), ", e *edge.Server, preHandlers ...edge.Handler) {")
	g.g.P("w := ", tools.ToLowerCamel(serviceName), "Wrapper{")
	g.g.P("h: h,")
	g.g.P("}")
	g.g.P("w.Register(e, func(c int64) []edge.Handler {")
	g.g.P("return preHandlers")
	g.g.P("})")
	g.g.P("}")
	g.g.P()
	g.g.P("func Register", s.Desc.Name(), "WithFunc(h I", s.Desc.Name(), ", e *edge.Server, handlerFunc func(c int64) []edge.Handler) {")
	g.g.P("w := ", tools.ToLowerCamel(serviceName), "Wrapper{")
	g.g.P("h: h,")
	g.g.P("}")
	g.g.P("w.Register(e, handlerFunc)")
	g.g.P("}")
	g.g.P()

}
func (g *Generator) genServerWrapper(s *protogen.Service) {
	serviceName := string(s.Desc.Name())
	g.g.P("type ", tools.ToLowerCamel(serviceName), "Wrapper struct {")
	g.g.P("h I", s.Desc.Name())
	g.g.P("}")
	g.g.P()
	for _, m := range s.Methods {
		methodName := string(m.Desc.Name())
		inputPkg, inputType := z.DescParts(g.f, g.g, m.Desc.Input())
		outputPkg, outputType := z.DescParts(g.f, g.g, m.Desc.Output())

		g.g.P("func (sw *", tools.ToLowerCamel(serviceName), "Wrapper) ", tools.ToLowerCamel(methodName), "Wrapper (ctx *edge.RequestCtx, in *rony.MessageEnvelope) {")
		if inputPkg == "" {
			g.g.P("req := Pool", inputType, ".Get()")
			g.g.P("defer Pool", inputType, ".Put(req)")
		} else {
			g.g.P("req := ", inputPkg, ".Pool", inputType, ".Get()")
			g.g.P("defer ", inputPkg, ".Pool", inputType, ".Put(req)")
		}
		if outputPkg == "" {
			g.g.P("res := Pool", outputType, ".Get()")
			g.g.P("defer Pool", outputType, ".Put(res)")
		} else {
			g.g.P("res := ", outputPkg, ".Pool", outputType, ".Get()")
			g.g.P("defer ", outputPkg, ".Pool", outputType, ".Put(res)")
		}

		g.g.P("err := proto.UnmarshalOptions{Merge:true}.Unmarshal(in.Message, req)")
		g.g.P("if err != nil {")
		g.g.P("ctx.PushError(errors.ErrInvalidRequest)")
		g.g.P("return")
		g.g.P("}")
		g.g.P()
		g.g.P("sw.h.", methodName, "(ctx, req, res)")
		g.g.P("if !ctx.Stopped() {")
		if outputPkg == "" {
			g.g.P("ctx.PushMessage(C_", outputType, ", res)")
		} else {
			g.g.P("ctx.PushMessage(", outputPkg, ".C_", outputType, ", res)")
		}
		g.g.P("}") // end of if block
		g.g.P("}") // end of func block
		g.g.P()
	}
	g.g.P("func (sw *", tools.ToLowerCamel(serviceName), "Wrapper) Register (e *edge.Server, handlerFunc func(c int64) []edge.Handler) {")
	g.g.P("if handlerFunc == nil {")
	g.g.P("handlerFunc = func(c int64) []edge.Handler {")
	g.g.P("return nil")
	g.g.P("}")
	g.g.P("}")
	g.g.P()
	for _, m := range s.Methods {
		methodName := string(m.Desc.Name())
		sb := strings.Builder{}
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		if proto.GetExtension(opt, rony.E_RonyInternal).(bool) {
			sb.WriteString(".TunnelOnly()")
		}
		methodConstructor := fmt.Sprintf("C_%s%s", serviceName, methodName)
		g.g.P(
			"e.SetHandler(",
			"edge.NewHandlerOptions().SetConstructor(", methodConstructor, ").",
			"SetHandler(handlerFunc(", methodConstructor, ")...).Append(sw.", tools.ToLowerCamel(methodName), "Wrapper)",
			sb.String(), ")",
		)
		if restOpt := proto.GetExtension(opt, rony.E_RonyRest).(*rony.RestOpt); restOpt != nil {
			// TODO:: sanitize input: check method is valid, path is valid , ...
			path := fmt.Sprintf("/%s", strings.Trim(restOpt.GetPath(), "/"))
			g.g.P(
				"e.SetRestProxy(\"", restOpt.GetMethod(), "\",\"", path, "\"",
				",edge.NewRestProxy(sw.", tools.ToLowerCamel(methodName), "RestClient, sw.", tools.ToLowerCamel(methodName), "RestServer))",
			)
		}
	}
	g.g.P("}")
	g.g.P()
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
func (g *Generator) genClient(s *protogen.Service) {
	g.g.P("type ", s.Desc.Name(), "Client struct {")
	g.g.P("c edgec.Client")
	g.g.P("}")
	g.g.P()
	g.g.P("func New", s.Desc.Name(), "Client (ec edgec.Client) *", s.Desc.Name(), "Client {")
	g.g.P("return &", s.Desc.Name(), "Client{")
	g.g.P("c: ec,")
	g.g.P("}")
	g.g.P("}")
	g.g.P()
	for _, m := range s.Methods {
		methodName := fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())
		methodLocalName := string(m.Desc.Name())
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		if proto.GetExtension(opt, rony.E_RonyInternal).(bool) {
			continue
		}
		inputName := z.Name(g.f, g.g, m.Desc.Input())
		outputName := z.Name(g.f, g.g, m.Desc.Output())
		outputPkg, outputType := z.DescParts(g.f, g.g, m.Desc.Output())

		g.g.P("func (c *", s.Desc.Name(), "Client) ", methodLocalName, "(req *", inputName, ", kvs ...*rony.KeyValue) (*", outputName, ", error) {")
		g.g.P("out := rony.PoolMessageEnvelope.Get()")
		g.g.P("defer rony.PoolMessageEnvelope.Put(out)")
		g.g.P("in := rony.PoolMessageEnvelope.Get()")
		g.g.P("defer rony.PoolMessageEnvelope.Put(in)")
		g.g.P("out.Fill(c.c.GetRequestID(), C_", methodName, ", req, kvs...)")
		g.g.P("err := c.c.Send(out, in)")
		g.g.P("if err != nil {")
		g.g.P("return nil, err")
		g.g.P("}")
		g.g.P("switch in.GetConstructor() {")
		if outputPkg != "" {
			g.g.P("case ", outputPkg, ".C_", outputType, ":")
			g.g.P("x := &", outputPkg, ".", outputType, "{}")
		} else {
			g.g.P("case C_", outputType, ":")
			g.g.P("x := &", outputType, "{}")
		}

		g.g.P("_ = proto.Unmarshal(in.Message, x)")
		g.g.P("return x, nil")
		g.g.P("case rony.C_Error:")
		g.g.P("x := &rony.Error{}")
		g.g.P("_ = proto.Unmarshal(in.Message, x)")
		g.g.P("return nil, fmt.Errorf(\"%s:%s\", x.GetCode(), x.GetItems())")
		g.g.P("default:")
		g.g.P("return nil, fmt.Errorf(\"unknown message: %d\", in.GetConstructor())")
		g.g.P("}")
		g.g.P("}")
		g.g.P()
	}
}
func (g *Generator) genTunnelCommand(s *protogen.Service) {
	for _, m := range s.Methods {
		methodName := fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())
		inputName := z.Name(g.f, g.g, m.Desc.Input())
		outputName := z.Name(g.f, g.g, m.Desc.Output())
		outputC := z.Constructor(g.f, g.g, m.Desc.Output())

		g.g.P("func TunnelRequest", methodName, "(ctx *edge.RequestCtx, replicaSet uint64, req *", inputName, ", res *", outputName, ", kvs ...*rony.KeyValue) error {")
		g.g.P("out := rony.PoolMessageEnvelope.Get()")
		g.g.P("defer rony.PoolMessageEnvelope.Put(out)")
		g.g.P("in := rony.PoolMessageEnvelope.Get()")
		g.g.P("defer rony.PoolMessageEnvelope.Put(in)")
		g.g.P("out.Fill(ctx.ReqID(), C_", methodName, ", req, kvs...)")
		g.g.P("err := ctx.TunnelRequest(replicaSet, out, in)")
		g.g.P("if err != nil {")
		g.g.P("return err")
		g.g.P("}")
		g.g.P("")
		g.g.P("switch in.GetConstructor() {")
		g.g.P("case ", outputC, ":")
		g.g.P("_ = res.Unmarshal(in.GetMessage())")
		g.g.P("return nil")
		g.g.P("case rony.C_Error:")
		g.g.P("x := &rony.Error{}")
		g.g.P("_ = x.Unmarshal(in.GetMessage())")
		g.g.P("return x")
		g.g.P("default:")
		g.g.P("return errors.ErrUnexpectedTunnelResponse")
		g.g.P("}")
		g.g.P("}")
		g.g.P()
	}
}
func (g *Generator) genCobraCmd(s *protogen.Service) {
	g.createPrepareFunc(s)
	g.createMethodGenerator(s)
	g.createClientCli(s)
}
func (g *Generator) createPrepareFunc(s *protogen.Service) {
	serviceName := string(s.Desc.Name())
	g.g.P("func prepare", serviceName, "Command(cmd *cobra.Command, c edgec.Client) (*", serviceName, "Client, error) {")
	g.g.P("// Bind the current flags to registered flags in config package")
	g.g.P("err := config.BindCmdFlags(cmd)")
	g.g.P("if err != nil {")
	g.g.P("return nil, err")
	g.g.P("}")
	g.g.P()
	g.g.P("return New", serviceName, "Client(c), nil")
	g.g.P("}")
}
func (g *Generator) createMethodGenerator(s *protogen.Service) {
	serviceName := string(s.Desc.Name())
	for _, m := range s.Methods {
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		if proto.GetExtension(opt, rony.E_RonyInternal).(bool) {
			continue
		}
		methodName := fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())
		methodLocalName := string(m.Desc.Name())
		g.g.P("var gen", methodName, "Cmd = func(h I", serviceName, "Cli, c edgec.Client) *cobra.Command {")
		g.g.P("cmd := &cobra.Command {")
		g.g.P("Use: \"", tools.ToKebab(methodLocalName), "\",")
		g.g.P("RunE: func(cmd *cobra.Command, args []string) error {")
		g.g.P("cli, err := prepare", serviceName, "Command(cmd, c)")
		g.g.P("if err != nil {")
		g.g.P("return err")
		g.g.P("}") // end if if clause
		g.g.P("return h.", methodLocalName, "(cli, cmd, args)")
		g.g.P("},") // end of RunE func block
		g.g.P("}")  // end of cobra.Command
		g.g.P("config.SetFlags(cmd,")
		for _, f := range m.Input.Fields {
			fieldName := string(f.Desc.Name())
			switch z.GoKind(g.f, g.g, f.Desc) {
			case "string", "[]byte":
				g.g.P("config.StringFlag(\"", tools.ToLowerCamel(fieldName), "\",\"\", \"\"),")
			case "int64":
				g.g.P("config.Int64Flag(\"", tools.ToLowerCamel(fieldName), "\",0, \"\"),")
			case "uint64":
				g.g.P("config.Uint64Flag(\"", tools.ToLowerCamel(fieldName), "\",0, \"\"),")
			case "int32":
				g.g.P("config.Int32Flag(\"", tools.ToLowerCamel(fieldName), "\",0, \"\"),")
			case "uint32":
				g.g.P("config.Uint32Flag(\"", tools.ToLowerCamel(fieldName), "\",0, \"\"),")
			case "bool":
				g.g.P("config.BoolFlag(\"", tools.ToLowerCamel(fieldName), "\",false, \"\"),")
			default:
			}
		}
		g.g.P(")") // end of SetFlags
		g.g.P("return cmd")
		g.g.P("}") // end of function
		g.g.P()
	}
}
func (g *Generator) createClientCli(s *protogen.Service) {
	g.g.P("type I", s.Desc.Name(), "Cli interface {")
	for _, m := range s.Methods {
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		if proto.GetExtension(opt, rony.E_RonyInternal).(bool) {
			continue
		}
		g.g.P(m.Desc.Name(), "(cli *", s.Desc.Name(), "Client, cmd *cobra.Command, args []string) error")
	}
	g.g.P("}")
	g.g.P()
	g.g.P("func Register", s.Desc.Name(), "Cli (h I", s.Desc.Name(), "Cli, c edgec.Client, rootCmd *cobra.Command) {")
	g.g.P("rootCmd.AddCommand(")
	var names []string
	for _, m := range s.Methods {
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		if proto.GetExtension(opt, rony.E_RonyInternal).(bool) {
			continue
		}
		methodName := fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())
		names = append(names, fmt.Sprintf("gen%sCmd(h, c)", methodName))
		if len(names) == 3 {
			sb := strings.Builder{}
			for _, name := range names {
				sb.WriteString(name)
				sb.WriteRune(',')
			}
			g.g.P(sb.String())
			names = names[:0]
		}
	}
	if len(names) > 0 {
		sb := strings.Builder{}
		for _, name := range names {
			sb.WriteString(name)
			sb.WriteRune(',')
		}
		g.g.P(sb.String())
	}
	g.g.P(")") // end of rootCmd.AddCommand
	g.g.P("}") // end of Register func block
}
