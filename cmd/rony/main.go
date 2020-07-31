package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/cmd/rony/plugins/rpc"
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
				protoParseCmd := getProtocCmd([]string{fmt.Sprintf(
					"-I=%s", dir),
					"--include_imports",
					"--include_source_info",
					"--descriptor_set_out=/dev/stdout",
					path,
				})
				protoCompileCmd := getProtocCmd([]string{"--descriptor_set_in=/dev/stdin"})

				protoParseCmd.Stdout = os.Stdin
				protoCompileCmd.Stdin = os.Stdout

				err = protoParseCmd.Start()
				panicOnErr(err, "ProtoParser Start")
				err = protoCompileCmd.Start()
				panicOnErr(err, "ProtoCompiler Start")

				go generate()

				_ = protoParseCmd.Wait()
				_ = protoCompileCmd.Wait()

				return nil
			})
		}

	},
}

func panicOnErr(err error, txt string) {
	if err != nil {
		panic(errors.Wrap(err, txt))
	}
}

func getProtocCmd(args []string) *exec.Cmd {
	var protocArgs []string
	protocArgs = append(protocArgs, args...)
	protocCmd := exec.Command("protoc", protocArgs...)
	fmt.Println("Check! ", protocCmd.String())
	protocCmd.Env = os.Environ()
	protocCmd.Stdout = nil
	return protocCmd
}

func generate() {
	req := command.Read()
	files := req.GetProtoFile()
	files = vanity.FilterFiles(files, vanity.NotGoogleProtobufDescriptorProto)

	vanity.ForEachFile(files, vanity.TurnOnMarshalerAll)
	vanity.ForEachFile(files, vanity.TurnOnSizerAll)
	vanity.ForEachFile(files, vanity.TurnOnUnmarshalerAll)

	command.Write(command.GeneratePlugin(req, &rpc.Generator{}, ".rony.go"))
}

// func generate(input io.Reader, output io.Writer, fileSuffix string) error {
// 	g := generator.New()
// 	data, err := ioutil.ReadAll(input)
// 	if err != nil {
// 		g.Error(err, "reading input")
// 	}
//
// 	if err := proto.Unmarshal(data, g.Request); err != nil {
// 		g.Error(err, "parsing input proto")
// 	}
//
// 	if len(g.Request.FileToGenerate) == 0 {
// 		g.Fail("no files to generate")
// 	}
//
// 	g.WrapTypes()
// 	g.SetPackageNames()
// 	g.BuildTypeNameMap()
// 	g.GeneratePlugin(&pools.Generator{})
//
// 	for i := 0; i < len(g.Response.File); i++ {
// 		g.Response.File[i].Name = proto.String(
// 			strings.Replace(*g.Response.File[i].Name, ".pb.go", fileSuffix, -1),
// 		)
// 	}
//
// 	if err := goFormat(g.Response); err != nil {
// 		g.Error(err)
// 	}
//
// 	// Send back the results.
// 	data, err = proto.Marshal(g.Response)
// 	if err != nil {
// 		g.Error(err, "failed to marshal output proto")
// 	}
//
// 	_, err = output.Write(data)
// 	if err != nil {
// 		g.Error(err, "failed to write output proto")
// 	}
//
// 	return err
// }

// func goFormat(resp *plugin.CodeGeneratorResponse) error {
// 	for i := 0; i < len(resp.File); i++ {
// 		formatted, err := format.Source([]byte(resp.File[i].GetContent()))
// 		if err != nil {
// 			return fmt.Errorf("go format error: %v", err)
// 		}
// 		fmts := string(formatted)
// 		resp.File[i].Content = &fmts
// 	}
// 	return nil
// }
