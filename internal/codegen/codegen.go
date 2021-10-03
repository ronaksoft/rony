package codegen

import "google.golang.org/protobuf/compiler/protogen"

//go:generate go run update_version.go

type (
	Generator interface {
		Generate(out *protogen.GeneratedFile, in ...*protogen.File) error
	}
	GeneratorFunc func(out *protogen.GeneratedFile, in ...*protogen.File) error
)

func (f GeneratorFunc) Generate(out *protogen.GeneratedFile, in ...*protogen.File) error {
	return f(out, in...)
}
