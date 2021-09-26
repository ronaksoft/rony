package di

import (
	"go.uber.org/dig"
)

var (
	di = dig.New()
)

func New() *Container {
	return dig.New()
}

func Provide(constructor interface{}) error {
	return di.Provide(constructor)
}

func MustProvide(constructor interface{}) {
	if err := di.Provide(constructor); err != nil {
		panic(err)
	}
}

func NamedProvide(constructor interface{}, name string) error {
	return di.Provide(constructor, dig.Name(name))
}

func MustNamedProvide(constructor interface{}, name string) {
	if err := di.Provide(constructor, dig.Name(name)); err != nil {
		panic(err)
	}
}

func Invoke(function interface{}) error {
	return di.Invoke(function)
}

func MustInvoke(function interface{}) {
	if err := di.Invoke(function); err != nil {
		panic(err)
	}
}
