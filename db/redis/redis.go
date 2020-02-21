package redis

import (
	"fmt"
	log "git.ronaksoftware.com/ronak/rony/logger"
	"github.com/mediocregopher/radix/v3"
	"net"
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
	_PermCache    *Cache
	_TempCache    *Cache
	_CounterCache *Cache
)

// InitTempCache
// Initializes Redis as an in-memory cache, this function must be called before creating a new Model instance
func InitTempCache(redisConn *Cache) {
	_TempCache = redisConn
}

// InitPermCache
// Initializes Redis as an in-memory data storage, this function must be called before creating a new Model instance
func InitPermCache(redisConn *Cache) {
	_PermCache = redisConn
}

// InitCounterCache
// Initializes Redis as an in-memory data storage, this function must be called before creating a new Model instance
func InitCounterCache(redisConn *Cache) {
	_CounterCache = redisConn
}

func PermCache() *Cache {
	return _PermCache
}

func TempCache() *Cache {
	return _TempCache
}

func CounterCache() *Cache {
	return _CounterCache
}

var (
	_Caches map[string]*Cache
)

func Init(name string, redisConn *Cache) {
	_Caches[name] = redisConn
}

func Get(name string) *Cache {
	return _Caches[name]
}

type connType int

const (
	_ connType = iota
	connTypePool
	connTypeCluster
)

// Config
type Config struct {
	PoolSize           int
	MaxPoolSize        int
	NewConnOnEmpty     bool
	DialTimeout        time.Duration
	ConnReadTimeout    time.Duration
	ConnWriteTimeout   time.Duration
	OnEmptyPoolTimeout time.Duration
	PingTime           time.Duration
	RefillInterval     time.Duration
	Password           string
	Host               string
	ClusterHosts       []string
	Db                 int
}

type (
	CmdAction radix.CmdAction
	Action    radix.Action
)

var (
	DefaultConfig = Config{
		PoolSize:           10,
		MaxPoolSize:        100,
		NewConnOnEmpty:     false,
		DialTimeout:        3 * time.Second,
		ConnReadTimeout:    3 * time.Second,
		ConnWriteTimeout:   3 * time.Second,
		OnEmptyPoolTimeout: 500 * time.Millisecond,
		PingTime:           time.Second,
		RefillInterval:     time.Second,
		Db:                 0,
	}
)

// Cache
type Cache struct {
	cluster  *radix.Cluster
	pool     *radix.Pool
	connType connType
	conn     Conn
	scripts  map[string]radix.EvalScript
}

type Conn interface {
	Do(action radix.Action) error
	Close() error
}

type Scanner interface {
	Next(*string) bool
	Close() error
}

// New
// This is the constructor of Cache, it accepts Config as input, you can use
// DefaultConfig for quick initialization, but make sure to add 'Conn' and 'Password' to it
//
// example:
// conf := redis.DefaultConfig
// conf.Conn = "your-host.com"
// conf.Password = "password123"
// c := New(conf)
func New(conf Config) *Cache {
	r := new(Cache)
	r.scripts = make(map[string]radix.EvalScript)

	PoolOpts := make([]radix.PoolOpt, 0)
	PoolOpts = append(PoolOpts,
		radix.PoolConnFunc(func(network, addr string) (conn radix.Conn, e error) {
			return radix.Dial(network, addr,
				radix.DialAuthPass(conf.Password),
				radix.DialSelectDB(conf.Db),
				radix.DialConnectTimeout(conf.DialTimeout),
				radix.DialReadTimeout(conf.ConnReadTimeout),
				radix.DialWriteTimeout(conf.ConnWriteTimeout),
			)
		}),
		radix.PoolPingInterval(conf.PingTime),
		radix.PoolRefillInterval(conf.RefillInterval),
	)
	if conf.NewConnOnEmpty {
		PoolOpts = append(PoolOpts,
			radix.PoolOnFullBuffer(conf.MaxPoolSize, time.Second),
			radix.PoolOnEmptyCreateAfter(time.Millisecond*100),
		)
	} else {
		if conf.OnEmptyPoolTimeout > 0 {
			PoolOpts = append(PoolOpts, radix.PoolOnEmptyErrAfter(conf.OnEmptyPoolTimeout))
		} else {
			PoolOpts = append(PoolOpts, radix.PoolOnEmptyWait())
		}
	}

	pool, err := radix.NewPool("tcp", conf.Host, conf.PoolSize, PoolOpts...)
	if err != nil {
		log.Fatal(err.Error())
	}
	r.conn = pool
	r.connType = connTypePool
	r.pool = pool
	return r
}

func NewCluster(conf Config) *Cache {
	r := new(Cache)
	r.scripts = make(map[string]radix.EvalScript)
	var ClusterOpt radix.ClusterOpt = nil
	if len(conf.Password) > 0 {
		radix.PoolConnFunc(func(network, addr string) (radix.Conn, error) {
			c, err := net.Dial(network, addr)
			if err != nil {
				return nil, err
			}
			conn := radix.NewConn(c)
			conn.Do(radix.Cmd(nil, "AUTH", conf.Password))
			return conn, nil
		})
	}
	cluster, err := radix.NewCluster(conf.ClusterHosts, ClusterOpt)
	if err != nil {
		log.Fatal(err.Error())
	}
	r.conn = cluster
	r.connType = connTypeCluster
	r.cluster = cluster
	time.AfterFunc(time.Minute, func() {
		cluster.Sync()
	})

	return r

}

