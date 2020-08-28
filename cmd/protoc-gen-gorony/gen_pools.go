package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
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

func GenPools(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("package ", file.GoPackageName)
	g.P("import (")
	g.P("\"sync\"")
	g.P("\"git.ronaksoft.com/ronak/rony/registry\"")
	g.P(")")
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
