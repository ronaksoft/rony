package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func main() {

	plugins := make(map[string]struct{})
	pgo := protogen.Options{
		ParamFunc: func(name, value string) error {
			switch name {
			case "plugin":
				for _, p := range strings.Split(value, "|") {
					plugins[p] = struct{}{}
				}
			}
			return nil
		},
		ImportRewriteFunc: func(path protogen.GoImportPath) protogen.GoImportPath {
			// TODO:: this is a hack for bug in Golang/Protobuf which does not support go module versions
			switch path {
			case "github.com/scylladb/gocqlx":
				return "github.com/scylladb/gocqlx/v2"
			}
			return path
		},
	}
	pgo.Run(func(plugin *protogen.Plugin) error {
		for _, f := range plugin.Files {
			// reset the global model and fill with the new data
			resetModels()
			for _, m := range f.Messages {
				fillModel(m)
			}

			// Generate Pools
			g1 := plugin.NewGeneratedFile(fmt.Sprintf("%s.rony.go", f.GeneratedFilenamePrefix), f.GoImportPath)
			GenPools(f, g1)
			GenDeepCopy(f, g1)
			if len(getModels()) > 0 {
				// Generate Model's repo functionality
				GenCql(f, g1)
			}
			if len(f.Services) > 0 {
				// Generate RPCs
				GenRPC(f, g1)
			}

		}
		return nil
	})
	return
}
