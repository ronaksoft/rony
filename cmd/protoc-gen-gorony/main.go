package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/go-openapi/spec"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/helper"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/repo"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/rpc"
	"github.com/ronaksoft/rony/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	pathRegEx = regexp.MustCompile(``)
)

func main() {
	pgo.Run(
		func(plugin *protogen.Plugin) error {
			if pluginOpt.CRC32 {
				codegen.CrcBits = 32
			}

			switch pluginOpt.ConstructorFormat {
			case codegen.StringJSON:
				return jsonStr(plugin)
			case codegen.Int64JSON:
				return jsonInt(plugin)
			}

			if pluginOpt.OpenAPI {
				return exportOpenAPI(plugin)
			}

			if pluginOpt.ExportCleanProto {
				return clearRonyTags(plugin)
			}

			err := normalMode(plugin)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func normalMode(plugin *protogen.Plugin) error {
	protocVer := plugin.Request.GetCompilerVersion()
	for _, protoFile := range plugin.Files {
		if !protoFile.Generate {
			continue
		}

		// Create the generator func
		generatedFile := plugin.NewGeneratedFile(
			fmt.Sprintf("%s.rony.go", protoFile.GeneratedFilenamePrefix), protoFile.GoImportPath,
		)
		generatedFile.P("// Code generated by Rony's protoc plugin; DO NOT EDIT.")
		generatedFile.P(
			"// ProtoC ver. v",
			protocVer.GetMajor(), ".", protocVer.GetMinor(), ".", protocVer.GetPatch(),
		)
		generatedFile.P("// Rony ver. ", codegen.Version)
		generatedFile.P("// Source: ", protoFile.Proto.GetName())
		generatedFile.P()

		// Generate all the helper functions
		_ = helper.GenFunc(generatedFile, pluginOpt, protoFile)

		// Generate rpc helper functions (Server, Client and CLI)
		_ = rpc.GenFunc(generatedFile, pluginOpt, protoFile)

		// Generate Repository functionalities
		g3 := repo.New(plugin, protoFile, generatedFile)
		g3.Generate()
	}

	return nil
}

func jsonStr(plugin *protogen.Plugin) error {
	var (
		importPath protogen.GoImportPath
		filePrefix string
		cn         = map[string]uint64{}
		cs         = map[uint64]string{}
	)
	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}
		importPath = f.GoImportPath
		filePrefix = f.GeneratedFilenamePrefix
		// reset the global model and fill with the new data
		for _, mt := range f.Messages {
			constructor := codegen.CrcHash([]byte(mt.Desc.Name()))
			cn[string(mt.Desc.Name())] = constructor
			cs[constructor] = string(mt.Desc.Name())
		}
		for _, s := range f.Services {
			for _, m := range s.Methods {
				methodName := fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())
				constructor := codegen.CrcHash([]byte(methodName))
				cn[methodName] = constructor
				cs[constructor] = methodName
			}
		}
	}

	t := template.Must(template.New("t1").Parse(`
	{
	    "ConstructorsByName": {
	    {{range $k,$v := .}}    "{{$k}}": "{{$v}}",
		{{end -}}
		},
		"ConstructorsByValue": {
		{{range $k,$v := .}}    "{{$v}}": "{{$k}}",
		{{end -}}
		}
	}
	`))

	out := &bytes.Buffer{}
	err := t.Execute(out, cn)
	if err != nil {
		panic(err)
	}

	gf := plugin.NewGeneratedFile(filepath.Join(filepath.Dir(filePrefix), "constructors.json"), importPath)
	_, err = gf.Write(out.Bytes())

	return err
}

func jsonInt(plugin *protogen.Plugin) error {
	var (
		importPath protogen.GoImportPath
		filePrefix string
		cn         = map[string]int64{}
		cs         = map[int64]string{}
	)
	for _, f := range plugin.Files {
		if !f.Generate {
			continue
		}
		importPath = f.GoImportPath
		filePrefix = f.GeneratedFilenamePrefix
		// reset the global model and fill with the new data
		for _, mt := range f.Messages {
			constructor := int64(codegen.CrcHash([]byte(mt.Desc.Name())))
			cn[string(mt.Desc.Name())] = constructor
			cs[constructor] = string(mt.Desc.Name())
		}
		for _, s := range f.Services {
			for _, m := range s.Methods {
				methodName := fmt.Sprintf("%s%s", s.Desc.Name(), m.Desc.Name())
				constructor := int64(codegen.CrcHash([]byte(methodName)))
				cn[methodName] = constructor
				cs[constructor] = methodName
			}
		}
	}

	t := template.Must(template.New("t1").Parse(`
	{
	    "ConstructorsByName": {
	    {{range $k,$v := .}}    "{{$k}}": "{{$v}}",
		{{end -}}
		},
		"ConstructorsByValue": {
		{{range $k,$v := .}}    "{{$v}}": "{{$k}}",
		{{end -}}
		}
	}
	`))

	out := &bytes.Buffer{}
	err := t.Execute(out, cn)
	if err != nil {
		panic(err)
	}

	gf := plugin.NewGeneratedFile(filepath.Join(filepath.Dir(filePrefix), "constructors.json"), importPath)
	_, err = gf.Write(out.Bytes())

	return err
}

