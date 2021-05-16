package helper

import (
	"fmt"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"hash/crc32"
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
	g.g.P("package ", g.f.GoPackageName)
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "sync"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/registry"})
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})
	if g.f.GoPackageName != "rony" {
		g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/edge"})
	}

	initFunc := &strings.Builder{}
	initFunc.WriteString("func init() {\n")
	for _, m := range g.f.Messages {
		g.genPool(m, initFunc)
		g.genDeepCopy(m)
		g.genMarshal(m)
		g.genUnmarshal(m)

		if g.f.GoPackageName != "rony" {
			g.genPushToContext(m)
		}

	}
	for _, st := range g.f.Services {
		for _, m := range st.Methods {
			methodName := fmt.Sprintf("%s%s", st.Desc.Name(), m.Desc.Name())
			constructor := crc32.ChecksumIEEE([]byte(methodName))
			initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, methodName))
			g.g.P("const C_", methodName, " int64 = ", fmt.Sprintf("%d", constructor))
		}
	}

	initFunc.WriteString("}")
	g.g.P("")
	g.g.P(initFunc.String())
	g.g.P()
}

func (g *Generator) genPool(m *protogen.Message, initFunc *strings.Builder) {
	messageName := m.Desc.Name()
	constructor := crc32.ChecksumIEEE([]byte(messageName))
	g.g.P(fmt.Sprintf("const C_%s int64 = %d", messageName, constructor))
	initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, messageName))
	g.g.P(fmt.Sprintf("type pool%s struct{", messageName))
	g.g.P("pool sync.Pool")
	g.g.P("}") // end of pool struct
	g.g.P(fmt.Sprintf("func (p *pool%s) Get() *%s {", messageName, messageName))
	g.g.P(fmt.Sprintf("x, ok := p.pool.Get().(*%s)", messageName))
	g.g.P("if !ok {")
	g.g.P(fmt.Sprintf("x = &%s{}", messageName))
	g.g.P("}") // end of if clause
	g.g.P("return x")
	g.g.P("}") // end of func Get()
	g.g.P()
	g.g.P(fmt.Sprintf("func (p *pool%s) Put(x *%s) {", messageName, messageName))
	g.g.P("if x == nil {")
	g.g.P("return")
	g.g.P("}")
	for _, ft := range m.Fields {
		ftName := ft.Desc.Name()
		switch ft.Desc.Cardinality() {
		case protoreflect.Repeated:
			switch ft.Desc.Kind() {
			case protoreflect.BytesKind:
				g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/pools"})
				g.g.P("for _, z := range x.", ftName, "{")
				g.g.P("pools.Bytes.Put(z)")
				g.g.P("}") // end of for/range
				g.g.P("x.", ftName, " = x.", ftName, "[:0]")
			case protoreflect.MessageKind:
				// If it is message we check if is nil then we leave it
				// If it is from same package use Pool
				g.g.P("for _, z := range x.", ftName, "{")
				ftPkg := string(ft.Desc.Message().FullName().Parent())
				if ftPkg != string(g.f.GoPackageName) {
					g.g.P(ftPkg, ".Pool", ft.Desc.Message().Name(), ".Put(z)")
				} else {
					g.g.P("Pool", ft.Desc.Message().Name(), ".Put(z)")
				}
				g.g.P("}") // end of for/range
				g.g.P("x.", ftName, " = x.", ftName, "[:0]")

			default:
				g.g.P("x.", ftName, " = x.", ftName, "[:0]")
			}
		default:
			switch ft.Desc.Kind() {
			case protoreflect.BytesKind:
				g.g.P(fmt.Sprintf("x.%s = x.%s[:0]", ftName, ftName))
			case protoreflect.MessageKind:
				// If it is message we check if is nil then we leave it
				// If it is from same package use Pool
				ftPkg := string(ft.Desc.Message().FullName().Parent())
				if ftPkg != string(g.f.GoPackageName) {
					g.g.P(ftPkg, ".Pool", ft.Desc.Message().Name(), ".Put(x.", ftName, ")")
				} else {
					g.g.P("Pool", ft.Desc.Message().Name(), ".Put(x.", ftName, ")")
				}
				g.g.P("x.", ftName, " = nil")
			default:
				g.g.P(fmt.Sprintf("x.%s = %s", ftName, z.ZeroValue(ft.Desc)))
			}
		}
	}
	g.g.P("p.pool.Put(x)")
	g.g.P("}") // end of func Put()
	g.g.P()
	g.g.P(fmt.Sprintf("var Pool%s = pool%s{}", messageName, messageName))
	g.g.P()
}
func (g *Generator) genDeepCopy(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") DeepCopy(z *", mtName, ") {")
	for _, ft := range m.Fields {
		ftName := ft.Desc.Name()
		ftPkg, ftType := z.DescName(g.f, g.g, ft.Desc.Message())
		switch ft.Desc.Cardinality() {
		case protoreflect.Repeated:
			switch ft.Desc.Kind() {
			case protoreflect.MessageKind:
				g.g.P("for idx := range x.", ftName, "{")
				g.g.P(fmt.Sprintf("if x.%s[idx] != nil {", ftName))
				if ftPkg == "" {
					g.g.P("xx := Pool", ftType, ".Get()")
				} else {
					g.g.P("xx := ", ftPkg, ".Pool", ftType, ".Get()")
				}
				g.g.P("x.", ftName, "[idx].DeepCopy(xx)")
				g.g.P("z.", ftName, " = append(z.", ftName, ", xx)")
				g.g.P("}")
				g.g.P("}")
			default:
				g.g.P("z.", ftName, " = append(z.", ftName, "[:0], x.", ftName, "...)")
			}
		default:
			switch ft.Desc.Kind() {
			case protoreflect.BytesKind:
				g.g.P("z.", ftName, " = append(z.", ftName, "[:0], x.", ftName, "...)")
			case protoreflect.MessageKind:
				// If it is message we check if is nil then we leave it
				// If it is from same package use Pool
				g.g.P("if x.", ftName, " != nil {")
				g.g.P("if z.", ftName, " == nil {")
				if ftPkg == "" {
					g.g.P("z.", ftName, " = Pool", ftType, ".Get()")
				} else {
					g.g.P("z.", ftName, " = ", ftPkg, ".Pool", ftType, ".Get()")
				}
				g.g.P("}")
				g.g.P("x.", ftName, ".DeepCopy(z.", ftName, ")")
				g.g.P("} else {")
				g.g.P("z.", ftName, "= nil")
				g.g.P("}")
			default:
				g.g.P(fmt.Sprintf("z.%s = x.%s", ftName, ftName))

			}
		}
	}
	g.g.P("}")
	g.g.P()
}
func (g *Generator) genPushToContext(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") PushToContext(ctx *edge.RequestCtx) {")
	g.g.P("ctx.PushMessage(C_", mtName, ", x)")
	g.g.P("}")
	g.g.P()
}
func (g *Generator) genUnmarshal(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") Unmarshal(b []byte) error {")
	g.g.P("return proto.UnmarshalOptions{}.Unmarshal(b, x)")
	g.g.P("}")
	g.g.P()
}
func (g *Generator) genMarshal(m *protogen.Message) {
	mtName := m.Desc.Name()
	g.g.P("func (x *", mtName, ") Marshal() ([]byte, error) {")
	g.g.P("return proto.Marshal(x)")
	g.g.P("}")
	g.g.P()
}
