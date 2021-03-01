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

func Generate(file *protogen.File, g *protogen.GeneratedFile) {
	for _, m := range file.Messages {
		opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)
		singleton := proto.GetExtension(opt, rony.E_RonySingleton).(bool)
		if !singleton {
			continue
		}
		g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})

		funcSave(g, m)
		funcRead(g, m)
		funcDelete(g, m)
	}
}
func genDbKey(m *protogen.Message) string {
	return fmt.Sprintf("'S', C_%s",
		m.Desc.Name(),
	)
}
func funcSave(g *protogen.GeneratedFile, m *protogen.Message) {
	// SaveWithTxn func
	g.P("func Save", m.Desc.Name(), "WithTxn (txn *store.Txn, alloc *store.Allocator, m *", m.Desc.Name(), ") (err error) {")
	g.P("if alloc == nil {")
	g.P("alloc = store.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("}") // end of if block
	g.P()
	g.P("err = store.Marshal(txn, alloc, m,", genDbKey(m), ")")
	g.P("if err != nil {")
	g.P("return")
	g.P("}")
	g.P("return nil")
	g.P("}") // end of SaveWithTxn func
	g.P()
	g.P("func Save", m.Desc.Name(), "(m *", m.Desc.Name(), ") (err error) {")
	g.P("alloc := store.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("return store.Update(func(txn *store.Txn) error {")
	g.P("return Save", m.Desc.Name(), "WithTxn(txn, alloc, m)")
	g.P("})") // end of store.Update func
	g.P("}")  // end of Save func
	g.P()
}
func funcRead(g *protogen.GeneratedFile, m *protogen.Message) {
	g.P("func Read", m.Desc.Name(), "WithTxn (txn *store.Txn, alloc *store.Allocator, m *", m.Desc.Name(), ") (*", m.Desc.Name(), ",error) {")
	g.P("if alloc == nil {")
	g.P("alloc = store.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("}") // end of if block
	g.P()
	g.P("err := store.Unmarshal(txn, alloc, m, ", genDbKey(m), ")")
	g.P("if err != nil {")
	g.P("return nil, err")
	g.P("}")
	g.P("return m, err")
	g.P("}") // end of ReadWithTxn func
	g.P()
	g.P("func Read", m.Desc.Name(), "(m *", m.Desc.Name(), ") (*", m.Desc.Name(), ",error) {")
	g.P("alloc := store.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P()
	g.P("if m == nil {")
	g.P("m = &", m.Desc.Name(), "{}")
	g.P("}")
	g.P()
	g.P("err := store.View(func(txn *store.Txn) (err error) {")
	g.P("m, err = Read", m.Desc.Name(), "WithTxn(txn, alloc,  m)")
	g.P("return")
	g.P("})") // end of View func
	g.P("return m, err")
	g.P("}") // end of Read func
	g.P()
}
func funcDelete(g *protogen.GeneratedFile, m *protogen.Message) {}
