package cmd

import (
	"embed"
	"os"

	"github.com/ronaksoft/rony/config"
	"github.com/spf13/cobra"
)

/*
   Creation Time: 2020 - Aug - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var Skeleton embed.FS

func init() {
	workingDir, _ := os.Getwd()

	config.SetPersistentFlags(RootCmd,
		config.BoolFlag("dry-run", false, "don't actually change anything"),
	)

	config.SetFlags(CreateProjectCmd,
		config.StringFlag("project.dir", workingDir, "the root path of the project"),
		config.StringFlag("project.name", "sample-project", "the name of the project"),
		config.StringFlag("package.path", "github.com/sample", "the full path of go package in go.mod file"),
	)

	config.SetFlags(GenProtoCmd,
		config.StringFlag("project.dir", workingDir, "the root path of the project"),
	)

	config.SetFlags(ExportProtoCmd,
		config.StringFlag("lang", "js", "possible values: js, dart"),
		config.StringFlag("c.format", "str", "possible values: str, int64"),
	)
	RootCmd.AddCommand(CreateProjectCmd, GenProtoCmd, ExportProtoCmd)
}

var RootCmd = &cobra.Command{
	Use: "rony",
}