// NewScanner
func (r *Cache) NewScanner(opts radix.ScanOpts) Scanner {
	switch r.connType {
	case connTypePool:
		return radix.NewScanner(r.conn, opts)
	case connTypeCluster:
		return r.cluster.NewScanner(opts)
	}
	return nil
}

// RegisterScript
func (r *Cache) RegisterScript(name string, numKeys int, script string) {
	r.scripts[name] = radix.NewEvalScript(numKeys, script)
}

// RunScript
func (r *Cache) RunScript(name string, result *interface{}, args ...string) error {
	return r.Do(r.scripts[name].Cmd(result, args...))
}

func (r *Cache) ErrCh() chan error {
	return r.pool.ErrCh
}

func (r *Cache) Do(action radix.Action) error {
	return r.conn.Do(action)
}

func (r *Cache) Pipeline(commands ...radix.CmdAction) error {
	return r.Do(radix.Pipeline(commands...))
}

func (r *Cache) Cmd(rcv interface{}, cmd string, key string, args ...interface{}) radix.CmdAction {
	return radix.FlatCmd(rcv, cmd, key, args...)
}

func (r *Cache) Multi() error {
	return r.conn.Do(radix.Cmd(nil, "MULTI"))
}

func (r *Cache) Exec() error {
	return r.conn.Do(radix.Cmd(nil, "EXEC"))
}

func (r *Cache) Close() error {
	return r.conn.Close()
}

func (r *Cache) Exists(keyName string) (reply bool, err error) {
	err = r.Do(radix.Cmd(&reply, "EXISTS", keyName))
	return
}

func (r *Cache) Expire(keyName string, ttl int) (reply bool, err error) {
	err = r.Do(radix.Cmd(&reply, "EXPIRE", keyName, fmt.Sprintf("%d", ttl)))
	return
}

func (r *Cache) Del(keyName ...string) (err error) {
	err = r.Do(radix.Cmd(nil, "DEL", keyName...))
	return
}

func (r *Cache) Set(keyName string, value interface{}) (err error) {
	err = r.Do(radix.FlatCmd(nil, "SET", keyName, value))
	return
}

func (r *Cache) SetNx(keyName string, value interface{}) (reply bool, err error) {
	err = r.Do(radix.FlatCmd(&reply, "SETNX", keyName, value))
	return
}

func (r *Cache) SetEx(keyName string, ttl, value interface{}) (err error) {
	err = r.Do(radix.FlatCmd(nil, "SETEX", keyName, ttl, value))
	return
}

func (r *Cache) GetString(keyName string) (reply string, err error) {
	err = r.Do(radix.Cmd(&reply, "GET", keyName))
	return
}

func (r *Cache) GetInt(keyName string) (reply int, err error) {
	err = r.Do(radix.Cmd(&reply, "GET", keyName))
	return
}

func (r *Cache) GetInt32(keyName string) (reply int32, err error) {
	err = r.Do(radix.Cmd(&reply, "GET", keyName))
	return
}

func (r *Cache) GetInt64(keyName string) (reply int64, err error) {
	err = r.Do(radix.Cmd(&reply, "GET", keyName))
	return
}

func (r *Cache) GetUInt64(keyName string) (reply uint64, err error) {
	err = r.Do(radix.Cmd(&reply, "GET", keyName))
	return
}

func (r *Cache) GetBytes(keyName string) (reply []byte, err error) {
	err = r.Do(radix.Cmd(&reply, "GET", keyName))
	return
}

func (r *Cache) GetByteSlice(keyName string) (reply [][]byte, err error) {
	err = r.Do(radix.Cmd(&reply, "GET", keyName))
	return
}

func (r *Cache) MGetBytes(keyNames ...string) (reply [][]byte, err error) {
	err = r.Do(radix.Cmd(&reply, "MGET", keyNames...))
	return
}

func (r *Cache) Inc(keyName string) (reply interface{}, err error) {
	err = r.Do(radix.Cmd(&reply, "INCR", keyName))
	return
}

func (r *Cache) IncInt64(keyName string) (reply int64, err error) {
	err = r.Do(radix.Cmd(&reply, "INCR", keyName))
	return
}

func (r *Cache) IncBy(keyName string, n int64) (reply int64, err error) {
	err = r.Do(radix.Cmd(&reply, "INCRBY", keyName, fmt.Sprintf("%d", n)))
	return
}

func (r *Cache) HSet(keyName string, fieldName interface{}, value interface{}) (reply bool, err error) {
	err = r.Do(radix.FlatCmd(&reply, "HSET", keyName, fieldName, value))
	return
}

func (r *Cache) HGetAllStringMap(keyName string) (reply map[string]string, err error) {
	err = r.Do(radix.Cmd(&reply, "HGETALL", keyName))
	return
}

func (r *Cache) HGetAllInt64Map(keyName string) (reply map[string]int64, err error) {
	err = r.Do(radix.Cmd(&reply, "HGETALL", keyName))
	return
}

func (r *Cache) HGetAllInt32Map(keyName string) (reply map[string]int32, err error) {
	err = r.Do(radix.Cmd(&reply, "HGETALL", keyName))
	return
}

func (r *Cache) ZCard(keyName string) (reply int, err error) {
	err = r.Do(radix.Cmd(&reply, "ZCARD", keyName))
	return
}
