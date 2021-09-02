package main

import (
	"fmt"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/helper"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/module"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/repo"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/rpc"
	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	pluginOpt = &codegen.PluginOptions{}
	pgo       = protogen.Options{
		ParamFunc: pluginOpt.ParamFunc,
		ImportRewriteFunc: func(path protogen.GoImportPath) protogen.GoImportPath {
			// TODO:: this is a hack for bug in Golang/Protobuf which does not support go module versions
			switch path {
			case "github.com/scylladb/gocqlx":
				return "github.com/scylladb/gocqlx/v2"
			}

			return path
		},
	}
)

func main() {
	pgo.Run(
		func(plugin *protogen.Plugin) error {
			if pluginOpt.CRC32 {
				codegen.CrcBits = 32
			}
			err := normalMode(plugin)
			if err != nil {
				return err
			}
			if pluginOpt.Module {
				return moduleMode(plugin)
			}
			return nil
		},
	)
}

func normalMode(plugin *protogen.Plugin) error {
	protocVer := plugin.Request.GetCompilerVersion()
	for _, protoFile := range plugin.Files {
		if !protoFile.Generate || protoFile.Proto.GetPackage() == "google.protobuf" {
			continue
		}

		// Create the generator func
		generatedFile := plugin.NewGeneratedFile(fmt.Sprintf("%s.rony.go", protoFile.GeneratedFilenamePrefix), protoFile.GoImportPath)
		generatedFile.P("// Code generated by Rony's protoc plugin; DO NOT EDIT.")
		generatedFile.P("// ProtoC ver. v", protocVer.GetMajor(), ".", protocVer.GetMinor(), ".", protocVer.GetPatch())
		generatedFile.P("// Rony ver. ", codegen.Version)
		generatedFile.P("// Source: ", protoFile.Proto.GetName())
		generatedFile.P()

		// Generate all the helper functions
		g1 := helper.New(protoFile, generatedFile, pluginOpt)
		g1.Generate()

		// Generate rpc helper functions (Server, Client and CLI)
		g2 := rpc.New(protoFile, generatedFile)
		g2.Generate()

		// Generate Repository functionalities
		g3 := repo.New(plugin, protoFile, generatedFile)
		g3.Generate()

	}

	return nil
}

func moduleMode(plugin *protogen.Plugin) error {
	g := module.New(plugin)
	return g.Generate()
}
