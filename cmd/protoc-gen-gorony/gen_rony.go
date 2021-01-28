package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"hash/crc32"
	"strings"
)

/*
   Creation Time: 2020 - Aug - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// GenHelpers generates codes related for pooling of the messages
func GenHelpers(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("package ", file.GoPackageName)
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "",
		GoImportPath: "sync",
	})
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "",
		GoImportPath: "github.com/ronaksoft/rony/registry",
	})
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "",
		GoImportPath: "google.golang.org/protobuf/proto",
	})
	if file.GoPackageName != "rony" {
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "",
			GoImportPath: "github.com/ronaksoft/rony/edge",
		})
	}

	initFunc := &strings.Builder{}
	initFunc.WriteString("func init() {\n")
	for _, mt := range file.Messages {
		genPool(file, mt, initFunc, g)
		genDeepCopy(file, mt, g)
		genMarshal(mt, g)
		genUnmarshal(mt, g)
		if file.GoPackageName != "rony" {
			genPushToContext(mt, g)
		}

	}
	for _, st := range file.Services {
		for _, m := range st.Methods {
			constructor := crc32.ChecksumIEEE([]byte(m.Desc.Name()))
			initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, m.Desc.Name()))
			g.P("const C_", m.Desc.Name(), " int64 = ", fmt.Sprintf("%d", constructor))
		}
	}

	initFunc.WriteString("}")
	g.P("")
	g.P(initFunc.String())
	g.P()
}
func genPool(file *protogen.File, mt *protogen.Message, initFunc *strings.Builder, g *protogen.GeneratedFile) {
	messageName := mt.Desc.Name()
	constructor := crc32.ChecksumIEEE([]byte(messageName))
	g.P(fmt.Sprintf("const C_%s int64 = %d", messageName, constructor))
	initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, messageName))
	g.P(fmt.Sprintf("type pool%s struct{", messageName))
	g.P("pool sync.Pool")
	g.P("}") // end of pool struct
	g.P(fmt.Sprintf("func (p *pool%s) Get() *%s {", messageName, messageName))
	g.P(fmt.Sprintf("x, ok := p.pool.Get().(*%s)", messageName))
	g.P("if !ok {")
	g.P(fmt.Sprintf("return &%s{}", messageName))
	g.P("}") // end of if clause
	g.P("return x")
	g.P("}") // end of func Get()
	g.P()
	g.P(fmt.Sprintf("func (p *pool%s) Put(x *%s) {", messageName, messageName))
	for _, ft := range mt.Fields {
		ftName := ft.Desc.Name()
		switch ft.Desc.Cardinality() {
		case protoreflect.Repeated:
			g.P(fmt.Sprintf("x.%s = x.%s[:0]", ftName, ftName))
		default:
			switch ft.Desc.Kind() {
			case protoreflect.BytesKind:
				g.P(fmt.Sprintf("x.%s = x.%s[:0]", ftName, ftName))
			case protoreflect.MessageKind:
				// If it is message we check if is nil then we leave it
				// If it is from same package use Pool
				g.P(fmt.Sprintf("if x.%s != nil {", ftName))
				ftPkg := string(ft.Desc.Message().FullName().Parent())
				if ftPkg != string(file.GoPackageName) {
					g.P(ftPkg, ".Pool", ft.Desc.Message().Name(), ".Put(x.", ftName, ")")
				} else {
					g.P("Pool", ft.Desc.Message().Name(), ".Put(x.", ftName, ")")
				}
				g.P("x.", ftName, " = nil")
				g.P("}") // end of if
			default:
				g.P(fmt.Sprintf("x.%s = %s", ftName, z.ZeroValue(ft.Desc)))

			}
		}
	}
	g.P("p.pool.Put(x)")
	g.P("}") // end of func Put()
	g.P()
	g.P(fmt.Sprintf("var Pool%s = pool%s{}", messageName, messageName))
	g.P()
}
func genDeepCopy(file *protogen.File, mt *protogen.Message, g *protogen.GeneratedFile) {
	mtName := mt.Desc.Name()
	g.P("func (x *", mtName, ") DeepCopy(z *", mtName, ") {")
	for _, ft := range mt.Fields {
		ftName := ft.Desc.Name()
		ftPkg, ftType := z.DescName(file, g, ft.Desc.Message())
		switch ft.Desc.Cardinality() {
		case protoreflect.Repeated:
			switch ft.Desc.Kind() {
			case protoreflect.MessageKind:
				g.P("for idx := range x.", ftName, "{")
				g.P(fmt.Sprintf("if x.%s[idx] != nil {", ftName))
				if ftPkg == "" {
					g.P("xx := Pool", ftType, ".Get()")
				} else {
					g.P("xx := ", ftPkg, ".Pool", ftType, ".Get()")
				}
				g.P("x.", ftName, "[idx].DeepCopy(xx)")
				g.P("z.", ftName, " = append(z.", ftName, ", xx)")
				g.P("}")
				g.P("}")
			default:
				g.P("z.", ftName, " = append(z.", ftName, "[:0], x.", ftName, "...)")
			}
		default:
			switch ft.Desc.Kind() {
			case protoreflect.BytesKind:
				g.P("z.", ftName, " = append(z.", ftName, "[:0], x.", ftName, "...)")
			case protoreflect.MessageKind:
				// If it is message we check if is nil then we leave it
				// If it is from same package use Pool
				g.P(fmt.Sprintf("if x.%s != nil {", ftName))
				if ftPkg == "" {
					g.P("z.", ftName, " = Pool", ftType, ".Get()")
				} else {
					g.P("z.", ftName, " = ", ftPkg, ".Pool", ftType, ".Get()")
				}
				g.P("x.", ftName, ".DeepCopy(z.", ftName, ")")
				g.P("}")
			default:
				g.P(fmt.Sprintf("z.%s = x.%s", ftName, ftName))

			}
		}
	}
	g.P("}")
	g.P()
}
func genPushToContext(mt *protogen.Message, g *protogen.GeneratedFile) {
	mtName := mt.Desc.Name()
	g.P("func (x *", mtName, ") PushToContext(ctx *edge.RequestCtx) {")
	g.P("ctx.PushMessage(C_", mtName, ", x)")
	g.P("}")
	g.P()
}
func genUnmarshal(mt *protogen.Message, g *protogen.GeneratedFile) {
	mtName := mt.Desc.Name()
	g.P("func (x *", mtName, ") Unmarshal(b []byte) error {")
	g.P("return proto.UnmarshalOptions{}.Unmarshal(b, x)")
	g.P("}")
	g.P()
}
func genMarshal(mt *protogen.Message, g *protogen.GeneratedFile) {
	mtName := mt.Desc.Name()
	g.P("func (x *", mtName, ") Marshal() ([]byte, error) {")
	g.P("return proto.Marshal(x)")
	g.P("}")
	g.P()
}

// GenRPC generates the server and client interfaces if any proto service has been defined
func GenRPC(file *protogen.File, s *protogen.Service, g *protogen.GeneratedFile) {
	genServerRPC(file, g, s)
	genExecuteRemoteRPC(file, g, s)
	genClientRPC(file, g, s)
}
func genServerRPC(file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	serviceName := string(s.Desc.Name())
	g.P("type I", s.Desc.Name(), " interface {")
	for _, m := range s.Methods {
		inputName := z.Name(file, g, m.Desc.Input())
		outputName := z.Name(file, g, m.Desc.Output())
		g.P(m.Desc.Name(), "(ctx *edge.RequestCtx, req *", inputName, ", res *", outputName, ")")
	}
	g.P("}")
	g.P()
	g.P()
	g.P("type ", tools.ToLowerCamel(serviceName), "Wrapper struct {")
	g.P("h I", s.Desc.Name())
	g.P("}")
	g.P()
	for _, m := range s.Methods {
		methodName := string(m.Desc.Name())
		inputPkg, inputType := z.DescName(file, g, m.Desc.Input())
		outputPkg, outputType := z.DescName(file, g, m.Desc.Output())

		g.P("func (sw *", tools.ToLowerCamel(serviceName), "Wrapper) ", tools.ToLowerCamel(methodName), "Wrapper (ctx *edge.RequestCtx, in *rony.MessageEnvelope) {")
		if inputPkg == "" {
			g.P("req := Pool", inputType, ".Get()")
			g.P("defer Pool", inputType, ".Put(req)")
		} else {
			g.P("req := ", inputPkg, ".Pool", inputType, ".Get()")
			g.P("defer ", inputPkg, ".Pool", inputType, ".Put(req)")
		}
		if outputPkg == "" {
			g.P("res := Pool", outputType, ".Get()")
			g.P("defer Pool", outputType, ".Put(res)")
		} else {
			g.P("res := ", outputPkg, ".Pool", outputType, ".Get()")
			g.P("defer ", outputPkg, ".Pool", outputType, ".Put(res)")
		}

		g.P("err := proto.UnmarshalOptions{Merge:true}.Unmarshal(in.Message, req)")
		g.P("if err != nil {")
		g.P("ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)")
		g.P("return")
		g.P("}")
		g.P()
		g.P("sw.h.", methodName, "(ctx, req, res)")
		g.P("if !ctx.Stopped() {")
		if outputPkg == "" {
			g.P("ctx.PushMessage(C_", outputType, ", res)")
		} else {
			g.P("ctx.PushMessage(", outputPkg, ".C_", outputType, ", res)")
		}
		g.P("}") // end of if block
		g.P("}") // end of func block
		g.P()
	}
	g.P("func (sw *", tools.ToLowerCamel(serviceName), "Wrapper) Register (e *edge.Server, ho *edge.HandlerOptions) {")
	for _, m := range s.Methods {
		methodName := string(m.Desc.Name())
		leaderOnlyText := "true"
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		leaderOnly := proto.GetExtension(opt, rony.E_FollowerOk).(bool)
		if leaderOnly {
			leaderOnlyText = "false"
		}
		g.P("e.SetHandlers(C_", methodName, ", ", leaderOnlyText, ", ho.ApplyTo(sw.", tools.ToLowerCamel(methodName), "Wrapper)...)")
	}
	g.P("}")
	g.P()
	g.P("func Register", s.Desc.Name(), "(h I", s.Desc.Name(), ", e *edge.Server, ho *edge.HandlerOptions) {")
	g.P("w := ", tools.ToLowerCamel(serviceName), "Wrapper{")
	g.P("h: h,")
	g.P("}")
	g.P("w.Register(e, ho)")
	g.P("}")
	g.P()

}
func genClientRPC(file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	g.P("type ", s.Desc.Name(), "Client struct {")
	g.P("c edgec.Client")
	g.P("}")
	g.P()
	g.P("func New", s.Desc.Name(), "Client (ec edgec.Client) *", s.Desc.Name(), "Client {")
	g.P("return &", s.Desc.Name(), "Client{")
	g.P("c: ec,")
	g.P("}")
	g.P("}")
	g.P()
	for _, m := range s.Methods {
		inputName := z.Name(file, g, m.Desc.Input())
		outputName := z.Name(file, g, m.Desc.Output())
		outputPkg, outputType := z.DescName(file, g, m.Desc.Output())

		leaderOnlyText := "true"
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		leaderOnly := proto.GetExtension(opt, rony.E_FollowerOk).(bool)
		if leaderOnly {
			leaderOnlyText = "false"
		}
		// constructor := crc32.ChecksumIEEE([]byte(*m.Name))
		g.P("func (c *", s.Desc.Name(), "Client) ", m.Desc.Name(), "(req *", inputName, ", kvs ...*rony.KeyValue) (*", outputName, ", error) {")
		g.P("out := rony.PoolMessageEnvelope.Get()")
		g.P("defer rony.PoolMessageEnvelope.Put(out)")
		g.P("in := rony.PoolMessageEnvelope.Get()")
		g.P("defer rony.PoolMessageEnvelope.Put(in)")
		g.P("out.Fill(c.c.GetRequestID(), C_", m.Desc.Name(), ", req, kvs...)")
		g.P("err := c.c.Send(out, in, ", leaderOnlyText, ")")
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("switch in.GetConstructor() {")
		if outputPkg != "" {
			g.P("case ", outputPkg, ".C_", outputType, ":")
			g.P("x := &", outputPkg, ".", outputType, "{}")
		} else {
			g.P("case C_", outputType, ":")
			g.P("x := &", outputType, "{}")
		}

		g.P("_ = proto.Unmarshal(in.Message, x)")
		g.P("return x, nil")
		g.P("case rony.C_Error:")
		g.P("x := &rony.Error{}")
		g.P("_ = proto.Unmarshal(in.Message, x)")
		g.P("return nil, fmt.Errorf(\"%s:%s\", x.GetCode(), x.GetItems())")
		g.P("default:")
		g.P("return nil, fmt.Errorf(\"unknown message: %d\", in.GetConstructor())")
		g.P("}")
		g.P("}")
		g.P()
	}
}
func genExecuteRemoteRPC(file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	for _, m := range s.Methods {
		leaderOnlyText := "true"
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		leaderOnly := proto.GetExtension(opt, rony.E_FollowerOk).(bool)
		if leaderOnly {
			leaderOnlyText = "false"
		}

		inputName := z.Name(file, g, m.Desc.Input())
		outputName := z.Name(file, g, m.Desc.Output())
		outputC := z.Constructor(file, g, m.Desc.Output())

		g.P("func ExecuteRemote", m.Desc.Name(), "(ctx *edge.RequestCtx, replicaSet uint64, req *", inputName, ", res *", outputName, ", kvs ...*rony.KeyValue) error {")
		g.P("out := rony.PoolMessageEnvelope.Get()")
		g.P("defer rony.PoolMessageEnvelope.Put(out)")
		g.P("in := rony.PoolMessageEnvelope.Get()")
		g.P("defer rony.PoolMessageEnvelope.Put(in)")
		g.P("out.Fill(ctx.ReqID(), C_", m.Desc.Name(), ", req, kvs...)")
		g.P("err := ctx.ExecuteRemote(replicaSet, ", leaderOnlyText, ", out, in)")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("")
		g.P("switch in.GetConstructor() {")
		g.P("case ", outputC, ":")
		g.P("_ = res.Unmarshal(in.GetMessage())")
		g.P("return nil")
		g.P("case rony.C_Error:")
		g.P("x := &rony.Error{}")
		g.P("_ = x.Unmarshal(in.GetMessage())")
		g.P("return x")
		g.P("default:")
		g.P("return edge.ErrUnexpectedTunnelResponse")
		g.P("}")
		g.P("}")
		g.P()
	}
}

// GenCobraCmd
func GenCobraCmd(file *protogen.File, s *protogen.Service, g *protogen.GeneratedFile) {
	genPrepareFunc(s, g)
	genMethodGenerators(file, g, s)
	genClientCliInterface(g, s)
}
func genPrepareFunc(s *protogen.Service, g *protogen.GeneratedFile) {
	serviceName := string(s.Desc.Name())
	opt, _ := s.Desc.Options().(*descriptorpb.ServiceOptions)
	clientProto := proto.GetExtension(opt, rony.E_RonyCobraCmdProtocol).(string)
	g.P("func prepare", serviceName, "Command(cmd *cobra.Command, header map[string]string) (*", serviceName, "Client, error) {")
	g.P("// Bind the current flags to registered flags in config package")
	g.P("err := config.BindCmdFlags(cmd)")
	g.P("if err != nil {")
	g.P("return nil, err")
	g.P("}")
	g.P()
	switch strings.ToLower(clientProto) {
	case "ws":
		g.P("wsC := edgec.NewWebsocket(edgec.WebsocketConfig{")
		g.P("SeedHostPort: fmt.Sprintf(\"%s:%d\", config.GetString(\"host\"), config.GetInt(\"port\")),")
		g.P("Header: header,")
		g.P("})")
		g.P("err = wsC.Start()")
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("return New", serviceName, "Client(wsC), nil")
	default:
		g.P("httpC := edgec.NewHttp(edgec.HttpConfig{")
		g.P("Name: \"Rony Client\",")
		g.P("SeedHostPort: fmt.Sprintf(\"%s:%d\", config.GetString(\"host\"), config.GetInt(\"port\")),")
		g.P("Header: header,")
		g.P("})")
		g.P()
		g.P("err = httpC.Start()")
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("return New", serviceName, "Client(httpC), nil")
	}
	g.P("}")
}
func genMethodGenerators(file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	serviceName := string(s.Desc.Name())
	for _, m := range s.Methods {
		methodName := string(m.Desc.Name())
		g.P("var gen", methodName, "Cmd = func(h I", serviceName, "Cli, header map[string]string) *cobra.Command {")
		g.P("cmd := &cobra.Command {")
		g.P("Use: \"", tools.ToKebab(methodName), "\",")
		g.P("RunE: func(cmd *cobra.Command, args []string) error {")
		g.P("cli, err := prepare", serviceName, "Command(cmd, header)")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}") // end if if clause
		g.P("return h.", methodName, "(cli, cmd, args)")
		g.P("},") // end of RunE func block
		g.P("}")  // end of cobra.Command
		g.P("config.SetFlags(cmd,")
		for _, f := range m.Input.Fields {
			fieldName := string(f.Desc.Name())
			switch z.GoKind(file, g, f.Desc) {
			case "string", "[]byte":
				g.P("config.StringFlag(\"", tools.ToLowerCamel(fieldName), "\",\"\", \"\"),")
			case "int64":
				g.P("config.Int64Flag(\"", tools.ToLowerCamel(fieldName), "\",0, \"\"),")
			case "uint64":
				g.P("config.Uint64Flag(\"", tools.ToLowerCamel(fieldName), "\",0, \"\"),")
			case "int32":
				g.P("config.Int32Flag(\"", tools.ToLowerCamel(fieldName), "\",0, \"\"),")
			case "uint32":
				g.P("config.Uint32Flag(\"", tools.ToLowerCamel(fieldName), "\",0, \"\"),")
			case "bool":
				g.P("config.BoolFlag(\"", tools.ToLowerCamel(fieldName), "\",false, \"\"),")
			default:
			}
		}
		g.P(")") // end of SetFlags
		g.P("return cmd")
		g.P("}") // end of function
		g.P()
	}
}
func genClientCliInterface(g *protogen.GeneratedFile, s *protogen.Service) {
	g.P("type I", s.Desc.Name(), "Cli interface {")
	for _, m := range s.Methods {
		g.P(m.Desc.Name(), "(cli *", s.Desc.Name(), "Client, cmd *cobra.Command, args []string) error")
	}
	g.P("}")
	g.P()
	g.P("func Register", s.Desc.Name(), "Cli (h I", s.Desc.Name(), "Cli, header map[string]string, rootCmd *cobra.Command) {")
	g.P("config.SetPersistentFlags(rootCmd, ")
	g.P("config.StringFlag(\"host\", \"127.0.0.1\", \"the seed host's address\"),")
	g.P("config.StringFlag(\"port\", \"80\", \"the seed host's port\"),")
	g.P(")") // end of SetPersistentFlags
	g.P("rootCmd.AddCommand(")
	var names []string
	for _, m := range s.Methods {
		methodName := string(m.Desc.Name())
		names = append(names, fmt.Sprintf("gen%sCmd(h, header)", methodName))
		if len(names) == 3 {
			sb := strings.Builder{}
			for _, name := range names {
				sb.WriteString(name)
				sb.WriteRune(',')
			}
			g.P(sb.String())
			names = names[:0]
		}
	}
	if len(names) > 0 {
		sb := strings.Builder{}
		for _, name := range names {
			sb.WriteString(name)
			sb.WriteRune(',')
		}
		g.P(sb.String())
		names = names[:0]
	}
	g.P(")") // end of rootCmd.AddCommand
	g.P("}") // end of Register func block
}
