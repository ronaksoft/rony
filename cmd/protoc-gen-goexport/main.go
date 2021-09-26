package main

import (
	"fmt"
	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	pluginOpt = &codegen.PluginOptions{}
	pgo       = protogen.Options{
		ParamFunc: pluginOpt.ParamFunc,
	}
)

func main() {
	pgo.Run(
		func(plugin *protogen.Plugin) error {
			return clearRonyTags(plugin)
		},
	)
}

func clearRonyTags(plugin *protogen.Plugin) error {
	for _, protoFile := range plugin.Files {
		if !protoFile.Generate || protoFile.Proto.GetPackage() == "google.protobuf" {
			continue
		}

		// Create the generator func
		gFile := plugin.NewGeneratedFile(fmt.Sprintf("%s.clean.proto", protoFile.GeneratedFilenamePrefix), protoFile.GoImportPath)
		gFile.P("syntax = \"", protoFile.Proto.GetSyntax(), "\";")
		gFile.P()
		gFile.P("package ", protoFile.Proto.GetPackage(), ";")
		gFile.P()
		for _, dep := range protoFile.Proto.Dependency {
			for _, f := range plugin.Request.FileToGenerate {
				if f == dep {
					gFile.P("import \"", dep, "\";")
				}
			}
		}
		for _, s := range protoFile.Services {
			gFile.P()
			for _, c := range s.Comments.LeadingDetached {
				gFile.P(s.Comments.Leading, " ", c, " ", s.Comments.Trailing)
			}
			gFile.P("service ", s.Desc.Name(), "{")
			for _, m := range s.Methods {
				gFile.P("\t rpc ", m.Desc.Name(), "(", m.Desc.Input().Name(), ") returns (", m.Desc.Output().Name(), ");")
			}
			gFile.P("}")
		}
		for _, m := range protoFile.Messages {
			gFile.P()
			for _, c := range m.Comments.LeadingDetached {
				gFile.P(m.Comments.Leading, " ", c, " ", m.Comments.Trailing)
			}
			gFile.P("message ", m.Desc.Name(), "{")
			for _, f := range m.Fields {
				switch f.Desc.Kind() {
				case protoreflect.MessageKind:
					gFile.P("\t", f.Desc.Cardinality().String(), " ", f.Desc.Message().Name(), " ", f.Desc.Name(), " = ", f.Desc.Number(), ";")
				case protoreflect.EnumKind:
					gFile.P("\t", f.Desc.Cardinality().String(), " ", f.Desc.Enum().Name(), " ", f.Desc.Name(), " = ", f.Desc.Number(), ";")
				default:
					gFile.P("\t", f.Desc.Cardinality().String(), " ", f.Desc.Kind(), " ", f.Desc.Name(), " = ", f.Desc.Number(), ";")
				}

			}
			gFile.P("}")

		}
		for _, m := range protoFile.Enums {
			gFile.P()
			for _, c := range m.Comments.LeadingDetached {
				gFile.P(m.Comments.Leading, " ", c, " ", m.Comments.Trailing)
			}
			gFile.P("enum ", m.Desc.Name(), "{")
			for _, f := range m.Values {
				gFile.P("\t", f.Desc.Name(), " = ", f.Desc.Number(), ";")
			}
			gFile.P("}")
		}


	}

	return nil
}
