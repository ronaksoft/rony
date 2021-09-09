package main

import (
	"embed"
	"github.com/ronaksoft/rony/cmd/rony/cmd"
)

//go:embed skel/*
var skeleton embed.FS

func main() {
	cmd.Skeleton = skeleton
	_ = cmd.RootCmd.Execute()
}
