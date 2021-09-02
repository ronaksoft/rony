package di

import (
	"go.uber.org/dig"
)

type (
	In  = dig.In
	Out = dig.Out
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
