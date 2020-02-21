package config

import (
	"path"
	"runtime"
	"strings"
	"time"

	"git.ronaksoftware.com/ronak/rony/gateway"
	"github.com/spf13/viper"
)

/*
   Creation Time: 2018 - Apr - 16
   Created by:  (ehsan)
   Maintainers:
       1.  (ehsan)
   Auditor: Ehsan N. Moosa
   Copyright Ronak Software Group 2018
*/

//go:generate go run update_version.go
const (
	// General Config Parameters
	BundleID             = "bundle.id"
	InstanceID           = "instance.id"
	TestMode             = "test.mode"
	LogLevel             = "log.level"
	LogFilePath          = "log.file.path"
	SentryDSN            = "log.sentry.dsn"
	ProfilerEnabled      = "profiler.enabled"
	ProfilerPort         = "profiler.port"
	MetricsPort          = "metrics.port"
	MaxConcurrency       = "max.concurrency"
	NewConnectionWorkers = "new.connection.workers"
	// PrivateKeyPoolSize defines the number of prepared generated private keys used for InitConnect
	// and InitCompleteAuth
	PrivateKeyPoolSize = "private.key.pool.size"

	// Gateways
	GatewayClient      = "gateway.client"
	GatewayProtocol    = "gateway.protocol"
	GatewayListenAddr  = "gateway.listen"
	GatewayTlsCertPath = "gateway.tls.cert.path"
	GatewayTlsKeyPath  = "gateway.tls.key.path"

	// External Connections
	NatsURL                 = "nats.url"
	NatsUser                = "nats.user"
	NatsPass                = "nats.pass"
	NatsClusterID           = "nats.cluster.id"
	NatsTimeout             = "nats.timeout"
	NatsRetries             = "nats.retries"
	RedisPermDsn            = "redis.perm.dsn"
	RedisPermPass           = "redis.perm.pass"
	RedisPermPoolSize       = "redis.perm.pool.size"
	RedisPermPoolMaxSize    = "redis.perm.pool.max.size"
	RedisTempDsn            = "redis.temp.dsn"
	RedisTempPass           = "redis.temp.pass"
	RedisTempPoolSize       = "redis.temp.pool.size"
	RedisTempPoolMaxSize    = "redis.temp.pool.max.size"
	RedisCounterDsn         = "redis.counter.dsn"
	RedisCounterPass        = "redis.counter.pass"
	RedisCounterPoolSize    = "redis.counter.pool.size"
	RedisCounterPoolMaxSize = "redis.counter.pool.max.size"
	DbScyllaDsn             = "db.scylla.dsn"
	DbScyllaUser            = "db.scylla.user"
	DbScyllaPass            = "db.scylla.pass"
	DbScyllaKeyspace        = "db.scylla.keyspace"
	DbScyllaConnections     = "db.scylla.connections"
	DbScyllaTracer          = "db.scylla.tracer"
	DbScyllaTracerFilePath  = "db.scylla.tracer.file.path"
	MongoDSN                = "mongo.dsn"
	MongoUser               = "mongo.user"
	MongoPass               = "mongo.path"
	MongoDBName             = "mongo.db.name"
	MongoTLS                = "mongo.tls"

	// Salt Manager
	SaltManagerSafetyMargin     = "salt.manager.safety.margin"
	SaltManagerRetrieveInterval = "salt.manager.retrieve.interval"

	// SMS Panel
	SmsEnabled                   = "sms.enabled"
	AdpMessageUrl                = "adp.message.url"
	AdpUsername                  = "adp.username"
	AdpPassword                  = "adp.password"
	AdpPhone                     = "adp.phone"
	AsanPardakhtDomesticHost     = "ap.domestic.host"
	AsanPardakhtDomesticUsername = "ap.domestic.username"
	AsanPardakhtDomesticPassword = "ap.domestic.password"
	AsanPardakhtIntlHost         = "ap.intl.host"
	AsanPardakhtIntlUsername     = "ap.intl.username"
	AsanPardakhtIntlPassword     = "ap.intl.password"

	// Misc
	NotificationServer   = "notification.server"
	StorageClusterID     = "storage.cluster.id"
	StorageBackend       = "storage.backend"
	StorageBackendDBName = "storage.backend.db.name"
	QuarkPath            = "quark.path"

	// Bcrypt
	BcryptCost = "bcrypt.cost"

	// JWT
	JWTSecret   = "jwt.secret"
	JWTIssuer   = "jwt.issuer"
	JWTDuration = "jwt.duration"

	// Bot
	BotEdgeAddr = "bot.edge.addr"
	BotAuthKey  = "bot.auth.key"
	BotID       = "bot.id"

	// Admin
	AdminServerAddr      = "admin.server.addr"
	AdminServerURL       = "admin.server.url"
	AdminServerStaticDir = "admin.server.static.dir"
	AdminRootUser        = "admin.root.user"
	AdminRootPass        = "admin.root.pass"

	// Notification
	NotificationServerAddr  = "notification.server.addr"
	NotificationSendDelay   = "notification.send.delay"
	NotificationConcurrency = "notification.concurrency"
)

