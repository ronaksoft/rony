package main

import (
	parse "git.ronaksoftware.com/ronak/rony/internal/parser"
	"github.com/scylladb/go-reflectx"
	"google.golang.org/protobuf/compiler/protogen"
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

func GenCql(file *protogen.File, g *protogen.GeneratedFile) {
	g.P("package ", file.GoPackageName)
	g.P("import (")
	g.P("\"git.ronaksoftware.com/ronak/rony\"")
	g.P("\"git.ronaksoftware.com/ronak/rony/pools\"")
	g.P("\"git.ronaksoftware.com/ronak/rony/repo/cql\"")
	g.P("\"github.com/scylladb/gocqlx/v2\"")
	g.P("\"github.com/scylladb/gocqlx/v2/qb\"")
	g.P(")")
	g.P()

	g.P("const (")
	g.P("colData = \"data\"")
	g.P("colVersion = \"ver\"")
	g.P(")")
	g.P()

	g.P("// Tables")
	g.P("const (")

	fields := map[string]struct{}{}
	for _, m := range file.Messages {
		t, err := parse.Parse(string(m.Desc.Name()), string(m.Comments.Leading))
		if err != nil {
			panic(err)
		}
		for _, n := range t.Root.Nodes {
			switch n.Type() {
			case parse.NodeTable:
				nn := n.(*parse.TableNode)
				for _, k := range nn.PartitionKeys {
					fields[k] = struct{}{}
				}
				for _, k := range nn.ClusteringKeys {
					fields[k] = struct{}{}
				}
				g.P("Table", m.Desc.Name(), "= \"t_", reflectx.CamelToSnakeASCII(string(m.Desc.Name())), "\"")
			case parse.NodeView:
				nn := n.(*parse.ViewNode)
				sb := strings.Builder{}
				for idx, k := range nn.PartitionKeys {
					fields[k] = struct{}{}
					if idx != 0 {
						sb.WriteString("_")
					}
					sb.WriteString(reflectx.CamelToSnakeASCII(k))
				}
				for _, k := range nn.ClusteringKeys {
					fields[k] = struct{}{}
				}

				g.P("View", m.Desc.Name(), "= \"mv_", reflectx.CamelToSnakeASCII(string(m.Desc.Name())), "_by_", sb.String(), "\"")
			}
		}
	}
	g.P(")")
	g.P()

	g.P("// Columns")
	g.P("const (")
	for n := range fields {
		g.P("Col", n, " = \"", reflectx.CamelToSnakeASCII(n), "\"")
	}
	g.P(")")
}
