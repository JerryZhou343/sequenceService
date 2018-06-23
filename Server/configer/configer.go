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
        DBHostIP: etcd.getStringValue(key + "sysDB/ip","127.0.0.1"),
        DBHostPort: etcd.getStringValue(key+"sysDB/port","3306"),
        DBName : etcd.getStringValue(key + "sysDB/name","sysDB"),
        DBUser: etcd.getStringValue(key + "sysDB/user","root"),
        DBPasswd: etcd.getStringValue(key + "sysDB/passwd","root"),
        DBMaxCon: etcd.getIntVale(key + "sysDB/maxcon",10),
        DBMinCon: etcd.getIntVale(key + "sysDB/mincon",20),
    }
    
    
    
    cfg.cacheConfig = &CacheConfig{
        Hostip : etcd.getStringValue(key + "redis/ip","127.0.0.1"),
        Hostport:etcd.getStringValue(key + "redis/port","6379"),
        Passwd: etcd.getStringValue(key + "redis/passwd","redis"),
        PoolSize:etcd.getIntVale(key + "redis/poolSize",10),
        DBNum: etcd.getIntVale(key + "redis/dbNum",0),
    }
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
