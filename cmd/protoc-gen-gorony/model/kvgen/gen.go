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
		funcSave(mm, g)
		funcRead(mm, g)
		funcDelete(mm, g)
	}
}
func funcSave(mm *model.Model, g *protogen.GeneratedFile) {
	g.P("func Save", mm.Name, "(m *", mm.Name, ") error {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("return kv.Update(func(txn *badger.Txn) error {")
	g.P("b := alloc.GenValue(m)")
	g.P("err := txn.Set(alloc.GenKey(C_", mm.Name, ",", mm.Table.String("m.", false, false), "), b)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	for _, pk := range mm.Views {
		g.P(
			"err = txn.Set(alloc.GenKey(C_", mm.Name, ",",
			crc32.ChecksumIEEE(tools.StrToByte(pk.String("m.", false, false))),
			",",
			pk.String("m.", false, false),
			"), b)",
		)
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
	g.P("func Read", mm.Name, "(", mm.FuncArgs(mm.Table, false), ", m *", mm.Name, ") (*", mm.Name, ",error) {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("if m == nil {")
	g.P("m = &", mm.Name, "{}")
	g.P("}")
	g.P("err := kv.View(func(txn *badger.Txn) error {")
	g.P("item, err := txn.Get(alloc.GenKey(C_", mm.Name, ",", mm.Table.String("", false, true), "))")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P("return item.Value(func (val []byte) error {")
	g.P("return m.Unmarshal(val)")
	g.P("})")
	g.P("})")
	g.P("return m, err")
	g.P("}")
	g.P()
	for _, pk := range mm.Views {
		g.P("func Read", mm.Name, pk.FuncName("By"), "(", mm.FuncArgs(pk, false), ", m *", mm.Name, ") ( *", mm.Name, ", error) {")
		g.P("alloc := kv.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P("if m == nil {")
		g.P("m = &", mm.Name, "{}")
		g.P("}")
		g.P("err := kv.View(func(txn *badger.Txn) error {")
		g.P("item, err := txn.Get(alloc.GenKey(C_", mm.Name, ",",
			crc32.ChecksumIEEE(tools.StrToByte(pk.String("", false, false))),
			",",
			pk.String("m.", false, false),
			"))",
		)
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("return item.Value(func (val []byte) error {")
		g.P("return m.Unmarshal(val)")
		g.P("})") // end of item.Value
		g.P("})")
		g.P("return m, err")
		g.P("}")
		g.P()
	}
}
func funcDelete(mm *model.Model, g *protogen.GeneratedFile) {
	g.P("func Delete", mm.Name, "(", mm.FuncArgs(mm.Table, false), ") error {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("return kv.Update(func(txn *badger.Txn) error {")
	if len(mm.Views) > 0 {
		g.P("m := &", mm.Name, "{}")
		g.P("item, err := txn.Get(alloc.GenKey(C_", mm.Name, ", ", mm.Table.String("", false, true), "))")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("err = item.Value(func(val []byte) error {")
		g.P("return m.Unmarshal(val)")
		g.P("})")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("err = txn.Delete(alloc.GenKey(C_", mm.Name, ",", mm.Table.String("", false, true), "))")
	} else {
		g.P("err := txn.Delete(alloc.GenKey(C_", mm.Name, ",", mm.Table.String("", false, true), "))")
	}

	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	for _, pk := range mm.Views {
		g.P("err = txn.Delete(alloc.GenKey(C_", mm.Name, ",",
			crc32.ChecksumIEEE(tools.StrToByte(pk.String("m.", false, false))),
			",",
			pk.String("m.", false, false),
			"))",
		)
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
