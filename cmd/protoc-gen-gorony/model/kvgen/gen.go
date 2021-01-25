package kvgen

import (
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/model"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
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
		funcList(mm, g)
		funcHasField(m, mm, g)
	}
}
func funcSave(mm *model.Model, g *protogen.GeneratedFile) {
	g.P("func Save", mm.Name, "WithTxn (txn *badger.Txn, alloc *kv.Allocator, m*", mm.Name, ") error {")
	g.P("if alloc == nil {")
	g.P("alloc = kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("}") // end of if block
	g.P()
	g.P("b := alloc.GenValue(m)")
	g.P("err := txn.Set(alloc.GenKey(C_", mm.Name, ",",
		crc32.ChecksumIEEE(tools.StrToByte(mm.Table.String("m.", false, false))),
		",",
		mm.Table.String("m.", false, false),
		"), b)",
	)
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
	g.P()
	g.P("}") // end of SaveWithTxn func
	g.P()
	g.P("func Save", mm.Name, "(m *", mm.Name, ") error {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("return kv.Update(func(txn *badger.Txn) error {")
	g.P("return Save", mm.Name, "WithTxn (txn, alloc, m)")
	g.P("})") // end of Update func
	g.P("}")  // end of Save func
	g.P()
}
func funcRead(mm *model.Model, g *protogen.GeneratedFile) {
	g.P("func Read", mm.Name, "WithTxn (txn *badger.Txn, alloc *kv.Allocator,", mm.FuncArgs(mm.Table, false), ", m *", mm.Name, ") (*", mm.Name, ",error) {")
	g.P("if alloc == nil {")
	g.P("alloc = kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("}") // end of if block
	g.P()
	g.P("item, err := txn.Get(alloc.GenKey(C_", mm.Name, ",",
		crc32.ChecksumIEEE(tools.StrToByte(mm.Table.String("m.", false, false))),
		",",
		mm.Table.String("", false, true),
		"))",
	)
	g.P("if err != nil {")
	g.P("return nil, err")
	g.P("}")
	g.P("err = item.Value(func (val []byte) error {")
	g.P("return m.Unmarshal(val)")
	g.P("})")
	g.P("return m, err")
	g.P("}") // end of Read func
	g.P()
	g.P("func Read", mm.Name, "(", mm.FuncArgs(mm.Table, false), ", m *", mm.Name, ") (*", mm.Name, ",error) {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("if m == nil {")
	g.P("m = &", mm.Name, "{}")
	g.P("}")
	g.P("err := kv.View(func(txn *badger.Txn) (err error) {")
	g.P("m, err = Read", mm.Name, "WithTxn(txn, alloc, ", mm.Table.String("", false, true), ", m)")
	g.P("return err")
	g.P("})") // end of View func
	g.P("return m, err")
	g.P("}") // end of Read func
	g.P()
	for _, pk := range mm.Views {
		g.P("func Read", mm.Name, pk.FuncName("By"), "WithTxn(txn *badger.Txn, alloc *kv.Allocator,", mm.FuncArgs(pk, false), ", m *", mm.Name, ") ( *", mm.Name, ", error) {")
		g.P("if alloc == nil {")
		g.P("alloc = kv.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P("}") // end of if block
		g.P()
		g.P("item, err := txn.Get(alloc.GenKey(C_", mm.Name, ",",
			crc32.ChecksumIEEE(tools.StrToByte(pk.String("m.", false, false))),
			",",
			pk.String("", false, true),
			"))",
		)
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("err = item.Value(func (val []byte) error {")
		g.P("return m.Unmarshal(val)")
		g.P("})") // end of item.Value
		g.P("return m, err")
		g.P("}")  // end of Read func
		g.P()
		g.P("func Read", mm.Name, pk.FuncName("By"), "(", mm.FuncArgs(pk, false), ", m *", mm.Name, ") ( *", mm.Name, ", error) {")
		g.P("alloc := kv.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P("if m == nil {")
		g.P("m = &", mm.Name, "{}")
		g.P("}")
		g.P("err := kv.View(func(txn *badger.Txn) (err error) {")
		g.P("m, err = Read", mm.Name, pk.FuncName("By"), "WithTxn (txn, alloc,", pk.String("", false, true), ", m)")
		g.P("return err")
		g.P("})") // end of View func
		g.P("return m, err")
		g.P("}") // end of Read func
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
		g.P("err = txn.Delete(alloc.GenKey(C_", mm.Name, ",",
			crc32.ChecksumIEEE(tools.StrToByte(mm.Table.String("m.", false, false))),
			",",
			mm.Table.String("", false, true),
			"))",
		)
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
func funcList(mm *model.Model, g *protogen.GeneratedFile) {
	g.P("func List", mm.Name, "(")
	g.P(mm.FuncArgsWithPrefix("offset", mm.Table, true), ", offset int32, limit int32, cond func(m *", mm.Name, ") bool, ")
	g.P(") ([]*", mm.Name, ", error) {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P()
	g.P("res := make([]*", mm.Name, ", 0, limit)")
	g.P("err := kv.View(func(txn *badger.Txn) error {")
	g.P("opt := badger.DefaultIteratorOptions")
	g.P("opt.Prefix = alloc.GenKey(C_", mm.Name, ",",
		crc32.ChecksumIEEE(tools.StrToByte(mm.Table.String("m.", false, false))),
		")",
	)
	g.P("osk := alloc.GenKey(C_", mm.Name, ",",
		crc32.ChecksumIEEE(tools.StrToByte(mm.Table.String("m.", false, false))),
		",",
		mm.Table.String("offset", true, false),
		")",
	)
	g.P("iter := txn.NewIterator(opt)")
	g.P("for iter.Seek(osk); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
	g.P("if offset--; offset >= 0 {")
	g.P("continue")
	g.P("}")
	g.P("if limit--; limit < 0 {")
	g.P("break")
	g.P("}")
	g.P("_ = iter.Item().Value(func (val []byte) error {")
	g.P("m := &", mm.Name, "{}")
	g.P("err := m.Unmarshal(val)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}") // end of if
	g.P("if cond(m) {")
	g.P("res = append(res, m)")
	g.P("}") // end of if cond
	g.P("return nil")
	g.P("})") // end of iter.Value func
	g.P("}")  // end of for
	g.P("iter.Close()")
	g.P("return nil")
	g.P("})") // end of View
	g.P("return res, err")
	g.P("}") // end of func List
	g.P()
}
func funcHasField(m *protogen.Message, mm *model.Model, g *protogen.GeneratedFile) {
	for _, f := range m.Fields {
		switch f.Desc.Cardinality() {
		case protoreflect.Repeated:
			if f.Desc.Kind() == protoreflect.MessageKind {
				break
			}
			mtName := m.Desc.Name()
			g.P("func (x *", mtName, ") Has", f.Desc.Name(),"(xx ",mm.FieldsGo[f.GoName], ") bool {")
			g.P("for idx := range x.", f.Desc.Name(), "{")
			g.P("if x.", f.Desc.Name(), "[idx] == xx {")
			g.P("return true")
			g.P("}")	// end of if
			g.P("}")	// end of for
			g.P("return false")
			g.P("}")	// end of func
			g.P()
		}
	}
}
