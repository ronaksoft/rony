package tools

import (
	"fmt"
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




/*
	Flags
 */

func StringFlag(name, value, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.String(name, value, usage)
	}
}

func GetString(cmd *cobra.Command, name string) string {
	f, err := cmd.Flags().GetString(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}

func IntFlag(name string, value int, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int(name, value, usage)
	}
}

func GetInt(cmd *cobra.Command, name string) int {
	f, err := cmd.Flags().GetInt(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}

func Int32Flag(name string, value int32, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int32(name, value, usage)
	}
}

func GetInt32(cmd *cobra.Command, name string) int32 {
	f, err := cmd.Flags().GetInt32(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}

func UInt32Flag(name string, value uint32, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Uint32(name, value, usage)
	}
}

func GetUInt32(cmd *cobra.Command, name string) uint32 {
	f, err := cmd.Flags().GetUint32(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}

func Int64Flag(name string, value int64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int64(name, value, usage)
	}
}

func GetInt64(cmd *cobra.Command, name string) int64 {
	f, err := cmd.Flags().GetInt64(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}

func UInt64Flag(name string, value uint64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Uint64(name, value, usage)
	}
}

func GetUInt64(cmd *cobra.Command, name string) uint64 {
	f, err := cmd.Flags().GetUint64(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}

func BoolFlag(name string, value bool, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Bool(name, value, usage)
	}
}

func GetBool(cmd *cobra.Command, name string) bool {
	f, err := cmd.Flags().GetBool(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}


func DurationFlag(name string, value time.Duration, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Duration(name, value, usage)
	}
}

func GetDuration(cmd *cobra.Command, name string) time.Duration {
	f, err := cmd.Flags().GetDuration(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}

func StringSliceFlag(name string, value []string, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.StringSlice(name, value, usage)
	}
}

func GetStringSlice(cmd *cobra.Command, name string) []string {
	f, err := cmd.Flags().GetStringSlice(name)
	if err != nil {
		panic(fmt.Sprintf("flag is not registered: %s for %s", name, cmd.Name()))
	}
	return f
}

func Int64SliceFlag(name string, value []int64, usage string) FlagOption {
	return func(fs *pflag.FlagSet) {
		fs.Int64Slice(name, value, usage)
	}
}


