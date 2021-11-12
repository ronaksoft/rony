package c

import "github.com/ronaksoft/rony/config"

var (
	globalC *config.Config
)

func init() {
	globalC = config.Global()
}

func Conf() *config.Config {
	return globalC
}