func Read(configName string, configPaths ...string) {
	viper.SetEnvPrefix("RIVER")
	viper.AutomaticEnv()

	// Use . for accessing nested objects,
	// so we can read config from environment variables
	// and config files the same way.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Check config dir no matter who imports this package.
	_, filename, _, _ := runtime.Caller(0)
	viper.AddConfigPath(path.Dir(filename))
	viper.AddConfigPath("/ronak/config")
	viper.AddConfigPath(".")

	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}
	if configName == "" {
		configName = "config"
	}
	viper.SetConfigName(configName)
	viper.SetDefault(ProfilerEnabled, true)
	viper.SetDefault(ProfilerPort, 6060)
	viper.SetDefault(MetricsPort, 7070)
	viper.SetDefault(BundleID, "API-EDGE")
	viper.SetDefault(InstanceID, "01")
	viper.SetDefault(TestMode, true)
	viper.SetDefault(LogLevel, 0)
	viper.SetDefault(LogFilePath, ".")
	viper.SetDefault(SentryDSN, "https://8c32490103354222ae89a015b1a1a555@sentry.ronaksoftware.com/8")
	viper.SetDefault(MaxConcurrency, 1000)
	viper.SetDefault(NewConnectionWorkers, 10)
	viper.SetDefault(PrivateKeyPoolSize, 100)

	// Gateways
	viper.SetDefault(GatewayProtocol, string(gateway.Websocket))
	viper.SetDefault(GatewayListenAddr, "0.0.0.0:8080")

	// DB and Cache and NATS
	viper.SetDefault(NatsURL, "nats://broker-nats:4222")
	viper.SetDefault(NatsUser, "ehsan")
	viper.SetDefault(NatsPass, "ehsan2374")
	viper.SetDefault(NatsClusterID, "test-cluster")
	viper.SetDefault(NatsTimeout, "5s")
	viper.SetDefault(NatsRetries, 10)

	viper.SetDefault(DbScyllaDsn, "db-cass.river")
	viper.SetDefault(DbScyllaUser, "ehsan")
	viper.SetDefault(DbScyllaPass, "ehsan2374")
	viper.SetDefault(DbScyllaKeyspace, "river")
	viper.SetDefault(DbScyllaConnections, 5)
	viper.SetDefault(DbScyllaTracer, false)
	viper.SetDefault(DbScyllaTracerFilePath, "/ronak/logs/scylla-trace.log")

	viper.SetDefault(MongoDSN, "localhost:27001")
	viper.SetDefault(MongoUser, "")
	viper.SetDefault(MongoPass, "")
	viper.SetDefault(MongoDBName, "river")
	viper.SetDefault(MongoTLS, true)

	viper.SetDefault(RedisPermDsn, "cache-redis.river:6379")
	viper.SetDefault(RedisPermPass, "ehsan2374")
	viper.SetDefault(RedisPermPoolSize, 10)
	viper.SetDefault(RedisPermPoolMaxSize, 500)
	viper.SetDefault(RedisTempDsn, "cache-redis.river:6379")
	viper.SetDefault(RedisTempPass, "ehsan2374")
	viper.SetDefault(RedisTempPoolSize, 10)
	viper.SetDefault(RedisTempPoolMaxSize, 500)
	viper.SetDefault(RedisCounterDsn, "cache-redis.river:6379")
	viper.SetDefault(RedisCounterPass, "ehsan2374")
	viper.SetDefault(RedisCounterPoolSize, 10)
	viper.SetDefault(RedisCounterPoolMaxSize, 500)

	// Salt Manager
	viper.SetDefault(SaltManagerSafetyMargin, 1*time.Minute)
	viper.SetDefault(SaltManagerRetrieveInterval, 10*time.Second)

	// SMS Panels
	viper.SetDefault(SmsEnabled, false)
	viper.SetDefault(AdpUsername, "")
	viper.SetDefault(AdpPassword, "E2e2374k19743")
	viper.SetDefault(AdpMessageUrl, "https://ws.adpdigital.com/url/send")
	viper.SetDefault(AdpPhone, "98200049112")
	viper.SetDefault(AsanPardakhtDomesticHost, "http://192.168.72.12:8021/v1/SMS")
	viper.SetDefault(AsanPardakhtDomesticUsername, "ron@kkz")
	viper.SetDefault(AsanPardakhtDomesticPassword, "mezebGG@er")
	viper.SetDefault(AsanPardakhtIntlUsername, "")
	viper.SetDefault(AsanPardakhtIntlPassword, "")
	viper.SetDefault(AsanPardakhtIntlHost, "")

	// Misc
	viper.SetDefault(NotificationServer, "http://notification-001:8080")
	viper.SetDefault(StorageClusterID, 1)
	viper.SetDefault(StorageBackend, "gridfs")
	viper.SetDefault(StorageBackendDBName, "river")
	viper.SetDefault(QuarkPath, "/ronak/plugins/quark.so")

	// Bcrypt
	viper.SetDefault(BcryptCost, 10)

	// JWT
	viper.SetDefault(JWTSecret, "secret")
	viper.SetDefault(JWTIssuer, "river")
	viper.SetDefault(JWTDuration, 8*time.Hour)

	// Bot
	viper.SetDefault(BotEdgeAddr, "192.168.1.22:8080")
	viper.SetDefault(BotAuthKey, "authkey")
	viper.SetDefault(BotID, "id")

	// Admin
	viper.SetDefault(AdminServerAddr, ":8000")
	viper.SetDefault(AdminServerURL, "localhost:8000")
	viper.SetDefault(AdminServerStaticDir, "static")
	viper.SetDefault(AdminRootUser, "root")
	viper.SetDefault(AdminRootPass, "pwd")

	// Notification
	viper.SetDefault(NotificationServerAddr, ":8085")
	viper.SetDefault(NotificationSendDelay, "5s")
	viper.SetDefault(NotificationConcurrency, 10)

	_ = viper.ReadInConfig()
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func GetInt32(key string) int32 {
	return viper.GetInt32(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
