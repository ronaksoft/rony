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

func Invoke(function interface{}) error {
	return di.Invoke(function)
}

func MustInvoke(function interface{}) {
	if err := di.Invoke(function); err != nil {
		panic(err)
	}
}
