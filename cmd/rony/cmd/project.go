package cmd

import (
	"context"
	"git.ronaksoft.com/ronak/rony/tools"
	"github.com/gobuffalo/genny/v2"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
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
	workingDir, _ := os.Getwd()
	tools.SetFlags(newCmd,
		tools.StringFlag("projectPath", workingDir, "the root path of the project"),
		tools.StringFlag("goPackage", "main", "the full path of go package in go.mod file"),
	)
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(newCmd)
}

var projectCmd = &cobra.Command{
	Use: "project",
}

var newCmd = &cobra.Command{
	Use: "new",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dryRun")
		r := genny.WetRunner(context.Background())
		if dryRun {
			r = genny.DryRunner(context.Background())
		}

		projectPath, _ := cmd.Flags().GetString("projectPath")
		goPackage, _ := cmd.Flags().GetString("goPackage")

		g := genny.New()
		createFolders(g, projectPath)
		goModuleInit(g, projectPath, goPackage)
		copyFiles(g, projectPath)

		// Create a Runner with the Generator customized by command's arguments
		err := r.With(g)
		if err != nil {
			return err
		}
		return r.Run()
	},
}

func createFolders(g *genny.Generator, projectPath string) {
	cmd := exec.Command("rm", "-r", "*")
	cmd.Env = os.Environ()
	g.Command(cmd)

	cmd = exec.Command("mkdir", "-p", filepath.Join(projectPath, "model"))
	cmd.Env = os.Environ()
	g.Command(cmd)

	cmd = exec.Command("mkdir", "-p", filepath.Join(projectPath, "service"))
	cmd.Env = os.Environ()
	g.Command(cmd)

	cmd = exec.Command("mkdir", "-p", filepath.Join(projectPath, "pkg/client"))
	cmd.Env = os.Environ()
	g.Command(cmd)

	cmd = exec.Command("mkdir", "-p", filepath.Join(projectPath, "cmd/cli"))
	cmd.Env = os.Environ()
	g.Command(cmd)
}
func goModuleInit(g *genny.Generator, projectPath, goPackage string) {
	cmd := exec.Command("go", "mod", "init", goPackage)
	cmd.Env = os.Environ()
	cmd.Dir = projectPath
	g.Command(cmd)
}
func copyFiles(g *genny.Generator, projectPath string) {
	f, err := pkger.Open("git.ronaksoft.com/ronak/rony:/internal/templates/main.go")
	if err != nil {
		panic(err)
	}
	g.File(genny.NewFile(filepath.Join(projectPath, "main.go"), f))
	_ = f.Close()
	f, err = pkger.Open("git.ronaksoft.com/ronak/rony:/internal/templates/server.go")
	if err != nil {
		panic(err)
	}
	g.File(genny.NewFile(filepath.Join(projectPath, "server.go"), f))
	_ = f.Close()
}
