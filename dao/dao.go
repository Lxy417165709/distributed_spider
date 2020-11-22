package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	mysqlDB *gorm.DB
)

var (
	AddressDB = new(addressDao)
	ImageDB   = new(ImageDao)
)

//func LogSet(isOpen bool) {
//	MysqlCon.LogMode(isOpen)
//	MysqlLogCon.LogMode(isOpen)
//	MysqlGameCon.LogMode(isOpen)
//}

func InitDB(link, dbName string, maxCon int) {
	mysqlDB = NewDB(GenDSN(link, dbName), maxCon)
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
