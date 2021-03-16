package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
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
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
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

// BindCmdFlags binds a full flag set to the configuration, using each flag's long
// name as the config key.
func BindCmdFlags(cmd *cobra.Command) error {
	return _Viper.BindPFlags(cmd.Flags())
}

// SetCmdFlags sets the flags defined by FlagOption to the command 'cmd'.
func SetCmdFlags(cmd *cobra.Command, opts ...FlagOption) {
	for _, fo := range opts {
		fo(cmd.Flags())
	}
}

// SetCmdPersistentFlags sets the persistent flags defined by FlagOption to the command 'cmd'.
func SetCmdPersistentFlags(cmd *cobra.Command, opts ...FlagOption) {
	for _, fo := range opts {
		fo(cmd.PersistentFlags())
	}
}

// Set sets the value for the key in the override register.
// Set is case-insensitive for a key.
// Will be used instead of values obtained via
// flags, config file, ENV, default, or key/value store.
func Set(key string, val interface{}) {
	_Viper.Set(key, val)
}

func GetBool(key string) bool {
	return _Viper.GetBool(key)
}

func GetString(key string) string {
	return _Viper.GetString(key)
}

func GetInt64(key string) int64 {
	return _Viper.GetInt64(key)
}

func GetUint64(key string) uint64 {
	return _Viper.GetUint64(key)
}

func GetInt32(key string) int32 {
	return _Viper.GetInt32(key)
}

func GetUint32(key string) uint32 {
	return _Viper.GetUint32(key)
}

func GetInt(key string) int {
	return _Viper.GetInt(key)
}

func GetUint(key string) uint {
	return _Viper.GetUint(key)
}

func GetIntSlice(key string) []int {
	return _Viper.GetIntSlice(key)
}

func GetInt32Slice(key string) []int32 {
	ints := _Viper.GetIntSlice(key)
	int32s := make([]int32, 0, len(ints))
	for _, i := range ints {
		int32s = append(int32s, int32(i))
	}
	return int32s
}

func GetInt64Slice(key string) []int64 {
	ints := _Viper.GetIntSlice(key)
	int64s := make([]int64, 0, len(ints))
	for _, i := range ints {
		int64s = append(int64s, int64(i))
	}
	return int64s
}

func GetStringSlice(key string) []string {
	return _Viper.GetStringSlice(key)
}

func GetDuration(key string) time.Duration {
	return _Viper.GetDuration(key)
}

func GetTime(key string) time.Time {
	return _Viper.GetTime(key)
}
