package tools

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
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

// FlagOption applies some predefined configurations on a FlagSet.
type FlagOption func(fs *pflag.FlagSet)

// SetFlags applies 'opts' on 'cmd' flags in order.
func SetFlags(cmd *cobra.Command, opts ...FlagOption) {
	for _, opt := range opts {
		opt(cmd.Flags())
	}
}

// SetPersistentFlags applies 'opts' on 'cmd' persistent flags.
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

/*

	Helper function for cli programs

*/

// PromptCompleter returns a completer function used by prompt package.
// This function is useful to create an interactive shell based on the rootCmd as the root commands.
func PromptCompleter(rootCmd *cobra.Command) func(d prompt.Document) []prompt.Suggest {
	return func(d prompt.Document) []prompt.Suggest {
		suggests := make([]prompt.Suggest, 0, 10)
		cols := d.TextBeforeCursor()
		currCmd := rootCmd
		for _, col := range strings.Fields(cols) {
			for _, cmd := range currCmd.Commands() {
				if cmd.Name() == col {
					currCmd = cmd

					break
				}
			}
		}

		currWord := d.GetWordBeforeCursor()
		if strings.HasPrefix(currWord, "--") {
			// Search in Flags
			rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
				if strings.HasPrefix(flag.Name, currWord[2:]) {
					suggests = append(suggests, prompt.Suggest{
						Text:        fmt.Sprintf("--%s", flag.Name),
						Description: flag.Usage,
					})
				}
			})
			currCmd.Flags().VisitAll(func(flag *pflag.Flag) {
				if strings.HasPrefix(flag.Name, currWord[2:]) {
					suggests = append(suggests, prompt.Suggest{
						Text:        fmt.Sprintf("--%s", flag.Name),
						Description: flag.DefValue,
					})
				}
			})
		} else {
			for _, cmd := range currCmd.Commands() {
				if strings.HasPrefix(cmd.Name(), currWord) {
					suggests = append(suggests, prompt.Suggest{
						Text:        cmd.Name(),
						Description: cmd.Short,
					})
				}
			}
		}

		return suggests
	}
}

// PromptExecutor returns an executor function used by prompt package.
// This function is useful to create an interactive shell based on the rootCmd as the root commands.
func PromptExecutor(rootCmd *cobra.Command) func(s string) {
	return func(s string) {
		if strings.TrimSpace(s) == "" {
			return
		}
		rootCmd.SetArgs(strings.Fields(s))
		_ = rootCmd.Execute()
	}
}

// RunShell runs an interactive shell
func RunShell(cmd *cobra.Command) {
	prompt.New(PromptExecutor(cmd), PromptCompleter(cmd)).Run()
}
