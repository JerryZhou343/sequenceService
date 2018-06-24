package cacheSession

import (
	"sync"
	"github.com/go-redis/redis"
	"github.com/mfslog/sequenceService/Server/configer"
	"github.com/mfslog/sequenceService/Server/log"
	"fmt"
	"os"
)


type redisConIns struct{
	client *redis.Client
}

var redisIns *redisConIns

var once sync.Once


//单例返回 redisCon Instance
func GetInstance() *redisConIns {
	once.Do(func() {
		redisIns = &redisConIns{}
	})
	return redisIns
}

func (rds *redisConIns)InitCon(){
	cfg := configer.GetInstance()
	redisCfg := cfg.GetRedisInfo()
	rds.client = redis.NewClient(&redis.Options{
		Addr: redisCfg.Hostip +":" + redisCfg.Hostport,
		Password: redisCfg.Passwd,
		DB:redisCfg.DBNum,
		PoolSize:redisCfg.PoolSize,
	})

	_,err := rds.client.Ping().Result()
	if err != nil{
		log.Error(fmt.Sprintf("%s",err))
		os.Exit(1)
	}
}


func (rds *redisConIns)GetClient()*redis.Client{
	return rds.client
}

func (rds *redisConIns)Close(){
	rds.Close()
}