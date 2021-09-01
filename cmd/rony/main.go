package main

import (
	"github.com/ronaksoft/rony/cmd/rony/cmd"
	"github.com/ronaksoft/rony/log"
)

func main() {
	log.Init(log.DefaultConfig)
	_ = cmd.RootCmd.Execute()
}
