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

			arg := GetArg(g, m, g.m(m))
			g.g.P(g.Exec(template.Must(template.New("genCreate").Parse(genCreate)), arg))
			g.g.P(g.Exec(template.Must(template.New("genUpdate").Parse(genUpdate)), arg))
			g.g.P(g.Exec(template.Must(template.New("genSave").Parse(genSave)), arg))
			g.g.P(g.Exec(template.Must(template.New("genRead").Parse(genRead)), arg))
			g.g.P(g.Exec(template.Must(template.New("genDelete").Parse(genDelete)), arg))

			g.genOrderByConstants(m)
			g.genHasField(m)
			g.genIter(m)
			g.genList(m)
			g.genIterByPK(m)
			g.genListByPK(m)
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

const genCreate = `
func Create{{.Name}} (m *{{.Type}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return store.Update(func(txn *rony.StoreTxn) error {
		return Create{{.Name}}WithTxn (txn, alloc, m)
	})
}

func Create{{.Name}}WithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{.Name}}) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	if store.Exists(txn, alloc, {{.DBKey "m."}}) {
		return store.ErrAlreadyExists
	}
	
	// save entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, {{(.DBKey "m.")}})
	if err != nil {
		return
	}

	{{- range .Views }}

	// save view {{.Keys}}
	err = store.Set(txn, alloc, val, {{(.DBKey "m.")}})
	if err != nil {
		return err 
	}
	{{- end }}
	
	{{- if .HasIndex }}

		key := alloc.Gen({{(.DBKey "m.")}})
		{{- range .Fields }}
			{{- if .HasIndex }}
				// update field index by saving new value: {{.Name}}
				{{- if eq .Kind "repeated" }}
					for idx := range m.{{.Name}} {
						err = store.Set(txn, alloc, key, {{.DBKey "m." "[idx]"}})
						if err != nil {
							return
						}
					}
				{{- else }}
					err = store.Set(txn, alloc, key, {{.DBKey "m." ""}})
					if err != nil {
						return
					}
				{{- end }}
			{{- end }}
		{{- end }}
	{{- end }}
	
	return
}

`

const genUpdate = `
func Update{{.Name}}WithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{.Name}}) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	
	err := Delete{{.Name}}WithTxn(txn, alloc, {{.String "m." "," false}})
	if err != nil {
		return err
	}
	
	return Create{{.Name}}WithTxn(txn, alloc, m)
}

func Update{{.Name}} ({{.FuncArgs ""}}, m *{{.Name}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		return store.ErrEmptyObject
	}
	
	err := store.View(func(txn *rony.StoreTxn) (err error) {
		return Update{{.Name}}WithTxn(txn, alloc, m)
	})
	return err 
}
`

const genSave = `
func Save{{.Name}}WithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{.Name}}) (err error) {
	if store.Exists(txn, alloc, {{.DBKey "m."}}) {
		return Update{{.Name}}WithTxn(txn, alloc, m)
	} else {
		return Create{{.Name}}WithTxn(txn, alloc, m)
	}
}

func Save{{.Name}} (m *{{.Name}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *rony.StoreTxn) error {
		return Save{{.Name}}WithTxn(txn, alloc, m)
	})
}

`

const genRead = `
func Read{{.Name}}WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, {{.FuncArgs ""}}, m *{{.Name}}) (*{{.Name}}, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, {{.DBKey ""}})
	if err != nil {
		return nil, err 
	}
	return m, nil
}

func Read{{.Name}} ({{.FuncArgs ""}}, m *{{.Name}}) (*{{.Name}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{.Name}}{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = Read{{.Name}}WithTxn(txn, alloc, {{.String "" "," true}}, m)
		return err 
	})
	return m, err
}

{{- range .Views }}

func Read{{.Name}}By{{.String "" "And" false}}WithTxn (
	txn *rony.StoreTxn, alloc *tools.Allocator,
	{{.FuncArgs ""}}, m *{{.Name}},
) (*{{.Name}}, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, {{.DBKey ""}})
	if err != nil {
		return nil, err
	}
	return m, err
}

func Read{{.Name}}By{{.String "" "And" false}}({{.FuncArgs ""}}, m *{{.Name}}) (*{{.Name}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{.Name}}{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = Read{{.Name}}By{{.String "" "And" false}}WithTxn (txn, alloc, {{.String "" "," true}}, m)
		return err 
	})
	return m, err
}

{{- end }}
`

