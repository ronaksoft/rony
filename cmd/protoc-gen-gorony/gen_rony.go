package main

import (
	"fmt"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"hash/crc32"
	"strings"
)

/*
   Creation Time: 2020 - Aug - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// GenCql generates the repo functions for messages which are identified as model with {{@model cql}}
func GenCql(file *protogen.File, g *protogen.GeneratedFile) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/pools"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/repo/cql"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx/v2/qb"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})

	constTables(file, g)
	initCqlQueries(file, g)
	funcsAndFactories(file, g)
}
func constTables(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// Tables")
	g.P("const (")
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		g.P("Table", m.Desc.Name(), "= \"", tools.ToSnake(string(m.Desc.Name())), "\"")
		for _, v := range mm.ViewParams {
			g.P("View", m.Desc.Name(), "By", v, "= \"", fmt.Sprintf("%s_by_%s", tools.ToSnake(string(m.Desc.Name())), tools.ToSnake(v)), "\"")
		}
	}
	g.P(")")
	g.P()
}
func initCqlQueries(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("func init() {")
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		g.P("cql.AddCqlQuery(`")
		g.P("CREATE TABLE IF NOT EXISTS ", tools.ToSnake(mm.Name), " (")
		for _, fn := range mm.FieldNames {
			g.P(fmt.Sprintf("%s \t %s,", tools.ToSnake(fn), mm.FieldsCql[fn]))
		}
		g.P("data \t blob,")
		pksb := strings.Builder{}
		pksb.WriteRune('(')
		switch {
		case len(mm.Table.PKs)+len(mm.Table.CKs) == 1:
			pksb.WriteString(tools.ToSnake(mm.Table.PKs[0]))
		case len(mm.Table.PKs) == 1:
			pksb.WriteString(tools.ToSnake(mm.Table.PKs[0]))
		default:
			pksb.WriteRune('(')
			for idx, pk := range mm.Table.PKs {
				if idx != 0 {
					pksb.WriteString(", ")
				}
				pksb.WriteString(tools.ToSnake(pk))
			}
			pksb.WriteRune(')')

		}
		for _, ck := range mm.Table.CKs {
			pksb.WriteString(", ")
			pksb.WriteString(tools.ToSnake(ck))
		}
		pksb.WriteRune(')')

		orders := strings.Builder{}
		for idx, k := range mm.Table.Orders {
			kWithoutSign := strings.TrimLeft(k, "-")
			if idx > 0 {
				orders.WriteString(", ")
			}
			if strings.HasPrefix(k, "-") {
				orders.WriteString(fmt.Sprintf("%s DESC", tools.ToSnake(kWithoutSign)))
			} else {
				orders.WriteString(fmt.Sprintf("%s ASC", tools.ToSnake(kWithoutSign)))
			}
		}
		g.P("PRIMARY KEY ", pksb.String())
		if len(mm.Table.Orders) > 0 {
			g.P(") WITH CLUSTERING ORDER BY (", orders.String(), ");")
		} else {
			g.P(");")
		}
		g.P("`)")

		// Create Materialized Views
		for idx, v := range mm.Views {
			g.P("cql.AddCqlQuery(`")
			g.P("CREATE MATERIALIZED VIEW IF NOT EXISTS ",
				fmt.Sprintf("%s_by_%s", tools.ToSnake(string(m.Desc.Name())), tools.ToSnake(mm.ViewParams[idx])),
				" AS ",
			)
			g.P("SELECT *")
			g.P("FROM ", tools.ToSnake(mm.Name))
			pksb := strings.Builder{}
			pksb.WriteRune('(')
			switch {
			case len(v.PKs)+len(v.CKs) == 1:
				g.P("WHERE ", tools.ToSnake(v.PKs[0]), " IS NOT null")
				pksb.WriteString(tools.ToSnake(v.PKs[0]))
			case len(v.PKs) == 1:
				g.P("WHERE ", tools.ToSnake(v.PKs[0]), " IS NOT null")
				pksb.WriteString(tools.ToSnake(v.PKs[0]))
			default:
				pksb.WriteRune('(')
				for idx, pk := range v.PKs {
					if idx != 0 {
						g.P("AND ", tools.ToSnake(v.PKs[idx]), " IS NOT null")
						pksb.WriteString(", ")
					} else {
						g.P("WHERE ", tools.ToSnake(v.PKs[idx]), " IS NOT null")
					}
					pksb.WriteString(tools.ToSnake(pk))
				}
				pksb.WriteRune(')')

			}
			for _, ck := range v.CKs {
				g.P("AND ", tools.ToSnake(ck), " IS NOT null")
				pksb.WriteString(", ")
				pksb.WriteString(tools.ToSnake(ck))
			}
			pksb.WriteRune(')')
			g.P("PRIMARY KEY ", pksb.String())
			g.P("`)")
		}
	}
	g.P("}")

}
func funcsAndFactories(file *protogen.File, g *protogen.GeneratedFile) {
	for _, m := range file.Messages {
		mm := _Models[string(m.Desc.Name())]
		if mm == nil {
			continue
		}
		funcInsert(mm, g)
		funcGet(mm, g)
		funcListBy(mm, g)
	}
}
func funcInsert(mm *Model, g *protogen.GeneratedFile) {
	g.P("var _", mm.Name, "InsertFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
	g.P("return qb.Insert(Table", mm.Name, ").")
	columns := strings.Builder{}
	binds := strings.Builder{}
	for _, f := range mm.FieldNames {
		columns.WriteString(fmt.Sprintf("%q, ", tools.ToSnake(f)))
		binds.WriteString(fmt.Sprintf("x.%s, ", f))
	}
	columns.WriteString("\"data\"")
	binds.WriteString("b")

	g.P("Columns(", columns.String(), ").")
	g.P("Query(cql.Session())")
	g.P("})")
	g.P("func ", mm.Name, "Insert (x *", mm.Name, ") (err error) {")
	g.P("q := _", mm.Name, "InsertFactory.GetQuery()")
	g.P("defer _", mm.Name, "InsertFactory.Put(q)")
	g.P()
	g.P("mo := proto.MarshalOptions{UseCachedSize: true}")
	g.P("b := pools.Bytes.GetCap(mo.Size(x))")
	g.P("defer pools.Bytes.Put(b)")
	g.P()
	g.P("b, err = mo.MarshalAppend(b, x)")
	g.P("if err != nil {")
	g.P("return err")
	g.P("}")
	g.P()
	g.P("q.Bind(", binds.String(), ")")
	g.P("err = cql.Exec(q)")
	g.P("return err")
	g.P("}")
}
func funcGet(mm *Model, g *protogen.GeneratedFile) {
	// Generate Factory
	g.P("var _", mm.Name, "GetFactory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
	g.P("return qb.Select(Table", mm.Name, ").")
	g.P("Columns(\"data\").")
	where := strings.Builder{}
	args := strings.Builder{}
	bind := strings.Builder{}
	for idx, f := range mm.Table.PKs {
		if idx != 0 {
			where.WriteString(", ")
			args.WriteString(", ")
			bind.WriteString(", ")
		}
		args.WriteString(fmt.Sprintf("%s %s", tools.ToLowerCamel(f), mm.FieldsGo[f]))
		where.WriteString(fmt.Sprintf("qb.Eq(%q)", tools.ToSnake(f)))
		bind.WriteString(tools.ToLowerCamel(f))
	}
	g.P("Where(", where.String(), ").")
	g.P("Query(cql.Session())")
	g.P("})")
	g.P()

	// Generate Func
	g.P("func ", mm.Name, "Get (", args.String(), ", x *", mm.Name, ") (*", mm.Name, ", error) {")
	g.P("if x == nil {")
	g.P("x = &", mm.Name, "{}")
	g.P("}")
	g.P("q := _", mm.Name, "GetFactory.GetQuery()")
	g.P("defer _", mm.Name, "GetFactory.Put(q)")
	g.P()
	g.P("b := pools.Bytes.GetCap(512)")
	g.P("defer pools.Bytes.Put(b)")
	g.P()
	g.P("q.Bind(", bind.String(), ")")
	g.P("err := cql.Scan(q, &b)")
	g.P("if err != nil {")
	g.P("return x, err")
	g.P("}")
	g.P()
	g.P("err = proto.UnmarshalOptions{Merge:true}.Unmarshal(b, x)")
	g.P("return x, err")
	g.P("}")
	g.P()
}
func funcListBy(mm *Model, g *protogen.GeneratedFile) {
	for idx, v := range mm.ViewParams {
		// Generate Factory
		g.P("var _", mm.Name, "ListBy", v, "Factory = cql.NewQueryFactory(func() *gocqlx.Queryx {")
		g.P("return qb.Select(View", mm.Name, "By", v, ").")
		g.P("Columns(\"data\").")
		where := strings.Builder{}
		args := strings.Builder{}
		bind := strings.Builder{}
		for idx, f := range mm.Views[idx].PKs {
			if idx != 0 {
				bind.WriteString(", ")
				where.WriteString(", ")
				args.WriteString(", ")
			}
			args.WriteString(fmt.Sprintf("%s %s", tools.ToLowerCamel(f), mm.FieldsGo[f]))
			where.WriteString(fmt.Sprintf("qb.Eq(%q)", tools.ToSnake(f)))
			bind.WriteString(tools.ToLowerCamel(f))
		}
		g.P("Where(", where.String(), ").")
		g.P("Query(cql.Session())")
		g.P("})")
		g.P()

		// Generate Function
		g.P("func ", mm.Name, "ListBy", v, " (", args.String(), ", limit int32, f func(x *", mm.Name, ") bool) error {")
		g.P("q := _", mm.Name, "ListBy", v, "Factory.GetQuery()")
		g.P("defer _", mm.Name, "ListBy", v, "Factory.Put(q)")
		g.P()
		g.P("b := pools.Bytes.GetCap(512)")
		g.P("defer pools.Bytes.Put(b)")
		g.P()
		g.P("q.Bind(", bind.String(), ")")
		g.P("iter := q.Iter()")
		g.P("for iter.Scan(&b) {")
		g.P("x := Pool", mm.Name, ".Get()")
		g.P("err := proto.UnmarshalOptions{Merge:true}.Unmarshal(b, x)")
		g.P("if err != nil {")
		g.P("Pool", mm.Name, ".Put(x)")
		g.P("return err")
		g.P("}")
		g.P("if limit--; limit <= 0 || !f(x) {")
		g.P("Pool", mm.Name, ".Put(x)")
		g.P("break")
		g.P("}")
		g.P("Pool", mm.Name, ".Put(x)")
		g.P("}")
		g.P("return nil")
		g.P("}")
	}
}

// GenPools generates codes related for pooling of the messages
func GenPools(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("package ", file.GoPackageName)
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "",
		GoImportPath: "sync",
	})
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "",
		GoImportPath: "github.com/ronaksoft/rony/registry",
	})

	initFunc := strings.Builder{}
	initFunc.WriteString("func init() {\n")

	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		constructor := crc32.ChecksumIEEE([]byte(mtName))
		g.P(fmt.Sprintf("const C_%s int64 = %d", mtName, constructor))
		initFunc.WriteString(fmt.Sprintf("registry.RegisterConstructor(%d, %q)\n", constructor, mtName))
		g.P(fmt.Sprintf("type pool%s struct{", mtName))
		g.P("pool sync.Pool")
		g.P("}")
		g.P(fmt.Sprintf("func (p *pool%s) Get() *%s {", mtName, mtName))
		g.P(fmt.Sprintf("x, ok := p.pool.Get().(*%s)", mtName))
		g.P("if !ok {")
		g.P(fmt.Sprintf("return &%s{}", mtName))
		g.P("}")
		g.P("return x")
		g.P("}")
		g.P("", "")
		g.P(fmt.Sprintf("func (p *pool%s) Put(x *%s) {", mtName, mtName))
		for _, ft := range mt.Fields {
			ftName := ft.Desc.Name()
			ftPkg, _ := descName(file, g, ft.Desc.Message())
			switch ft.Desc.Cardinality() {
			case protoreflect.Repeated:
				g.P(fmt.Sprintf("x.%s = x.%s[:0]", ftName, ftName))
			default:
				switch ft.Desc.Kind() {
				case protoreflect.BytesKind:
					g.P(fmt.Sprintf("x.%s = x.%s[:0]", ftName, ftName))
				case protoreflect.MessageKind:
					// If it is message we check if is nil then we leave it
					// If it is from same package use Pool
					g.P(fmt.Sprintf("if x.%s != nil {", ftName))
					if ftPkg != "" {
						g.P(ftPkg, ".Pool", ft.Desc.Message().Name(), ".Put(x.", ftName, ")")
					} else {
						g.P("Pool", ft.Desc.Message().Name(), ".Put(x.", ftName, ")")
					}
					g.P("x.", ftName, " = nil")
					g.P("}")
				default:
					g.P(fmt.Sprintf("x.%s = %s", ftName, zeroValue(ft.Desc.Kind())))

				}
			}
		}
		g.P("p.pool.Put(x)")
		g.P("}")
		g.P("")
		g.P(fmt.Sprintf("var Pool%s = pool%s{}", mtName, mtName))
		g.P("")
	}
	initFunc.WriteString("}")
	g.P("")
	g.P(initFunc.String())
	g.P()
}

// GenDeepCopy generates codes which deep copy a message
func GenDeepCopy(file *protogen.File, g *protogen.GeneratedFile) {
	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		g.P("func (x *", mtName, ") DeepCopy(z *", mtName, ") {")
		for _, ft := range mt.Fields {
			ftName := ft.Desc.Name()
			ftPkg, ftType := descName(file, g, ft.Desc.Message())
			switch ft.Desc.Cardinality() {
			case protoreflect.Repeated:
				switch ft.Desc.Kind() {
				case protoreflect.MessageKind:
					g.P("for idx := range x.", ftName, "{")
					g.P(fmt.Sprintf("if x.%s[idx] != nil {", ftName))
					if ftPkg == "" {
						g.P("xx := Pool", ftType, ".Get()")
					} else {
						g.P("xx := ", ftPkg, ".Pool", ftType, ".Get()")
					}
					g.P("x.", ftName, "[idx].DeepCopy(xx)")
					g.P("z.", ftName, " = append(z.", ftName, ", xx)")
					g.P("}")
					g.P("}")
				default:
					g.P("z.", ftName, " = append(z.", ftName, "[:0], x.", ftName, "...)")
				}
			default:
				switch ft.Desc.Kind() {
				case protoreflect.BytesKind:
					g.P("z.", ftName, " = append(z.", ftName, "[:0], x.", ftName, "...)")
				case protoreflect.MessageKind:
					// If it is message we check if is nil then we leave it
					// If it is from same package use Pool
					g.P(fmt.Sprintf("if x.%s != nil {", ftName))
					if ftPkg == "" {
						g.P("z.", ftName, " = Pool", ftType, ".Get()")
					} else {
						g.P("z.", ftName, " = ", ftPkg, ".Pool", ftType, ".Get()")
					}
					g.P("x.", ftName, ".DeepCopy(z.", ftName, ")")
					g.P("}")
				default:
					g.P(fmt.Sprintf("z.%s = x.%s", ftName, ftName))

				}
			}
		}
		g.P("}")
		g.P()
	}
}

// GenPushToContext generates codes related for pooling of the messages
func GenPushToContext(file *protogen.File, g *protogen.GeneratedFile) {
	for _, mt := range file.Messages {
		mtName := mt.Desc.Name()
		g.P("func (x *", mtName, ") PushToContext(ctx *edge.RequestCtx) {")
		g.P("ctx.PushMessage(C_", mtName, ", x)")
		g.P("}")
		g.P()
	}
}

// GenRPC generates the server and client interfaces if any proto service has been defined
func GenRPC(file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) > 0 {
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "",
			GoImportPath: "github.com/ronaksoft/rony/edge",
		})
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "",
			GoImportPath: "google.golang.org/protobuf/proto",
		})
		g.QualifiedGoIdent(protogen.GoIdent{
			GoName:       "",
			GoImportPath: "fmt",
		})
		if file.GoPackageName != "rony" {
			g.QualifiedGoIdent(protogen.GoIdent{
				GoName:       "MessageEnvelope",
				GoImportPath: "github.com/ronaksoft/rony",
			})
		}
	}

	// Generate Server
	for _, s := range file.Services {
		for _, m := range s.Methods {
			constructor := crc32.ChecksumIEEE([]byte(m.Desc.Name()))
			g.P("const C_", m.Desc.Name(), " int64 = ", fmt.Sprintf("%d", constructor))
		}
		g.P()
		g.P("type I", s.Desc.Name(), " interface {")
		for _, m := range s.Methods {
			inputPkg, inputType := descName(file, g, m.Desc.Input())
			outputPkg, outputType := descName(file, g, m.Desc.Output())
			inputName := inputType
			if inputPkg != "" {
				inputName = fmt.Sprintf("%s.%s", inputPkg, inputType)
			}
			outputName := outputType
			if outputPkg != "" {
				outputName = fmt.Sprintf("%s.%s", outputPkg, outputType)
			}
			g.P(m.Desc.Name(), "(ctx *edge.RequestCtx, req *", inputName, ", res *", outputName, ")")
		}
		g.P("}")
		g.P()
		g.P()
		g.P("type ", s.Desc.Name(), "Wrapper struct {")
		g.P("h I", s.Desc.Name())
		g.P("}")
		g.P()
		g.P("func Register", s.Desc.Name(), "(h I", s.Desc.Name(), ", e *edge.Server) {")
		g.P("w := ", s.Desc.Name(), "Wrapper{")
		g.P("h: h,")
		g.P("}")
		g.P("w.Register(e)")
		g.P("}")
		g.P()
		g.P("func (sw *", s.Desc.Name(), "Wrapper) Register (e *edge.Server) {")
		for _, m := range s.Methods {
			g.P("e.SetHandlers(C_", m.Desc.Name(), ", true, sw.", m.Desc.Name(), "Wrapper)")
		}
		g.P("}")
		g.P()
		for _, m := range s.Methods {
			inputPkg, inputType := descName(file, g, m.Desc.Input())
			outputPkg, outputType := descName(file, g, m.Desc.Output())

			g.P("func (sw *", s.Desc.Name(), "Wrapper) ", m.Desc.Name(), "Wrapper (ctx *edge.RequestCtx, in *rony.MessageEnvelope) {")
			if inputPkg == "" {
				g.P("req := Pool", inputType, ".Get()")
				g.P("defer Pool", inputType, ".Put(req)")
			} else {
				g.P("req := ", inputPkg, ".Pool", inputType, ".Get()")
				g.P("defer ", inputPkg, ".Pool", inputType, ".Put(req)")
			}
			if outputPkg == "" {
				g.P("res := Pool", outputType, ".Get()")
				g.P("defer Pool", outputType, ".Put(res)")
			} else {
				g.P("res := ", outputPkg, ".Pool", outputType, ".Get()")
				g.P("defer ", outputPkg, ".Pool", outputType, ".Put(res)")
			}

			g.P("err := proto.UnmarshalOptions{Merge:true}.Unmarshal(in.Message, req)")
			g.P("if err != nil {")
			g.P("ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)")
			g.P("return")
			g.P("}")
			g.P()
			g.P("sw.h.", m.Desc.Name(), "(ctx, req, res)")
			g.P("if !ctx.Stopped() {")
			if outputPkg == "" {
				g.P("ctx.PushMessage(C_", outputType, ", res)")
			} else {
				g.P("ctx.PushMessage(", outputPkg, ".C_", outputType, ", res)")
			}

			g.P("}")
			g.P("}")
			g.P()
		}
	}

	// Generate Client
	g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "Client",
		GoImportPath: "github.com/ronaksoft/rony/edgec",
	})
	for _, s := range file.Services {
		g.P("type ", s.Desc.Name(), "Client struct {")
		g.P("c edgec.Client")
		g.P("}")
		g.P()
		g.P("func New", s.Desc.Name(), "Client (ec edgec.Client) *", s.Desc.Name(), "Client {")
		g.P("return &", s.Desc.Name(), "Client{")
		g.P("c: ec,")
		g.P("}")
		g.P("}")
		g.P()
		for _, m := range s.Methods {
			inputPkg, inputType := descName(file, g, m.Desc.Input())
			outputPkg, outputType := descName(file, g, m.Desc.Output())
			inputName := inputType
			if inputPkg != "" {
				inputName = fmt.Sprintf("%s.%s", inputPkg, inputType)
			}
			outputName := outputType
			if outputPkg != "" {
				outputName = fmt.Sprintf("%s.%s", outputPkg, outputType)
			}
			// constructor := crc32.ChecksumIEEE([]byte(*m.Name))
			g.P("func (c *", s.Desc.Name(), "Client) ", m.Desc.Name(), "(req *", inputName, ", kvs ...*rony.KeyValue) (*", outputName, ", error) {")
			g.P("out := rony.PoolMessageEnvelope.Get()")
			g.P("defer rony.PoolMessageEnvelope.Put(out)")
			g.P("in := rony.PoolMessageEnvelope.Get()")
			g.P("defer rony.PoolMessageEnvelope.Put(in)")
			g.P("out.Fill(c.c.GetRequestID(), C_", m.Desc.Name(), ", req, kvs...)")
			g.P("err := c.c.Send(out, in)")
			g.P("if err != nil {")
			g.P("return nil, err")
			g.P("}")
			g.P("switch in.GetConstructor() {")
			if outputPkg != "" {
				g.P("case ", outputPkg, ".C_", outputType, ":")
				g.P("x := &", outputPkg, ".", outputType, "{}")
			} else {
				g.P("case C_", outputType, ":")
				g.P("x := &", outputType, "{}")
			}

			g.P("_ = proto.Unmarshal(in.Message, x)")
			g.P("return x, nil")
			g.P("case rony.C_Error:")
			g.P("x := &rony.Error{}")
			g.P("_ = proto.Unmarshal(in.Message, x)")
			g.P("return nil, fmt.Errorf(\"%s:%s\", x.GetCode(), x.GetItems())")
			g.P("default:")
			g.P("return nil, fmt.Errorf(\"unknown message: %d\", in.GetConstructor())")
			g.P("}")
			g.P("}")
			g.P()
		}
	}
}

func zeroValue(t protoreflect.Kind) string {
	switch t {
	case protoreflect.BoolKind:
		return "false"
	case protoreflect.StringKind:
		return "\"\""
	case protoreflect.MessageKind:
		return "nil"
	case protoreflect.BytesKind:
		return "nil"
	default:
		return "0"
	}
}
