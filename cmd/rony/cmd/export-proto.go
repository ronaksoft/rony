package cmd

import (
	"context"
	"fmt"
	"github.com/gobuffalo/genny/v2"
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/internal/codegen"
	"github.com/ronaksoft/rony/log"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	protoExt = ".proto"
)

var ExportProtoCmd = &cobra.Command{
	Use: "export-proto",
	RunE: func(cmd *cobra.Command, args []string) error {
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
		err := r.With(g)
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
		// generate protoc-gen-goexport
		args1 := []string{
			fmt.Sprintf("-I=%s", projectPathAbs),
			fmt.Sprintf("-I=%s", folderPathAbs),
			fmt.Sprintf("-I=%s/vendor", projectPathAbs),
			fmt.Sprintf("--goexport_out=paths=source_relative:%s", filepath.Join(projectPathAbs, "exports/proto")),
		}
		args1 = append(args1, files...)
		cmd1 := exec.Command(
			"protoc", args1...,
		)
		cmd1.Env = os.Environ()
		cmd1.Dir = filepath.Dir(projectPathAbs)
		cmd1.Stderr = os.Stderr
		g.Command(cmd1)

		// generate protoc-gen-gorony
		args2 := []string{
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
	}
}

func collectProto(g *genny.Generator, folders []string) {
	var (
		files             []string
		projectPathAbs, _ = filepath.Abs(".")
		dstTempFolder     = "_exported-proto"
	)

	f, _ := Skeleton.Open("skel/proto/msg.proto")
	g.File(genny.NewFile(filepath.Join(dstTempFolder, "msg.proto"), f))
	fmt.Println("run in: ", projectPathAbs)
	for _, folder := range folders {
		files = files[:0]
		_ = filepath.Walk(filepath.Join(".", folder), func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return nil
			}
			if filepath.Ext(info.Name()) == ".proto" {
				files = append(files, path)
				cmd0 := exec.Command("mkdir", "-p", filepath.Join(dstTempFolder, filepath.Dir(path)))
				cmd0.Dir = projectPathAbs
				g.Command(cmd0)
				cmd1 := exec.Command("cp", path, filepath.Join(dstTempFolder, filepath.Dir(path)))
				cmd1.Dir = projectPathAbs
				g.Command(cmd1)
			}

			return nil
		})
	}
	cmd3 := exec.Command("tar", "-czf", fmt.Sprintf("proto.tar.gz"), dstTempFolder)
	cmd3.Env = os.Environ()
	cmd3.Dir = projectPathAbs
	cmd3.Stderr = os.Stderr
	g.Command(cmd3)
	fmt.Println(cmd3.String())

	cmd4 := exec.Command("rm", "-r", dstTempFolder)
	cmd4.Dir = projectPathAbs
	cmd4.Stderr = os.Stderr
	g.Command(cmd4)
	fmt.Println(cmd4.String())
}
