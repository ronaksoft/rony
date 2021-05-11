package rpc

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
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
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "fmt"})
		if g.f.GoPackageName != "rony" {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
		}

		for _, s := range g.f.Services {
			g.genServer(s)
			g.genExecuteRemote(s)

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
	g.g.P("type I", s.Desc.Name(), " interface {")
	for _, m := range s.Methods {
		inputName := z.Name(g.f, g.g, m.Desc.Input())
		outputName := z.Name(g.f, g.g, m.Desc.Output())
		g.g.P(m.Desc.Name(), "(ctx *edge.RequestCtx, req *", inputName, ", res *", outputName, ")")
	}
	g.g.P("}")
	g.g.P()
	g.g.P()
	g.g.P("type ", tools.ToLowerCamel(serviceName), "Wrapper struct {")
	g.g.P("h I", s.Desc.Name())
	g.g.P("}")
	g.g.P()
	for _, m := range s.Methods {
		methodName := string(m.Desc.Name())
		inputPkg, inputType := z.DescName(g.f, g.g, m.Desc.Input())
		outputPkg, outputType := z.DescName(g.f, g.g, m.Desc.Output())

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
		g.g.P("ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)")
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
		if proto.GetExtension(opt, rony.E_RonyInconsistentRead).(bool) {
			sb.WriteString(".InconsistentRead()")
		}
		if proto.GetExtension(opt, rony.E_RonyInternal).(bool) {
			sb.WriteString(".TunnelOnly()")
		}
		constructorName := fmt.Sprintf("C_%s%s", serviceName, methodName)
		g.g.P(
			"e.SetHandler(",
			"edge.NewHandlerOptions().SetConstructor(", constructorName, ").",
			"SetHandler(handlerFunc(", constructorName, ")...).Append(sw.", tools.ToLowerCamel(methodName), "Wrapper)",
			sb.String(), ")")
	}
	g.g.P("}")
	g.g.P()
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
		outputPkg, outputType := z.DescName(g.f, g.g, m.Desc.Output())

		leaderOnlyText := "true"
		leaderOnly := proto.GetExtension(opt, rony.E_RonyInconsistentRead).(bool)
		if leaderOnly {
			leaderOnlyText = "false"
		}
		g.g.P("func (c *", s.Desc.Name(), "Client) ", methodLocalName, "(req *", inputName, ", kvs ...*rony.KeyValue) (*", outputName, ", error) {")
		g.g.P("out := rony.PoolMessageEnvelope.Get()")
		g.g.P("defer rony.PoolMessageEnvelope.Put(out)")
		g.g.P("in := rony.PoolMessageEnvelope.Get()")
		g.g.P("defer rony.PoolMessageEnvelope.Put(in)")
		g.g.P("out.Fill(c.c.GetRequestID(), C_", methodName, ", req, kvs...)")
		g.g.P("err := c.c.Send(out, in, ", leaderOnlyText, ")")
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
func (g *Generator) genExecuteRemote(s *protogen.Service) {
	for _, m := range s.Methods {
		leaderOnlyText := "true"
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		leaderOnly := proto.GetExtension(opt, rony.E_RonyInconsistentRead).(bool)
		if leaderOnly {
			leaderOnlyText = "false"
		}
		methodName := fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())
		inputName := z.Name(g.f, g.g, m.Desc.Input())
		outputName := z.Name(g.f, g.g, m.Desc.Output())
		outputC := z.Constructor(g.f, g.g, m.Desc.Output())

		g.g.P("func ExecuteRemote", methodName, "(ctx *edge.RequestCtx, replicaSet uint64, req *", inputName, ", res *", outputName, ", kvs ...*rony.KeyValue) error {")
		g.g.P("out := rony.PoolMessageEnvelope.Get()")
		g.g.P("defer rony.PoolMessageEnvelope.Put(out)")
		g.g.P("in := rony.PoolMessageEnvelope.Get()")
		g.g.P("defer rony.PoolMessageEnvelope.Put(in)")
		g.g.P("out.Fill(ctx.ReqID(), C_", methodName, ", req, kvs...)")
		g.g.P("err := ctx.ExecuteRemote(replicaSet, ", leaderOnlyText, ", out, in)")
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
		g.g.P("return edge.ErrUnexpectedTunnelResponse")
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