func exportOpenAPI(plugin *protogen.Plugin) error {
	var (
		importPath protogen.GoImportPath
		filePrefix string
	)
	swag := &spec.Swagger{}
	swag.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Description: "Rony Swagger Service Description",
			Title:       "Rony",
			Version:     "0.0.1",
		},
	}
	swag.Schemes = []string{"http", "https"}
	swag.Swagger = "2.0"

	for _, protoFile := range plugin.Files {
		if !protoFile.Generate || protoFile.Proto.GetPackage() == "google.protobuf" {
			continue
		}
		importPath = protoFile.GoImportPath
		filePrefix = protoFile.GeneratedFilenamePrefix

		arg := codegen.GenTemplateArg(protoFile)
		for _, s := range arg.Services {
			addTag(swag, s)
			for _, m := range s.Methods {
				if !m.RestEnabled {
					continue
				}
				addOperation(swag, s, m)
				addDefinition(swag, m.Input, m.Rest.Json)
				addDefinition(swag, m.Output, m.Rest.Json)
			}
		}
	}

	out, err := swag.MarshalJSON()
	if err != nil {
		return err
	}

	gf := plugin.NewGeneratedFile(filepath.Join(filepath.Dir(filePrefix), "swagger.json"), importPath)

	_, err = gf.Write(out)

	return err
}
func addTag(swag *spec.Swagger, s codegen.ServiceArg) {
	swag.Tags = append(
		swag.Tags,
		spec.NewTag(s.Name(), s.Comments, nil),
	)
}
func addDefinition(swag *spec.Swagger, m codegen.MessageArg, jsonEncode bool) {
	if swag.Definitions == nil {
		swag.Definitions = map[string]spec.Schema{}
	}

	def := spec.Schema{}
	def.Description = m.Comments
	def.Typed("object", "")
	for _, f := range m.Fields {
		fName := f.DescName()
		if jsonEncode {
			fName = f.JSONName()
		}
		switch f.GoKind {
		case "string":
			def.SetProperty(fName, *spec.StringProperty())
		case "[]byte":
			def.SetProperty(fName, *spec.ArrayProperty(spec.Int8Property()))
		case "int8", "uint8":
			def.SetProperty(fName, *spec.Int8Property())
		case "int16", "uint16":
			def.SetProperty(fName, *spec.Int16Property())
		case "int32", "uint32":
			def.SetProperty(fName, *spec.Int32Property())
		case "int64", "uint64":
			def.SetProperty(fName, *spec.Int64Property())
		case "float32":
			def.SetProperty(fName, *spec.Float32Property())
		case "float64":
			def.SetProperty(fName, *spec.Float64Property())
		default:
			def.SetProperty(fName, *spec.StringProperty())
		}
	}

	swag.Definitions[m.Name()] = def
}
func addOperation(swag *spec.Swagger, s codegen.ServiceArg, m codegen.MethodArg) {
	if swag.Paths == nil {
		swag.Paths = &spec.Paths{
			Paths: map[string]spec.PathItem{},
		}
	}

	opID := fmt.Sprintf("%s%s", s.NameCC(), m.Name())
	op := spec.NewOperation(opID).
		RespondsWith(200,
			spec.NewResponse().
				WithSchema(
					spec.RefProperty(fmt.Sprintf("#/definitions/%s", m.Output.Name())),
				),
		).
		WithTags(s.Name())

	if m.Rest.Json {
		op.WithProduces("application/json").
			WithConsumes("application/json")
	} else {
		op.WithProduces("application/protobuf").
			WithConsumes("application/protobuf")
	}

	for name, kind := range m.Rest.PathParams {
		op.AddParam(
			setParamType(
				spec.PathParam(name).
					AsRequired().
					NoEmptyValues(),
				kind,
			),
		)
	}

	for name, kind := range m.Rest.QueryParams {
		op.AddParam(
			setParamType(
				spec.QueryParam(name).
					AsRequired().
					NoEmptyValues(),
				kind,
			),
		)
	}

	if m.Rest.Unmarshal {
		op.AddParam(
			spec.BodyParam(
				m.Input.Name(),
				spec.RefProperty(fmt.Sprintf("#/definitions/%s", m.Input.Name())),
			),
		)
	}

	restPath := replacePath(m.Rest.Path)
	pathItem := swag.Paths.Paths[restPath]
	switch strings.ToLower(m.Rest.Method) {
	case "get":
		pathItem.Get = op
	case "post":
		pathItem.Post = op
	case "put":
		pathItem.Put = op
	case "delete":
		pathItem.Delete = op
	case "patch":
		pathItem.Patch = op
	}
	swag.Paths.Paths[restPath] = pathItem
}
func setParamType(p *spec.Parameter, kind protoreflect.Kind) *spec.Parameter {
	switch kind {
	case protoreflect.StringKind:
		p.Typed("string", kind.String())
	case protoreflect.BytesKind:
		p.Typed("array", "int8")
	case protoreflect.DoubleKind, protoreflect.FloatKind:
		p.Typed("number", kind.String())
	case protoreflect.Int32Kind, protoreflect.Sint32Kind,
		protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		p.Typed("integer", "int32")
	case protoreflect.Int64Kind, protoreflect.Sint64Kind,
		protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		p.Typed("integer", "int64")
	default:
		p.Typed("integer", kind.String())
	}

	return p
}

