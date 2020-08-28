package main

import (
	"git.ronaksoft.com/ronak/rony/cmd/rony/cmd"
)

func main() {
	_ = cmd.RootCmd.Execute()
	// var (
	// 	projectName string = "proj1"
	// )
	//
	// // Make Dir
	// g := genny.New()
	// cmd := exec.Command("mkdir", "-p", filepath.Join(projectName, "models"))
	// cmd.Env = os.Environ()
	// g.Command(cmd)
	//
	//
	// // Main
	// f, err := pkger.Open("git.ronaksoft.com/ronak/rony:/internal/templates/main.go")
	// if err != nil {
	// 	panic(fmt.Sprintf("Main: %v", err))
	// }
	// g.File(genny.NewFile(filepath.Join(projectName, "main.go"), f))
	// _ = f.Close()
	//
	// // Server File
	// f, err = pkger.Open("git.ronaksoft.com/ronak/rony:/internal/templates/server.go")
	// if err != nil {
	// 	panic(fmt.Sprintf("Server: %v", err))
	// }
	// g.File(genny.NewFile(filepath.Join(projectName, "server.go"), f))
	// _ = f.Close()
	//
	// // Wet Runner
	// r := genny.WetRunner(context.Background())
	//
	// err = r.With(g)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// err = r.Run()
	// if err != nil {
	// 	panic(err)
	// }

}
