package cmd

import (
	"context"
	"git.ronaksoft.com/ronak/rony/tools"
	"github.com/gobuffalo/genny"
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

func init() {
	tools.SetFlags(newCmd,
		tools.StringFlag("projectPath", "", "the project path in go.mod"),
	)
	RootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use: "new",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dryRun")
		r := genny.WetRunner(context.Background())
		if dryRun {
			r = genny.DryRunner(context.Background())
		}

		// Create a Runner with the Generator customized by command's arguments
		err := r.With(gen(cmd))
		if err != nil {
			return err
		}
		return r.Run()
	},
}

func gen(cmd *cobra.Command) *genny.Generator {
	projectPath, _ := cmd.Flags().GetString("projectPath")
	projectName, _ := cmd.Flags().GetString("projectName")

	g := genny.New()
	createFolders(g, projectName)
	goModuleInit(g, projectPath)
	return g
}
func createFolders(g *genny.Generator, projectPath string) {}
func goModuleInit(g *genny.Generator, projectPath string)  {}
func copyFiles(g *genny.Generator)                         {}
