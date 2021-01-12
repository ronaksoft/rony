package kvgen

import (
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/model"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"hash/crc32"
)

/*
   Creation Time: 2021 - Jan - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Generate generates the repo functions for messages which are identified as model with {{@entity cql}}
func Generate(file *protogen.File, g *protogen.GeneratedFile) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/repo/kv"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "badger", GoImportPath: "github.com/dgraph-io/badger"})

	genFuncs(file, g)
}
func genFuncs(file *protogen.File, g *protogen.GeneratedFile) {
	for _, m := range file.Messages {
		mm := model.GetModels()[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		funcCreate(mm, g)
		funcRead(mm, g)
		funcDelete(mm, g)
	}
}
func funcCreate(mm *model.Model, g *protogen.GeneratedFile) {
	g.P("func Create", mm.Name, "(m *", mm.Name, ") error {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("return kv.Update(func(txn *badger.Txn) error {")
	g.P("b := alloc.GenValue(m)")
	g.P("err := txn.Set(alloc.GenKey(C_", mm.Name, ",", mm.Table.StringKeys("m."), "), b)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	for _, pk := range mm.Views {
		g.P("err = txn.Set(alloc.GenKey(C_", mm.Name, ",", crc32.ChecksumIEEE(tools.StrToByte(pk.StringKeys("m."))), ",", pk.StringKeys("m."), "), b)")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P()
	}
	g.P("return nil")
	g.P("})")
	g.P("}")
	g.P()
}
func funcRead(mm *model.Model, g *protogen.GeneratedFile) {
	g.P("func Read", mm.Name, "(m *", mm.Name, ") error {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("return kv.View(func(txn *badger.Txn) error {")
	g.P("item, err := txn.Get(alloc.GenKey(C_", mm.Name, ",", mm.Table.StringKeys("m."), "))")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P("return item.Value(func (val []byte) error {")
	g.P("return m.Unmarshal(val)")
	g.P("})")
	g.P("})")
	g.P("}")
	g.P()
}
func funcDelete(mm *model.Model, g *protogen.GeneratedFile) {
	g.P("func Delete", mm.Name, "(m *", mm.Name, ") error {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("return kv.Update(func(txn *badger.Txn) error {")
	g.P("err := txn.Delete(alloc.GenKey(C_", mm.Name, ",", mm.Table.StringKeys("m."), "))")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	for _, pk := range mm.Views {
		g.P("err = txn.Delete(alloc.GenKey(C_", mm.Name, ",", crc32.ChecksumIEEE(tools.StrToByte(pk.StringKeys("m."))), ",", pk.StringKeys("m."), "))")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P()
	}
	g.P("return nil")
	g.P("})")
	g.P("}")
	g.P()
}
