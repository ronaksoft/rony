package tools

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"time"
)

/*
   Creation Time: 2020 - Jun - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type FlagOption func(fs *pflag.FlagSet)

func SetFlags(cmd *cobra.Command, opts ...FlagOption) {
	for _, opt := range opts {
		opt(cmd.Flags())
	}
}

func SetPersistentFlags(cmd *cobra.Command, opts ...FlagOption) {
	for _, opt := range opts {
		opt(cmd.PersistentFlags())
	}
}

func RegisterStringFlag(name, value, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.String(name, value, usage)
	}
}

func RegisterIntFlag(name string, value int, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int(name, value, usage)
	}
}

func RegisterInt32Flag(name string, value int32, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int32(name, value, usage)
	}
}

func RegisterUInt64Flag(name string, value uint64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Uint64(name, value, usage)
	}
}

func RegisterBoolFlag(name string, value bool, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Bool(name, value, usage)
	}
}

func RegisterDurationFlag(name string, value time.Duration, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Duration(name, value, usage)
	}
}

func RegisterStringSliceFlag(name string, value []string, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.StringSlice(name, value, usage)
	}
}

func RegisterInt64SliceFlag(name string, value []int64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int64Slice(name, value, usage)
	}
}

func RegisterInt64Flag(name string, value int64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int64(name, value, usage)
	}
}
