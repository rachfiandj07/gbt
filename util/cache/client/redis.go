package client

import (
	"sync"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/pkg/errors"
	"gopkg.in/redis.v5"

	"github.com/golang-base-template/util/config"
)

type (
	//RedisConnsMap holds maps of redis connections
	RedisConnsMap map[string]*RedisConnInfo

	//RedisConnInfo holds conn infos
	RedisConnInfo struct {
		Addr string
		Pass string
		Conn *redis.Client
	}

	RedisList struct {
		CriticalRedis    []string
		NonCriticalRedis []string
	}
)

var (
	//RedisClients holds all redis clients
	RedisClients RedisConnsMap

	onReconnect    = &sync.Once{}
	isUnitTest     = false
	ErrRedisClosed = errors.New("redis: client is closed")
	mxRc           = sync.Mutex{}
)

func InitRedis(connString []string) (err error) {
	if RedisClients == nil {
		RedisClients = make(map[string]*RedisConnInfo)
	}
	connList := config.Get().Redis
	for _, name := range connString {
		addstruct, exist := connList[name]
		if !exist {
			continue
		}
		RedisClients[name] = NewConnection(name, addstruct.Address, "")
	}

	err = pingRedis()

	return
}

// NewConnection for given connection string.
func NewConnection(name, connection, password string) *RedisConnInfo {
	// init with default config
	redisOpt := &redis.Options{
		Addr:               connection,
		Password:           password,
		DB:                 0,
		PoolSize:           500,
		MaxRetries:         2,
		PoolTimeout:        time.Millisecond * 1200,
		DialTimeout:        time.Millisecond * 1000,
		ReadTimeout:        time.Millisecond * 1000,
		WriteTimeout:       time.Millisecond * 1000,
		IdleTimeout:        time.Second * 60,
		IdleCheckFrequency: time.Second * 10,
	}

	if rdsConfig, exist := config.Get().Redis[name]; exist {
		redisOpt = &redis.Options{
			Addr:               connection,
			Password:           password, // set password if gcp environment, no need for aliyun
			DB:                 0,        // use default DB
			PoolSize:           rdsConfig.PoolSize,
			MaxRetries:         2,
			PoolTimeout:        time.Millisecond * time.Duration(rdsConfig.PoolTimeoutMS),
			DialTimeout:        time.Millisecond * time.Duration(rdsConfig.DialTimeoutMS),
			ReadTimeout:        time.Millisecond * time.Duration(rdsConfig.ReadTimeoutMS),
			WriteTimeout:       time.Millisecond * time.Duration(rdsConfig.WriteTimeoutMS),
			IdleTimeout:        time.Second * time.Duration(rdsConfig.IdleTimeoutSec),
			IdleCheckFrequency: time.Second * time.Duration(rdsConfig.IdleFreqCheckSec),
		}
	}
	rds := redis.NewClient(redisOpt)

	// mock redis in unit test
	if isUnitTest {
		mr, _ := miniredis.Run()
		rds = redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})
	}

	return &RedisConnInfo{
		Addr: connection,
		Pass: password,
		Conn: rds,
	}
}

// PingRedis connection.
func pingRedis() error {
	// prevent data race
	mxRc.Lock()
	defer mxRc.Unlock()
	for _, val := range RedisClients {
		_, err := val.Conn.Ping().Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// InitMock for mocking redis in unit test
func InitMock(address, password string, list []string) (conn RedisConnsMap, err error) {
	conn = make(RedisConnsMap)
	if len(list) == 0 {
		err = errors.New("[Redis Init][InitMock] No connection inputs")
		return conn, err
	}

	for _, name := range list {
		conn[name] = NewConnection(name, address, password)
	}
	return conn, nil
}
