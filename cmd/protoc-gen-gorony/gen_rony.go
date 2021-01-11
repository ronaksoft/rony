package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
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

// GenPools generates codes related for pooling of the messages
func GenPools(file *protogen.File, g *protogen.GeneratedFile) {
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

	initFunc := strings.Builder{}
	initFunc.WriteString("func init() {\n")
	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		constructor := crc32.ChecksumIEEE([]byte(mt.Desc.Name()))
		g.P(fmt.Sprintf("const C_%s int64 = %d", mtName, constructor))
		initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, mt.Desc.Name()))
		g.P(fmt.Sprintf("type pool%s struct{", mtName))
		g.P("pool sync.Pool")
		g.P("}")
		g.P(fmt.Sprintf("func (p *pool%s) Get() *%s {", mtName, mtName))
		g.P(fmt.Sprintf("x, ok := p.pool.Get().(*%s)", mtName))
		g.P("if !ok {")
		g.P(fmt.Sprintf("return &%s{}", mtName))
		g.P("}")
		g.P("return x")
		g.P("}")
		g.P("", "")
		g.P(fmt.Sprintf("func (p *pool%s) Put(x *%s) {", mtName, mtName))
		for _, ft := range mt.Fields {
			ftName := ft.Desc.Name()
			ftPkg, _ := descName(file, g, ft.Desc.Message())
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
					if ftPkg != "" {
						g.P(ftPkg, ".Pool", ft.Desc.Message().Name(), ".Put(x.", ftName, ")")
					} else {
						g.P("Pool", ft.Desc.Message().Name(), ".Put(x.", ftName, ")")
					}
					g.P("x.", ftName, " = nil")
					g.P("}")
				default:
					g.P(fmt.Sprintf("x.%s = %s", ftName, zeroValue(ft.Desc.Kind())))

				}
			}
		}
		g.P("p.pool.Put(x)")
		g.P("}")
		g.P("")
		g.P(fmt.Sprintf("var Pool%s = pool%s{}", mtName, mtName))
		g.P("")
	}
	for _, st := range file.Services {
		for _, m := range st.Methods {
			constructor := crc32.ChecksumIEEE([]byte(m.Desc.Name()))
			initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, m.Desc.Name()))
		}
	}

	initFunc.WriteString("}")
	g.P("")
	g.P(initFunc.String())
	g.P()
}

// GenDeepCopy generates codes which deep copy a message
func GenDeepCopy(file *protogen.File, g *protogen.GeneratedFile) {
	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		g.P("func (x *", mtName, ") DeepCopy(z *", mtName, ") {")
		for _, ft := range mt.Fields {
			ftName := ft.Desc.Name()
			ftPkg, ftType := descName(file, g, ft.Desc.Message())
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
}

// GenPushToContext generates codes related for pooling of the messages
func GenPushToContext(file *protogen.File, g *protogen.GeneratedFile) {
	if file.GoPackageName == "rony" {
		return
	}
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "",
		GoImportPath: "github.com/ronaksoft/rony/edge",
	})
	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		g.P("func (x *", mtName, ") PushToContext(ctx *edge.RequestCtx) {")
		g.P("ctx.PushMessage(C_", mtName, ", x)")
		g.P("}")
		g.P()
	}
}

// GenUnmarshal generates codes related for pooling of the messages
func GenUnmarshal(file *protogen.File, g *protogen.GeneratedFile) {
	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		g.P("func (x *", mtName, ") Unmarshal(b []byte) error {")
		g.P("return proto.UnmarshalOptions{}.Unmarshal(b, x)")
		g.P("}")
		g.P()
	}
}

// GenUnmarshal generates codes related for pooling of the messages
func GenMarshal(file *protogen.File, g *protogen.GeneratedFile) {
	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		g.P("func (x *", mtName, ") Marshal() ([]byte, error) {")
		g.P("return proto.Marshal(x)")
		g.P("}")
		g.P()
	}
}

