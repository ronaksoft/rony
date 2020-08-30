package cmd

import (
	"context"
	"fmt"
	"git.ronaksoft.com/ronak/rony/tools"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
   Creation Time: 2020 - Aug - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	skeletonPath = "git.ronaksoft.com/ronak/rony:/internal/templates"
)

func init() {
	workingDir, _ := os.Getwd()
	tools.SetFlags(newCmd,
		tools.StringFlag("projectPath", workingDir, "the root path of the project"),
		tools.StringFlag("goPackage", "main", "the full path of go package in go.mod file"),
		tools.BoolFlag("withExample", false, "use an example project instead of empty project"),
	)
	tools.SetFlags(buildProtoCmd,
		tools.StringFlag("projectPath", workingDir, "the root path of the project"),
	)

	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(newCmd, buildProtoCmd)

	_ = pkger.Walk(skeletonPath, func(path string, info os.FileInfo, err error) error {
		return nil
	})
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
		withExample, _ := cmd.Flags().GetBool("withExample")

		g := genny.New()
		setupSkeleton(g, projectPath, goPackage, withExample)
		goModuleInit(g, projectPath, goPackage)

		// Create a Runner with the Generator customized by command's arguments
		err := r.With(g)
		if err != nil {
			return err
		}
		return r.Run()
	},
}

func setupSkeleton(g *genny.Generator, projectPath, goPackage string, withExample bool) {
	cmd := exec.Command("mkdir", "-p", projectPath)
	cmd.Env = os.Environ()
	g.Command(cmd)

	pathPrefix := skeletonPath + "/skel"
	if withExample {
		pathPrefix += "_example"
	}
	err := pkger.Walk(pathPrefix, func(path string, info os.FileInfo, err error) error {
		realPath := strings.TrimPrefix(path, pathPrefix)
		if info.IsDir() {
			g.File(genny.NewDir(filepath.Join(projectPath, realPath), os.ModeDir|0744))
		} else {
			f, err := pkger.Open(path)
			if err != nil {
				return err
			}
			tCtx := plush.NewContext()
			tCtx.Set("goPackage", func() string {
				return goPackage
			})
			s, err := plush.RenderR(f, tCtx)
			if err != nil {
				return err
			}
			g.File(genny.NewFileS(filepath.Join(projectPath, realPath), s))
			_ = f.Close()
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
func goModuleInit(g *genny.Generator, projectPath, goPackage string) {
	cmd := exec.Command("go", "mod", "init", goPackage)
	cmd.Env = os.Environ()
	cmd.Dir = projectPath
	g.Command(cmd)
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Env = os.Environ()
	cmd.Dir = projectPath
	g.Command(cmd)
}

var buildProtoCmd = &cobra.Command{
	Use: "build-proto",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dryRun")
		r := genny.WetRunner(context.Background())
		if dryRun {
			r = genny.DryRunner(context.Background())
		}

		projectPath, _ := cmd.Flags().GetString("projectPath")

		g := genny.New()
		compileProto(g, projectPath)

		// Create a Runner with the Generator customized by command's arguments
		err := r.With(g)
		if err != nil {
			return err
		}
		return r.Run()
	},
}

func compileProto(g *genny.Generator, projectPath string) {
	// Compile proto files
	folders := []string{"service", "model"}
	for _, f := range folders {
		_ = filepath.Walk(filepath.Join(projectPath, f), func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) == ".proto" {
				projectPathAbs, _ := filepath.Abs(projectPath)
				cmd1 := exec.Command(
					"protoc",
					fmt.Sprintf("-I=%s", projectPathAbs),
					fmt.Sprintf("--go_out=%s", projectPathAbs),
					fmt.Sprintf(path),
				)
				cmd1.Env = os.Environ()
				cmd1.Dir = filepath.Dir(projectPathAbs)
				g.Command(cmd1)

				cmd2 := exec.Command(
					"protoc",
					fmt.Sprintf("-I=%s", projectPathAbs),
					fmt.Sprintf("--gorony_out=%s", projectPathAbs),
					fmt.Sprintf(path),
				)
				cmd2.Env = os.Environ()
				cmd2.Dir = filepath.Dir(projectPathAbs)
				g.Command(cmd2)
			}
			return nil
		})
	}

}
