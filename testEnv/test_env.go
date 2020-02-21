package testEnv

import (
	"git.ronaksoftware.com/ronak/rony/config"
	"git.ronaksoftware.com/ronak/rony/db/redis"
	"git.ronaksoftware.com/ronak/rony/db/scylla"
	log "git.ronaksoftware.com/ronak/rony/logger"
	"git.ronaksoftware.com/ronak/rony/metrics"
	"path"
	"runtime"
)

/*
   Creation Time: 2019 - Oct - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var _Initialized bool

func Init() {
	if !_Initialized {
		_Initialized = true
	} else {
		return
	}
	_, filename, _, _ := runtime.Caller(0)
	config.Read("testConfig", path.Dir(filename))
	log.InitLogger(log.WarnLevel, "")
	metrics.Run("Test", "01", 2374)
	dbConf := scylla.DefaultConfig
	dbConf.Host = config.GetString(config.DbScyllaDsn)
	dbConf.Keyspace = "river"
	scylla.Init(dbConf)

	redisConf := redis.DefaultConfig
	redisConf.Host = config.GetString(config.RedisPermDsn)
	redisConf.Password = "ehsan2374"
	redis.InitPermCache(redis.New(redisConf))
	redis.InitTempCache(redis.New(redisConf))
	redis.InitCounterCache(redis.New(redisConf))

}
