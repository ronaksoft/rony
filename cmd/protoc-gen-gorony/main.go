package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	pgo := protogen.Options{
		ParamFunc:         nil,
		ImportRewriteFunc: nil,
	}
	pgo.Run(func(plugin *protogen.Plugin) error {
		for _, f := range plugin.Files {
			resetModels()
			g1 := plugin.NewGeneratedFile(fmt.Sprintf("%s.pools.go", f.GeneratedFilenamePrefix), f.GoImportPath)
			GenPools(f, g1)
			if len(f.Services) > 0 {
				g2 := plugin.NewGeneratedFile(fmt.Sprintf("%s.rony.go", f.GeneratedFilenamePrefix), f.GoImportPath)
				GenRPC(f, g2)
			}

			for _, m := range f.Messages {
				fillModel(m)
			}
			if len(getModels()) > 0 {
				g3 := plugin.NewGeneratedFile(fmt.Sprintf("%s.cql.go", f.GeneratedFilenamePrefix), f.GoImportPath)
				GenCql(f, g3)
			}
		}
		return nil
	})
	return
}
