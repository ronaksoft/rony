package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
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

func descName(file *protogen.File, g *protogen.GeneratedFile, desc protoreflect.MessageDescriptor) (string, string) {
	if string(desc.FullName().Parent()) == string(file.GoPackageName) {
		return "", string(desc.Name())
	} else {
		fd, ok := desc.ParentFile().Options().(*descriptorpb.FileOptions)
		if ok {
			g.QualifiedGoIdent(protogen.GoImportPath(fd.GetGoPackage()).Ident(string(desc.Name())))
		}
		return string(desc.ParentFile().Package()), string(desc.Name())
	}
}

func GenPools(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("package ", file.GoPackageName)
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "",
		GoImportPath: "sync",
	})
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "",
		GoImportPath: "git.ronaksoft.com/ronak/rony/registry",
	})

	initFunc := strings.Builder{}
	initFunc.WriteString("func init() {\n")

	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		constructor := crc32.ChecksumIEEE([]byte(mtName))
		g.P(fmt.Sprintf("const C_%s int64 = %d", mtName, constructor))
		initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, mtName))
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
		g.P("x.Reset()")
		// for _, ft := range mt.Fields {
		// 	ftName := ft.Desc.Name()
		// 	switch ft.Desc.Cardinality() {
		// 	case protoreflect.Repeated:
		// 		g.P(fmt.Sprintf("x.%s = x.%s[:0]", ftName, ftName))
		// 	default:
		// 		switch ft.Desc.Kind() {
		// 		case protoreflect.BytesKind:
		// 			g.P(fmt.Sprintf("x.%s = x.%s[:0]", ftName, ftName))
		// 		case protoreflect.MessageKind:
		// 			g.P(fmt.Sprintf("if x.%s != nil {", ftName))
		// 			g.P(fmt.Sprintf("*x.%s = %s{}", ftName, ft.Desc.Message().Name()))
		// 			g.P("}")
		// 		default:
		// 			g.P(fmt.Sprintf("x.%s = %s", ftName, zeroValue(ft.Desc.Kind())))
		//
		// 		}
		// 	}
		// }
		g.P("p.pool.Put(x)")
		g.P("}")
		g.P("")
		g.P(fmt.Sprintf("var Pool%s = pool%s{}", mtName, mtName))
		g.P("")
	}
	initFunc.WriteString("}")
	g.P("")
	g.P(initFunc.String())
}
func GenRPC(file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) > 0 {
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "",
			GoImportPath: "git.ronaksoft.com/ronak/rony/edge",
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
				GoName:       "MessageEnvelope",
				GoImportPath: "git.ronaksoft.com/ronak/rony",
			})
		}
	}

	// Generate Server
	for _, s := range file.Services {
		for _, m := range s.Methods {
			constructor := crc32.ChecksumIEEE([]byte(m.Desc.Name()))
			g.P("const C_", m.Desc.Name(), " int64 = ", fmt.Sprintf("%d", constructor))
		}
		g.P()
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
		g.P("func New", s.Desc.Name(), "Server(h I", s.Desc.Name(), ") *", s.Desc.Name(), "Wrapper {")
		g.P("return &", s.Desc.Name(), "Wrapper{")
		g.P("h: h,")
		g.P("}")
		g.P("}")
		g.P()
		g.P("func (sw *", s.Desc.Name(), "Wrapper) Register (e *edge.Server) {")

		for _, m := range s.Methods {
			g.P("e.AddHandler(C_", m.Desc.Name(), ", sw.", m.Desc.Name(), "Wrapper)")
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

			g.P("err := proto.Unmarshal(in.Message, req)")
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

	// Generate Client
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "Client",
		GoImportPath: "git.ronaksoft.com/ronak/rony/edgeClient",
	})
	for _, s := range file.Services {
		g.P("type ", s.Desc.Name(), "Client struct {")
		g.P("c edgeClient.Client")
		g.P("}")
		g.P()
		g.P("func New", s.Desc.Name(), "Client (ec edgeClient.Client) *", s.Desc.Name(), "Client {")
		g.P("return &", s.Desc.Name(), "Client{")
		g.P("c: ec,")
		g.P("}")
		g.P("}")

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
			// constructor := crc32.ChecksumIEEE([]byte(*m.Name))
			g.P("func (c *", s.Desc.Name(), "Client) ", m.Desc.Name(), "(req *", inputName, ") (*", outputName, ", error) {")
			g.P("out := rony.PoolMessageEnvelope.Get()")
			g.P("defer rony.PoolMessageEnvelope.Put(out)")
			g.P("in := rony.PoolMessageEnvelope.Get()")
			g.P("defer rony.PoolMessageEnvelope.Put(in)")
			g.P("out.Fill(c.c.GetRequestID(), C_", m.Desc.Name(), ", req)")
			g.P("err := c.c.Send(out, in)")
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

		}
	}
}

func zeroValue(t protoreflect.Kind) string {
	switch t {
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.StringKind:
		return "\"\""
	case protoreflect.MessageKind:
		return "nil"
	case protoreflect.BytesKind:
		return "nil"
	default:
		return "0"
	}
}
