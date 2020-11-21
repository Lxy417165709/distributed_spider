package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

//func LogSet(isOpen bool) {
//	MysqlCon.LogMode(isOpen)
//	MysqlLogCon.LogMode(isOpen)
//	MysqlGameCon.LogMode(isOpen)
//}

func InitDB(link, dbName string, maxCon int) {
	db = NewDB(GenDSN(link, dbName), maxCon)
}

func GenDSN(link, db string) string {
	return fmt.Sprintf("%v/%v?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local", link, db)
}

func NewDB(dsn string, maxCon int) *gorm.DB {
	con, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("Got error when connecting database, the error is '%s'", err))
	}
	idle := maxCon
	con.DB().SetMaxOpenConns(maxCon)
	con.DB().SetMaxIdleConns(idle)
	con.LogMode(true) // 开启sql日志
	con.BlockGlobalUpdate(true)
	return con
}
