package datasource

import (
	"fmt"
	"github.com/fuloge/basework/api"
	cfg "github.com/fuloge/basework/configs"
	"github.com/fuloge/basework/pkg/log"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	redisURL       string
	redisPassword  string
	db             int
	timeout        = 2
	maxActive      int
	maxIdle        int
	idleTimeoutSec int
)

func init() {
	redisURL = cfg.EnvConfig.Redis.Hosts[0]
	redisPassword = cfg.EnvConfig.Redis.Password
	db = cfg.EnvConfig.Redis.DB
	maxActive = cfg.EnvConfig.Redis.MaxActive
	maxIdle = cfg.EnvConfig.Redis.MaxIdle
	idleTimeoutSec = cfg.EnvConfig.Redis.IdleTimeoutSec
}

var once sync.Once
var redisPl *redis.Pool

func newRedisPool() *redis.Pool {
	logger := log.New()

	once.Do(func() {
		redisPl = &redis.Pool{
			MaxIdle:     maxIdle,
			MaxActive:   maxActive,
			IdleTimeout: time.Duration(idleTimeoutSec) * time.Second,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialURL(redisURL, redis.DialDatabase(db), redis.DialPassword(redisPassword),
					redis.DialConnectTimeout(time.Duration(timeout)*time.Second),
					redis.DialReadTimeout(time.Duration(timeout)*time.Second),
					redis.DialWriteTimeout(time.Duration(timeout)*time.Second))
				if err != nil {
					logger.Error("RedisPool", zap.String(api.RedisConnErr.Message, err.Error()))
					return nil, fmt.Errorf("redis connection error: %s", err)
				}

				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) (err error) {
				if time.Since(t) < time.Minute {
					return nil
				}

				_, err = c.Do("PING")
				if err != nil {
					logger.Error("RedisPool", zap.String(api.RedisConnErr.Message, err.Error()))
				}
				return err
			},
		}
	})

	return redisPl
}

func GetRedisConn() (redis.Conn, *redis.Pool) {
	pool := newRedisPool()
	return pool.Get(), pool
}

//key:"lock_uid"
//uid: user_id
func AddLock(val string) bool {
	msg, _ := redis.String(
		RedisExec("set", "lock:LOCK_"+val, val, "nx", "ex", 4),
	)

	if msg == "OK" {
		return true
	}

	return false
}

func DelLock(val string) {
	_, err := RedisExec("del", "lock:LOCK_"+val)
	if err != nil {
		fmt.Println(api.RedisConnErr, err.Error())
	}
}

//func GetLock(conn redis.Conn, val string) string {
//	defer conn.Close()
//
//	msg, _ := redis.String(conn.Do("get", "lock:LOCK_"+val))
//	return msg
//}

func RedisExec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con, _ := GetRedisConn()
	defer con.Close()
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}

func CloseReidsPool() {
	pool := newRedisPool()
	pool.Close()
}