const genDelete = `
func Delete{{.Name}}WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, {{.FuncArgs ""}}) error {
	{{- if or (gt (len .Views) 0) (.HasIndex) }}
		m := &{{.Name}}{}
		err := store.Unmarshal(txn, alloc, m, {{.DBKey ""}})
		if err != nil {
			return err 
		}
		err = store.Delete(txn, alloc, {{.DBKey "m."}})
	{{- else }}
		err := store.Delete(txn, alloc, {{.DBKey ""}})
	{{- end }}
	if err != nil {
		return err 
	}
	{{- range .Fields }}
		
		{{ if .HasIndex }}
			// delete field index
			{{- if eq .Kind "repeated" }}
				for idx := range m.{{.Name}} {
					err = store.Delete(txn, alloc, {{.DBKey "m." "[idx]"}})
					if err != nil {
						return err 
					}
				}
			{{- else }}
				err = store.Delete(txn, alloc, {{.DBKey "m." ""}})
				if err != nil {
					return err
				}
			{{- end }}
		{{- end }}
	{{- end }}
	{{- range .Views }}
		err = store.Delete(txn, alloc, {{.DBKey "m."}})
		if err != nil {
			return err 
		}

	{{- end }}
	
	return nil
}

func Delete{{.Name}}({{.FuncArgs ""}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *rony.StoreTxn) error {
		return Delete{{.Name}}WithTxn(txn, alloc, {{.String "" "," true}})
	})
}

`

// genHasField generate helper function for repeated fields of the aggregate.
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

