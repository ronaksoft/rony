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

type GenRony struct {
	g *generator.Generator
}

func (g *GenRony) Name() string {
	return "gen-server"
}

func (g *GenRony) Init(gen *generator.Generator) {
	if gen == nil {
		g.g = generator.New()
	} else {
		g.g = gen
	}
}

func (g *GenRony) Generate(file *generator.FileDescriptor) {
	if len(file.Service) > 0 {
		g.g.AddImport("git.ronaksoftware.com/ronak/rony/edge")
		g.g.AddImport("git.ronaksoftware.com/ronak/rony/edgeClient")
		g.g.AddImport("git.ronaksoftware.com/ronak/rony/pools")
		if *file.Package != "rony" {
			g.g.AddImport("git.ronaksoftware.com/ronak/rony")
		}

	}

	// Generate Server
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
			g.g.P("if !ctx.Stopped() {")
			g.g.In()
			g.g.P("ctx.PushMessage(C_", outputType, ", res)")
			g.g.Out()
			g.g.P("}")
			g.g.Out()
			g.g.P("}")
			g.g.P()
		}
	}

	// Generate Client
	for _, s := range file.Service {
		g.g.P("type ", s.Name, "Client struct {")
		g.g.In()
		g.g.P("c edgeClient.Client")
		g.g.Out()
		g.g.P("}")
		g.g.P()
		g.g.P("func New", s.Name, "Client (ec edgeClient.Client) *", s.Name, "Client {")
		g.g.In()
		g.g.P("return &", s.Name, "Client{")
		g.g.In()
		g.g.P("c: ec,")
		g.g.Out()
		g.g.P("}")
		g.g.Out()
		g.g.P("}")

		for _, m := range s.Method {
			inputType := strings.Split(*m.InputType, ".")[2]
			outputType := strings.Split(*m.OutputType, ".")[2]
			// constructor := crc32.ChecksumIEEE([]byte(*m.Name))
			g.g.P("func (c *", s.Name, "Client) ", m.Name, "(req *", inputType, ") (*", outputType, ", error) {")
			g.g.In()
			g.g.P("out := rony.PoolMessageEnvelope.Get()")
			g.g.P("defer rony.PoolMessageEnvelope.Put(out)")
			g.g.P("in := rony.PoolMessageEnvelope.Get()")
			g.g.P("defer rony.PoolMessageEnvelope.Put(in)")
			g.g.P("b := pools.Bytes.GetLen(req.Size())")
			g.g.P("req.MarshalToSizedBuffer(b)")
			g.g.P("out.RequestID = c.c.GetRequestID()")
			g.g.P("out.Constructor = C_", m.Name)
			g.g.P("out.Message = append(out.Message[:0], b...)")
			g.g.P("err := c.c.Send(out, in)")
			g.g.P("if err != nil {")
			g.g.In()
			g.g.P("return nil, err")
			g.g.Out()
			g.g.P("}")
			g.g.P("switch in.Constructor {")
			g.g.In()
			g.g.P("case C_", outputType, ":")
			g.g.In()
			g.g.P("x := &", outputType, "{}")
			g.g.P("_ = x.Unmarshal(in.Message)")
			g.g.P("return x, nil")
			g.g.Out()
			g.g.P("case rony.C_Error:")
			g.g.In()
			g.g.P("x := &rony.Error{}")
			g.g.P("_ = x.Unmarshal(in.Message)")
			g.g.P("return nil, fmt.Errorf(\"%s:%s\", x.Code, x.Items)")
			g.g.Out()
			g.g.P("default:")
			g.g.In()
			g.g.P("return nil, fmt.Errorf(\"unknown message: %d\", in.Constructor)")
			g.g.Out()
			g.g.Out()
			g.g.P("}")
			g.g.Out()
			g.g.P("}")

		}
	}
}

func (g *GenRony) GenerateImports(file *generator.FileDescriptor) {

}
