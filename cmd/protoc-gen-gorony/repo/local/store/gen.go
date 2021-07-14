package store

import (
	"fmt"
	"github.com/jinzhu/inflection"
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
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "bytes"})
	g.g.P("var _ = bytes.MinRead")

	singletonFuncs := map[string]interface{}{
		"DBKey": func(arg codegen.MessageArg) string {
			return fmt.Sprintf("'S', C_%s",
				arg.Name,
			)
		},
	}
	aggregateFuncs := map[string]interface{}{
		"Singular": func(x string) string {
			return inflection.Singular(x)
		},
		"Plural": func(x string) string {
			return inflection.Plural(x)
		},
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
		"FuncArgsPKs": func(m codegen.ModelKey, prefix string) string {
			nc := codegen.None
			if prefix == "" {
				nc = codegen.LowerCamelCase
			}
			return m.NameTypes(codegen.PropFilterPKs, prefix, nc, codegen.LangGo)
		},
		"FuncArgsCKs": func(m codegen.ModelKey, prefix string) string {
			nc := codegen.None
			if prefix == "" {
				nc = codegen.LowerCamelCase
			}
			return m.NameTypes(codegen.PropFilterCKs, prefix, nc, codegen.LangGo)
		},
		"DBKey": func(m codegen.ModelKey, prefix string) string {
			nc := codegen.None
			if prefix == "" {
				nc = codegen.LowerCamelCase
			}
			return fmt.Sprintf("'M', C_%s, %d, %s",
				m.Name(),
				crc32.ChecksumIEEE(tools.StrToByte(m.Names(codegen.PropFilterALL, "", "", ",", codegen.None))),
				m.Names(codegen.PropFilterALL, prefix, "", ",", nc),
			)
		},
		"DBKeyPKs": func(m codegen.ModelKey, prefix string) string {
			nc := codegen.None
			if prefix == "" {
				nc = codegen.LowerCamelCase
			}
			return fmt.Sprintf("'M', C_%s, %d, %s",
				m.Name(),
				crc32.ChecksumIEEE(tools.StrToByte(m.Names(codegen.PropFilterALL, "", "", ",", codegen.None))),
				m.Names(codegen.PropFilterPKs, prefix, "", ",", nc),
			)
		},
		"DBPrefix": func(m codegen.ModelKey) string {
			return fmt.Sprintf("'M', C_%s, %d",
				m.Name(),
				crc32.ChecksumIEEE(tools.StrToByte(m.Names(codegen.PropFilterALL, "", "", ",", codegen.None))),
			)
		},
		"IndexDBKey": func(m codegen.MessageArg, f codegen.FieldArg, prefix, postfix string) string {
			nc := codegen.None
			if prefix == "" {
				nc = codegen.LowerCamelCase
			}
			return fmt.Sprintf("'I', C_%s, uint64(%d), %s%s%s, %s",
				m.Name, crc64.Checksum([]byte(f.Name), codegen.CrcTab),
				prefix, f.Name, postfix,
				m.Table.Names(codegen.PropFilterALL, prefix, "", ",", nc),
			)
		},
		"IndexDBPrefix": func(m codegen.MessageArg, f codegen.FieldArg, prefix, postfix string) string {
			return fmt.Sprintf("'I', C_%s, uint64(%d), %s%s%s",
				m.Name, crc64.Checksum([]byte(f.Name), codegen.CrcTab),
				prefix, inflection.Singular(f.Name), postfix,
			)
		},
		"String": func(m codegen.ModelKey, prefix, sep string, lcc bool) string {
			nc := codegen.None
			if lcc {
				nc = codegen.LowerCamelCase
			}
			return m.Names(codegen.PropFilterALL, prefix, "", sep, nc)
		},
		"StringPKs": func(m codegen.ModelKey, prefix, sep string, lcc bool) string {
			nc := codegen.None
			if lcc {
				nc = codegen.LowerCamelCase
			}
			return m.Names(codegen.PropFilterPKs, prefix, "", sep, nc)
		},
		"StringCKs": func(m codegen.ModelKey, prefix, sep string, lcc bool) string {
			nc := codegen.None
			if lcc {
				nc = codegen.LowerCamelCase
			}
			return m.Names(codegen.PropFilterCKs, prefix, "", sep, nc)
		},
		"OrderTypes": func(m codegen.MessageArg) map[string]int {
			var (
				uniqueOrders = make(map[string]int)
			)
			for idx, v := range m.Views {
				uniqueOrders[v.Names(codegen.PropFilterPKs, "", "", "", codegen.None)] = idx
			}

			return uniqueOrders
		},
	}
	for _, m := range g.f.Messages {
		arg := codegen.GetMessageArg(g.f, g.g, m)
		if arg.IsSingleton {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
			g.g.P(g.Exec(template.Must(template.New("genSingleton").Funcs(singletonFuncs).Parse(genSingleton)), arg))
		} else if arg.IsAggregate {
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/store"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony"})
			g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/tools"})
			g.g.P(g.Exec(template.Must(template.New("genHelpers").Funcs(aggregateFuncs).Parse(genHelpers)), arg))
			g.g.P(g.Exec(template.Must(template.New("genLocalRepo").Funcs(aggregateFuncs).Parse(genLocalRepo)), arg))
			g.g.P(g.Exec(template.Must(template.New("genCreate").Funcs(aggregateFuncs).Parse(genCreate)), arg))
			g.g.P(g.Exec(template.Must(template.New("genUpdate").Funcs(aggregateFuncs).Parse(genUpdate)), arg))
			g.g.P(g.Exec(template.Must(template.New("genSave").Funcs(aggregateFuncs).Parse(genSave)), arg))
			g.g.P(g.Exec(template.Must(template.New("genRead").Funcs(aggregateFuncs).Parse(genRead)), arg))
			g.g.P(g.Exec(template.Must(template.New("genDelete").Funcs(aggregateFuncs).Parse(genDelete)), arg))

		}
	}
}

