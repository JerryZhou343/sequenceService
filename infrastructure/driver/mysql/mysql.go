package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/mfslog/sequenceService/infrastructure/config"
	"github.com/pkg/errors"
)

func NewMySQLDB(conf *config.Config) (db *gorm.DB, err error) {

	DBUrl := conf.MySQLUser + ":" + conf.MySQLPassword + "@tcp(" + conf.MySQLAddr + ")/" + conf.MySQLDBName + "?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", DBUrl)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	db.DB().SetMaxIdleConns(1)
	db.DB().SetMaxOpenConns(conf.MySQLConLimit)
	db.SingularTable(true) //全局禁用表名复数形式
	db.LogMode(false)
	return
}
