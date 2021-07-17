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

var singletonFuncs = map[string]interface{}{
	"DBKey": func(arg codegen.MessageArg) string {
		return fmt.Sprintf("'S', C_%s",
			arg.Name,
		)
	},
}
var aggregateFuncs = map[string]interface{}{
	"Singular": func(x string) string {
		return inflection.Singular(x)
	},
	"Plural": func(x string) string {
		return inflection.Plural(x)
	},
	"MVName": func(m codegen.ModelKey) string {
		alias := m.Alias()
		if alias == "" {
			alias = m.Names(codegen.PropFilterALL, "", "", "", codegen.None)
		}
		sb := strings.Builder{}
		sb.WriteString(m.Name())
		sb.WriteString(alias)
		return sb.String()
	},
	"MVAlias": func(m codegen.ModelKey, prefix string) string {
		if m.Alias() != "" {
			return m.Alias()
		}
		sb := strings.Builder{}
		sb.WriteString(prefix)
		sb.WriteString(m.Names(codegen.PropFilterALL, "", "", "", codegen.None))

		return sb.String()
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
			prefix, inflection.Singular(f.NameCC()), postfix,
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

func Generate(g *Generator, arg codegen.MessageArg) {
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
		g.g.P(g.Exec(template.Must(template.New("genPrimaryKey").Funcs(aggregateFuncs).Parse(genPrimaryKey)), arg))
		g.g.P(g.Exec(template.Must(template.New("genLocalRepo").Funcs(aggregateFuncs).Parse(genLocalRepo)), arg))
		g.g.P(g.Exec(template.Must(template.New("genCreate").Funcs(aggregateFuncs).Parse(genCreate)), arg))
		g.g.P(g.Exec(template.Must(template.New("genUpdate").Funcs(aggregateFuncs).Parse(genUpdate)), arg))
		g.g.P(g.Exec(template.Must(template.New("genSave").Funcs(aggregateFuncs).Parse(genSave)), arg))
		g.g.P(g.Exec(template.Must(template.New("genRead").Funcs(aggregateFuncs).Parse(genRead)), arg))
		g.g.P(g.Exec(template.Must(template.New("genDelete").Funcs(aggregateFuncs).Parse(genDelete)), arg))
		g.g.P(g.Exec(template.Must(template.New("genList").Funcs(aggregateFuncs).Parse(genList)), arg))
		g.g.P(g.Exec(template.Must(template.New("genIter").Funcs(aggregateFuncs).Parse(genIter)), arg))
		g.g.P(g.Exec(template.Must(template.New("genListByIndex").Funcs(aggregateFuncs).Parse(genListByIndex)), arg))
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

func (r *{{$repoName}}) Create(m *{{$modelName}}) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	return r.s.Update(func(txn *rony.StoreTxn) error {
		return r.CreateWithTxn (txn, alloc, m)
	})
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
func (r *{{$repoName}}) ReadBy{{MVAlias . ""}}WithTxn (
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

func (r *{{$repoName}}) ReadBy{{MVAlias . ""}}({{FuncArgs . ""}}, m *{{$modelName}}) (*{{$modelName}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	if m == nil {
		m = &{{$modelName}}{}
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		m, err = r.ReadBy{{MVAlias . ""}}WithTxn (txn, alloc, {{String . "" "," true}}, m)
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
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	{{ if or (gt (len .Views) 0) (HasIndex .) }}
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

const genPrimaryKey = `
{{$modelName := .Name}}
type {{$modelName}}PrimaryKey interface {
	makeItPrivate()
}

type {{$modelName}}PK struct {
{{- range .Table.Keys }}
{{.Name}}  {{.GoType}}
{{- end }}
}

func ({{$modelName}}PK) makeItPrivate() {}

{{- range .Views }}

type {{MVName .}}PK struct {
{{- range .Keys }}
{{.Name}}  {{.GoType}}
{{- end }}
}

func ({{MVName .}}PK) makeItPrivate() {}
{{ end }}
`

const genIter = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) IterWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator, offset {{$modelName}}PrimaryKey, ito *store.IterOption, cb func(m *{{$modelName}}) bool,
) error {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	
	var seekKey []byte
	opt := store.DefaultIteratorOptions
	opt.Reverse = ito.Backward()

	switch offset := offset.(type) {
	case {{$modelName}}PK:
		opt.Prefix = alloc.Gen({{DBKeyPKs .Table "offset."}})
		seekKey = alloc.Gen({{DBKey .Table "offset."}})
{{ range .Views }}
	case {{MVName .}}PK:
		opt.Prefix = alloc.Gen({{DBKeyPKs . "offset."}})
		seekKey = alloc.Gen({{DBKey . "offset."}})
{{ end }}
	default:
		opt.Prefix = alloc.Gen({{DBPrefix .Table}})
		seekKey = opt.Prefix
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		iter := txn.NewIterator(opt)
		if ito.OffsetKey() == nil {
			iter.Seek(seekKey)
		} else {
			iter.Seek(ito.OffsetKey())
		}
		exitLoop := false
		for ; iter.ValidForPrefix(opt.Prefix); iter.Next() {
			err = iter.Item().Value(func(val []byte) error {
				m := &{{$modelName}}{}
				err := m.Unmarshal(val)
				if err != nil {
					return err
				}
				if !cb(m) {
					exitLoop = true
				} 
				return nil
			})
			if err != nil || exitLoop {
				break
			}
		}
		iter.Close()

		return
	})

	return err
}

func (r *{{$repoName}}) Iter(
	pk {{$modelName}}PrimaryKey, ito *store.IterOption, cb func(m *{{$modelName}}) bool,
) error {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	

	return r.s.View(func(txn *rony.StoreTxn) error {
		return r.IterWithTxn(txn, alloc, pk, ito, cb)
	})
}
`

const genList = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
func (r *{{$repoName}}) ListWithTxn(
	txn *rony.StoreTxn, alloc *tools.Allocator, offset {{$modelName}}PrimaryKey, lo *store.ListOption, cond func(m *{{$modelName}}) bool,
) ([]*{{$modelName}}, error) {
	if alloc == nil {
		alloc = tools.NewAllocator()
		defer alloc.ReleaseAll()
	}
	
	var seekKey []byte
	opt := store.DefaultIteratorOptions
	opt.Reverse = lo.Backward()
	res := make([]*{{$modelName}}, 0, lo.Limit())
	
	switch offset := offset.(type) {
	case {{$modelName}}PK:
		opt.Prefix = alloc.Gen({{DBKeyPKs .Table "offset."}})
		seekKey = alloc.Gen({{DBKey .Table "offset."}})
{{ range .Views }}
	case {{MVName .}}PK:
		opt.Prefix = alloc.Gen({{DBKeyPKs . "offset."}})
		seekKey = alloc.Gen({{DBKey . "offset."}})
{{ end }}
	default:
		opt.Prefix = alloc.Gen({{DBPrefix .Table}})
		seekKey = opt.Prefix
	}

	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
		for iter.Seek(seekKey); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			err = iter.Item().Value(func(val []byte) error {
				m := &{{$modelName}}{}
				err := m.Unmarshal(val)
				if err != nil {
					return err
				}
				if cond == nil || cond(m) {
					res = append(res, m)
				} else {
					limit++
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		iter.Close()
		return
	})

	return res, err
}

func (r *{{$repoName}}) List(
	pk {{$modelName}}PrimaryKey, lo *store.ListOption, cond func(m *{{$modelName}}) bool,
) ([]*{{$modelName}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()
	
	var (
		res []*{{$modelName}}
		err error
	)
	err = r.s.View(func(txn *rony.StoreTxn) error {
		res, err = r.ListWithTxn(txn, alloc, pk, lo, cond)
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
`

const genListByIndex = `
{{$model := .}}
{{$repoName := print .Name "LocalRepo"}}
{{$modelName := .Name}}
{{ range .Fields }}
{{ if .HasIndex }}
func (r *{{$repoName}}) ListBy{{Singular .Name}} ({{Singular .NameCC}} {{.GoKind}}, lo *store.ListOption, cond func(*{{$modelName}}) bool) ([]*{{$modelName}}, error) {
	alloc := tools.NewAllocator()
	defer alloc.ReleaseAll()

	opt := store.DefaultIteratorOptions
	opt.Reverse = lo.Backward()
	opt.Prefix = alloc.Gen({{IndexDBPrefix $model . "" ""}})
	res := make([]*{{$modelName}}, 0, lo.Limit())
	err := r.s.View(func(txn *rony.StoreTxn) (err error) {
		iter := txn.NewIterator(opt)
		offset := lo.Skip()
		limit := lo.Limit()
		for iter.Seek(opt.Prefix); iter.ValidForPrefix(opt.Prefix); iter.Next() {
			if offset--; offset >= 0 {
				continue
			}
			if limit--; limit < 0 {
				break
			}
			err = iter.Item().Value(func(val []byte) error {
				b, err := store.GetByKey(txn, alloc, val)
				if err != nil {
					return err
				}
				m := &{{$modelName}}{}
				err = m.Unmarshal(b)
				if err != nil {
					return err
				}
				if cond == nil || cond(m) {
					res = append(res, m)
				} else {
					limit++
				}
				return nil
			})
			if err != nil {
				break
			}
		}
		iter.Close()
		return err
	})
	if err != nil {
		return nil, err
	}

	return res, nil 
}
{{ end }}
{{ end }}

`
