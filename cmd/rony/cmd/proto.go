package cmd

import (
	"context"
	"fmt"
	"git.ronaksoft.com/ronak/rony/tools"
	"github.com/gobuffalo/genny/v2"
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
	tools.SetFlags(protoBuildCmd,
		tools.StringFlag("protoPath", "", "the dir path of the proto files"),
		tools.StringFlag("importPath", "", "the dir path of the import proto files"),
		tools.StringFlag("outPath", "", "the dir path of the go files generated from proto files"),
	)
	RootCmd.AddCommand(protoCmd)
	protoCmd.AddCommand(protoBuildCmd)
}

var protoCmd = &cobra.Command{
	Use: "proto",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var protoBuildCmd = &cobra.Command{
	Use: "build",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dryRun")
		r := genny.WetRunner(context.Background())
		if dryRun {
			r = genny.DryRunner(context.Background())
		}

		protoPath, _ := cmd.Flags().GetString("protoPath")
		importPath, _ := cmd.Flags().GetString("importPath")
		outPath, _ := cmd.Flags().GetString("outPath")
		if outPath == "" {
			outPath = protoPath
		}

		g := genny.New()
		compileProto(g, importPath, protoPath, outPath)

		// Create a Runner with the Generator customized by command's arguments
		err := r.With(g)
		if err != nil {
			return err
		}
		return r.Run()
	},
}

func compileProto(g *genny.Generator, importPath, protoPath, outPath string) {
	// protoc  -I=./testdata  --go_out=./testdata ./testdata/*.proto
	// protoc  -I=./testdata  --gorony_out=./testdata ./testdata/*.proto
	s, err := os.Stat(protoPath)
	if err != nil {
		panic(err)
	}
	if s.IsDir() {
		filepath.Walk(protoPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) == ".proto" {
				cmd1 := exec.Command(
					"protoc",
					fmt.Sprintf("-I=%s", importPath),
					fmt.Sprintf("--go_out=%s", outPath),
					fmt.Sprintf(protoPath),
				)
				cmd1.Env = os.Environ()
				cmd1.Dir = filepath.Dir(protoPath)
				g.Command(cmd1)

				cmd2 := exec.Command(
					"protoc",
					fmt.Sprintf("-I=%s", importPath),
					fmt.Sprintf("--gorony_out=%s", outPath),
					fmt.Sprintf(protoPath),
				)
				cmd2.Env = os.Environ()
				cmd2.Dir = filepath.Dir(protoPath)
				g.Command(cmd2)

			}
			return nil
		})
	} else {
		if filepath.Ext(protoPath) == ".proto" {
			cmd := exec.Command(
				"protoc",
				fmt.Sprintf("-I=%s", importPath),
				fmt.Sprintf("--go_out=%s", outPath),
				fmt.Sprintf(protoPath),
			)
			cmd.Env = os.Environ()
			cmd.Dir = filepath.Dir(protoPath)
			g.Command(cmd)
		}

	}

}
