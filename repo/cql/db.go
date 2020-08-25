package cql

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"strings"
)

/*
   Creation Time: 2020 - Aug - 13
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	_Session gocqlx.Session
)

func Init(config Config) {
	scyllaCluster := gocql.NewCluster(strings.Split(config.Host, ",")...)
	retryPolicy := &gocql.ExponentialBackoffRetryPolicy{
		NumRetries: config.Retries,
		Min:        config.RetryMinBackOff,
		Max:        config.RetryMaxBackOff,
	}

	scyllaCluster.RetryPolicy = retryPolicy
	scyllaCluster.ConnectTimeout = config.ConnectTimeout
	scyllaCluster.Timeout = config.Timeout
	scyllaCluster.ReconnectInterval = config.ReconnectInterval
	scyllaCluster.DefaultIdempotence = config.DefaultIdempotence
	scyllaCluster.QueryObserver = config.QueryObserver
	scyllaCluster.Compressor = gocql.SnappyCompressor{}
	scyllaCluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}

	scyllaCluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	scyllaCluster.NumConns = config.Concurrency
	scyllaCluster.PageSize = config.PageSize
	scyllaCluster.WriteCoalesceWaitTime = config.WriteCoalesceWaitTime
	if len(config.Keyspace) == 0 {
		panic("keyspace is not set")
	}
	scyllaCluster.QueryObserver = config.QueryObserver
	scyllaCluster.Keyspace = config.Keyspace
	scyllaCluster.Consistency = config.Consistency
	scyllaCluster.SerialConsistency = config.SerialConsistency
	scyllaCluster.CQLVersion = config.CqlVersion
	session, err := scyllaCluster.CreateSession()
	if err != nil {
		panic(err)
	}
	_Session = gocqlx.NewSession(session)
}

func Session() gocqlx.Session {
	return _Session
}

func Exec(q *gocqlx.Queryx) error {
	return q.Exec()
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
