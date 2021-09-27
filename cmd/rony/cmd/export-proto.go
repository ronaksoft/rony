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

		g := genny.New()
		exportProto(g, args)
		// Create a Runner with the Generator customized by command's arguments
		err := r.With(g)
		if err != nil {
			return err
		}

		return r.Run()
	},
}

func exportProto(g *genny.Generator, folders []string) {
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
			if filepath.Ext(info.Name()) == ".proto" {
				files = append(files, path)
			}

			return nil
		})
		// generate protoc-gen-goexport
		args := []string{
			fmt.Sprintf("-I=%s", projectPathAbs),
			fmt.Sprintf("-I=%s", folderPathAbs),
			fmt.Sprintf("-I=%s/vendor", projectPathAbs),
			fmt.Sprintf("--goexport_out=paths=source_relative:%s", filepath.Join(projectPathAbs, "exports/proto")),
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
