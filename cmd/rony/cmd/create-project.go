package cmd

import (
	"context"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/markbates/pkger"
	"github.com/ronaksoft/rony/config"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
   Creation Time: 2021 - Jan - 30
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var CreateProjectCmd = &cobra.Command{
	Use: "create-project",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.BindCmdFlags(cmd)
		if err != nil {
			return err
		}
		r := genny.WetRunner(context.Background())
		if config.GetBool("dry-run") {
			r = genny.DryRunner(context.Background())
		}

		g := genny.New()
		setupSkeleton(g)
		goModuleInit(g)
		goModuleTidy(g)
		goModuleVendor(g)
		gofmt(g)

		// Create a Runner with the Generator customized by command's arguments
		err = r.With(g)
		if err != nil {
			return err
		}

		return r.Run()
	},
}

func setupSkeleton(g *genny.Generator) {
	projectPath := config.GetString("project.dir")
	projectName := config.GetString("project.name")

	_ = os.Mkdir(projectPath, os.ModePerm)

	pathPrefix := skeletonPath + "/skel"
	err := pkger.Walk(pathPrefix, func(path string, info os.FileInfo, err error) error {
		realPath := strings.TrimSuffix(strings.TrimPrefix(path, pathPrefix), ".tpl")
		if info.IsDir() {
			g.File(genny.NewDir(filepath.Join(projectPath, realPath), os.ModeDir|0744))
		} else {
			f, err := pkger.Open(path)
			if err != nil {
				return err
			}
			tCtx := plush.NewContext()
			tCtx.Set("projectName", func() string {
				return projectName
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
func goModuleInit(g *genny.Generator) {
	projectPath := config.GetString("project.dir")
	packageName := config.GetString("package.name")

	cmd := exec.Command("go", "mod", "init", packageName)
	cmd.Env = os.Environ()
	cmd.Dir = projectPath
	g.Command(cmd)
}
func goModuleTidy(g *genny.Generator) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Env = os.Environ()
	cmd.Dir = config.GetString("project.dir")
	g.Command(cmd)
}
func goModuleVendor(g *genny.Generator) {
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Env = os.Environ()
	cmd.Dir = config.GetString("project.dir")
	g.Command(cmd)
}
func gofmt(g *genny.Generator) {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Env = os.Environ()
	cmd.Dir = config.GetString("project.dir")

	g.Command(cmd)
}
