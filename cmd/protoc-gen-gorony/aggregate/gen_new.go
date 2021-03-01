package aggregate

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	parse "github.com/ronaksoft/rony/internal/parser"
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
	g.g.P("func Read", mn, "WithTxn (txn *store.Txn, alloc *store.Allocator,", g.m(m).FuncArgs("", g.m(m).Table), ", m *", mn, ") (*", mn, ",error) {")
	g.blockAlloc()
	g.g.P("err := store.Unmarshal(txn, alloc, m,", tableKey(g.m(m), ""), ")")
	g.g.P("if err != nil {")
	g.g.P("return nil, err")
	g.g.P("}")
	g.g.P("return m, err")
	g.g.P("}") // end of Read func
	g.g.P()
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
	for _, pk := range g.m(m).Views {
		g.g.P(
			"func Read", mn, "By", pk.String("", "And", false), "WithTxn",
			"(txn *store.Txn, alloc *store.Allocator,", g.m(m).FuncArgs("", pk), ", m *", mn, ")",
			"( *", mn, ", error) {")
		g.blockAlloc()
		g.g.P("err := store.Unmarshal(txn, alloc, m,", tableKey(g.m(m), ""), " )")
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
	g.g.P("om := &", mn, "{}")
	g.g.P("err := store.Unmarshal(txn, alloc, om,", tableKey(g.m(m), "m."), ")")
	g.g.P("if err != nil {")
	g.g.P("return err")
	g.g.P("}")
	g.g.P()
	g.g.P("err = Delete", mn, "WithTxn(txn, alloc, ", g.m(m).Table.String("om.", ",", false), ")")
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
	if len(g.m(m).Views) > 0 {
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
				g.g.P("func (x *", mtName, ") Has", f.Desc.Name(), "(xx ", g.m(m).FieldsGo[ftName], ") bool {")
				g.g.P("for idx := range x.", f.Desc.Name(), "{")
				g.g.P("if bytes.Equal(x.", f.Desc.Name(), "[idx], xx) {")
				g.g.P("return true")
				g.g.P("}") // end of if
				g.g.P("}") // end of for
				g.g.P("return false")
				g.g.P("}") // end of func
				g.g.P()
			case protoreflect.EnumKind:
				mtName := m.Desc.Name()
				g.g.P("func (x *", mtName, ") Has", f.Desc.Name(), "(xx ", f.Enum.Desc.Name(), ") bool {")
				g.g.P("for idx := range x.", f.Desc.Name(), "{")
				g.g.P("if x.", f.Desc.Name(), "[idx] == xx {")
				g.g.P("return true")
				g.g.P("}") // end of if
				g.g.P("}") // end of for
				g.g.P("return false")
				g.g.P("}") // end of func
				g.g.P()
			default:
				mtName := m.Desc.Name()
				g.g.P("func (x *", mtName, ") Has", f.Desc.Name(), "(xx ", g.m(m).FieldsGo[ftName], ") bool {")
				g.g.P("for idx := range x.", f.Desc.Name(), "{")
				g.g.P("if x.", f.Desc.Name(), "[idx] == xx {")
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
	// Save func
	g.g.P("func Save", g.model(m).Name, "(m *", g.model(m).Name, ") error {")
	g.g.P("alloc := store.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P("return store.Update(func(txn *store.Txn) error {")
	g.g.P("return Save", g.model(m).Name, "WithTxn (txn, alloc, m)")
	g.g.P("})") // end of Update func
	g.g.P("}")  // end of Save func
	g.g.P()

	// SaveWithTxn func
	g.g.P("func Save", g.model(m).Name, "WithTxn (txn *store.Txn, alloc *store.Allocator, m*", g.model(m).Name, ") (err error) {")
	g.g.P("return nil")
	g.g.P("}")
}
func (g *Generator) genIter(m *protogen.Message) {}
func (g *Generator) genIterByPK(m *protogen.Message) {}
func (g *Generator) genList(m *protogen.Message) {}
func (g *Generator) genListByPK(m *protogen.Message) {}

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

func viewKey(mm *Aggregate, prefixKey string, idx int) string {
	return fmt.Sprintf("'M', C_%s, %d, %s",
		mm.Name,
		mm.Views[idx].Checksum(),
		mm.Views[idx].String(prefixKey, ",", prefixKey == ""),
	)
}

func indexKey(mm *Aggregate, fieldName string, prefix string, postfix string) string {
	lower := prefix == ""
	return fmt.Sprintf("'I', C_%s, %d, %s%s%s, %s",
		mm.Name, crc64.Checksum([]byte(fieldName), crcTab),
		prefix, fieldName, postfix,
		mm.Table.String(prefix, ",", lower),
	)
}

var (
	crcTab = crc64.MakeTable(crc64.ISO)
)

func genDbPrefixPKs(mm *Aggregate, key Key, keyPrefix string) string {
	lowerCamel := keyPrefix == ""
	return fmt.Sprintf("'M', C_%s, %d, %s",
		mm.Name,
		key.Checksum(),
		key.StringPKs(keyPrefix, ",", lowerCamel),
	)
}
func genDbPrefixCKs(mm *Aggregate, key Key, keyPrefix string) string {
	lowerCamel := keyPrefix == ""
	return fmt.Sprintf("'M', C_%s, %d, %s",
		mm.Name,
		key.Checksum(),
		key.StringCKs(keyPrefix, ",", lowerCamel),
	)
}
