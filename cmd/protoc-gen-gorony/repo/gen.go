package repo

import (
	"fmt"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/repo/cql"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/repo/store"
	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

/*
   Creation Time: 2021 - Jul - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	localRepoPrefix  = "Local"
	GlobalRepoPrefix = "Global"
)

type Generator struct {
	p         *protogen.Plugin
	f         *protogen.File
	g         *protogen.GeneratedFile
	initBlock *strings.Builder
}

func New(p *protogen.Plugin, f *protogen.File, g *protogen.GeneratedFile) *Generator {
	return &Generator{
		p:         p,
		f:         f,
		g:         g,
		initBlock: &strings.Builder{},
	}
}

func (g *Generator) Generate() {
	g.g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "bytes"})
	g.g.P("var _ = bytes.MinRead")

	var cqlGen *protogen.GeneratedFile
	for _, m := range g.f.Messages {
		arg := codegen.GetMessageArg(m).With(g.f)
		if !arg.IsAggregate && !arg.IsSingleton {
			continue
		}

		switch arg.LocalRepo {
		case "store":
			store.Generate(store.New(g.f, g.g, localRepoPrefix), arg)
		case "cql":
			if cqlGen == nil {
				cqlGen = g.p.NewGeneratedFile(fmt.Sprintf("%s.cql", g.f.GeneratedFilenamePrefix), g.f.GoImportPath)
			}
			cql.GenerateCQL(cql.New(g.f, cqlGen, localRepoPrefix), arg)
			cql.GenerateGo(cql.New(g.f, g.g, localRepoPrefix), arg)
		default:
		}

		switch arg.GlobalRepo {
		case "cql":
			if cqlGen == nil {
				cqlGen = g.p.NewGeneratedFile(fmt.Sprintf("%s.cql", g.f.GeneratedFilenamePrefix), g.f.GoImportPath)
			}
			cql.GenerateCQL(cql.New(g.f, cqlGen, GlobalRepoPrefix), arg)
			cql.GenerateGo(cql.New(g.f, g.g, GlobalRepoPrefix), arg)
		default:
		}
	}
}
