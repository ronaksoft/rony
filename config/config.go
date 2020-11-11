package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

/*
   Creation Time: 2020 - Nov - 11
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	_Viper *viper.Viper
)

func MustInit(configFileName string, configSearchPaths ...string) {
	err := Init(configFileName, configSearchPaths...)
	if err != nil {
		panic(err)
	}
}

func Init(configFileName string, configSearchPaths ...string) error {
	_Viper = viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)
	for _, p := range configSearchPaths {
		_Viper.AddConfigPath(p)
	}
	if configFileName == "" {
		configFileName = "config"
	}

	_Viper.SetConfigName(configFileName)
	return _Viper.ReadInConfig()
}

// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use.
// // E.g. if your prefix is "spf", the env registry will look for env
// // variables that start with "SPF_".
func SetEnvPrefix(prefix string) {
	_Viper.SetEnvPrefix(prefix)
	_Viper.AutomaticEnv()
}

// Set sets the value for the key in the override register.
// Set is case-insensitive for a key.
// Will be used instead of values obtained via
// flags, config file, ENV, default, or key/value store.
func Set(key string, val interface{}) {
	_Viper.Set(key, val)
}

// BindCommandFlags binds a full flag set to the configuration, using each flag's long
// name as the config key.
func BindCommandFlags(cmd *cobra.Command) error {
	return _Viper.BindPFlags(cmd.Flags())
}

// EnableCommandFlags sets the flags defined by FlagOption to the command 'cmd'.
func EnableCommandFlags(cmd *cobra.Command, opts ...FlagOption) {
	for _, fo := range opts {
		fo(cmd.Flags())
	}
}
