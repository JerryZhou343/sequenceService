
package DBSession

import (
    "github.com/jinzhu/gorm"
    "os"
    _ "github.com/go-sql-driver/mysql"
    "sync"
    "github.com/mfslog/sequenceService/Server/configer"
    "github.com/mfslog/sequenceService/Server/log"
)



var DBIns *DBConnectInstance

var once sync.Once


//单例返回 DBCconnect Instance
func GetInstance() *DBConnectInstance {
    once.Do(func() {
        DBIns = &DBConnectInstance {}
    })
    return DBIns
}


type DBConnectInstance struct{
    SYSDB *gorm.DB
}


func (db *DBConnectInstance)connectSYSDB(cfg *configer.DBConfig){
    var err error
    sysDBUrl := cfg.DBUser+":" + cfg.DBPasswd + "@tcp("+ cfg.DBHostIP + ":" + cfg.DBHostPort + ")/" + cfg.DBName + "?charset=utf8&parseTime=True&loc=Local"
    db.SYSDB, err = gorm.Open("mysql",sysDBUrl)
    if err != nil{
        log.Error(0,"connect sys DB Failed [" + sysDBUrl +"]")
        os.Exit(1)
    }
    db.SYSDB.DB().SetMaxIdleConns(10)
    db.SYSDB.DB().SetMaxOpenConns(100)
}




func (db *DBConnectInstance)InitDBCon(){
    dbConfiger := configer.GetInstance()
    //init db
    //1.连接parkDB
    db.connectSYSDB(dbConfiger.GetSYSDBInfo())
    //2.连接userDB
}


//断开连接
func (db *DBConnectInstance)UninitDBCon(){
    db.SYSDB.Close()
}