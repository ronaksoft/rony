package codegen

import (
	"strings"
)

/*
   Creation Time: 2021 - Jul - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	optKeyRony                 = "rony_opt"
	optValEdgeDep              = "no_edge_dep"
	optValCrc32                = "crc32"
	optValConstructorJSONStr   = "json_str"
	optValConstructorJSONInt64 = "json_int64"
	optOpenAPI                 = "open_api"
	optCleanProto              = "clean_proto"
)

type ConstructorFormat string

const (
	StringJSON ConstructorFormat = "str"
	Int64JSON  ConstructorFormat = "int64"
)

type PluginOptions struct {
	NoEdgeDependency  bool
	CRC32             bool
	ConstructorFormat ConstructorFormat
	OpenAPI           bool
	ExportCleanProto  bool
}

func (op *PluginOptions) ParamFunc(key, value string) error {
	switch key {
	case optKeyRony:
		for _, v := range strings.Split(value, "|") {
			switch strings.TrimSpace(v) {
			case optValEdgeDep:
				op.NoEdgeDependency = true
			case optValCrc32:
				op.CRC32 = true
			case optValConstructorJSONStr:
				op.ConstructorFormat = StringJSON
			case optValConstructorJSONInt64:
				op.ConstructorFormat = Int64JSON
			case optOpenAPI:
				op.OpenAPI = true
			case optCleanProto:
				op.ExportCleanProto = true
			}
		}
	}

	return nil
}
