package cql

import (
	"github.com/gocql/gocql"
	"time"
)

/*
   Creation Time: 2019 - Oct - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Config struct {
	Host                  string
	Username              string
	Password              string
	Keyspace              string
	Retries               int
	RetryMinBackOff       time.Duration
	RetryMaxBackOff       time.Duration
	ConnectTimeout        time.Duration
	Timeout               time.Duration
	ReconnectInterval     time.Duration
	Concurrency           int
	Consistency           gocql.Consistency
	SerialConsistency     gocql.SerialConsistency
	ReplicationClass      string
	ReplicationFactor     int
	CqlVersion            string
	DefaultIdempotence    bool
	QueryObserver         gocql.QueryObserver
	PageSize              int
	WriteCoalesceWaitTime time.Duration
}

var (
	DefaultConfig = Config{
		Concurrency:           5,
		Timeout:               600 * time.Millisecond,
		ConnectTimeout:        600 * time.Millisecond,
		Retries:               100,
		RetryMinBackOff:       100 * time.Millisecond,
		RetryMaxBackOff:       time.Second,
		ReconnectInterval:     10 * time.Second,
		Consistency:           gocql.LocalQuorum,
		SerialConsistency:     gocql.LocalSerial,
		ReplicationClass:      "SimpleStrategy",
		ReplicationFactor:     1,
		CqlVersion:            "3.4.4",
		DefaultIdempotence:    true,
		QueryObserver:         nil,
		PageSize:              1000,
		WriteCoalesceWaitTime: time.Millisecond,
	}
)
