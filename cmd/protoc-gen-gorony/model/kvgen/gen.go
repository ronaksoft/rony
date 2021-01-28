package kvgen

import (
	"fmt"
	"github.com/jinzhu/inflection"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/model"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
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
func genDbKey(mm *model.Model, pk model.Key, keyPrefix string) string {
	lowerCamel := keyPrefix == ""
	return fmt.Sprintf("C_%s, %d, %s",
		mm.Name,
		pk.Checksum(),
		pk.String(keyPrefix, ",", lowerCamel),
	)
}
func genDbPrefixPKs(mm *model.Model, key model.Key, keyPrefix string) string {
	lowerCamel := keyPrefix == ""
	return fmt.Sprintf("C_%s, %d, %s",
		mm.Name,
		key.Checksum(),
		key.StringPKs(keyPrefix, ",", lowerCamel),
	)
}
func genDbPrefixCKs(mm *model.Model, key model.Key, keyPrefix string) string {
	lowerCamel := keyPrefix == ""
	return fmt.Sprintf("C_%s, %d, %s",
		mm.Name,
		key.Checksum(),
		key.StringCKs(keyPrefix, ",", lowerCamel),
	)
}
func genDbIndexKey(mm *model.Model, fieldName string, prefix string, postfix string) string {
	lower := prefix == ""
	return fmt.Sprintf("\"IDX\", C_%s, %d, %s%s%s, %s",
		mm.Name, crc32.ChecksumIEEE([]byte(fieldName)), prefix, fieldName, postfix, mm.Table.String(prefix, ",", lower),
	)
}
func genFuncs(file *protogen.File, g *protogen.GeneratedFile) {
	for _, m := range file.Messages {
		mm := model.GetModels()[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		funcSave(file, g, m, mm)
		funcRead(g, mm)
		funcDelete(g, mm)
		funcList(g, mm)
		funcHasField(g, m, mm)
		funcListByIndex(g, m, mm)
		if len(mm.Table.CKs) > 0 {
			funcListByPartitionKey(g, m, mm)
		}
	}
}
func funcSave(file *protogen.File, g *protogen.GeneratedFile, mt *protogen.Message, mm *model.Model) {
	// SaveWithTxn func
	g.P("func Save", mm.Name, "WithTxn (txn *badger.Txn, alloc *kv.Allocator, m*", mm.Name, ") (err error) {")
	g.P("if alloc == nil {")
	g.P("alloc = kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("}") // end of if block
	g.P()

	var hasIndexedField bool
	for _, f := range mt.Fields {
		opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
		hasIndexedField = proto.GetExtension(opt, rony.E_RonyIndex).(bool)
		if hasIndexedField {
			break
		}
	}

	if hasIndexedField {
		g.P("// Try to read old value")
		g.P("om := &", mm.Name, "{}")
		g.P("om, err = Read", mm.Name, "WithTxn(txn, alloc, ", mm.Table.String("m.", ",", false), ", om)")
		g.P("if err != nil && err != badger.ErrKeyNotFound {")
		g.P("return")
		g.P("}")
		g.P()
		g.P("if om != nil {")
		for _, f := range mt.Fields {
			ftName := string(f.Desc.Name())
			opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
			index := proto.GetExtension(opt, rony.E_RonyIndex).(bool)
			if index {
				g.P("// update field index by deleting old values")
				switch f.Desc.Cardinality() {
				case protoreflect.Repeated:
					g.P("for i := 0; i < len(om.", ftName, "); i++ {")
					g.P("found := false")
					g.P("for j := 0; j < len(m.", ftName, "); j++ {")
					switch f.Desc.Kind() {
					case protoreflect.MessageKind:
						panic("rony_index on MessageKind field is not valid")
					case protoreflect.BytesKind:
						g.Import("bytes")
						g.P("if bytes.Equal(om.", ftName, "[i], m.", ftName, "[j])) {")
						g.P("found = true")
						g.P("break")
						g.P("}")
					default:
						g.P("if om.", ftName, "[i] == m.", ftName, "[j] {")
						g.P("found = true")
						g.P("break")
						g.P("}")
					}
					g.P("}") // end of for (j)
					g.P("if !found {")
					g.P("err = txn.Delete(alloc.GenKey(", genDbIndexKey(mm, ftName, "om.", "[i]"), "))")
					g.P("if err != nil {")
					g.P("return")
					g.P("}")
					g.P("}")
					g.P("}") // end of for (i)
				default:
					switch f.Desc.Kind() {
					case protoreflect.MessageKind:
						panic("rony_index on MessageKind field is not valid")
					case protoreflect.BytesKind:
						g.P("if !bytes.Equal(om.", ftName, ", m.", ftName, "){")
						g.P("err = txn.Delete(alloc.GenKey(", genDbIndexKey(mm, ftName, "om.", ""), "))")
						g.P("if err != nil {")
						g.P("return")
						g.P("}")
						g.P("}")
					default:
						g.P("if om.", ftName, " != m.", ftName, "{")
						g.P("err = txn.Delete(alloc.GenKey(", genDbIndexKey(mm, ftName, "om.", ""), "))")
						g.P("if err != nil {")
						g.P("return")
						g.P("}")
						g.P("}")
					}
				}
				g.P()
			}
		}
		g.P("}")
		g.P()
	}

	g.P("// save entry")
	g.P("b := alloc.GenValue(m)")
	g.P("key := alloc.GenKey(", genDbKey(mm, mm.Table, "m."), ")")
	g.P("err = txn.Set(key, b)")
	g.P("if err != nil {")
	g.P("return")
	g.P("}")
	g.P()
	for idx := range mm.Views {
		g.P("// save entry for view", mm.Views[idx].Keys())
		g.P("err = txn.Set(alloc.GenKey(", genDbKey(mm, mm.Views[idx], "m."), "), b)")
		g.P("if err != nil {")
		g.P("return")
		g.P("}")
		g.P()
	}
	for _, f := range mt.Fields {
		ftName := string(f.Desc.Name())
		opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
		index := proto.GetExtension(opt, rony.E_RonyIndex).(bool)
		if index {
			g.P("// update field index by saving new values")
			switch f.Desc.Kind() {
			case protoreflect.MessageKind:
				// TODO:: support index on message fields
			default:
				switch f.Desc.Cardinality() {
				case protoreflect.Repeated:
					g.P("for idx := range m.", ftName, "{")
					g.P("err = txn.Set(alloc.GenKey(", genDbIndexKey(mm, ftName, "m.", "[idx]"), "), key)")
					g.P("if err != nil {")
					g.P("return")
					g.P("}")
					g.P("}") // end of for
				default:
					g.P("err = txn.Set(alloc.GenKey(", genDbIndexKey(mm, ftName, "m.", ""), "), key)")
					g.P("if err != nil {")
					g.P("return")
					g.P("}")
				}
			}
			g.P()
		}
	}
	g.P("return")
	g.P()
	g.P("}") // end of SaveWithTxn func
	g.P()

	// Save func
	g.P("func Save", mm.Name, "(m *", mm.Name, ") error {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("return kv.Update(func(txn *badger.Txn) error {")
	g.P("return Save", mm.Name, "WithTxn (txn, alloc, m)")
	g.P("})") // end of Update func
	g.P("}")  // end of Save func
	g.P()
}
func funcRead(g *protogen.GeneratedFile, mm *model.Model) {
	g.P("func Read", mm.Name, "WithTxn (txn *badger.Txn, alloc *kv.Allocator,", mm.FuncArgs("", mm.Table), ", m *", mm.Name, ") (*", mm.Name, ",error) {")
	g.P("if alloc == nil {")
	g.P("alloc = kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("}") // end of if block
	g.P()
	g.P("item, err := txn.Get(alloc.GenKey(", genDbKey(mm, mm.Table, ""), "))")
	g.P("if err != nil {")
	g.P("return nil, err")
	g.P("}")
	g.P("err = item.Value(func (val []byte) error {")
	g.P("return m.Unmarshal(val)")
	g.P("})")
	g.P("return m, err")
	g.P("}") // end of Read func
	g.P()
	g.P("func Read", mm.Name, "(", mm.FuncArgs("", mm.Table), ", m *", mm.Name, ") (*", mm.Name, ",error) {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P()
	g.P("if m == nil {")
	g.P("m = &", mm.Name, "{}")
	g.P("}")
	g.P()
	g.P("err := kv.View(func(txn *badger.Txn) (err error) {")
	g.P("m, err = Read", mm.Name, "WithTxn(txn, alloc, ", mm.Table.String("", ",", true), ", m)")
	g.P("return err")
	g.P("})") // end of View func
	g.P("return m, err")
	g.P("}") // end of Read func
	g.P()
	for _, pk := range mm.Views {
		g.P(
			"func Read", mm.Name, "By", pk.String("", "And", false),
			"WithTxn(txn *badger.Txn, alloc *kv.Allocator,", mm.FuncArgs("", pk), ", m *", mm.Name, ") ( *",
			mm.Name,
			", "+
				"error) {")
		g.P("if alloc == nil {")
		g.P("alloc = kv.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P("}") // end of if block
		g.P()
		g.P("item, err := txn.Get(alloc.GenKey(", genDbKey(mm, pk, ""), "))")
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("err = item.Value(func (val []byte) error {")
		g.P("return m.Unmarshal(val)")
		g.P("})") // end of item.Value
		g.P("return m, err")
		g.P("}") // end of Read func
		g.P()
		g.P(
			"func Read", mm.Name, "By",
			pk.String("", "And", false),
			"(", mm.FuncArgs("", pk), ", m *", mm.Name, ") ( *", mm.Name, ", error) {",
		)
		g.P("alloc := kv.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P("if m == nil {")
		g.P("m = &", mm.Name, "{}")
		g.P("}")
		g.P("err := kv.View(func(txn *badger.Txn) (err error) {")
		g.P("m, err = Read", mm.Name, "By", pk.String("", "And", false), "WithTxn (txn, alloc,", pk.String("", ",", true), ", m)")
		g.P("return err")
		g.P("})") // end of View func
		g.P("return m, err")
		g.P("}") // end of Read func
		g.P()
	}
}
func funcDelete(g *protogen.GeneratedFile, mm *model.Model) {
	g.P("func Delete", mm.Name, "WithTxn(txn *badger.Txn, alloc *kv.Allocator, ", mm.FuncArgs("", mm.Table), ") error {")
	if len(mm.Views) > 0 {
		g.P("m := &", mm.Name, "{}")
		g.P("item, err := txn.Get(alloc.GenKey(", genDbKey(mm, mm.Table, ""), "))")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("err = item.Value(func(val []byte) error {")
		g.P("return m.Unmarshal(val)")
		g.P("})")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P("err = txn.Delete(alloc.GenKey(", genDbKey(mm, mm.Table, "m."), "))")
	} else {
		g.P("err := txn.Delete(alloc.GenKey(", genDbKey(mm, mm.Table, ""), "))")
	}

	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	for _, pk := range mm.Views {
		g.P("err = txn.Delete(alloc.GenKey(", genDbKey(mm, pk, "m."), "))")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}")
		g.P()
	}
	g.P("return nil")
	g.P("}")	// end of DeleteWithTxn
	g.P()
	g.P("func Delete", mm.Name, "(", mm.FuncArgs("", mm.Table), ") error {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P()
	g.P("return kv.Update(func(txn *badger.Txn) error {")
	g.P("return Delete", mm.Name, "WithTxn(txn, alloc, ", mm.Table.String("", ",", true), ")")
	g.P("})")	// end of Update func
	g.P("}")	// end of Delete func
	g.P()
}
func funcList(g *protogen.GeneratedFile, mm *model.Model) {
	g.P("func List", mm.Name, "(")
	g.P(mm.FuncArgs("offset", mm.Table), ", lo *kv.ListOption, cond func(m *", mm.Name, ") bool, ")
	g.P(") ([]*", mm.Name, ", error) {")
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P()
	g.P("res := make([]*", mm.Name, ", 0, lo.Limit())")
	g.P("err := kv.View(func(txn *badger.Txn) error {")
	g.P("opt := badger.DefaultIteratorOptions")
	g.P("opt.Prefix = alloc.GenKey(C_", mm.Name, ",", mm.Table.Checksum(), ")")
	g.P("opt.Reverse = lo.Backward()")
	g.P("osk := alloc.GenKey(", genDbPrefixPKs(mm, mm.Table, "offset"), ")")
	g.P("iter := txn.NewIterator(opt)")
	g.P("offset := lo.Skip()")
	g.P("limit := lo.Limit()")
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
func funcListByPartitionKey(g *protogen.GeneratedFile, m *protogen.Message, mm *model.Model) {
	g.P(
		"func List", mm.Name, "By",
		mm.Table.StringPKs("", "And", false),
		"(", mm.FuncArgsPKs("", mm.Table), ",",
		mm.FuncArgsCKs("offset", mm.Table), ", lo *kv.ListOption) ([]*", mm.Name, ", error) {",
	)
	g.P("alloc := kv.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P()
	g.P("res := make([]*", mm.Name, ", 0, lo.Limit())")
	g.P("err := kv.View(func(txn *badger.Txn) error {")
	g.P("opt := badger.DefaultIteratorOptions")
	g.P("opt.Prefix = alloc.GenKey(", genDbPrefixPKs(mm, mm.Table, ""), ")")
	g.P("opt.Reverse = lo.Backward()")
	g.P("osk := alloc.GenKey(", genDbPrefixPKs(mm, mm.Table, ""), ",", mm.Table.StringCKs("offset", ",", false), ")")
	g.P("iter := txn.NewIterator(opt)")
	g.P("offset := lo.Skip()")
	g.P("limit := lo.Limit()")
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
	g.P("res = append(res, m)")
	g.P("return nil")
	g.P("})") // end of item.Value
	g.P("}")  // end of for
	g.P("iter.Close()")
	g.P("return nil")
	g.P("})") // end of View
	g.P("return res, err")
	g.P("}") // end of func List
	g.P()
}
func funcListByIndex(g *protogen.GeneratedFile, m *protogen.Message, mm *model.Model) {
	for _, f := range m.Fields {
		ftName := string(f.Desc.Name())
		opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
		index := proto.GetExtension(opt, rony.E_RonyIndex).(bool)
		if index {
			switch f.Desc.Kind() {
			case protoreflect.MessageKind:
				// TODO:: support index on message fields
			default:
				ftNameS := inflection.Singular(ftName)
				g.P("func List", mm.Name, "By", ftNameS, "(", tools.ToLowerCamel(ftNameS), " ", mm.FieldsGo[ftName], ", lo *kv.ListOption) ([]*", mm.Name, ", "+
					"error) {")
				g.P("alloc := kv.NewAllocator()")
				g.P("defer alloc.ReleaseAll()")
				g.P()
				g.P("res := make([]*", mm.Name, ", 0, lo.Limit())")
				g.P("err := kv.View(func(txn *badger.Txn) error {")
				g.P("opt := badger.DefaultIteratorOptions")
				g.P("opt.Prefix = alloc.GenKey(\"IDX\", C_", mm.Name, ",", crc32.ChecksumIEEE([]byte(ftName)), ",", tools.ToLowerCamel(ftNameS), ")")
				g.P("opt.Reverse = lo.Backward()")
				g.P("iter := txn.NewIterator(opt)")
				g.P("offset := lo.Skip()")
				g.P("limit := lo.Limit()")
				g.P("for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
				g.P("if offset--; offset >= 0 {")
				g.P("continue")
				g.P("}")
				g.P("if limit--; limit < 0 {")
				g.P("break")
				g.P("}")
				g.P("_ = iter.Item().Value(func (val []byte) error {")
				g.P("item, err := txn.Get(val)")
				g.P("if err != nil {")
				g.P("return err")
				g.P("}") // end of if
				g.P("return item.Value(func (val []byte) error {")
				g.P("m := &", mm.Name, "{}")
				g.P("err := m.Unmarshal(val)")
				g.P("if err != nil {")
				g.P("return err")
				g.P("}") // end of if
				g.P("res = append(res, m)")
				g.P("return nil")
				g.P("})") // end of item.Value
				g.P("})") // end of iter.Value func
				g.P("}")  // end of for
				g.P("iter.Close()")
				g.P("return nil")
				g.P("})") // end of View
				g.P("return res, err")
				g.P("}") // end of func List
				g.P()
			}
		}
	}
}
func funcHasField(g *protogen.GeneratedFile, m *protogen.Message, mm *model.Model) {
	for _, f := range m.Fields {
		switch f.Desc.Cardinality() {
		case protoreflect.Repeated:
			if f.Desc.Kind() == protoreflect.MessageKind {
				break
			}
			mtName := m.Desc.Name()
			g.P("func (x *", mtName, ") Has", f.Desc.Name(), "(xx ", mm.FieldsGo[f.GoName], ") bool {")
			g.P("for idx := range x.", f.Desc.Name(), "{")
			g.P("if x.", f.Desc.Name(), "[idx] == xx {")
			g.P("return true")
			g.P("}") // end of if
			g.P("}") // end of for
			g.P("return false")
			g.P("}") // end of func
			g.P()
		}
	}
}
