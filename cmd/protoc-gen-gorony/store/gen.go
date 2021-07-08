package store

import (
	"fmt"
	"github.com/ronaksoft/rony/internal/codegen"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"hash/crc32"
	"hash/crc64"
	"strings"
	"text/template"
)

/*
   Creation Time: 2021 - Jul - 08
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	crcTab = crc64.MakeTable(crc64.ISO)
)

type Generator struct {
	f *protogen.File
	g *protogen.GeneratedFile
}

func New(f *protogen.File, g *protogen.GeneratedFile) *Generator {
	return &Generator{
		f: f,
		g: g,
	}
}

func (g *Generator) Generate() {
	singletonFuncs := map[string]interface{}{
		"DBKey": func(arg codegen.MessageArg) string {
			return fmt.Sprintf("'S', C_%s",
				arg.Name,
			)
		},
	}
	aggregateFuncs := map[string]interface{}{
		"HasIndex": func(m codegen.MessageArg) bool {
			for _, f := range m.Fields {
				if f.HasIndex {
					return true
				}
			}
			return false
		},
		"FuncArgs": func(m codegen.ModelKey, prefix string) string {
			nc := codegen.None
			if prefix == "" {
				nc = codegen.LowerCamelCase
			}
			return m.NameTypes(codegen.PropFilterALL, prefix, nc, codegen.LangGo)
		},
		"DBKey": func(m codegen.ModelKey, prefix string) string {
			nc := codegen.None
			if prefix == "" {
				nc = codegen.LowerCamelCase
			}
			return fmt.Sprintf("'M', C_%s, %d, %s",
				m.Name(),
				crc32.ChecksumIEEE(tools.StrToByte(m.Names(codegen.PropFilterALL, "", ",", codegen.None))),
				m.Names(codegen.PropFilterALL, prefix, ",", nc),
			)
		},
		"IndexDBKey": func(m codegen.MessageArg, f codegen.FieldArg, prefix, postfix string) string {
			nc := codegen.None
			if prefix == "" {
				nc = codegen.LowerCamelCase
			}
			return fmt.Sprintf("'I', C_%s, uint64(%d), %s%s%s, %s",
				m.Name, crc64.Checksum([]byte(f.Name), crcTab),
				prefix, f.Name, postfix,
				m.Table.Names(codegen.PropFilterALL, prefix, ",", nc),
			)
		},
		"String": func(m codegen.ModelKey, prefix, sep string, lcc bool) string {
			nc := codegen.None
			if lcc {
				nc = codegen.LowerCamelCase
			}
			return m.Names(codegen.PropFilterALL, prefix, sep, nc)
		},
	}
	for _, m := range g.f.Messages {
		arg := codegen.GetMessageArg(g.f, g.g, m)
		if arg.IsSingleton {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
			g.g.P(g.Exec(template.Must(template.New("genSingletonSave").Funcs(singletonFuncs).Parse(genSingletonSave)), arg))
			g.g.P(g.Exec(template.Must(template.New("genSingletonRead").Funcs(singletonFuncs).Parse(genSingletonRead)), arg))
		} else if arg.IsAggregate {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})
			g.g.P(g.Exec(template.Must(template.New("genAggregateCreate").Funcs(aggregateFuncs).Parse(genAggregateCreate)), arg))
			g.g.P(g.Exec(template.Must(template.New("genAggregateUpdate").Funcs(aggregateFuncs).Parse(genAggregateUpdate)), arg))
			g.g.P(g.Exec(template.Must(template.New("genAggregateSave").Funcs(aggregateFuncs).Parse(genAggregateSave)), arg))
			g.g.P(g.Exec(template.Must(template.New("genAggregateRead").Funcs(aggregateFuncs).Parse(genAggregateRead)), arg))
			g.g.P(g.Exec(template.Must(template.New("genAggregateDelete").Funcs(aggregateFuncs).Parse(genAggregateDelete)), arg))
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

const genSingletonSave = `
func Save{{.Name}}WithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{.Name}}) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	
	err = store.Marshal(txn, alloc, m, {{ DBKey . }})
	if err != nil {
		return 
	}
	return nil
}

func Save{{.Name}} (m *{{.Name}}) (err error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	
	return store.Update(func(txn *rony.StoreTxn) error {
		return Save{{.Name}}WithTxn(txn, alloc, m)
	})
}
`

const genSingletonRead = `
func Read{{.Name}}WithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{.Name}}) (*{{.Name}}, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	
	err := store.Unmarshal(txn, alloc, m, {{ DBKey . }})
	if err != nil {
		return nil, err 
	}
	return m, err
}

func Read{{.Name}} (m *{{.Name}}) (*{{.Name}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{.Name}}{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = Read{{.Name}}WithTxn(txn, alloc, m)
		return
	})
	return m, err 
}
`

const genAggregateCreate = `
{{$model := .}}
func Create{{.Name}} (m *{{.Name}}) error {
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
	if store.Exists(txn, alloc, {{DBKey .Table "m."}}) {
		return store.ErrAlreadyExists
	}
	
	// save entry
	val := alloc.Marshal(m)
	err = store.Set(txn, alloc, val, {{(DBKey .Table "m.")}})
	if err != nil {
		return
	}

	{{- range .Views }}
	// save view {{.Keys}}
	err = store.Set(txn, alloc, val, {{(DBKey . "m.")}})
	if err != nil {
		return err 
	}
	{{- end }}
	
	
	{{- if HasIndex . }}

		key := alloc.Gen({{(DBKey .Table "m.")}})
		{{- range .Fields }}
			{{- if .HasIndex }}
				// update field index by saving new value: {{.Name}}
				{{- if eq .Kind "repeated" }}
					for idx := range m.{{.Name}} {
						err = store.Set(txn, alloc, key, {{IndexDBKey $model . "m." "[idx]"}})
						if err != nil {
							return
						}
					}
				{{- else }}
					err = store.Set(txn, alloc, key, {{IndexDBKey $model . "m." ""}})
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

const genAggregateUpdate = `
func Update{{.Name}}WithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{.Name}}) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	
	err := Delete{{.Name}}WithTxn(txn, alloc, {{String .Table "m." "," false}})
	if err != nil {
		return err
	}
	
	return Create{{.Name}}WithTxn(txn, alloc, m)
}

func Update{{.Name}} ({{FuncArgs .Table ""}}, m *{{.Name}}) error {
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

const genAggregateSave = `
{{$model := .}}
func Save{{.Name}}WithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{.Name}}) (err error) {
	if store.Exists(txn, alloc, {{DBKey .Table "m."}}) {
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

const genAggregateRead = `
{{$model := .}}
func Read{{.Name}}WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, {{FuncArgs .Table ""}}, m *{{.Name}}) (*{{.Name}}, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, {{DBKey .Table ""}})
	if err != nil {
		return nil, err 
	}
	return m, nil
}

func Read{{.Name}} ({{FuncArgs .Table ""}}, m *{{.Name}}) (*{{.Name}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{.Name}}{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = Read{{.Name}}WithTxn(txn, alloc, {{String .Table "" "," true}}, m)
		return err 
	})
	return m, err
}

{{- range .Views }}

func Read{{.Name}}By{{String . "" "And" false}}WithTxn (
	txn *rony.StoreTxn, alloc *tools.Allocator,
	{{FuncArgs . ""}}, m *{{.Name}},
) (*{{.Name}}, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}

	err := store.Unmarshal(txn, alloc, m, {{DBKey . ""}})
	if err != nil {
		return nil, err
	}
	return m, err
}

func Read{{.Name}}By{{String . "" "And" false}}({{FuncArgs . ""}}, m *{{.Name}}) (*{{.Name}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{.Name}}{}
	}

	err := store.View(func(txn *rony.StoreTxn) (err error) {
		m, err = Read{{.Name}}By{{String . "" "And" false}}WithTxn (txn, alloc, {{String . "" "," true}}, m)
		return err 
	})
	return m, err
}

{{- end }}
`

const genAggregateDelete = `
{{$model := .}}
func Delete{{.Name}}WithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, {{FuncArgs .Table ""}}) error {
	{{- if or (gt (len .Views) 0) (HasIndex .) }}
		m := &{{.Name}}{}
		err := store.Unmarshal(txn, alloc, m, {{DBKey .Table ""}})
		if err != nil {
			return err 
		}
		err = store.Delete(txn, alloc, {{DBKey .Table "m."}})
	{{- else }}
		err := store.Delete(txn, alloc, {{DBKey .Table ""}})
	{{- end }}
	if err != nil {
		return err 
	}
	{{- range .Fields }}
		
		{{ if .HasIndex }}
			// delete field index
			{{- if eq .Cardinality "repeated" }}
				for idx := range m.{{.Name}} {
					err = store.Delete(txn, alloc, {{IndexDBKey $model . "m." "[idx]"}})
					if err != nil {
						return err 
					}
				}
			{{- else }}
				err = store.Delete(txn, alloc, {{IndexDBKey $model . "m." ""}})
				if err != nil {
					return err
				}
			{{- end }}
		{{- end }}
	{{- end }}
	{{- range .Views }}
		err = store.Delete(txn, alloc, {{DBKey . "m."}})
		if err != nil {
			return err 
		}

	{{- end }}
	
	return nil
}

func Delete{{.Name}}({{FuncArgs .Table ""}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return store.Update(func(txn *rony.StoreTxn) error {
		return Delete{{.Name}}WithTxn(txn, alloc, {{String .Table "" "," true}})
	})
}

`
