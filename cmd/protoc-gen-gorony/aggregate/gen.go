package aggregate

import (
	"fmt"
	"github.com/jinzhu/inflection"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	parse "github.com/ronaksoft/rony/internal/parser"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"hash/crc64"
	"strings"
)

/*
   Creation Time: 2021 - Mar - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Generator struct {
	savedModels   map[string]*Aggregate
	visitedFields map[string]struct{}
	f             *protogen.File
	g             *protogen.GeneratedFile
}

func New(f *protogen.File, g *protogen.GeneratedFile) *Generator {
	return &Generator{
		savedModels:   map[string]*Aggregate{},
		visitedFields: map[string]struct{}{},
		f:             f,
		g:             g,
	}
}

func (g *Generator) Generate() {
	for _, m := range g.f.Messages {
		g.createModel(m)
		if g.m(m) != nil {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.genCreate(m)
			g.genRead(m)
			g.genUpdate(m)
			g.genDelete(m)
			g.genHasField(m)
			g.genSave(m)
			g.genIter(m)
			g.genList(m)
			g.genIterByPK(m)
			g.genListByPK(m)
			g.genListByIndex(m)
		}
	}
}

func (g *Generator) genCreate(m *protogen.Message) {
	// Create func
	g.g.P("func Create", g.model(m).Name, "(m *", g.model(m).Name, ") error {")
	g.g.P("alloc := store.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P("return store.Update(func(txn *store.Txn) error {")
	g.g.P("return Create", g.model(m).Name, "WithTxn (txn, alloc, m)")
	g.g.P("})") // end of Update func
	g.g.P("}")  // end of Save func
	g.g.P()

	// CreateWithTxn func
	g.g.P("func Create", g.model(m).Name, "WithTxn (txn *store.Txn, alloc *store.Allocator, m*", g.model(m).Name, ") (err error) {")
	g.blockAlloc()
	g.g.P("if store.Exists(txn, alloc, ", tableKey(g.m(m), "m."), ") {")
	g.g.P("return store.ErrAlreadyExists")
	g.g.P("}")

	g.g.P("// save entry")
	g.g.P("val := alloc.Marshal(m)")
	g.g.P("err = store.Set(txn, alloc, val, ", tableKey(g.m(m), "m."), ")")
	g.g.P("if err != nil {")
	g.g.P("return")
	g.g.P("}")
	g.g.P()

	g.createViews(m)
	g.createIndices(m)

	g.g.P("return")
	g.g.P()
	g.g.P("}") // end of CreateWithTxn func
	g.g.P()
	g.g.P()

}
func (g *Generator) createViews(m *protogen.Message) {
	if len(g.m(m).Views) > 0 {
		g.g.P("// save views")
		for idx := range g.m(m).Views {
			g.g.P("// save entry for view: ", g.m(m).Views[idx].Keys())
			g.g.P("err = store.Set(txn, alloc, val, ", viewKey(g.m(m), "m.", idx), ")")
			g.g.P("if err != nil {")
			g.g.P("return")
			g.g.P("}")
			g.g.P()
		}
	}
}
func (g *Generator) createIndices(m *protogen.Message) {
	if g.m(m).HasIndex {
		g.g.P("// save indices")
		g.g.P("key := alloc.Gen(", tableKey(g.m(m), "m."), ")")
		for _, f := range m.Fields {
			ftName := string(f.Desc.Name())
			opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
			index := proto.GetExtension(opt, rony.E_RonyIndex).(bool)
			if index {
				g.g.P("// update field index by saving new values")
				switch f.Desc.Kind() {
				case protoreflect.MessageKind:
					panicF("does not support index in Message: %s", m.Desc.Name())
				default:
					switch f.Desc.Cardinality() {
					case protoreflect.Repeated:
						g.g.P("for idx := range m.", ftName, "{")
						g.g.P("err = store.Set(txn, alloc, key,", indexKey(g.m(m), ftName, "m.", "[idx]"), " )")
						g.g.P("if err != nil {")
						g.g.P("return")
						g.g.P("}")
						g.g.P("}") // end of for
					default:
						g.g.P("err = store.Set(txn, alloc, key,", indexKey(g.m(m), ftName, "m.", ""), " )")
						g.g.P("if err != nil {")
						g.g.P("return")
						g.g.P("}")
					}
				}
				g.g.P()
			}
		}
	}
}
func (g *Generator) genRead(m *protogen.Message) {
	mn := g.m(m).Name

	// ReadWithTxn Func
	g.g.P("func Read", mn, "WithTxn (txn *store.Txn, alloc *store.Allocator,", g.m(m).FuncArgs("", g.m(m).Table), ", m *", mn, ") (*", mn, ",error) {")
	g.blockAlloc()
	g.g.P("err := store.Unmarshal(txn, alloc, m,", tableKey(g.m(m), ""), ")")
	g.g.P("if err != nil {")
	g.g.P("return nil, err")
	g.g.P("}")
	g.g.P("return m, err")
	g.g.P("}") // end of Read func
	g.g.P()

	// Read Func
	g.g.P("func Read", mn, "(", g.m(m).FuncArgs("", g.m(m).Table), ", m *", mn, ") (*", mn, ",error) {")
	g.g.P("alloc := store.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P()
	g.g.P("if m == nil {")
	g.g.P("m = &", mn, "{}")
	g.g.P("}")
	g.g.P()
	g.g.P("err := store.View(func(txn *store.Txn) (err error) {")
	g.g.P("m, err = Read", mn, "WithTxn(txn, alloc, ", g.m(m).Table.String("", ",", true), ", m)")
	g.g.P("return err")
	g.g.P("})") // end of View func
	g.g.P("return m, err")
	g.g.P("}") // end of Read func
	g.g.P()

	g.readViews(m)
}
func (g *Generator) readViews(m *protogen.Message) {
	mn := g.m(m).Name
	for idx, pk := range g.m(m).Views {
		g.g.P(
			"func Read", mn, "By", pk.String("", "And", false), "WithTxn",
			"(txn *store.Txn, alloc *store.Allocator,", g.m(m).FuncArgs("", pk), ", m *", mn, ")",
			"( *", mn, ", error) {")
		g.blockAlloc()
		g.g.P("err := store.Unmarshal(txn, alloc, m,", viewKey(g.m(m), "", idx), " )")
		g.g.P("if err != nil {")
		g.g.P("return nil, err")
		g.g.P("}")
		g.g.P("return m, err")
		g.g.P("}") // end of Read func
		g.g.P()
		g.g.P(
			"func Read", mn, "By", pk.String("", "And", false),
			"(", g.m(m).FuncArgs("", pk), ", m *", mn, ")",
			"( *", mn, ", error) {",
		)
		g.g.P("alloc := store.NewAllocator()")
		g.g.P("defer alloc.ReleaseAll()")
		g.g.P("if m == nil {")
		g.g.P("m = &", mn, "{}")
		g.g.P("}")
		g.g.P("err := store.View(func(txn *store.Txn) (err error) {")
		g.g.P("m, err = Read", mn, "By", pk.String("", "And", false), "WithTxn (txn, alloc,", pk.String("", ",", true), ", m)")
		g.g.P("return err")
		g.g.P("})") // end of View func
		g.g.P("return m, err")
		g.g.P("}") // end of Read func
		g.g.P()
	}
}
func (g *Generator) genUpdate(m *protogen.Message) {
	mn := g.m(m).Name
	g.g.P("func Update", mn, "WithTxn (txn *store.Txn, alloc *store.Allocator, m *", mn, ") error {")
	g.blockAlloc()
	g.g.P("err := Delete", mn, "WithTxn(txn, alloc, ", g.m(m).Table.String("m.", ",", false), ")")
	g.g.P("if err != nil {")
	g.g.P("return err")
	g.g.P("}")
	g.g.P()
	g.g.P("return Create", mn, "WithTxn(txn, alloc, m)")
	g.g.P("}") // end of UpdateWithTxn func
	g.g.P()
	g.g.P("func Update", mn, "(", g.m(m).FuncArgs("", g.m(m).Table), ", m *", mn, ") error {")
	g.g.P("alloc := store.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P()
	g.g.P("if m == nil {")
	g.g.P("return store.ErrEmptyObject")
	g.g.P("}")
	g.g.P()
	g.g.P("err := store.View(func(txn *store.Txn) (err error) {")
	g.g.P("return Update", mn, "WithTxn(txn, alloc, m)")
	g.g.P("})") // end of View func
	g.g.P("return err")
	g.g.P("}") // end of Read func
	g.g.P()
}
func (g *Generator) genDelete(m *protogen.Message) {
	mn := g.m(m).Name
	g.g.P("func Delete", mn, "WithTxn(txn *store.Txn, alloc *store.Allocator, ", g.m(m).FuncArgs("", g.m(m).Table), ") error {")

	if len(g.m(m).Views) > 0 || g.m(m).HasIndex {
		g.g.P("m := &", mn, "{}")
		g.g.P("err := store.Unmarshal(txn, alloc, m, ", tableKey(g.m(m), ""), ")")
		g.g.P("if err != nil {")
		g.g.P("return err")
		g.g.P("}")
		g.g.P("err = store.Delete(txn, alloc,", tableKey(g.m(m), "m."), ")")
	} else {
		g.g.P("err := store.Delete(txn, alloc,", tableKey(g.m(m), ""), ")")
	}
	g.g.P("if err != nil {")
	g.g.P("return err")
	g.g.P("}")
	g.g.P()

	for _, f := range m.Fields {
		ftName := string(f.Desc.Name())
		opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
		index := proto.GetExtension(opt, rony.E_RonyIndex).(bool)
		if index {
			g.g.P("// delete field index ")
			switch f.Desc.Kind() {
			case protoreflect.MessageKind:
				panicF("does not support index in Message: %s", m.Desc.Name())
			default:
				switch f.Desc.Cardinality() {
				case protoreflect.Repeated:
					g.g.P("for idx := range m.", ftName, "{")
					g.g.P("err = store.Delete(txn, alloc, ", indexKey(g.m(m), ftName, "m.", "[idx]"), " )")
					g.g.P("if err != nil {")
					g.g.P("return err")
					g.g.P("}")
					g.g.P("}") // end of for
				default:
					g.g.P("err = store.Delete(txn, alloc, ", indexKey(g.m(m), ftName, "m.", ""), " )")
					g.g.P("if err != nil {")
					g.g.P("return err")
					g.g.P("}")
				}
			}
			g.g.P()
		}
	}

	for idx := range g.m(m).Views {
		g.g.P("err = store.Delete(txn, alloc,", viewKey(g.m(m), "m.", idx), ")")
		g.g.P("if err != nil {")
		g.g.P("return err")
		g.g.P("}")
		g.g.P()
	}

	g.g.P("return nil")
	g.g.P("}") // end of DeleteWithTxn
	g.g.P()

	g.g.P("func Delete", mn, "(", g.m(m).FuncArgs("", g.m(m).Table), ") error {")
	g.g.P("alloc := store.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P()
	g.g.P("return store.Update(func(txn *store.Txn) error {")
	g.g.P("return Delete", mn, "WithTxn(txn, alloc, ", g.m(m).Table.String("", ",", true), ")")
	g.g.P("})") // end of Update func
	g.g.P("}")  // end of Delete func
	g.g.P()
}

func (g *Generator) genHasField(m *protogen.Message) {
	for _, f := range m.Fields {
		ftName := string(f.Desc.Name())
		switch f.Desc.Cardinality() {
		case protoreflect.Repeated:
			switch f.Desc.Kind() {
			case protoreflect.MessageKind, protoreflect.GroupKind:
			case protoreflect.BytesKind:
				mtName := m.Desc.Name()
				g.g.P("func (x *", mtName, ") Has", inflection.Singular(ftName), "(xx ", g.m(m).FieldsGo[ftName], ") bool {")
				g.g.P("for idx := range x.", ftName, "{")
				g.g.P("if bytes.Equal(x.", ftName, "[idx], xx) {")
				g.g.P("return true")
				g.g.P("}") // end of if
				g.g.P("}") // end of for
				g.g.P("return false")
				g.g.P("}") // end of func
				g.g.P()
			case protoreflect.EnumKind:
				mtName := m.Desc.Name()
				g.g.P("func (x *", mtName, ") Has", inflection.Singular(ftName), "(xx ", f.Enum.Desc.Name(), ") bool {")
				g.g.P("for idx := range x.", ftName, "{")
				g.g.P("if x.", ftName, "[idx] == xx {")
				g.g.P("return true")
				g.g.P("}") // end of if
				g.g.P("}") // end of for
				g.g.P("return false")
				g.g.P("}") // end of func
				g.g.P()
			default:
				mtName := m.Desc.Name()
				g.g.P("func (x *", mtName, ") Has", inflection.Singular(ftName), "(xx ", g.m(m).FieldsGo[ftName], ") bool {")
				g.g.P("for idx := range x.", ftName, "{")
				g.g.P("if x.", ftName, "[idx] == xx {")
				g.g.P("return true")
				g.g.P("}") // end of if
				g.g.P("}") // end of for
				g.g.P("return false")
				g.g.P("}") // end of func
				g.g.P()

			}

		}
	}
}
func (g *Generator) genSave(m *protogen.Message) {
	mn := g.m(m).Name
	// SaveWithTxn func
	g.g.P("func Save", mn, "WithTxn (txn *store.Txn, alloc *store.Allocator, m*", mn, ") (err error) {")
	g.g.P("if store.Exists(txn, alloc, ", tableKey(g.m(m), "m."), ") {")
	g.g.P("return Update", mn, "WithTxn(txn, alloc, m)")
	g.g.P("} else {")
	g.g.P("return Create", mn, "WithTxn(txn, alloc, m)")
	g.g.P("}")
	g.g.P("}")
	g.g.P()

	// Save func
	g.g.P("func Save", g.model(m).Name, "(m *", g.model(m).Name, ") error {")
	g.g.P("alloc := store.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P("return store.Update(func(txn *store.Txn) error {")
	g.g.P("return Save", g.model(m).Name, "WithTxn (txn, alloc, m)")
	g.g.P("})") // end of Update func
	g.g.P("}")  // end of Save func
	g.g.P()
}
func (g *Generator) genIter(m *protogen.Message) {
	mn := g.m(m).Name
	g.g.P("func Iter", inflection.Plural(mn), "(txn *store.Txn, alloc *store.Allocator, cb func(m *", mn, ") bool, )  error {")
	g.blockAlloc()

	g.g.P("exitLoop := false")
	g.g.P("iterOpt := store.DefaultIteratorOptions")
	g.g.P("iterOpt.Prefix = alloc.Gen('M', C_", mn, ",", g.m(m).Table.Checksum(), ")")
	g.g.P("iter := txn.NewIterator(iterOpt)")
	g.g.P("for iter.Rewind(); iter.ValidForPrefix(iterOpt.Prefix); iter.Next() {")
	g.g.P("_ = iter.Item().Value(func (val []byte) error {")
	g.g.P("m := &", mn, "{}")
	g.g.P("err := m.Unmarshal(val)")
	g.g.P("if err != nil {")
	g.g.P("return err")
	g.g.P("}") // end of if
	g.g.P("if !cb(m) {")
	g.g.P("exitLoop = true")
	g.g.P("}") // end of if callback
	g.g.P("return nil")
	g.g.P("})") // end of iter.Value func
	g.g.P("if exitLoop {")
	g.g.P("break")
	g.g.P("}")
	g.g.P("}") // end of for
	g.g.P("iter.Close()")
	g.g.P("return nil")
	g.g.P("}") // end of func List
	g.g.P()
}
func (g *Generator) genIterByPK(m *protogen.Message) {
	mn := g.m(m).Name
	if len(g.m(m).Table.CKs) > 0 {
		g.g.P(
			"func Iter", mn, "By",
			g.m(m).Table.StringPKs("", "And", false),
			"(txn *store.Txn, alloc *store.Allocator,", g.m(m).FuncArgsPKs("", g.m(m).Table), ", cb func(m *", mn, ") bool) error {",
		)
		g.blockAlloc()
		g.g.P("exitLoop := false")
		g.g.P("opt := store.DefaultIteratorOptions")
		g.g.P("opt.Prefix = alloc.Gen(", tablePrefix(g.m(m), ""), ")")
		g.g.P("iter := txn.NewIterator(opt)")
		g.g.P("for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
		g.g.P("_ = iter.Item().Value(func (val []byte) error {")
		g.g.P("m := &", mn, "{}")
		g.g.P("err := m.Unmarshal(val)")
		g.g.P("if err != nil {")
		g.g.P("return err")
		g.g.P("}") // end of if
		g.g.P("if !cb(m) {")
		g.g.P("exitLoop = true")
		g.g.P("}")
		g.g.P("return nil")
		g.g.P("})") // end of item.Value
		g.g.P("if exitLoop {")
		g.g.P("break")
		g.g.P("}")
		g.g.P("}") // end of for
		g.g.P("iter.Close()")
		g.g.P("return nil")
		g.g.P("}") // end of func Iter
		g.g.P()
	}
	for idx := range g.m(m).Views {
		if len(g.m(m).Views[idx].CKs) == 0 {
			continue
		}
		g.g.P(
			"func Iter", mn, "By",
			g.m(m).Views[idx].StringPKs("", "And", false),
			"(txn *store.Txn, alloc *store.Allocator,", g.m(m).FuncArgsPKs("", g.m(m).Views[idx]), ", cb func(m *", mn, ") bool) error {",
		)
		g.g.P("if alloc == nil {")
		g.g.P("alloc = store.NewAllocator()")
		g.g.P("defer alloc.ReleaseAll()")
		g.g.P("}")
		g.g.P()
		g.g.P("exitLoop := false")
		g.g.P("opt := store.DefaultIteratorOptions")
		g.g.P("opt.Prefix = alloc.Gen(", viewPrefix(g.m(m), "", idx), ")")
		g.g.P("iter := txn.NewIterator(opt)")
		g.g.P("for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
		g.g.P("_ = iter.Item().Value(func (val []byte) error {")
		g.g.P("m := &", mn, "{}")
		g.g.P("err := m.Unmarshal(val)")
		g.g.P("if err != nil {")
		g.g.P("return err")
		g.g.P("}") // end of if
		g.g.P("if !cb(m) {")
		g.g.P("exitLoop = true")
		g.g.P("}")
		g.g.P("return nil")
		g.g.P("})") // end of item.Value
		g.g.P("if exitLoop {")
		g.g.P("break")
		g.g.P("}")
		g.g.P("}") // end of for
		g.g.P("iter.Close()")
		g.g.P("return nil")
		g.g.P("}") // end of func List
		g.g.P()
	}
}
func (g *Generator) genList(m *protogen.Message) {
	mn := g.m(m).Name
	g.g.P("func List", mn, "(")
	g.g.P(g.m(m).FuncArgs("offset", g.m(m).Table), ", lo *store.ListOption, cond func(m *", mn, ") bool, ")
	g.g.P(") ([]*", mn, ", error) {")
	g.g.P("alloc := store.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P()
	g.g.P("res := make([]*", mn, ", 0, lo.Limit())")
	g.g.P("err := store.View(func(txn *store.Txn) error {")
	g.g.P("opt := store.DefaultIteratorOptions")
	g.g.P("opt.Prefix = alloc.Gen('M', C_", mn, ",", g.m(m).Table.Checksum(), ")")
	g.g.P("opt.Reverse = lo.Backward()")
	g.g.P("osk := alloc.Gen(", tablePrefix(g.m(m), "offset"), ")")
	g.g.P("iter := txn.NewIterator(opt)")
	g.g.P("offset := lo.Skip()")
	g.g.P("limit := lo.Limit()")
	g.g.P("for iter.Seek(osk); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
	g.g.P("if offset--; offset >= 0 {")
	g.g.P("continue")
	g.g.P("}")
	g.g.P("if limit--; limit < 0 {")
	g.g.P("break")
	g.g.P("}")
	g.g.P("_ = iter.Item().Value(func (val []byte) error {")
	g.g.P("m := &", mn, "{}")
	g.g.P("err := m.Unmarshal(val)")
	g.g.P("if err != nil {")
	g.g.P("return err")
	g.g.P("}") // end of if
	g.g.P("if cond == nil || cond(m) {")
	g.g.P("res = append(res, m)")
	g.g.P("} else {") // end of if cond
	g.g.P("limit++")
	g.g.P("}")
	g.g.P("return nil")
	g.g.P("})") // end of iter.Value func
	g.g.P("}")  // end of for
	g.g.P("iter.Close()")
	g.g.P("return nil")
	g.g.P("})") // end of View
	g.g.P("return res, err")
	g.g.P("}") // end of func List
	g.g.P()
}
func (g *Generator) genListByPK(m *protogen.Message) {
	mn := g.m(m).Name
	if len(g.m(m).Table.CKs) > 0 {
		g.g.P(
			"func List", mn, "By", g.m(m).Table.StringPKs("", "And", false), "(",
		)
		g.g.P(
			g.m(m).FuncArgsPKs("", g.m(m).Table), ",",
			g.m(m).FuncArgsCKs("offset", g.m(m).Table), ", lo *store.ListOption, cond func(m *", mn, ") bool,",
		)
		g.g.P(
			") ([]*", mn, ", error) {",
		)
		g.g.P("alloc := store.NewAllocator()")
		g.g.P("defer alloc.ReleaseAll()")
		g.g.P()
		g.g.P("res := make([]*", mn, ", 0, lo.Limit())")
		g.g.P("err := store.View(func(txn *store.Txn) error {")
		g.g.P("opt := store.DefaultIteratorOptions")
		g.g.P("opt.Prefix = alloc.Gen(", tablePrefix(g.m(m), ""), ")")
		g.g.P("opt.Reverse = lo.Backward()")
		g.g.P("osk := alloc.Gen(", tablePrefix(g.m(m), ""), ",", g.m(m).Table.StringCKs("offset", ",", false), ")")
		g.g.P("iter := txn.NewIterator(opt)")
		g.g.P("offset := lo.Skip()")
		g.g.P("limit := lo.Limit()")
		g.g.P("for iter.Seek(osk); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
		g.g.P("if offset--; offset >= 0 {")
		g.g.P("continue")
		g.g.P("}")
		g.g.P("if limit--; limit < 0 {")
		g.g.P("break")
		g.g.P("}")
		g.g.P("_ = iter.Item().Value(func (val []byte) error {")
		g.g.P("m := &", mn, "{}")
		g.g.P("err := m.Unmarshal(val)")
		g.g.P("if err != nil {")
		g.g.P("return err")
		g.g.P("}") // end of if
		g.g.P("if cond == nil || cond(m) {")
		g.g.P("res = append(res, m)")
		g.g.P("} else {") // end of if cond
		g.g.P("limit++")
		g.g.P("}")
		g.g.P("return nil")
		g.g.P("})") // end of item.Value
		g.g.P("}")  // end of for
		g.g.P("iter.Close()")
		g.g.P("return nil")
		g.g.P("})") // end of View
		g.g.P("return res, err")
		g.g.P("}") // end of func List
		g.g.P()
	}
	for idx := range g.m(m).Views {
		if len(g.m(m).Views[idx].CKs) == 0 {
			continue
		}
		g.g.P(
			"func List", mn, "By", g.m(m).Views[idx].StringPKs("", "And", false), "(",
		)
		g.g.P(
			g.m(m).FuncArgsPKs("", g.m(m).Views[idx]), ",",
			g.m(m).FuncArgsCKs("offset", g.m(m).Views[idx]), ", lo *store.ListOption, cond func(m *", mn, ") bool,",
		)
		g.g.P(
			") ([]*", mn, ", error) {",
		)
		g.g.P("alloc := store.NewAllocator()")
		g.g.P("defer alloc.ReleaseAll()")
		g.g.P()
		g.g.P("res := make([]*", mn, ", 0, lo.Limit())")
		g.g.P("err := store.View(func(txn *store.Txn) error {")
		g.g.P("opt := store.DefaultIteratorOptions")
		g.g.P("opt.Prefix = alloc.Gen(", viewPrefix(g.m(m), "", idx), ")")
		g.g.P("opt.Reverse = lo.Backward()")
		g.g.P("osk := alloc.Gen(", viewPrefix(g.m(m), "", idx), ",", g.m(m).Views[idx].StringCKs("offset", ",", false), ")")
		g.g.P("iter := txn.NewIterator(opt)")
		g.g.P("offset := lo.Skip()")
		g.g.P("limit := lo.Limit()")
		g.g.P("for iter.Seek(osk); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
		g.g.P("if offset--; offset >= 0 {")
		g.g.P("continue")
		g.g.P("}")
		g.g.P("if limit--; limit < 0 {")
		g.g.P("break")
		g.g.P("}")
		g.g.P("_ = iter.Item().Value(func (val []byte) error {")
		g.g.P("m := &", g.m(m).Name, "{}")
		g.g.P("err := m.Unmarshal(val)")
		g.g.P("if err != nil {")
		g.g.P("return err")
		g.g.P("}") // end of if
		g.g.P("if cond == nil || cond(m) {")
		g.g.P("res = append(res, m)")
		g.g.P("} else {") // end of if cond
		g.g.P("limit++")
		g.g.P("}")
		g.g.P("return nil")
		g.g.P("})") // end of item.Value
		g.g.P("}")  // end of for
		g.g.P("iter.Close()")
		g.g.P("return nil")
		g.g.P("})") // end of View
		g.g.P("return res, err")
		g.g.P("}") // end of func List
		g.g.P()
	}
}
func (g *Generator) genListByIndex(m *protogen.Message) {
	mn := g.m(m).Name
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
				g.g.P("func List", g.m(m).Name, "By", ftNameS, "(")
				g.g.P(
					tools.ToLowerCamel(ftNameS), " ", g.m(m).FieldsGo[ftName], ", lo *store.ListOption, cond func(m *", mn, ") bool,",
				)
				g.g.P(") ([]*", mn, ", error) {")
				g.g.P("alloc := store.NewAllocator()")
				g.g.P("defer alloc.ReleaseAll()")
				g.g.P()
				g.g.P("res := make([]*", mn, ", 0, lo.Limit())")
				g.g.P("err := store.View(func(txn *store.Txn) error {")
				g.g.P("opt := store.DefaultIteratorOptions")
				g.g.P("opt.Prefix = alloc.Gen(", indexPrefix(g.m(m), ftName, tools.ToLowerCamel(ftNameS)), ")")
				g.g.P("opt.Reverse = lo.Backward()")
				g.g.P("iter := txn.NewIterator(opt)")
				g.g.P("offset := lo.Skip()")
				g.g.P("limit := lo.Limit()")
				g.g.P("for iter.Rewind(); iter.ValidForPrefix(opt.Prefix); iter.Next() {")
				g.g.P("if offset--; offset >= 0 {")
				g.g.P("continue")
				g.g.P("}")
				g.g.P("if limit--; limit < 0 {")
				g.g.P("break")
				g.g.P("}")
				g.g.P("_ = iter.Item().Value(func (val []byte) error {")
				g.g.P("item, err := txn.Get(val)")
				g.g.P("if err != nil {")
				g.g.P("return err")
				g.g.P("}") // end of if
				g.g.P("return item.Value(func (val []byte) error {")
				g.g.P("m := &", mn, "{}")
				g.g.P("err := m.Unmarshal(val)")
				g.g.P("if err != nil {")
				g.g.P("return err")
				g.g.P("}") // end of if
				g.g.P("if cond == nil || cond(m) {")
				g.g.P("res = append(res, m)")
				g.g.P("} else {") // end of if cond
				g.g.P("limit++")
				g.g.P("}")
				g.g.P("return nil")
				g.g.P("})") // end of item.Value
				g.g.P("})") // end of iter.Value func
				g.g.P("}")  // end of for
				g.g.P("iter.Close()")
				g.g.P("return nil")
				g.g.P("})") // end of View
				g.g.P("return res, err")
				g.g.P("}") // end of func List
				g.g.P()
			}
		}
	}
}

func (g *Generator) createModel(m *protogen.Message) {
	var (
		isAggregate = false
		mm          = Aggregate{
			FieldsCql: make(map[string]string),
			FieldsGo:  make(map[string]string),
		}
	)

	// Generate the aggregate description from proto options
	aggregateDesc := strings.Builder{}
	opt, _ := m.Desc.Options().(*descriptorpb.MessageOptions)

	if entity := proto.GetExtension(opt, rony.E_RonyAggregate).(bool); entity {
		aggrType := proto.GetExtension(opt, rony.E_RonyAggregateType).(string)
		if aggrType == "" {
			panic("define rony_aggregate_type")
		}
		aggregateDesc.WriteString(fmt.Sprintf("{{@model %s}}\n", aggrType))
	}
	if tab := proto.GetExtension(opt, rony.E_RonyAggregateTable).(string); tab != "" {
		aggregateDesc.WriteString(fmt.Sprintf("{{@tab %s}}\n", tab))
	}
	if view := proto.GetExtension(opt, rony.E_RonyAggregateView).(string); view != "" {
		aggregateDesc.WriteString(fmt.Sprintf("{{@view %s}}\n", view))
	}

	// Parse the generated description
	t, err := parse.Parse(string(m.Desc.Name()), aggregateDesc.String())
	if err != nil {
		panic(err)
	}
	fields := make(map[string]struct{})
	for _, n := range t.Root.Nodes {
		switch n.Type() {
		case parse.NodeModel:
			nn := n.(*parse.ModelNode)
			mm.Type = nn.Text
			isAggregate = true
		case parse.NodeTable:
			pk := Key{}
			nn := n.(*parse.TableNode)
			for _, k := range nn.PartitionKeys {
				fields[k] = struct{}{}
				pk.PKs = append(pk.PKs, k)
			}
			for _, k := range nn.ClusteringKeys {
				kWithoutSign := strings.TrimLeft(k, "-")
				fields[kWithoutSign] = struct{}{}
				pk.Orders = append(pk.Orders, k)
				pk.CKs = append(pk.CKs, kWithoutSign)
			}
			mm.Table = pk
		case parse.NodeView:
			pk := Key{}
			nn := n.(*parse.ViewNode)
			sb := strings.Builder{}
			for _, k := range nn.PartitionKeys {
				fields[k] = struct{}{}
				pk.PKs = append(pk.PKs, k)
				sb.WriteString(k)
			}
			mm.ViewParams = append(mm.ViewParams, sb.String())
			for _, k := range nn.ClusteringKeys {
				kWithoutSign := strings.TrimLeft(k, "-")
				fields[kWithoutSign] = struct{}{}
				pk.Orders = append(pk.Orders, k)
				pk.CKs = append(pk.CKs, kWithoutSign)
			}
			mm.Views = append(mm.Views, pk)
		}
	}
	for f := range fields {
		g.visitedFields[f] = struct{}{}
		mm.FieldNames = append(mm.FieldNames, f)
	}
	if isAggregate {
		for _, f := range m.Fields {
			mm.FieldsCql[f.GoName] = z.CqlKind(f.Desc)
			mm.FieldsGo[f.GoName] = z.GoKind(g.f, g.g, f.Desc)
		}
		mm.Name = string(m.Desc.Name())

		// Check if Aggregator has indexed field
		for _, f := range m.Fields {
			opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
			mm.HasIndex = proto.GetExtension(opt, rony.E_RonyIndex).(bool)
			if mm.HasIndex {
				break
			}
		}

		g.savedModels[mm.Name] = &mm
	}

}
func (g *Generator) model(m *protogen.Message) *Aggregate {
	return g.savedModels[string(m.Desc.Name())]
}
func (g *Generator) m(m *protogen.Message) *Aggregate {
	return g.model(m)
}

func (g *Generator) blockAlloc() {
	g.g.P("if alloc == nil {")
	g.g.P("alloc = store.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P("}") // end of if block
	g.g.P()
}

func panicF(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}

func tableKey(mm *Aggregate, prefixKey string) string {
	return fmt.Sprintf("'M', C_%s, %d, %s",
		mm.Name,
		mm.Table.Checksum(),
		mm.Table.String(prefixKey, ",", prefixKey == ""),
	)
}

func tablePrefix(mm *Aggregate, prefixKey string) string {
	return fmt.Sprintf("'M', C_%s, %d, %s",
		mm.Name,
		mm.Table.Checksum(),
		mm.Table.StringPKs(prefixKey, ",", prefixKey == ""),
	)
}

func viewKey(mm *Aggregate, prefixKey string, idx int) string {
	return fmt.Sprintf("'M', C_%s, %d, %s",
		mm.Name,
		mm.Views[idx].Checksum(),
		mm.Views[idx].String(prefixKey, ",", prefixKey == ""),
	)
}

func viewPrefix(mm *Aggregate, prefixKey string, idx int) string {
	return fmt.Sprintf("'M', C_%s, %d, %s",
		mm.Name,
		mm.Views[idx].Checksum(),
		mm.Views[idx].StringPKs(prefixKey, ",", prefixKey == ""),
	)
}

func indexKey(mm *Aggregate, fieldName string, prefix string, postfix string) string {
	lower := prefix == ""
	return fmt.Sprintf("'I', C_%s, uint64(%d), %s%s%s, %s",
		mm.Name, crc64.Checksum([]byte(fieldName), crcTab),
		prefix, fieldName, postfix,
		mm.Table.String(prefix, ",", lower),
	)
}

func indexPrefix(mm *Aggregate, fieldName string, fieldVarName string) string {
	return fmt.Sprintf("'I', C_%s, uint64(%d), %s",
		mm.Name, crc64.Checksum([]byte(fieldName), crcTab), fieldVarName,
	)
}

var (
	crcTab = crc64.MakeTable(crc64.ISO)
)
