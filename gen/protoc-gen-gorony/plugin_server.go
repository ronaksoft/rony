package main

import (
	"fmt"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"hash/crc32"
	"strings"
)

/*
   Creation Time: 2020 - Jul - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type GenServer struct {
	g *generator.Generator
}

func (g *GenServer) Name() string {
	return "gen-server"
}

func (g *GenServer) Init(gen *generator.Generator) {
	if gen == nil {
		g.g = generator.New()
	} else {
		g.g = gen
	}
}

func (g *GenServer) Generate(file *generator.FileDescriptor) {
	if len(file.Service) > 0 {
		g.g.AddImport("git.ronaksoftware.com/ronak/rony/edge")
		if *file.Package != "rony" {
			g.g.AddImport("git.ronaksoftware.com/ronak/rony")
		}

	}
	for _, s := range file.Service {
		for _, m := range s.Method {
			constructor := crc32.ChecksumIEEE([]byte(*m.Name))
			g.g.P("const C_", m.Name, " int64 = ", fmt.Sprintf("%d", constructor))
		}
		g.g.P()
		g.g.P("type I", s.Name, " interface {")
		g.g.In()
		for _, m := range s.Method {
			inputType := strings.Split(*m.InputType, ".")[2]
			outputType := strings.Split(*m.OutputType, ".")[2]
			g.g.P(m.Name, "(ctx *edge.RequestCtx, req *", inputType, ", res *", outputType, ")")
		}
		g.g.Out()
		g.g.P("}")
		g.g.P()
		g.g.P()
		g.g.P("type ", s.Name, "Wrapper struct {")
		g.g.In()
		g.g.P("h I", s.Name)
		g.g.Out()
		g.g.P("}")
		g.g.P()
		g.g.P("func New", s.Name, "Server(h I", s.Name, ") ", s.Name, "Wrapper {")
		g.g.In()
		g.g.P("return ", s.Name, "Wrapper{")
		g.g.In()
		g.g.P("h: h,")
		g.g.Out()
		g.g.P("}")
		g.g.Out()
		g.g.P("}")
		g.g.P()
		g.g.P("func (sw *", s.Name, "Wrapper) Register (e *edge.Server) {")
		g.g.In()
		for _, m := range s.Method {
			g.g.P("e.AddHandler(C_", m.Name, ", sw.", m.Name, "Wrapper)")
		}
		g.g.Out()
		g.g.P("}")
		g.g.P()
		for _, m := range s.Method {
			inputType := strings.Split(*m.InputType, ".")[2]
			outputType := strings.Split(*m.OutputType, ".")[2]
			g.g.P("func (sw *", s.Name, "Wrapper) ", m.Name, "Wrapper (ctx *edge.RequestCtx, in *rony.MessageEnvelope) {")
			g.g.In()
			g.g.P("req := Pool", inputType, ".Get()")
			g.g.P("defer Pool", inputType, ".Put(req)")
			g.g.P("res := Pool", outputType, ".Get()")
			g.g.P("defer Pool", outputType, ".Put(res)")
			g.g.P("err := req.Unmarshal(in.Message)")
			g.g.P("if err != nil {")
			g.g.In()
			g.g.P("ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)")
			g.g.P("return")
			g.g.Out()
			g.g.P("}")
			g.g.P()
			g.g.P("sw.h.", m.Name, "(ctx, req, res)")
			g.g.P("ctx.PushMessage(C_", outputType, ", res)")
			g.g.Out()
			g.g.P("}")
			g.g.P()
		}
	}
}

func (g *GenServer) GenerateImports(file *generator.FileDescriptor) {

}
