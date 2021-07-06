package rpc

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/cmd/protoc-gen-gorony/z"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"hash/crc32"
	"strings"
)

/*
   Creation Time: 2021 - Jul - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Arg struct {
	Name    string
	NameCC  string // LowerCamelCase(Name)
	C       uint32
	Methods []struct {
		Name        string
		NameCC      string // LowerCamelCase(Name)
		InputName   string
		InputType   string
		InputPkg    string
		InputC      uint32
		OutputName  string
		OutputType  string
		OutputPkg   string
		OutputC     uint32
		RestEnabled bool
		TunnelOnly  bool
		Rest        struct {
			Method string
			Path   string
			Json   bool
		}
	}
}

func GetArg(g *Generator, s *protogen.Service) Arg {
	arg := Arg{
		Name:   string(s.Desc.Name()),
		NameCC: tools.ToLowerCamel(string(s.Desc.Name())),
		C:      crc32.ChecksumIEEE([]byte(s.Desc.Name())),
	}
	for _, m := range s.Methods {
		inputPkg, inputType := z.DescParts(g.f, g.g, m.Input.Desc)
		outputPkg, outputType := z.DescParts(g.f, g.g, m.Output.Desc)

		opt, _ := s.Desc.Options().(*descriptorpb.MethodOptions)
		restOpt := proto.GetExtension(opt, rony.E_RonyRest).(*rony.RestOpt)
		rest := struct {
			Method string
			Path   string
			Json   bool
		}{Method: "", Path: "", Json: false}
		if restOpt != nil {
			rest.Method = restOpt.GetMethod()
			rest.Path = fmt.Sprintf("/%s", strings.Trim(restOpt.GetPath(), "/"))
			rest.Json = restOpt.GetJsonEncode()
		}
		arg.Methods = append(
			arg.Methods,
			struct {
				Name        string
				NameCC      string
				InputName   string
				InputType   string
				InputPkg    string
				InputC      uint32
				OutputName  string
				OutputType  string
				OutputPkg   string
				OutputC     uint32
				RestEnabled bool
				TunnelOnly  bool
				Rest        struct {
					Method string
					Path   string
					Json   bool
				}
			}{
				Name:        string(m.Desc.Name()),
				NameCC: tools.ToLowerCamel(string(m.Desc.Name())),
				InputName:   z.DescName(g.f, g.g, m.Input.Desc),
				InputType:   inputType,
				InputPkg:    inputPkg,
				InputC:      crc32.ChecksumIEEE([]byte(m.Input.Desc.Name())),
				OutputName:  z.DescName(g.f, g.g, m.Output.Desc),
				OutputType:  outputType,
				OutputPkg:   outputPkg,
				OutputC:     crc32.ChecksumIEEE([]byte(m.Output.Desc.Name())),
				RestEnabled: restOpt != nil,
				TunnelOnly:  proto.GetExtension(opt, rony.E_RonyInternal).(bool),
				Rest:        rest,
			},
		)
	}

	return arg
}
