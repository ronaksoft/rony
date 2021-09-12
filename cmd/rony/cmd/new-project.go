package cmd

import (
	"context"
	"embed"
	"fmt"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/plushgen"
	"github.com/gobuffalo/plush/v4"
	"github.com/ronaksoft/rony/config"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
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
	Use: "new",
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

func addFiles(g *genny.Generator, tCtx *plush.Context, fs embed.FS, fsPath, path string, outputDir string) {
	files, _ := fs.ReadDir(filepath.Join(fsPath, path))
	for _, ent := range files {
		if ent.IsDir() {
			continue
		}

		f, err := Skeleton.Open(filepath.Join(fsPath, path, ent.Name()))
		if err != nil {
			panic(err)
		}
		ft, err := plushgen.Transformer(tCtx).Transform(genny.NewFile(filepath.Join(outputDir, ent.Name()), f))
		g.File(ft)
		_ = f.Close()
	}

}

func setupSkeleton(g *genny.Generator) {
	projectPath, err := filepath.Abs(config.GetString("project.dir"))
	if err != nil {
		panic(err)
	}
	projectName := config.GetString("project.name")
	packagePath := config.GetString("package.path")

	_ = os.MkdirAll(projectPath, os.ModePerm)

	tCtx := plush.NewContext()
	tCtx.Set("projectName", func() string {
		return projectName
	})
	tCtx.Set("packagePath", func() string {
		return packagePath
	})

	addFiles(g, tCtx, Skeleton, "skel", ".", filepath.Join(projectPath, "."))

	// create cmd folder
	g.File(genny.NewDir(filepath.Join(projectPath, fmt.Sprintf("cmd/cli-%s", projectName)), os.ModeDir|0744))
	addFiles(g, tCtx, Skeleton, "skel", "cmd/cli-project", filepath.Join(projectPath, fmt.Sprintf("cmd/cli-%s", projectName)))

	// create rpc folder
	g.File(genny.NewDir(filepath.Join(projectPath, "rpc"), os.ModeDir|0744))
	addFiles(g, tCtx, Skeleton, "skel", "rpc", filepath.Join(projectPath, "rpc"))

	g.File(genny.NewDir(filepath.Join(projectPath, "conf"), os.ModeDir|0744))
	addFiles(g, tCtx, Skeleton, "skel", "conf", filepath.Join(projectPath, "conf"))
}

func goModuleInit(g *genny.Generator) {
	projectPath := config.GetString("project.dir")
	packagePath := config.GetString("package.path")

	cmd := exec.Command("go", "mod", "init", packagePath)
	cmd.Env = os.Environ()
	cmd.Dir = projectPath
	g.Command(cmd)
}

func goModuleTidy(g *genny.Generator) {
	cmd := exec.Command("go", "mod", "tidy", "-go=1.16")
	cmd.Env = os.Environ()
	cmd.Dir = config.GetString("project.dir")
	g.Command(cmd)

	cmd = exec.Command("go", "mod", "tidy", "-go=1.17")
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
