//go:build ignore
// +build ignore

package main

import (
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

/*
   Creation Time: 2019 - Nov - 30
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func main() {
	f, err := os.Create("version.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var riverVer, gitVer string
	b, err := exec.Command("git", "--no-pager", "describe", "--abbrev=0").Output()
	riverVer = strings.TrimSpace(string(b))
	b, err = exec.Command("git", "--no-pager", "log", "--format=\"%H\"", "-n", "1").Output()
	gitVer = strings.TrimSpace(string(b))

	packageTemplate.Execute(f, struct {
		Timestamp     time.Time
		GitCommit     string
		GitVersionTag string
	}{
		Timestamp:     time.Now(),
		GitCommit:     gitVer,
		GitVersionTag: riverVer,
	})
}

var packageTemplate = template.Must(template.New("").Parse(`
// This is auto-generated code; DO NOT EDIT.
package codegen

var (
	Commit = {{ .GitCommit }}
	Version = "{{ .GitVersionTag }}"
)
`))
