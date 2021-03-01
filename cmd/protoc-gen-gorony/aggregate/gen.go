package aggregate

import (
	"github.com/jinzhu/inflection"
	"github.com/ronaksoft/rony"
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

func funcSave(g *protogen.GeneratedFile, mt *protogen.Message, mm *Aggregate) {

	// if hasIndexedField {
	// 	g.P("// Try to read old value")
	// 	g.P("om := &", mm.Name, "{}")
	// 	g.P("err := store.Unmarshal(txn alloc, om,", mm.Table.String("m.", ",", false), ")")
	// 	g.P("if err != nil && err != store.ErrKeyNotFound {")
	// 	g.P("return")
	// 	g.P("}")
	// 	g.P()
	// 	g.P("if om != nil {")
	// 	for _, f := range mt.Fields {
	// 		ftName := string(f.Desc.Name())
	// 		opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
	// 		index := proto.GetExtension(opt, rony.E_RonyIndex).(bool)
	// 		if index {
	// 			g.P("// update field index by deleting old values")
	// 			switch f.Desc.Cardinality() {
	// 			case protoreflect.Repeated:
	// 				g.P("for i := 0; i < len(om.", ftName, "); i++ {")
	// 				g.P("found := false")
	// 				g.P("for j := 0; j < len(m.", ftName, "); j++ {")
	// 				switch f.Desc.Kind() {
	// 				case protoreflect.MessageKind:
	// 					panic("rony_index on MessageKind field is not valid")
	// 				case protoreflect.BytesKind:
	// 					g.Import("bytes")
	// 					g.P("if bytes.Equal(om.", ftName, "[i], m.", ftName, "[j])) {")
	// 					g.P("found = true")
	// 					g.P("break")
	// 					g.P("}")
	// 				default:
	// 					g.P("if om.", ftName, "[i] == m.", ftName, "[j] {")
	// 					g.P("found = true")
	// 					g.P("break")
	// 					g.P("}")
	// 				}
	// 				g.P("}") // end of for (j)
	// 				g.P("if !found {")
	// 				g.P("err = store.Delete(txn, alloc, ", genDbIndexKey(mm, ftName, "om.", "[i]"), ")")
	// 				g.P("if err != nil {")
	// 				g.P("return")
	// 				g.P("}")
	// 				g.P("}")
	// 				g.P("}") // end of for (i)
	// 			default:
	// 				switch f.Desc.Kind() {
	// 				case protoreflect.MessageKind:
	// 					panic("rony_index on MessageKind field is not valid")
	// 				case protoreflect.BytesKind:
	// 					g.P("if !bytes.Equal(om.", ftName, ", m.", ftName, "){")
	// 					g.P("err = store.Delete(txn, alloc,", genDbIndexKey(mm, ftName, "om.", ""), ")")
	// 					g.P("if err != nil {")
	// 					g.P("return")
	// 					g.P("}")
	// 					g.P("}")
	// 				default:
	// 					g.P("if om.", ftName, " != m.", ftName, "{")
	// 					g.P("err = store.Delete(txn, alloc,", genDbIndexKey(mm, ftName, "om.", ""), ")")
	// 					g.P("if err != nil {")
	// 					g.P("return")
	// 					g.P("}")
	// 					g.P("}")
	// 				}
	// 			}
	// 			g.P()
	// 		}
	// 	}
	// 	g.P("}")
	// 	g.P()
	// }

}

func funcIter(g *protogen.GeneratedFile, mm *Aggregate) {
	g.P("func Iter", inflection.Plural(mm.Name), "(txn *store.Txn, alloc *store.Allocator, cb func(m *", mm.Name, ") bool, )  error {")
	g.P("if alloc == nil {")
	g.P("alloc = store.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P("}")
	g.P()
	g.P("exitLoop := false")
	g.P("iterOpt := store.DefaultIteratorOptions")
	g.P("iterOpt.Prefix = alloc.GenKey('M', C_", mm.Name, ",", mm.Table.Checksum(), ")")
	g.P("iter := txn.NewIterator(iterOpt)")
	g.P("for iter.Rewind(); iter.ValidForPrefix(iterOpt.Prefix); iter.Next() {")
	g.P("_ = iter.Item().Value(func (val []byte) error {")
	g.P("m := &", mm.Name, "{}")
	g.P("err := m.Unmarshal(val)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}") // end of if
	g.P("if !cb(m) {")
	g.P("exitLoop = true")
	g.P("}") // end of if callback
	g.P("return nil")
	g.P("})") // end of iter.Value func
	g.P("if exitLoop {")
	g.P("break")
	g.P("}")
	g.P("}") // end of for
	g.P("iter.Close()")
	g.P("return nil")
	g.P("}") // end of func List
	g.P()
}
func funcIterByPartitionKey(g *protogen.GeneratedFile, mm *Aggregate) {
	if len(mm.Table.CKs) > 0 {
		g.P(
			"func Iter", mm.Name, "By",
			mm.Table.StringPKs("", "And", false),
			"(txn *store.Txn, alloc *store.Allocator,", mm.FuncArgsPKs("", mm.Table), ", cb func(m *", mm.Name, ") bool) error {",
		)
		g.P("if alloc == nil {")
		g.P("alloc = store.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P("}")
		g.P()
		g.P("exitLoop := false")
		g.P("opt := store.DefaultIteratorOptions")
		g.P("opt.Prefix = alloc.GenKey(", genDbPrefixPKs(mm, mm.Table, ""), ")")
		g.P("iter := txn.NewIterator(opt)")
		g.P("for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
		g.P("_ = iter.Item().Value(func (val []byte) error {")
		g.P("m := &", mm.Name, "{}")
		g.P("err := m.Unmarshal(val)")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}") // end of if
		g.P("if !cb(m) {")
		g.P("exitLoop = true")
		g.P("}")
		g.P("return nil")
		g.P("})") // end of item.Value
		g.P("if exitLoop {")
		g.P("break")
		g.P("}")
		g.P("}") // end of for
		g.P("iter.Close()")
		g.P("return nil")
		g.P("}") // end of func Iter
		g.P()
	}

	for idx := range mm.Views {
		if len(mm.Views[idx].CKs) == 0 {
			continue
		}
		g.P(
			"func Iter", mm.Name, "By",
			mm.Views[idx].StringPKs("", "And", false),
			"(txn *store.Txn, alloc *store.Allocator,", mm.FuncArgsPKs("", mm.Views[idx]), ", cb func(m *", mm.Name, ") bool) error {",
		)
		g.P("if alloc == nil {")
		g.P("alloc = store.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P("}")
		g.P()
		g.P("exitLoop := false")
		g.P("opt := store.DefaultIteratorOptions")
		g.P("opt.Prefix = alloc.GenKey(", genDbPrefixPKs(mm, mm.Views[idx], ""), ")")
		g.P("iter := txn.NewIterator(opt)")
		g.P("for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
		g.P("_ = iter.Item().Value(func (val []byte) error {")
		g.P("m := &", mm.Name, "{}")
		g.P("err := m.Unmarshal(val)")
		g.P("if err != nil {")
		g.P("return err")
		g.P("}") // end of if
		g.P("if !cb(m) {")
		g.P("exitLoop = true")
		g.P("}")
		g.P("return nil")
		g.P("})") // end of item.Value
		g.P("if exitLoop {")
		g.P("break")
		g.P("}")
		g.P("}") // end of for
		g.P("iter.Close()")
		g.P("return nil")
		g.P("}") // end of func List
		g.P()
	}
}
func funcList(g *protogen.GeneratedFile, mm *Aggregate) {
	g.P("func List", mm.Name, "(")
	g.P(mm.FuncArgs("offset", mm.Table), ", lo *store.ListOption, cond func(m *", mm.Name, ") bool, ")
	g.P(") ([]*", mm.Name, ", error) {")
	g.P("alloc := store.NewAllocator()")
	g.P("defer alloc.ReleaseAll()")
	g.P()
	g.P("res := make([]*", mm.Name, ", 0, lo.Limit())")
	g.P("err := store.View(func(txn *store.Txn) error {")
	g.P("opt := store.DefaultIteratorOptions")
	g.P("opt.Prefix = alloc.GenKey('M', C_", mm.Name, ",", mm.Table.Checksum(), ")")
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
	g.P("if cond == nil || cond(m) {")
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
func funcListByPartitionKey(g *protogen.GeneratedFile, mm *Aggregate) {
	if len(mm.Table.CKs) > 0 {
		g.P(
			"func List", mm.Name, "By",
			mm.Table.StringPKs("", "And", false),
			"(", mm.FuncArgsPKs("", mm.Table), ",",
			mm.FuncArgsCKs("offset", mm.Table), ", lo *store.ListOption) ([]*", mm.Name, ", error) {",
		)
		g.P("alloc := store.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P()
		g.P("res := make([]*", mm.Name, ", 0, lo.Limit())")
		g.P("err := store.View(func(txn *store.Txn) error {")
		g.P("opt := store.DefaultIteratorOptions")
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

	for idx := range mm.Views {
		if len(mm.Views[idx].CKs) == 0 {
			continue
		}
		g.P(
			"func List", mm.Name, "By",
			mm.Views[idx].StringPKs("", "And", false),
			"(", mm.FuncArgsPKs("", mm.Views[idx]), ",",
			mm.FuncArgsCKs("offset", mm.Views[idx]), ", lo *store.ListOption) ([]*", mm.Name, ", error) {",
		)
		g.P("alloc := store.NewAllocator()")
		g.P("defer alloc.ReleaseAll()")
		g.P()
		g.P("res := make([]*", mm.Name, ", 0, lo.Limit())")
		g.P("err := store.View(func(txn *store.Txn) error {")
		g.P("opt := store.DefaultIteratorOptions")
		g.P("opt.Prefix = alloc.GenKey(", genDbPrefixPKs(mm, mm.Views[idx], ""), ")")
		g.P("opt.Reverse = lo.Backward()")
		g.P("osk := alloc.GenKey(", genDbPrefixPKs(mm, mm.Views[idx], ""), ",", mm.Views[idx].StringCKs("offset", ",", false), ")")
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
}
func funcListByIndex(g *protogen.GeneratedFile, m *protogen.Message, mm *Aggregate) {
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
				g.P("func List", mm.Name, "By", ftNameS, "(", tools.ToLowerCamel(ftNameS), " ", mm.FieldsGo[ftName], ", lo *store.ListOption) ([]*", mm.Name, ", "+
					"error) {")
				g.P("alloc := store.NewAllocator()")
				g.P("defer alloc.ReleaseAll()")
				g.P()
				g.P("res := make([]*", mm.Name, ", 0, lo.Limit())")
				g.P("err := store.View(func(txn *store.Txn) error {")
				g.P("opt := store.DefaultIteratorOptions")
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
