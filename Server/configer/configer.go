package configer

import (
    "strings"
    "github.com/spf13/viper"
    "sync"
    "os"
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
}

type CacheConfig struct{
    Hostip string
    Hostport string
    Passwd string
}


//加载etcd 中的配置信息
func (cfg *Configer)LocadConfig() {
    //local db config
    basePath := viper.GetString("service_register.base_path")
    name := viper.GetString("service_register.service_name")
    etcdEndPoints := viper.GetStringSlice("service_register.etcd_address")
    var key string = strings.TrimRight(basePath,"/") + "/" + name + "/config/"
    
    etcd := EtcdService{}
    
    //连接etcd
    err := etcd.connectEtcd(etcdEndPoints)
    if err != nil{
        os.Exit(1)
    }
    
    //读取配置信息
    cfg.SYSDBCfg =&DBConfig{
        DBHostIP: etcd.getStringValue(key + "sysDB/ip"),
        DBHostPort: etcd.getStringValue(key+"sysDB/port"),
        DBName : etcd.getStringValue(key + "sysDB/name"),
        DBUser: etcd.getStringValue(key + "sysDB/user"),
        DBPasswd: etcd.getStringValue(key + "sysDB/passwd"),
    }
    
    
    
    cfg.cacheConfig = &CacheConfig{
        Hostip : etcd.getStringValue(key + "redis/ip"),
        Hostport:etcd.getStringValue(key + "redis/port"),
        Passwd: etcd.getStringValue(key + "redis/passwd"),
    }
    
    //cfg.kafkaAddress = etcd.getStringValue(key + "kafka/hosts")
    //断开和etcd 的连接
    etcd.disconnect()
}

//返回sys库信息
func (cfg *Configer)GetSYSDBInfo()*DBConfig{
    return cfg.SYSDBCfg
}




//返回redis信息
func (cfg *Configer)GetRedisInfo()*CacheConfig{
    return cfg.cacheConfig;
}

//返回Kafka地址
//func (cfg *Configer)GetKafaAddress()string{
//    return cfg.kafkaAddress
//}