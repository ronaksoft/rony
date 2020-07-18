package main

import (
	"fmt"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"hash/crc32"
	"strings"
)

/*
   Creation Time: 2019 - Nov - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type GenPools struct {
	g *generator.Generator
}

func (g *GenPools) Name() string {
	return "gen-pools"
}

func (g *GenPools) Init(gen *generator.Generator) {
	if gen == nil {
		g.g = generator.New()
	} else {
		g.g = gen
	}

}

func (g *GenPools) Generate(file *generator.FileDescriptor) {
	g.g.AddImport("sync")
	initFunc := strings.Builder{}
	initFunc.WriteString("func init() {\n")
	for _, mt := range file.MessageType {
		constructor := crc32.ChecksumIEEE([]byte(*mt.Name))
		g.g.P(fmt.Sprintf("const C_%s int64 = %d", *mt.Name, constructor))
		initFunc.WriteString(fmt.Sprintf("ConstructorNames[%d] = \"%s\"\n", constructor, *mt.Name))
		g.g.P(fmt.Sprintf("type pool%s struct{", *mt.Name))
		g.g.In()
		g.g.P("pool sync.Pool")
		g.g.Out()
		g.g.P("}")
		g.g.P(fmt.Sprintf("func (p *pool%s) Get() *%s {", *mt.Name, *mt.Name))
		g.g.In()
		g.g.P(fmt.Sprintf("x, ok := p.pool.Get().(*%s)", *mt.Name))
		g.g.P("if !ok {")
		g.g.In()
		g.g.P(fmt.Sprintf("return &%s{}", *mt.Name))
		g.g.Out()
		g.g.P("}")
		g.g.P("return x")
		g.g.Out()
		g.g.P("}")
		g.g.P("", "")

		g.g.P(fmt.Sprintf("func (p *pool%s) Put(x *%s) {", *mt.Name, *mt.Name))
		g.g.In()
		for _, ft := range mt.Field {
			if *ft.Label == descriptor.FieldDescriptorProto_LABEL_REPEATED || *ft.Type == descriptor.FieldDescriptorProto_TYPE_BYTES {
				g.g.P(fmt.Sprintf("x.%s = x.%s[:0]", *ft.Name, *ft.Name))
			} else if *ft.Label == descriptor.FieldDescriptorProto_LABEL_OPTIONAL {
				g.g.P(fmt.Sprintf("x.%s = %s", *ft.Name, zeroValue(ft.Type)))
			}
		}
		g.g.P("p.pool.Put(x)")
		g.g.Out()
		g.g.P("}")
		g.g.P("")
		g.g.P(fmt.Sprintf("var Pool%s = pool%s{}", *mt.Name, *mt.Name))
		g.g.P("")
	}
	initFunc.WriteString("}")
	g.g.P("")
	g.g.P(initFunc.String())

}

func zeroValue(t *descriptor.FieldDescriptorProto_Type) string {
	switch *t {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "false"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "\"\""
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		return "nil"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "nil"
	default:
		return "0"
	}
}

func (g *GenPools) GenerateImports(file *generator.FileDescriptor) {
	g.g.AddImport("sync")
	// g.g.AddImport("git.ronaksoftware.com/river")
}
