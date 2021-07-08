package aggregate

import (
	"fmt"
	"github.com/jinzhu/inflection"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/codegen"
	parse "github.com/ronaksoft/rony/internal/parser"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"hash/crc64"
	"strings"
	"text/template"
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
	savedModels map[string]*Aggregate
	f           *protogen.File
	g           *protogen.GeneratedFile
}

func New(f *protogen.File, g *protogen.GeneratedFile) *Generator {
	return &Generator{
		savedModels: map[string]*Aggregate{},
		f:           f,
		g:           g,
	}
}

func (g *Generator) Generate() {
	for _, m := range g.f.Messages {
		g.createModel(m)
		if g.m(m) != nil {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})

			g.genListByIndex(m)
		}
	}
}

func (g *Generator) Exec(t *template.Template, v interface{}) string {
	sb := &strings.Builder{}
	err := t.Execute(sb, v)
	if err != nil {
		panic(err)
	}
	return sb.String()
}

func (g *Generator) createModel(m *protogen.Message) {
	var (
		isAggregate = false
		agg         = Aggregate{
			Name:      string(m.Desc.Name()),
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
	if views := proto.GetExtension(opt, rony.E_RonyAggregateView).([]string); len(views) > 0 {
		for _, view := range views {
			aggregateDesc.WriteString(fmt.Sprintf("{{@view %s}}\n", view))
		}
	}

	// Parse the generated description
	t, err := parse.Parse(string(m.Desc.Name()), aggregateDesc.String())
	if err != nil {
		panic(err)
	}

	// Generate Go and CQL kinds of the fields
	for _, f := range m.Fields {
		agg.FieldsCql[f.GoName] = codegen.CqlKind(f.Desc)
		agg.FieldsGo[f.GoName] = codegen.GoKind(g.f, g.g, f.Desc)
		opt, _ := f.Desc.Options().(*descriptorpb.FieldOptions)
		if proto.GetExtension(opt, rony.E_RonyIndex).(bool) {
			agg.HasIndex = true
		}
	}

	fields := make(map[string]struct{})
	for _, n := range t.Root.Nodes {
		switch n.Type() {
		case parse.NodeModel:
			agg.Type = n.(*parse.ModelNode).Text
			isAggregate = true
		case parse.NodeTable:
			pk := Key{
				aggregate: string(m.Desc.Name()),
			}
			nn := n.(*parse.TableNode)
			for _, k := range nn.PartitionKeys {
				fields[k] = struct{}{}
				pk.PKs = append(pk.PKs, k)
				pk.PKGoTypes = append(pk.PKGoTypes, agg.FieldsGo[k])
			}
			for _, k := range nn.ClusteringKeys {
				kWithoutSign := strings.TrimLeft(k, "-")
				fields[kWithoutSign] = struct{}{}
				pk.Orders = append(pk.Orders, k)
				pk.CKs = append(pk.CKs, kWithoutSign)
				pk.CKGoTypes = append(pk.CKGoTypes, agg.FieldsGo[kWithoutSign])
			}
			agg.Table = pk
		case parse.NodeView:
			pk := Key{
				aggregate: string(m.Desc.Name()),
			}
			nn := n.(*parse.ViewNode)
			for _, k := range nn.PartitionKeys {
				fields[k] = struct{}{}
				pk.PKs = append(pk.PKs, k)
				pk.PKGoTypes = append(pk.PKGoTypes, agg.FieldsGo[k])
			}

			for _, k := range nn.ClusteringKeys {
				kWithoutSign := strings.TrimLeft(k, "-")
				fields[kWithoutSign] = struct{}{}
				pk.Orders = append(pk.Orders, k)
				pk.CKs = append(pk.CKs, kWithoutSign)
				pk.CKGoTypes = append(pk.CKGoTypes, agg.FieldsGo[kWithoutSign])
			}
			agg.Views = append(agg.Views, pk)
		}
	}

	if !isAggregate {
		return
	}

	g.savedModels[agg.Name] = &agg

}
func (g *Generator) model(m *protogen.Message) *Aggregate {
	return g.savedModels[string(m.Desc.Name())]
}
func (g *Generator) m(m *protogen.Message) *Aggregate {
	return g.model(m)
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
				g.g.P("alloc := tools.NewAllocator()")
				g.g.P("defer alloc.ReleaseAll()")
				g.g.P()
				g.g.P("res := make([]*", mn, ", 0, lo.Limit())")
				g.g.P("err := store.View(func(txn *rony.StoreTxn) error {")
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


func panicF(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
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