func replacePath(path string) string {
	sb := strings.Builder{}
	for idx, p := range strings.Split(path, "/") {
		if idx > 0 {
			sb.WriteRune('/')
		}
		if strings.HasPrefix(p, ":") {
			sb.WriteRune('{')
			sb.WriteString(p[1:])
			sb.WriteRune('}')
		} else {
			sb.WriteString(p)
		}
	}

	return sb.String()
}
func clearRonyTags(plugin *protogen.Plugin) error {
	for _, protoFile := range plugin.Files {
		if !protoFile.Generate || protoFile.Proto.GetPackage() == "google.protobuf" {
			continue
		}

		// Create the generator func
		gFile := plugin.NewGeneratedFile(
			fmt.Sprintf(
				"%s.clean.proto",
				protoFile.GeneratedFilenamePrefix,
			),
			protoFile.GoImportPath,
		)
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
				for _, c := range m.Comments.LeadingDetached {
					gFile.P(m.Comments.Leading, " ", c, " ", m.Comments.Trailing)
				}
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
				for _, c := range f.Comments.LeadingDetached {
					gFile.P(f.Comments.Leading, " ", c, " ", f.Comments.Trailing)
				}
				switch protoFile.Proto.GetSyntax() {
				case "proto3":
					switch f.Desc.Cardinality() {
					case protoreflect.Optional, protoreflect.Required:
						switch f.Desc.Kind() {
						case protoreflect.MessageKind:
							gFile.P("\t", f.Desc.Message().Name(), " ", f.Desc.Name(), " = ", f.Desc.Number(), ";")
						case protoreflect.EnumKind:
							gFile.P("\t", f.Desc.Enum().Name(), " ", f.Desc.Name(), " = ", f.Desc.Number(), ";")
						default:
							gFile.P("\t", f.Desc.Kind(), " ", f.Desc.Name(), " = ", f.Desc.Number(), ";")
						}
					case protoreflect.Repeated:
						switch f.Desc.Kind() {
						case protoreflect.MessageKind:
							gFile.P(
								"\t",
								f.Desc.Cardinality().String(), " ", f.Desc.Message().Name(), " ",
								f.Desc.Name(), " = ", f.Desc.Number(), ";",
							)
						case protoreflect.EnumKind:
							gFile.P(
								"\t",
								f.Desc.Cardinality().String(), " ", f.Desc.Enum().Name(), " ",
								f.Desc.Name(), " = ", f.Desc.Number(), ";",
							)
						default:
							gFile.P(
								"\t",
								f.Desc.Cardinality().String(), " ", f.Desc.Kind(), " ",
								f.Desc.Name(), " = ", f.Desc.Number(), ";",
							)
						}
					}
				case "proto2":
					switch f.Desc.Kind() {
					case protoreflect.MessageKind:
						gFile.P(
							"\t",
							f.Desc.Cardinality().String(), " ", f.Desc.Message().Name(), " ",
							f.Desc.Name(), " = ", f.Desc.Number(), ";",
						)
					case protoreflect.EnumKind:
						gFile.P(
							"\t",
							f.Desc.Cardinality().String(), " ", f.Desc.Enum().Name(), " ",
							f.Desc.Name(), " = ", f.Desc.Number(), ";",
						)
					default:
						gFile.P(
							"\t",
							f.Desc.Cardinality().String(), " ", f.Desc.Kind(), " ",
							f.Desc.Name(), " = ", f.Desc.Number(), ";",
						)
					}
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
