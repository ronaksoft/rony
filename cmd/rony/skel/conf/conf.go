package c

import "github.com/ronaksoft/rony/config"

var (
	globalC *config.Config
)

func init() {
	globalC = config.New()
}

func Conf() *config.Config {
	return globalC
}
