package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/internal/codegen"
	"github.com/ronaksoft/rony/log"
	"github.com/spf13/cobra"
)

const (
	protoExt = ".proto"
)

var ExportProtoCmd = &cobra.Command{
	Use: "export-proto",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.BindCmdFlags(cmd)
		if err != nil {
			return err
		}

		r := genny.WetRunner(context.Background())
		if config.GetBool("dry-run") {
			r = genny.DryRunner(context.Background())
		}
		r.Logger = log.DefaultLogger.Sugared()

		if len(args) == 0 {
			args = append(args, "rpc", "model")
		}

		cFormat := codegen.ConstructorFormat(config.GetString("c.format"))
		switch cFormat {
		case codegen.Int64JSON, codegen.StringJSON:
		default:
			cFormat = codegen.StringJSON
		}

		cmd.Println("Folders:", args)
		cmd.Println("Constructor Format:", cFormat)

		g := genny.New()
		exportProto(g, args, cFormat)

		// Create a Runner with the Generator customized by command's arguments
		err = r.With(g)
		if err != nil {
			return err
		}

		return r.Run()
	},
}

func exportProto(g *genny.Generator, folders []string, cFormat codegen.ConstructorFormat) {
	var (
		files             []string
		projectPathAbs, _ = filepath.Abs(".")
		folderPathAbs     string
	)

	_ = os.MkdirAll(filepath.Join(projectPathAbs, "exports/proto"), os.ModePerm|0755)
	for _, folder := range folders {
		files = files[:0]
		folderPathAbs, _ = filepath.Abs(filepath.Join(".", folder))
		_ = filepath.Walk(filepath.Join(".", folder), func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) == protoExt {
				files = append(files, path)
			}

			return nil
		})
		// generate protoc-gen-gorony (Clean Proto Files)
		args1 := []string{
			fmt.Sprintf("-I=%s/src", os.Getenv("GOPATH")),
			fmt.Sprintf("-I=%s", projectPathAbs),
			fmt.Sprintf("-I=%s", folderPathAbs),
			fmt.Sprintf("-I=%s/vendor", projectPathAbs),
			fmt.Sprintf(
				"--gorony_out=paths=source_relative,rony_opt=clean_proto:%s",
				filepath.Join(projectPathAbs, "exports/proto"),
			),
		}
		args1 = append(args1, files...)
		cmd1 := exec.Command(
			"protoc", args1...,
		)
		cmd1.Env = os.Environ()
		cmd1.Dir = filepath.Dir(projectPathAbs)
		cmd1.Stderr = os.Stderr
		g.Command(cmd1)

		// generate protoc-gen-gorony (Constructors)
		args2 := []string{
			fmt.Sprintf("-I=%s/src", os.Getenv("GOPATH")),
			fmt.Sprintf("-I=%s", projectPathAbs),
			fmt.Sprintf("-I=%s", folderPathAbs),
			fmt.Sprintf("-I=%s/vendor", projectPathAbs),
			fmt.Sprintf(
				"--gorony_out=paths=source_relative,rony_opt=json_%s:%s",
				cFormat,
				filepath.Join(projectPathAbs, "exports/proto"),
			),
		}
		args2 = append(args2, files...)
		cmd2 := exec.Command(
			"protoc", args2...,
		)
		cmd2.Env = os.Environ()
		cmd2.Dir = filepath.Dir(projectPathAbs)
		cmd2.Stderr = os.Stderr
		g.Command(cmd2)

		// generate protoc-gen-gorony (Swagger APIs)
		args3 := []string{
			fmt.Sprintf("-I=%s/src", os.Getenv("GOPATH")),
			fmt.Sprintf("-I=%s", projectPathAbs),
			fmt.Sprintf("-I=%s", folderPathAbs),
			fmt.Sprintf("-I=%s/vendor", projectPathAbs),
			fmt.Sprintf(
				"--gorony_out=paths=source_relative,rony_opt=open_api:%s",
				filepath.Join(projectPathAbs, "exports/proto"),
			),
		}
		args3 = append(args3, files...)
		cmd3 := exec.Command(
			"protoc", args3...,
		)
		cmd3.Env = os.Environ()
		cmd3.Dir = filepath.Dir(projectPathAbs)
		cmd3.Stderr = os.Stderr
		g.Command(cmd3)
	}
}
