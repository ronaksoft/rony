package module

import (
	"google.golang.org/protobuf/compiler/protogen"
)

/*
   Creation Time: 2021 - Jul - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Generator struct {
	p      *protogen.Plugin
	f      *protogen.File
	g      *protogen.GeneratedFile
	remote string
	local  string
}

func New(p *protogen.Plugin) *Generator {
	return &Generator{
		p: p,
	}
}

func (g *Generator) checkFiles() error {
	return nil
}

func (g *Generator) filterFiles() []*protogen.File {
	var ret []*protogen.File
	for _, f := range g.p.Files {
		if f.Generate {
			ret = append(ret, f)
		}
	}
	return ret
}

func (g *Generator) Generate() error {
	err := g.checkFiles()
	if err != nil {
		return err
	}
	files := g.filterFiles()
	g.g = g.p.NewGeneratedFile("text.go", files[0].GoImportPath)
	g.g.P("package ", files[0].GoPackageName)
	for _, f := range files {
		for _, m := range f.Messages {
			g.g.P("// ", m.GoIdent.GoName, " ", m.GoIdent.GoImportPath, " ", m.Desc.ParentFile().Package())
		}
	}
	return nil
}
