package singleton

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

/*
   Creation Time: 2021 - Jan - 29
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
	for _, m := range g.f.Messages {
		opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)
		singleton := proto.GetExtension(opt, rony.E_RonySingleton).(bool)
		if singleton {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})
			g.genSave(m)
			g.genRead(m)
			continue
		}
	}
}

func (g *Generator) genSave(m *protogen.Message) {
	// SaveWithTxn func
	g.g.P("func Save", m.Desc.Name(), "WithTxn (txn *store.LTxn, alloc *tools.Allocator, m *", m.Desc.Name(), ") (err error) {")
	g.g.P("if alloc == nil {")
	g.g.P("alloc = tools.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P("}") // end of if block
	g.g.P()
	g.g.P("err = store.Marshal(txn, alloc, m,", genDbKey(m), ")")
	g.g.P("if err != nil {")
	g.g.P("return")
	g.g.P("}")
	g.g.P("return nil")
	g.g.P("}") // end of SaveWithTxn func
	g.g.P()
	g.g.P("func Save", m.Desc.Name(), "(m *", m.Desc.Name(), ") (err error) {")
	if g.f.GoPackageName != "rony" {
		g.g.P("alloc := tools.NewAllocator()")
	} else {
		g.g.P("alloc := tools.NewAllocator()")
	}
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P("return store.Update(func(txn *store.LTxn) error {")
	g.g.P("return Save", m.Desc.Name(), "WithTxn(txn, alloc, m)")
	g.g.P("})") // end of store.Update func
	g.g.P("}")  // end of Save func
	g.g.P()
}
func (g *Generator) genRead(m *protogen.Message) {
	g.g.P("func Read", m.Desc.Name(), "WithTxn (txn *store.LTxn, alloc *tools.Allocator, m *", m.Desc.Name(), ") (*", m.Desc.Name(), ",error) {")
	g.g.P("if alloc == nil {")
	g.g.P("alloc = tools.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P("}") // end of if block
	g.g.P()
	g.g.P("err := store.Unmarshal(txn, alloc, m, ", genDbKey(m), ")")
	g.g.P("if err != nil {")
	g.g.P("return nil, err")
	g.g.P("}")
	g.g.P("return m, err")
	g.g.P("}") // end of ReadWithTxn func
	g.g.P()
	g.g.P("func Read", m.Desc.Name(), "(m *", m.Desc.Name(), ") (*", m.Desc.Name(), ",error) {")
	if g.f.GoPackageName != "rony" {
		g.g.P("alloc := tools.NewAllocator()")
	} else {
		g.g.P("alloc := tools.NewAllocator()")
	}
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P()
	g.g.P("if m == nil {")
	g.g.P("m = &", m.Desc.Name(), "{}")
	g.g.P("}")
	g.g.P()
	g.g.P("err := store.View(func(txn *store.LTxn) (err error) {")
	g.g.P("m, err = Read", m.Desc.Name(), "WithTxn(txn, alloc,  m)")
	g.g.P("return")
	g.g.P("})") // end of View func
	g.g.P("return m, err")
	g.g.P("}") // end of Read func
	g.g.P()
}
func (g *Generator) genDelete(m *protogen.Message) {}

func genDbKey(m *protogen.Message) string {
	return fmt.Sprintf("'S', C_%s",
		m.Desc.Name(),
	)
}
