package scylla

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"strings"
	"time"
)

/*
   Creation Time: 2019 - Sep - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	_Session *gocql.Session
)

func Init(config Config) {
	cassCluster := gocql.NewCluster(strings.Split(config.Host, ",")...)
	retryPolicy := new(gocql.ExponentialBackoffRetryPolicy)
	retryPolicy.NumRetries = config.Retries
	retryPolicy.Min = config.RetryMinBackOff
	retryPolicy.Max = config.RetryMaxBackOff
	cassCluster.RetryPolicy = retryPolicy
	cassCluster.ConnectTimeout = config.ConnectTimeout
	cassCluster.Timeout = config.Timeout
	cassCluster.ReconnectInterval = config.ReconnectInterval
	cassCluster.DefaultIdempotence = config.DefaultIdempotence
	cassCluster.QueryObserver = config.QueryObserver
	cassCluster.Compressor = gocql.SnappyCompressor{}
	cassCluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}

	cassCluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cassCluster.NumConns = config.Concurrency
	cassCluster.PageSize = config.PageSize
	cassCluster.WriteCoalesceWaitTime = config.WriteCoalesceWaitTime
	if len(config.Keyspace) == 0 {
		panic("keyspace is not set")
	}
	cassCluster.QueryObserver = queryObserver{}
	cassCluster.Keyspace = config.Keyspace
	cassCluster.Consistency = config.Consistency
	cassCluster.SerialConsistency = config.SerialConsistency
	cassCluster.CQLVersion = config.CqlVersion
	session, err := cassCluster.CreateSession()
	if err != nil {
		panic(err)
	}
	_Session = session
	initDbQueries()
}

func Session() *gocql.Session {
	return _Session
}

const (
	speculativeAttempts = 3
	speculativeDelay    = time.Second
)

func Exec(q *gocqlx.Queryx) (err error) {
	tries := 5
	for tries > 0 {
		err = q.Exec()
		if err == nil {
			return
		}
		time.Sleep(time.Millisecond)
		tries--
	}
	return
}

func Scan(q *gocqlx.Queryx, dest ...interface{}) error {
	return q.Scan(dest...)
}

func Get(q *gocqlx.Queryx, dest interface{}) error {
	return q.Get(dest)
}

func Select(q *gocqlx.Queryx, dest interface{}) error {
	return q.Select(dest)
}

/*
	The following part is for backward compatibility of older packages, this must be removed when they
	got updated.
*/

type Consistency uint16
type SerialConsistency uint16

// Db
type Db struct {
	config  Config
	session *gocql.Session
}

// NewCassDB
// Returns CassDB struct which has a 'gocql' session object enclosed.
// You can use DefaultCassConfig for quick configuration but make sure to set
// Username, Password and KeySpace
//
// example :
//  conf := DefaultConfig
//  conf.Username = "username"
//  conf.Password = "password"
//  conf.KeySpace = "key-space"
//  db := NewCassDB(conf)
func NewDB(conf Config) *Db {
	db := new(Db)
	db.config = conf
	cassCluster := gocql.NewCluster(conf.Host)
	retryPolicy := new(gocql.ExponentialBackoffRetryPolicy)
	retryPolicy.NumRetries = conf.Retries
	retryPolicy.Min = conf.RetryMinBackOff
	retryPolicy.Max = conf.RetryMaxBackOff

	cassCluster.RetryPolicy = retryPolicy
	cassCluster.ConnectTimeout = conf.ConnectTimeout
	cassCluster.Timeout = conf.Timeout
	cassCluster.ReconnectInterval = conf.ReconnectInterval
	cassCluster.DefaultIdempotence = conf.DefaultIdempotence
	cassCluster.QueryObserver = conf.QueryObserver
	cassCluster.Compressor = gocql.SnappyCompressor{}
	cassCluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.Username,
		Password: conf.Password,
	}
	cassCluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cassCluster.NumConns = conf.Concurrency
	cassCluster.PageSize = conf.PageSize
	cassCluster.WriteCoalesceWaitTime = conf.WriteCoalesceWaitTime
	if len(conf.Keyspace) > 0 {
		cassCluster.Keyspace = conf.Keyspace
		cassCluster.Consistency = gocql.Consistency(conf.Consistency)
		cassCluster.SerialConsistency = gocql.SerialConsistency(conf.SerialConsistency)
		cassCluster.CQLVersion = conf.CqlVersion
		if session, err := cassCluster.CreateSession(); err != nil {
			panic(err)
		} else {
			db.session = session
		}
	}
	return db
}

func (db *Db) SetConsistency(c Consistency) {
	db.session.SetConsistency(gocql.Consistency(c))
}

func (db *Db) GetSession() *gocql.Session {
	return db.session
}

func (db *Db) CloseSession() {
	db.session.Close()
}
