package cmd

import (
	"context"
	"fmt"
	"github.com/gobuffalo/genny/v2"
	"github.com/ronaksoft/rony/config"
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

const (
	skeletonPath = "github.com/ronaksoft/rony:/internal/templates"
)

var GenProtoCmd = &cobra.Command{
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

func compileProto(g *genny.Generator) {
	projectPath := config.GetString("project.dir")
	// Compile proto files
	folders := []string{"rpc", "model"}
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
