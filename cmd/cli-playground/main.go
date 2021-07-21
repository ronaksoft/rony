package main

import (
	"github.com/c-bata/go-prompt"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/tools"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

func main() {
	log.SetLevel(log.InfoLevel)
	runtime.GOMAXPROCS(runtime.NumCPU())
	RootCmd.AddCommand(ExitCmd)
	p := prompt.New(tools.PromptExecutor(RootCmd), tools.PromptCompleter(RootCmd))
	p.Run()
}

var RootCmd = &cobra.Command{
	Use: "Root",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
}

var ExitCmd = &cobra.Command{
	Use: "exit",
	Run: func(cmd *cobra.Command, args []string) {
		for _, e := range Edges {
			_ = e.Cluster().Leave()
			e.Shutdown()
		}
		os.Exit(0)
	},
}
