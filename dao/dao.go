package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"spider/common/utils"
)

var (
	mysqlDB *gorm.DB
)

var (
	AddressDB = new(addressDao)
	ImageDB   = new(ImageDao)
)

func InitDB(link, dbName string, maxCon int) {
	mysqlDB = utils.NewDB(utils.GenDSN(link, dbName), maxCon)
}

func CloseLog() {
	mysqlDB.LogMode(false)
}
