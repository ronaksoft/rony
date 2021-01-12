package cmd

import (
	"context"
	"fmt"
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
   Creation Time: 2020 - Aug - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	skeletonPath = "github.com/ronaksoft/rony:/internal/templates"
)

func init() {
	workingDir, _ := os.Getwd()
	_ = config.Init("", workingDir)

	config.SetPersistentFlags(RootCmd,
		config.BoolFlag("dry-run", false, "don't actually change anything"),
	)

	config.SetFlags(createCmd,
		config.StringFlag("project.dir", workingDir, "the root path of the project"),
		config.StringFlag("project.name", "sample-project", "the name of the project"),
		config.StringFlag("package.name", "github.com/sample", "the full path of go package in go.mod file"),
	)

	config.SetFlags(genProtoCmd,
		config.StringFlag("project.dir", workingDir, "the root path of the project"),
	)

	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(createCmd, genProtoCmd)

	_ = pkger.Walk(skeletonPath, func(path string, info os.FileInfo, err error) error {
		return nil
	})
}

var projectCmd = &cobra.Command{
	Use: "project",
}

var createCmd = &cobra.Command{
	Use: "create",
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
		gofmt(g)

		// Create a Runner with the Generator customized by command's arguments
		err = r.With(g)
		if err != nil {
			return err
		}
		return r.Run()
	},
}

var genProtoCmd = &cobra.Command{
	Use: "gen-proto",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())
		if config.GetBool("dry-run") {
			r = genny.DryRunner(context.Background())
		}

		g := genny.New()
		compileProto(g)
		gofmt(g)
		goModuleTidy(g)

		// Create a Runner with the Generator customized by command's arguments
		err := r.With(g)
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

func gofmt(g *genny.Generator) {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Env = os.Environ()
	cmd.Dir = config.GetString("project.dir")

	g.Command(cmd)
}

func compileProto(g *genny.Generator) {
	projectPath := config.GetString("project.dir")
	// Compile proto files
	folders := []string{"service", "model"}
	for _, f := range folders {
		_ = filepath.Walk(filepath.Join(projectPath, f), func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) == ".proto" {
				projectPathAbs, _ := filepath.Abs(projectPath)
				cmd1 := exec.Command(
					"protoc",
					fmt.Sprintf("-I=%s", projectPathAbs),
					fmt.Sprintf("--go_out=paths=source_relative:%s", projectPathAbs),
					path,
				)
				cmd1.Env = os.Environ()
				cmd1.Dir = filepath.Dir(projectPathAbs)
				g.Command(cmd1)

				cmd2 := exec.Command(
					"protoc",
					fmt.Sprintf("-I=%s", projectPathAbs),
					fmt.Sprintf("--gorony_out=paths=source_relative,plugin=server:%s", projectPathAbs),
					path,
				)
				cmd2.Env = os.Environ()
				cmd2.Dir = filepath.Dir(projectPathAbs)
				g.Command(cmd2)
			}
			return nil
		})
	}

}
