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
	optKeyRony    = "rony_opt"
	optValModule  = "module"
	optValEdgeDep = "no_edge_dep"
)

type PluginOptions struct {
	NoEdgeDependency bool
	Module           bool
}

func (op *PluginOptions) ParamFunc(key, value string) error {
	switch key {
	case optKeyRony:
		for _, v := range strings.Split(value, "|") {
			switch strings.TrimSpace(v) {
			case optValModule:
				op.Module = true
			case optValEdgeDep:
				op.NoEdgeDependency = true
			}
		}
	}
	return nil
}