// genOrderByConstants generates constants used in Iter and List functions to identify the order.
func (g *Generator) genOrderByConstants(m *protogen.Message) {
	orderType := fmt.Sprintf("%sOrder", g.m(m).Name)
	g.g.P("type ", orderType, " string")
	for _, view := range g.m(m).Views {
		orderName := fmt.Sprintf("%sOrderBy%s", g.m(m).Name, strings.Join(view.PKs, ""))
		g.g.P("const ", orderName, " ", orderType, " =\"", strings.Join(view.PKs, ""), "\"")
	}
}
func (g *Generator) genIter(m *protogen.Message) {
	mn := g.m(m).Name
	orderType := fmt.Sprintf("%sOrder", g.m(m).Name)

	g.g.P("func Iter", inflection.Plural(mn), "(txn *rony.StoreTxn, alloc *tools.Allocator, cb func(m *", mn, ") bool, orderBy ...", orderType, ")  error {")
	g.blockAlloc()

	g.g.P("exitLoop := false")
	g.g.P("iterOpt := store.DefaultIteratorOptions")

	if len(g.m(m).Views) > 0 {
		g.g.P("if len(orderBy) == 0 {")
		g.g.P("iterOpt.Prefix = alloc.Gen(", g.m(m).Table.DBIterPrefix(), ")")
		g.g.P("} else {")
		g.g.P("switch orderBy[0] {")
		for idx, view := range g.m(m).Views {
			orderName := fmt.Sprintf("%sOrderBy%s", g.m(m).Name, strings.Join(view.PKs, ""))
			g.g.P("case ", orderName, ":")
			g.g.P("iterOpt.Prefix = alloc.Gen(", g.m(m).Views[idx].DBIterPrefix(), ")")
		}
		g.g.P("default:")
		g.g.P("iterOpt.Prefix = alloc.Gen(", g.m(m).Table.DBIterPrefix(), ")")
		g.g.P("}")
		g.g.P("}")
	} else {
		g.g.P("iterOpt.Prefix = alloc.Gen(", g.m(m).Table.DBIterPrefix(), ")")
	}

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
			"(txn *rony.StoreTxn, alloc *tools.Allocator,", g.m(m).Table.FuncArgsPKs(""), ", cb func(m *", mn, ") bool) error {",
		)
		g.blockAlloc()
		g.g.P("exitLoop := false")
		g.g.P("opt := store.DefaultIteratorOptions")
		g.g.P("opt.Prefix = alloc.Gen(", g.m(m).Table.DBKeyPrefix(""), ")")
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
			"(txn *rony.StoreTxn, alloc *tools.Allocator,", g.m(m).Views[idx].FuncArgsPKs(""), ", cb func(m *", mn, ") bool) error {",
		)
		g.g.P("if alloc == nil {")

		g.g.P("alloc = tools.NewAllocator()")
		g.g.P("defer alloc.ReleaseAll()")
		g.g.P("}")
		g.g.P()
		g.g.P("exitLoop := false")
		g.g.P("opt := store.DefaultIteratorOptions")
		g.g.P("opt.Prefix = alloc.Gen(", g.m(m).Views[idx].DBKeyPrefix(""), ")")
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
	orderType := fmt.Sprintf("%sOrder", g.m(m).Name)
	g.g.P("func List", mn, "(")
	g.g.P(g.m(m).Table.FuncArgs("offset"), ", lo *store.ListOption, cond func(m *", mn, ") bool, orderBy ...", orderType, ",")
	g.g.P(") ([]*", mn, ", error) {")
	g.g.P("alloc := tools.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P()
	g.g.P("res := make([]*", mn, ", 0, lo.Limit())")
	g.g.P("err := store.View(func(txn *rony.StoreTxn) error {")
	g.g.P("opt := store.DefaultIteratorOptions")
	if len(g.m(m).Views) > 0 {
		g.g.P("if len(orderBy) == 0 {")
		g.g.P("opt.Prefix = alloc.Gen(", g.m(m).Table.DBIterPrefix(), ")")
		g.g.P("} else {")
		g.g.P("switch orderBy[0] {")
		for idx, view := range g.m(m).Views {
			orderName := fmt.Sprintf("%sOrderBy%s", g.m(m).Name, strings.Join(view.PKs, ""))
			g.g.P("case ", orderName, ":")
			g.g.P("opt.Prefix = alloc.Gen(", g.m(m).Views[idx].DBIterPrefix(), ")")
		}
		g.g.P("default:")
		g.g.P("opt.Prefix = alloc.Gen(", g.m(m).Table.DBIterPrefix(), ")")
		g.g.P("}")
		g.g.P("}")
	} else {
		g.g.P("opt.Prefix = alloc.Gen(", g.m(m).Table.DBIterPrefix(), ")")
	}

	g.g.P("opt.Reverse = lo.Backward()")
	g.g.P("osk := alloc.Gen(", g.m(m).Table.DBKeyPrefix("offset"), ")")
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
			g.m(m).Table.FuncArgsPKs(""), ",",
			g.m(m).Table.FuncArgsCKs("offset"), ", lo *store.ListOption, cond func(m *", mn, ") bool,",
		)
		g.g.P(
			") ([]*", mn, ", error) {",
		)
		g.g.P("alloc := tools.NewAllocator()")
		g.g.P("defer alloc.ReleaseAll()")
		g.g.P()
		g.g.P("res := make([]*", mn, ", 0, lo.Limit())")
		g.g.P("err := store.View(func(txn *rony.StoreTxn) error {")
		g.g.P("opt := store.DefaultIteratorOptions")
		g.g.P("opt.Prefix = alloc.Gen(", g.m(m).Table.DBKeyPrefix(""), ")")
		g.g.P("opt.Reverse = lo.Backward()")
		g.g.P("osk := alloc.Gen(", g.m(m).Table.DBKeyPrefix(""), ",", g.m(m).Table.StringCKs("offset", ",", false), ")")
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
			g.m(m).Views[idx].FuncArgsPKs(""), ",",
			g.m(m).Views[idx].FuncArgsCKs("offset"), ", lo *store.ListOption, cond func(m *", mn, ") bool,",
		)
		g.g.P(
			") ([]*", mn, ", error) {",
		)
		g.g.P("alloc := tools.NewAllocator()")
		g.g.P("defer alloc.ReleaseAll()")
		g.g.P()
		g.g.P("res := make([]*", mn, ", 0, lo.Limit())")
		g.g.P("err := store.View(func(txn *rony.StoreTxn) error {")
		g.g.P("opt := store.DefaultIteratorOptions")
		g.g.P("opt.Prefix = alloc.Gen(", g.m(m).Views[idx].DBKeyPrefix(""), ")")
		g.g.P("opt.Reverse = lo.Backward()")
		g.g.P("osk := alloc.Gen(", g.m(m).Views[idx].DBKeyPrefix(""), ",", g.m(m).Views[idx].StringCKs("offset", ",", false), ")")
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

func (g *Generator) blockAlloc() {
	g.g.P("if alloc == nil {")
	g.g.P("alloc = tools.NewAllocator()")
	g.g.P("defer alloc.ReleaseAll()")
	g.g.P("}") // end of if block
	g.g.P()
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
