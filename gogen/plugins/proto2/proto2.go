package proto2

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gogen"
	"github.com/iancoleman/strcase"
	"strings"
)

/*
   Creation Time: 2020 - Apr - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type ProtoBuffer struct {
	g *gogen.Generator
}

func (g *ProtoBuffer) Name() string {
	return "ProtoBuffer2 Generator"
}

func (g *ProtoBuffer) Init(gen *gogen.Generator) {
	g.g = gen
}

func (g *ProtoBuffer) Generate(desc *gogen.Descriptor) {
	for _, m := range desc.Models {
		for _, c := range m.Comments {
			g.g.P("//", c)
		}
		g.g.P("message", m.Name, "{")
		g.g.In()
		for idx, p := range m.Properties {
			ro := "required"
			if p.CheckOption(gogen.Optional) {
				ro = "optional"
			}
			if p.CheckOption(gogen.Slice) {
				ro = "repeated"
			}

			tags := strings.Builder{}
			for _, t := range p.Tags {
				parts := strings.Split(t, ":")
				switch len(parts) {
				case 1:
					tags.WriteString("(gogoproto.moretags) = ")
					tags.WriteString(parts[0])
					tags.WriteRune(':')
					tags.WriteString(strcase.ToSnake(p.Name))
				case 2:
					tags.WriteString("(gogoproto.moretags) = ")
					tags.WriteString(parts[0])
					tags.WriteRune(':')
					tags.WriteString(parts[1])
				default:
					panic(fmt.Sprintf("invalid tag: %v", t))
				}

			}
			if tags.Len() > 0 {
				g.g.P(ro, p.Type, p.Name, "=", idx+1, fmt.Sprintf("[%s]", tags.String()), ";", "//", p.Comment)
			} else {
				g.g.P(ro, p.Type, p.Name, "=", idx+1, ";", "//", p.Comment)
			}

		}
		g.g.Out()
		g.g.P("}")
		g.g.Nl()

	}
}

func (g *ProtoBuffer) GeneratePrepend(desc *gogen.Descriptor) {
	g.g.P("syntax = \"proto2\";")
	g.g.P("package ", desc.Name, ";")
	g.g.Nl(3)
	return
}
