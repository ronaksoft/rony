package di

import (
	"go.uber.org/dig"
)

var (
	di = dig.New()
)

func Provide(constructor interface{}) error {
	return di.Provide(constructor)
}

func MustProvide(constructor interface{}) {
	err := di.Provide(constructor)
	if err != nil {
		panic(err)
	}
}

func Invoke(function interface{}) error {
	return di.Invoke(function)
}

func MustInvoke(function interface{}) {
	err := di.Invoke(function)
	if err != nil {
		panic(err)
	}
}