func (g *Generator) Exec(t *template.Template, v interface{}) string {
	sb := &strings.Builder{}
	if err := t.Execute(sb, v); err != nil {
		panic(err)
	}
	return sb.String()
}

const genSingleton = `
{{$model := .}}
{{$repoName := print .Name "LocalSingleton"}}
{{$modelName := .Name}}
type {{$repoName}} struct {
    s rony.Store
}

func New{{$repoName}}(s rony.Store) *{{$repoName}} {
	return &{{$repoName}}{
		s: s,
	}
}

func (r *{{$repoName}}) SaveWithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{$modelName}}) (err error) {
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

func (r *{{$repoName}}) Save (m *{{$modelName}}) (err error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	
	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.SaveWithTxn(txn, alloc, m)
	})
}

func (r *{{$repoName}}) ReadWithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{$modelName}}) (*{{$modelName}}, error) {
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

func (r *{{$repoName}}) Read (m *{{$modelName}}) (*{{$modelName}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{$modelName}}{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadWithTxn(txn, alloc, m)
		return
	})
	return m, err 
}

`

const genLocalRepo = `
{{$repoName := print .Name "LocalRepo"}}
type {{$repoName}} struct {
    s rony.Store
}

func New{{$repoName}}(s rony.Store) *{{$repoName}} {
	return &{{$repoName}}{
		s: s,
	}
}

`

const genCreate = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) Create(m *{{$modelName}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.CreateWithTxn (txn, alloc, m)
	})
}

func (r *{{$repoName}}) CreateWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, m *{{$modelName}}) (err error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	key := alloc.Gen({{DBKey .Table "m."}})
	if store.ExistsByKey(txn, alloc, key) {
		return store.ErrAlreadyExists
	}
	
	// save table entry
	val := alloc.Marshal(m)
	err = store.SetByKey(txn, val, key)
	if err != nil {
		return
	}
	{{range .Views }}
	// save view entry
	err = store.Set(txn, alloc, val, {{(DBKey . "m.")}})
	if err != nil {
		return err 
	}
	{{- end }}
	{{ if HasIndex . }}
		// key := alloc.Gen({{(DBKey .Table "m.")}})
		{{- range .Fields }}
			{{- if .HasIndex }}
				// update field index by saving new value: {{.Name}}
				{{- if eq .Cardinality "repeated" }}
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

const genUpdate = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) UpdateWithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{$modelName}}) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	
	err := r.DeleteWithTxn(txn, alloc, {{String .Table "m." "," false}})
	if err != nil {
		return err
	}
	
	return r.CreateWithTxn(txn, alloc, m)
}

