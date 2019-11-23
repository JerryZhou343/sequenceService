package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/mfslog/sequenceService/infrastructure/config"
	"time"
)

func NewRedisPool(conf *config.Config) *redis.Pool {
	readTimeout := redis.DialReadTimeout(time.Second * time.Duration(2))
	writeTimeout := redis.DialWriteTimeout(time.Second * time.Duration(2))
	conTimeout := redis.DialConnectTimeout(time.Second * time.Duration(5))
	redisPool := &redis.Pool{
		MaxIdle:     1,
		MaxActive:   128,
		IdleTimeout: 0,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.RedisAddr, readTimeout, writeTimeout, conTimeout)
			if err != nil {
				return nil, err
			}
			if len(conf.RedisPassword) > 0 {
				if _, err := c.Do("AUTH", conf.RedisPassword); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return redisPool
}
