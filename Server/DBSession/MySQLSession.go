
package DBSession

import (
    "github.com/jinzhu/gorm"
    "os"
    _ "github.com/go-sql-driver/mysql"
    "sync"
    "github.com/mfslog/sequenceService/Server/configer"
    "github.com/mfslog/sequenceService/Server/log"
)



var dbIns *dbConnectInstance

var once sync.Once


//单例返回 DBCconnect Instance
func GetInstance() *dbConnectInstance {
    once.Do(func() {
        dbIns = &dbConnectInstance{}
    })
    return dbIns
}


type dbConnectInstance struct{
    SYSDB *gorm.DB
}


func (db *dbConnectInstance)connectSYSDB(cfg *configer.DBConfig){
    var err error
    sysDBUrl := cfg.DBUser+":" + cfg.DBPasswd + "@tcp("+ cfg.DBHostIP + ":" + cfg.DBHostPort + ")/" + cfg.DBName + "?charset=utf8&parseTime=True&loc=Local"
    db.SYSDB, err = gorm.Open("mysql",sysDBUrl)
    if err != nil{
        log.Error("connect sys DB Failed [" + sysDBUrl +"]")
        os.Exit(1)
    }
    db.SYSDB.DB().SetMaxIdleConns(cfg.DBMinCon)
    db.SYSDB.DB().SetMaxOpenConns(cfg.DBMaxCon)
    db.SYSDB.SingularTable(true) //全局禁用表名复数形式
    db.SYSDB.LogMode(false)
}




func (db *dbConnectInstance)InitDBCon(){
    dbConfiger := configer.GetInstance()
    //init db
    //1.连接parkDB
    db.connectSYSDB(dbConfiger.GetSYSDBInfo())
    //2.连接userDB
}


//断开连接
func (db *dbConnectInstance)UninitDBCon(){
    db.SYSDB.Close()
}