func (r *{{$repoName}}) Update ({{FuncArgs .Table ""}}, m *{{$modelName}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		return store.ErrEmptyObject
	}
	
	err := r.s.Update(func(txn *rony.StoreTxn) (err error) {
		return r.UpdateWithTxn(txn, alloc, m)
	})

	return err 
}
`

const genSave = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) SaveWithTxn (txn *rony.StoreTxn, alloc *tools.Allocator, m *{{$modelName}}) (err error) {
	if store.Exists(txn, alloc, {{DBKey .Table "m."}}) {
		return r.UpdateWithTxn(txn, alloc, m)
	} else {
		return r.CreateWithTxn(txn, alloc, m)
	}
}

func (r *{{$repoName}}) Save (m *{{$modelName}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.SaveWithTxn(txn, alloc, m)
	})
}
`

const genRead = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) ReadWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, {{FuncArgs .Table ""}}, m *{{$modelName}}) (*{{$modelName}}, error) {
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

func (r *{{$repoName}}) Read ({{FuncArgs .Table ""}}, m *{{$modelName}}) (*{{$modelName}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{$modelName}}{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadWithTxn(txn, alloc, {{String .Table "" "," true}}, m)
		return err 
	})
	return m, err
}
{{ range .Views }}
func (r *{{$repoName}}) ReadBy{{String . "" "And" false}}WithTxn (
	txn *rony.StoreTxn, alloc *tools.Allocator,
	{{FuncArgs . ""}}, m *{{$modelName}},
) (*{{$modelName}}, error) {
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

func (r *{{$repoName}}) ReadBy{{String . "" "And" false}}({{FuncArgs . ""}}, m *{{$modelName}}) (*{{$modelName}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{$modelName}}{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadBy{{String . "" "And" false}}WithTxn (txn, alloc, {{String . "" "," true}}, m)
		return err 
	})
	return m, err
}
{{ end }}
`

const genDelete = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) DeleteWithTxn(txn *rony.StoreTxn, alloc *tools.Allocator, {{FuncArgs .Table ""}}) error {
	{{- if or (gt (len .Views) 0) (HasIndex .) }}
		m := &{{$modelName}}{}
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

func (r *{{$repoName}})  Delete({{FuncArgs .Table ""}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.DeleteWithTxn(txn, alloc, {{String .Table "" "," true}})
	})
}
`

const genHelpers = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
{{- range .Fields }}
	{{- if eq .Cardinality "repeated" }}
		{{- if not (or (eq .Kind "message") (eq .Kind "group")) }}
			{{- if eq .Kind "bytes" }}
				func (x *{{$modelName}}) Has{{Singular .Name}}(xx {{.GoKind}}) bool {
					for idx := range x.{{.Name}} {
						if bytes.Equal(x.{{.Name}}[idx], xx) {
							return true
						}
					}
					return false
				}

			{{- else if eq .Kind "enum" }}
				func (x *{{$modelName}}) Has{{Singular .Name}} (xx {{.GoKind}}) bool {
					for idx := range x.{{.Name}} {
						if x.{{.Name}}[idx] == xx {
							return true
						}
					}
					return false
				}

			{{- else }}
				func (x *{{$modelName}})Has{{Singular .Name}} (xx {{.GoKind}}) bool {
					for idx := range x.{{.Name}} {
						if x.{{.Name}}[idx] == xx {
							return true
						}
					}
					return false
				}

			{{- end }}
		{{- end }}
	{{- end }}
{{ end }}
`
