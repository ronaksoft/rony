package gen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

/*
   Creation Time: 2020 - Apr - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Generator struct {
	b1          []byte
	b2          []byte
	buf         *bytes.Buffer
	indent      string
	imports     map[string]struct{}
	outputDir   string
	packageName string
	prefix      string
	postfix     string
}

func (g *Generator) Pns(str ...interface{}) {
	g.p(0, str...)
}

func (g *Generator) P(str ...interface{}) {
	g.p(' ', str...)
}

func (g *Generator) p(sep rune, str ...interface{}) {
	g.buf.WriteString(g.indent)
	for _, s := range str {
		switch x := s.(type) {
		case int:
			g.buf.WriteString(strconv.Itoa(x))
		case uint:
			g.buf.WriteString(strconv.Itoa(int(x)))
		case int32:
			g.buf.WriteString(strconv.Itoa(int(x)))
		case int64:
			g.buf.WriteString(strconv.Itoa(int(x)))
		case uint32:
			g.buf.WriteString(strconv.Itoa(int(x)))
		case uint64:
			g.buf.WriteString(strconv.Itoa(int(x)))
		case string:
			g.buf.WriteString(x)
		case []byte:
			g.buf.Write(x)
		default:
			fmt.Println(x, reflect.TypeOf(x))
		}
		if sep > 0 {
			g.buf.WriteRune(sep)
		}
	}
	g.buf.WriteRune('\n')
}

func (g *Generator) Pf(format string, a ...interface{}) {
	g.buf.WriteString(g.indent)
	g.buf.WriteString(fmt.Sprintf(format, a...))
	g.buf.WriteRune('\n')
}

func (g *Generator) Nl(n ...int) {
	if len(n) > 0 {
		for i := 0; i < n[0]; i++ {
			g.buf.WriteRune('\n')
		}
	} else {
		g.buf.WriteRune('\n')
	}
}

func (g *Generator) In(n ...int) {
	if len(n) > 0 {
		for i := 0; i < n[0]; i++ {
			g.indent += "\t"
		}
	} else {
		g.indent += "\t"
	}

}

func (g *Generator) Out(n ...int) {
	if len(n) > 0 {
		for i := 0; i < n[0]; i++ {
			if len(g.indent) > 0 {
				g.indent = g.indent[1:]
			}
		}
	} else {
		if len(g.indent) > 0 {
			g.indent = g.indent[1:]
		}
	}

}

func (g *Generator) AddImport(importPath string) {
	g.imports[importPath] = struct{}{}
}

// GenerateModels takes the output directory and a package name which all the models will be generated in that
// package. User can also use prefix and postfix to create filenames i.e. the output file name will be '[prefix][model_name][postfix].go'
func Generate(p Plugin, outputDir string, packageName, prefix, postfix string) error {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}

	for _, desc := range descriptors {
		g := &Generator{
			b1:          make([]byte, 0, 1<<16),
			b2:          make([]byte, 0, 1<<16),
			outputDir:   outputDir,
			packageName: packageName,
			prefix:      prefix,
			postfix:     postfix,
		}
		g.buf = bytes.NewBuffer(g.b2)
		p.Init(g)
		p.Generate(&desc)
		g.b2 = g.buf.Bytes()
		g.buf = bytes.NewBuffer(g.b1)
		p.GeneratePrepend(&desc)
		g.b1 = g.buf.Bytes()
		fname := filepath.Join(outputDir, fmt.Sprintf("%s%s%s", prefix, strings.ToLower(desc.Name), postfix))
		f, err := os.Create(fname)
		if err != nil {
			return err
		}

		_, err = f.Write(g.b1)
		if err != nil {
			return err
		}
		_, err = f.Write(g.b2)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
