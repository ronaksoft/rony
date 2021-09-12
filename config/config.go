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

type Config struct {
	v *viper.Viper
}

var (
	globalConfig *Config
)

func init() {
	globalConfig = New()
}

func Global() *Config {
	return globalConfig
}

func New() *Config {
	return &Config{
		v: viper.NewWithOptions(
			viper.EnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
		),
	}
}

func (c *Config) ReadFile(configFileName string, configSearchPaths ...string) error {
	for _, p := range configSearchPaths {
		c.v.AddConfigPath(p)
	}
	if configFileName == "" {
		configFileName = "config"
	}

	c.v.SetConfigName(configFileName)

	return c.v.ReadInConfig()
}

func ReadFile(configFileName string, configSearchPaths ...string) error {
	return globalConfig.ReadFile(configFileName, configSearchPaths...)
}

func (c *Config) AllSettings() map[string]interface{} {
	return c.v.AllSettings()
}

func AllSettings() map[string]interface{} {
	return globalConfig.AllSettings()
}

// SetEnvPrefix defines a prefix that ENVIRONMENT variables will use.
// // E.g. if your prefix is "spf", the env registry will look for env
// // variables that start with "SPF_".
func (c *Config) SetEnvPrefix(prefix string) {
	c.v.SetEnvPrefix(prefix)
	c.v.AutomaticEnv()
}
func SetEnvPrefix(prefix string) {
	globalConfig.SetEnvPrefix(prefix)
}

// BindCmdFlags binds a full flag set to the configuration, using each flag's long
// name as the config key.
func (c *Config) BindCmdFlags(cmd *cobra.Command) error {
	return c.v.BindPFlags(cmd.Flags())
}
func BindCmdFlags(cmd *cobra.Command) error {
	return globalConfig.BindCmdFlags(cmd)
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
func (c *Config) Set(key string, val interface{}) {
	c.v.Set(key, val)
}

func Set(key string, val interface{}) {
	globalConfig.Set(key, val)
}

func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}
func GetBool(key string) bool {
	return globalConfig.GetBool(key)
}

func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}
func GetString(key string) string {
	return globalConfig.GetString(key)
}

func (c *Config) GetInt64(key string) int64 {
	return c.v.GetInt64(key)
}
func GetInt64(key string) int64 {
	return globalConfig.GetInt64(key)
}

func (c *Config) GetUint64(key string) uint64 {
	return c.v.GetUint64(key)
}
func GetUint64(key string) uint64 {
	return globalConfig.GetUint64(key)
}

func (c *Config) GetInt32(key string) int32 {
	return c.v.GetInt32(key)
}
func GetInt32(key string) int32 {
	return globalConfig.GetInt32(key)
}

func (c *Config) GetUint32(key string) uint32 {
	return c.v.GetUint32(key)
}
func GetUint32(key string) uint32 {
	return globalConfig.GetUint32(key)
}

func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}
func GetInt(key string) int {
	return globalConfig.GetInt(key)
}

func (c *Config) GetUint(key string) uint {
	return c.v.GetUint(key)
}
func GetUint(key string) uint {
	return globalConfig.GetUint(key)
}

func (c *Config) GetIntSlice(key string) []int {
	return c.v.GetIntSlice(key)
}
func GetIntSlice(key string) []int {
	return globalConfig.GetIntSlice(key)
}

func (c *Config) GetInt32Slice(key string) []int32 {
	ints := c.v.GetIntSlice(key)
	int32s := make([]int32, 0, len(ints))
	for _, i := range ints {
		int32s = append(int32s, int32(i))
	}

	return int32s
}
func GetInt32Slice(key string) []int32 {
	return globalConfig.GetInt32Slice(key)
}

func (c *Config) GetInt64Slice(key string) []int64 {
	ints := c.v.GetIntSlice(key)
	int64s := make([]int64, 0, len(ints))
	for _, i := range ints {
		int64s = append(int64s, int64(i))
	}

	return int64s
}
func GetInt64Slice(key string) []int64 {
	return globalConfig.GetInt64Slice(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}
func GetStringSlice(key string) []string {
	return globalConfig.GetStringSlice(key)
}

func (c *Config) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}
func GetDuration(key string) time.Duration {
	return globalConfig.GetDuration(key)
}

func (c *Config) GetTime(key string) time.Time {
	return c.v.GetTime(key)
}
func GetTime(key string) time.Time {
	return globalConfig.GetTime(key)
}
