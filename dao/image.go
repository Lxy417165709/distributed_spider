package dao

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"spider/common/logger"
	"spider/model"
	"time"
)

type ImageDao struct{}

func (*ImageDao) Create(dataUrl string, md5 string, addressId int) error {
	if err := mysqlDB.Create(&model.Image{
		MD5:       md5,
		Url:       dataUrl,
		AddressId: addressId,
		CreatedAt: time.Now(),
	}).Error; err != nil {
		//logger.Error("Fail to finish mysqlDB.Create",
		//	zap.String("dataUrl", dataUrl),
		//	zap.String("md5", md5),
		//	zap.Int("addressId", addressId),
		//	zap.Error(err))
		return err
	}
	return nil
}

func (*ImageDao) GetByMd5(imageMd5 string) (*model.Image, error) {
	var result model.Image
	if err := mysqlDB.First(&result, "md5 = ?", imageMd5).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logger.Error("Fail to finish mysqlDB.First", zap.String("md5", imageMd5), zap.Error(err))
		return nil, err
	}
	return &result, nil
}

func (*ImageDao) GetAll() ([]*model.Image, error) {
	var result []*model.Image
	if err := mysqlDB.Find(&result).Error; err != nil {
		logger.Error("Fail to finish mysqlDB.Find", zap.Error(err))
		return nil, err
	}
	return result, nil
}
