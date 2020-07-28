package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/cmd/rony/plugins/pools"
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

/*
	RONY Command-line executive
	This creates project skeleton with generated codes for cli, models (scylla, aerospike, redis, ...) and boilerplate code
	to make developing new projects with Rony framework easier and faster.
*/

func main() {
	tools.SetFlags(RootCmd,
		tools.StringSliceFlag(FlagInputDirs, []string{"."}, "directory where proto files are"),
	)
	_ = RootCmd.Execute()

}

var RootCmd = &cobra.Command{
	Use: "rony",
	Run: func(cmd *cobra.Command, args []string) {
		var protocArgs []string
		for _, arg := range tools.GetStringSlice(cmd, FlagInputDirs) {
			protocArgs = append(protocArgs, "-I", arg)
		}

		for _, dir := range tools.GetStringSlice(cmd, FlagInputDirs) {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Print("ERR!", err.Error())
					return err
				}
				if info.IsDir() {
					return nil
				}

				protocArgs = append(protocArgs[:0], " --descriptor_set_out", "/dev/stdout")
				err = runProtoc(os.Stdout, os.Stdin, protocArgs, path)
				if err != nil {
					return nil
				}

				// go func() {
				g := generator.New()

				fmt.Println("Lets read descriptors from stdin")
				data, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					g.Error(err, "reading input")
				}

				if err := proto.Unmarshal(data, g.Request); err != nil {
					g.Error(err, "parsing input proto")
				}

				if len(g.Request.FileToGenerate) == 0 {
					g.Fail("no files to generate")
				}

				files := g.Request.GetProtoFile()
				files = vanity.FilterFiles(files, vanity.NotGoogleProtobufDescriptorProto)

				vanity.ForEachFile(files, vanity.TurnOnMarshalerAll)
				vanity.ForEachFile(files, vanity.TurnOnSizerAll)
				vanity.ForEachFile(files, vanity.TurnOnUnmarshalerAll)

				fmt.Println("Lets write generated code to the stdout")
				command.Write(command.GeneratePlugin(g.Request, &pools.Generator{}, ".pools.go"))
				// }()

				protocArgs = append(protocArgs[:0], " --descriptor_set_in", "/dev/stdin")
				err = runProtoc(os.Stdout, os.Stdin, protocArgs, path)
				if err != nil {
					fmt.Println("ERR! On Write", err.Error())
				}

				return nil
			})
		}

	},
}

func runProtoc(stdin, stdout *os.File, args []string, path string) error {
	var protocArgs []string
	protocArgs = append(protocArgs, args...)
	protocArgs = append(protocArgs, path)
	protocCmd := exec.Command("protoc", protocArgs...)
	fmt.Println("Check! ", path, protocCmd.String())
	protocCmd.Env = os.Environ()
	protocCmd.Stdout = stdout
	protocCmd.Stdin = stdin
	err := protocCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
