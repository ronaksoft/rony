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

func Provide(constructor interface{}, opt ...ProvideOption) error {
	return di.Provide(constructor, opt...)
}

func MustProvide(constructor interface{}, opt ...ProvideOption) {
	if err := di.Provide(constructor, opt...); err != nil {
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