// GenRPC generates the server and client interfaces if any proto service has been defined
func GenRPC(file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) > 0 {
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "",
			GoImportPath: "github.com/ronaksoft/rony/edge",
		})
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "",
			GoImportPath: "google.golang.org/protobuf/proto",
		})
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "",
			GoImportPath: "fmt",
		})
		if file.GoPackageName != "rony" {
			g.QualifiedGoIdent(protogen.GoIdent{
				GoName:       "",
				GoImportPath: "github.com/ronaksoft/rony",
			})
		}
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "Client",
			GoImportPath: "github.com/ronaksoft/rony/edgec",
		})
	}

	// Generate Server
	for _, s := range file.Services {
		for _, m := range s.Methods {
			constructor := crc32.ChecksumIEEE([]byte(m.Desc.Name()))
			g.P("const C_", m.Desc.Name(), " int64 = ", fmt.Sprintf("%d", constructor))
		}
		g.P()

		genServerRPC(file, g, s)
		genExecuteRemoteRPC(file, g, s)
		genClientRPC(file, g, s)
	}
}
func genServerRPC(file *protogen.File, g *protogen.GeneratedFile, s *protogen.Service) {
	g.P("type I", s.Desc.Name(), " interface {")
	for _, m := range s.Methods {
		inputPkg, inputType := descName(file, g, m.Desc.Input())
		outputPkg, outputType := descName(file, g, m.Desc.Output())
		inputName := inputType
		if inputPkg != "" {
			inputName = fmt.Sprintf("%s.%s", inputPkg, inputType)
		}
		outputName := outputType
		if outputPkg != "" {
			outputName = fmt.Sprintf("%s.%s", outputPkg, outputType)
		}
		g.P(m.Desc.Name(), "(ctx *edge.RequestCtx, req *", inputName, ", res *", outputName, ")")
	}
	g.P("}")
	g.P()
	g.P()
	g.P("type ", s.Desc.Name(), "Wrapper struct {")
	g.P("h I", s.Desc.Name())
	g.P("}")
	g.P()
	g.P("func Register", s.Desc.Name(), "(h I", s.Desc.Name(), ", e *edge.Server) {")
	g.P("w := ", s.Desc.Name(), "Wrapper{")
	g.P("h: h,")
	g.P("}")
	g.P("w.Register(e)")
	g.P("}")
	g.P()
	g.P("func (sw *", s.Desc.Name(), "Wrapper) Register (e *edge.Server) {")
	for _, m := range s.Methods {
		leaderOnlyText := "true"
		opt, _ := m.Desc.Options().(*descriptorpb.MethodOptions)
		leaderOnly := proto.GetExtension(opt, rony.E_FollowerOk).(bool)
		if leaderOnly {
			leaderOnlyText = "false"
		}

		g.P("e.SetHandlers(C_", m.Desc.Name(), ", ", leaderOnlyText, ", sw.", m.Desc.Name(), "Wrapper)")
	}
	g.P("}")
	g.P()
	for _, m := range s.Methods {
		inputPkg, inputType := descName(file, g, m.Desc.Input())
		outputPkg, outputType := descName(file, g, m.Desc.Output())

		g.P("func (sw *", s.Desc.Name(), "Wrapper) ", m.Desc.Name(), "Wrapper (ctx *edge.RequestCtx, in *rony.MessageEnvelope) {")
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
		g.P("sw.h.", m.Desc.Name(), "(ctx, req, res)")
		g.P("if !ctx.Stopped() {")
		if outputPkg == "" {
			g.P("ctx.PushMessage(C_", outputType, ", res)")
		} else {
			g.P("ctx.PushMessage(", outputPkg, ".C_", outputType, ", res)")
		}

		g.P("}")
		g.P("}")
		g.P()
	}
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
		inputPkg, inputType := descName(file, g, m.Desc.Input())
		outputPkg, outputType := descName(file, g, m.Desc.Output())
		inputName := inputType
		if inputPkg != "" {
			inputName = fmt.Sprintf("%s.%s", inputPkg, inputType)
		}
		outputName := outputType
		if outputPkg != "" {
			outputName = fmt.Sprintf("%s.%s", outputPkg, outputType)
		}

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

		inputName := name(file, g, m.Desc.Input())
		inputC := constructor(file, g, m.Desc.Input())
		outputName := name(file, g, m.Desc.Output())
		outputC := constructor(file, g, m.Desc.Output())

		g.P("func ExecuteRemote", m.Desc.Name(), "(ctx *edge.RequestCtx, replicaSet uint64, req *", inputName, ", res *", outputName, ") error {")
		g.P("out := rony.PoolMessageEnvelope.Get()")
		g.P("defer rony.PoolMessageEnvelope.Put(out)")
		g.P("in := rony.PoolMessageEnvelope.Get()")
		g.P("defer rony.PoolMessageEnvelope.Put(in)")
		g.P("out.Fill(ctx.ReqID(), ", inputC, ", req)")
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
