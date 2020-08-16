package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"hash/crc32"
)

/*
   Creation Time: 2020 - Aug - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func GenRPC(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("package ", file.GoPackageName)
	g.P("import (")
	if len(file.Services) > 0 {
		g.P("\"git.ronaksoftware.com/ronak/rony/edge\"")
		g.P("\"git.ronaksoftware.com/ronak/rony/edgeClient\"")
		g.P("\"google.golang.org/protobuf/proto\"")
		g.P("\"fmt\"")
		if file.GoPackageName != "rony" {
			g.P("\"git.ronaksoftware.com/ronak/rony\"")
		}
	}
	g.P(")")

	// Generate Server
	for _, s := range file.Services {
		for _, m := range s.Methods {
			constructor := crc32.ChecksumIEEE([]byte(m.Desc.Name()))
			g.P("const C_", m.Desc.Name(), " int64 = ", fmt.Sprintf("%d", constructor))
		}
		g.P()
		g.P("type I", s.Desc.Name(), " interface {")
		for _, m := range s.Methods {
			// inputType := strings.Split(string(m.Desc.Input().FullName()), ".")[2]
			// outputType := strings.Split(string(m.Desc.Output().FullName()), ".")[2]
			inputType := m.Desc.Input().Name()
			outputType := m.Desc.Output().Name()
			g.P(m.Desc.Name(), "(ctx *edge.RequestCtx, req *", inputType, ", res *", outputType, ")")
		}
		g.P("}")
		g.P()
		g.P()
		g.P("type ", s.Desc.Name(), "Wrapper struct {")
		g.P("h I", s.Desc.Name())
		g.P("}")
		g.P()
		g.P("func New", s.Desc.Name(), "Server(h I", s.Desc.Name(), ") ", s.Desc.Name(), "Wrapper {")
		g.P("return ", s.Desc.Name(), "Wrapper{")
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
			inputType := m.Desc.Input().Name()
			outputType := m.Desc.Output().Name()
			g.P("func (sw *", s.Desc.Name(), "Wrapper) ", m.Desc.Name(), "Wrapper (ctx *edge.RequestCtx, in *rony.MessageEnvelope) {")
			g.P("req := Pool", inputType, ".Get()")
			g.P("defer Pool", inputType, ".Put(req)")
			g.P("res := Pool", outputType, ".Get()")
			g.P("defer Pool", outputType, ".Put(res)")
			g.P("err := proto.Unmarshal(in.Message, req)")
			g.P("if err != nil {")
			g.P("ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)")
			g.P("return")
			g.P("}")
			g.P()
			g.P("sw.h.", m.Desc.Name(), "(ctx, req, res)")
			g.P("if !ctx.Stopped() {")
			g.P("ctx.PushMessage(C_", outputType, ", res)")
			g.P("}")
			g.P("}")
			g.P()
		}
	}

	// Generate Client
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
			inputType := m.Desc.Input().Name()
			outputType := m.Desc.Output().Name()
			// constructor := crc32.ChecksumIEEE([]byte(*m.Name))
			g.P("func (c *", s.Desc.Name(), "Client) ", m.Desc.Name(), "(req *", inputType, ") (*", outputType, ", error) {")
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
			g.P("case C_", outputType, ":")
			g.P("x := &", outputType, "{}")
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
