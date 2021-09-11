package cmd

import (
	"context"
	"fmt"
	"github.com/gobuffalo/genny/v2"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/log"
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

var GenProtoCmd = &cobra.Command{
	Use: "gen-proto",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())
		if config.GetBool("dry-run") {
			r = genny.DryRunner(context.Background())
		}
		r.Logger = log.DefaultLogger.Sugared()

		if len(args) == 0 {
			args = append(args, "rpc", "model")
		}

		g := genny.New()
		compileProto(g, args)
		gofmt(g)
		goModuleTidy(g)
		goModuleVendor(g)

		// Create a Runner with the Generator customized by command's arguments
		err := r.With(g)
		if err != nil {
			return err
		}

		return r.Run()
	},
}

func compileProto(g *genny.Generator, folders []string) {
	// Compile proto files
	var (
		files             []string
		projectPathAbs, _ = filepath.Abs(".")
		folderPathAbs     string
	)
	for _, folder := range folders {
		files = files[:0]
		folderPathAbs, _ = filepath.Abs(filepath.Join(".", folder))
		_ = filepath.Walk(filepath.Join(".", folder), func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) == ".proto" {
				//files = append(files, filepath.Join(folderPathAbs, filepath.Base(path)))
				files = append(files, path)
				// call protoc-gen-go
				cmd1 := exec.Command(
					"protoc",
					fmt.Sprintf("-I=%s", projectPathAbs),
					fmt.Sprintf("-I=%s", folderPathAbs),
					fmt.Sprintf("-I=%s/vendor", projectPathAbs),
					fmt.Sprintf("--go_out=paths=source_relative:%s", projectPathAbs),
					path,
				)
				cmd1.Env = os.Environ()
				cmd1.Dir = filepath.Dir(projectPathAbs)
				g.Command(cmd1)

				// generate protoc-gen-gorony
				cmd2 := exec.Command(
					"protoc",
					fmt.Sprintf("-I=%s", projectPathAbs),
					fmt.Sprintf("-I=%s", folderPathAbs),
					fmt.Sprintf("-I=%s/vendor", projectPathAbs),
					fmt.Sprintf("--gorony_out=paths=source_relative:%s", projectPathAbs),
					path,
				)
				cmd2.Env = os.Environ()
				cmd2.Dir = filepath.Dir(projectPathAbs)
				g.Command(cmd2)

			}

			return nil
		})
		// generate protoc-gen-gorony
		args := []string{
			fmt.Sprintf("-I=%s", projectPathAbs),
			fmt.Sprintf("-I=%s", folderPathAbs),
			fmt.Sprintf("-I=%s/vendor", projectPathAbs),
			fmt.Sprintf("--gorony_out=paths=source_relative,rony_opt=json_c:%s", projectPathAbs),
		}
		args = append(args, files...)
		cmd3 := exec.Command(
			"protoc", args...,
		)
		cmd3.Env = os.Environ()
		cmd3.Dir = filepath.Dir(projectPathAbs)
		cmd3.Stderr = os.Stderr
		g.Command(cmd3)
	}
}
