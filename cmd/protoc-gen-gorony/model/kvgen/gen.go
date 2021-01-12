package kvgen

import (
	"google.golang.org/protobuf/compiler/protogen"
)

/*
   Creation Time: 2021 - Jan - 12
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Generate generates the repo functions for messages which are identified as model with {{@entity cql}}
func Generate(file *protogen.File, g *protogen.GeneratedFile) {
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/pools"})
	g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/ronaksoft/rony/repo/kv"})
	// g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx"})
	// g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "github.com/scylladb/gocqlx/v2/qb"})
	// g.QualifiedGoIdent(protogen.GoIdent{GoName: "", GoImportPath: "google.golang.org/protobuf/proto"})

	// constTables(file, g)
	// initCqlQueries(file, g)
	// funcsAndFactories(file, g)
}
