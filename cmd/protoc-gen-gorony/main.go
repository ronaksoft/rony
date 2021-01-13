package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/model"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/model/cqlgen"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/model/kvgen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
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
			case "github.com/dgraph-io/badger":
				return "github.com/dgraph-io/badger/v2"
			}
			return path
		},
	}
	pgo.Run(func(plugin *protogen.Plugin) error {
		for _, f := range plugin.Files {
			if !f.Generate {
				continue
			}

			if f.Proto.GetPackage() == "google.protobuf" {
				continue
			}

			// reset the global model and fill with the new data
			model.ResetModels()
			for _, m := range f.Messages {
				model.FillModel(m)
			}

			// Generate Pools
			g1 := plugin.NewGeneratedFile(fmt.Sprintf("%s.rony.go", f.GeneratedFilenamePrefix), f.GoImportPath)
			GenPools(f, g1)
			GenDeepCopy(f, g1)
			GenPushToContext(f, g1)
			GenMarshal(f, g1)
			GenUnmarshal(f, g1)

			// Generate Model's repo functionality
			if len(model.GetModels()) > 0 {
				opt, _ := f.Desc.Options().(*descriptorpb.FileOptions)
				storage := proto.GetExtension(opt, rony.E_RonyStorage).(string)
				switch strings.ToLower(storage) {
				case "local":
					kvgen.Generate(f, g1)
				case "cql":
					cqlgen.Generate(f, g1)
				}
			}

			// Generate RPCs if there is any service definition in the file
			if len(f.Services) > 0 {
				g1.QualifiedGoIdent(protogen.GoIdent{
					GoName:       "",
					GoImportPath: "github.com/ronaksoft/rony/edge",
				})
				g1.QualifiedGoIdent(protogen.GoIdent{
					GoName:       "",
					GoImportPath: "google.golang.org/protobuf/proto",
				})
				g1.QualifiedGoIdent(protogen.GoIdent{
					GoName:       "",
					GoImportPath: "fmt",
				})
				if f.GoPackageName != "rony" {
					g1.QualifiedGoIdent(protogen.GoIdent{
						GoName:       "",
						GoImportPath: "github.com/ronaksoft/rony",
					})
				}
				g1.QualifiedGoIdent(protogen.GoIdent{
					GoName:       "Client",
					GoImportPath: "github.com/ronaksoft/rony/edgec",
				})

				for _, s := range f.Services {
					GenRPC(f, s, g1)
					opt, _ := s.Desc.Options().(*descriptorpb.ServiceOptions)
					if proto.GetExtension(opt, rony.E_RonyCobraCmd).(bool) {
						g1.QualifiedGoIdent(protogen.GoIdent{
							GoName:       "",
							GoImportPath: "github.com/spf13/cobra",
						})
						g1.QualifiedGoIdent(protogen.GoIdent{
							GoName:       "",
							GoImportPath: "github.com/ronaksoft/rony/config",
						})

						GenCobraCmd(f, s, g1)
					}
				}

			}

		}
		return nil
	})
	return
}
