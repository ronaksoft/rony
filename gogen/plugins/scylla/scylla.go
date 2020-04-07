package scylla

import (
	"git.ronaksoftware.com/ronak/rony/gogen"
)

/*
   Creation Time: 2020 - Apr - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Plugin struct {
	imports map[string]struct{}
	g       *gogen.Generator
}

func (s *Plugin) Name() string {
	return "Scylla"
}

func (s *Plugin) Init(g *gogen.Generator) {
	s.g = g
	s.imports = make(map[string]struct{})
}

func (s *Plugin) Generate(desc *gogen.Descriptor) {
	for _, m := range desc.Models {
		for _, c := range m.Comments {
			s.g.P("//", c)
		}
		s.g.P("type", m.Name, "struct {")
		s.g.In()
		for _, p := range m.Properties {
			s.g.P(p.Name, p.Type, "//", p.Comment)
		}
		s.g.Out()
		s.g.P("}")
		s.g.P()
	}

}

func (s *Plugin) GeneratePrepend(desc *gogen.Descriptor) {
	panic("implement me")
}
