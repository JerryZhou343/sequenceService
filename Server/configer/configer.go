package configer

import (
    "strings"
    "github.com/spf13/viper"
    "sync"
    "github.com/mfslog/sequenceService/Server/serverPlugin"
)

var instance *Configer

var once sync.Once


//单例返回 configer instance
func GetInstance() *Configer {
    once.Do(func() {
        instance = &Configer {}
    })
    return instance
}


//configer 数据成员
type Configer struct{
    SYSDBCfg *DBConfig
    cacheConfig *CacheConfig
}



//数据库配置描述
type DBConfig struct{
    DBHostIP string
    DBHostPort string
    DBUser string
    DBPasswd string
    DBName string
    DBMinCon int
    DBMaxCon int
}

type CacheConfig struct{
    Hostip string
    Hostport string
    Passwd string
    PoolSize int
    DBNum int
}

type lastTimeStamp struct{
    SnowFlakeStamp int
    DisorderStamp int
}


//加载etcd 中的配置信息
func (cfg *Configer)LocadConfig() {
    basePath := viper.GetString("service_register.base_path")
    name := viper.GetString("service_register.service_name")
    var key string = strings.TrimRight(basePath,"/") + "/" + name + "/config/"
    etcd := serverplugin.GetEtcdConIns()
    //读取配置信息
    cfg.SYSDBCfg =&DBConfig{
        DBHostIP: etcd.GetStringValue(key + "sysDB/ip","127.0.0.1"),
        DBHostPort: etcd.GetStringValue(key+"sysDB/port","3306"),
        DBName : etcd.GetStringValue(key + "sysDB/name","sysDB"),
        DBUser: etcd.GetStringValue(key + "sysDB/user","root"),
        DBPasswd: etcd.GetStringValue(key + "sysDB/passwd","root"),
        DBMaxCon: etcd.GetIntVale(key + "sysDB/maxcon",10),
        DBMinCon: etcd.GetIntVale(key + "sysDB/mincon",20),
    }
    
    
    
    cfg.cacheConfig = &CacheConfig{
        Hostip : etcd.GetStringValue(key + "redis/ip","127.0.0.1"),
        Hostport:etcd.GetStringValue(key + "redis/port","6379"),
        Passwd: etcd.GetStringValue(key + "redis/passwd","redis"),
        PoolSize:etcd.GetIntVale(key + "redis/poolSize",10),
        DBNum: etcd.GetIntVale(key + "redis/dbNum",0),
    }
}

//返回sys库信息
func (cfg *Configer)GetSYSDBInfo()*DBConfig{
    return cfg.SYSDBCfg
}




//返回redis信息
func (cfg *Configer)GetRedisInfo()*CacheConfig{
    return cfg.cacheConfig;
}
