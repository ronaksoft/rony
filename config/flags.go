package config

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

/*
	Flags
*/

func StringFlag(name, value, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.String(name, value, usage)
	}
}

func IntFlag(name string, value int, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int(name, value, usage)
	}
}

func Int32Flag(name string, value int32, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int32(name, value, usage)
	}
}

func Uint32Flag(name string, value uint32, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Uint32(name, value, usage)
	}
}

func Int64Flag(name string, value int64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int64(name, value, usage)
	}
}

func Uint64Flag(name string, value uint64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Uint64(name, value, usage)
	}
}

func BoolFlag(name string, value bool, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Bool(name, value, usage)
	}
}

func DurationFlag(name string, value time.Duration, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Duration(name, value, usage)
	}
}

func StringSliceFlag(name string, value []string, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.StringSlice(name, value, usage)
	}
}

func Int64SliceFlag(name string, value []int64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int64Slice(name, value, usage)
	}
}
