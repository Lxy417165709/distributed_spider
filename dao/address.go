package dao

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"spider/common/logger"
	"spider/model"
	"time"
)

type addressDao struct{}

func (*addressDao) Create(url string, crawlSource model.CrawlSource, crawlNodeNum int) error {
	if err := mysqlDB.Create(&model.Address{
		CrawlStatus:  model.NotCrawl,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CrawlSource:  crawlSource,
		Url:          url,
		CrawlNodeNum: crawlNodeNum,
	}).Error; err != nil {
		//logger.Error("Fail to finish mysqlDB.Create", zap.String("url", url), zap.Error(err))
		return err
	}
	return nil
}

func (*addressDao) UpdateStatus(url string, status model.CrawlStatus) error {
	tableName := (&model.Address{}).TableName()
	if err := mysqlDB.Table(tableName).Where("url = ?", url).Update(map[string]interface{}{
		"crawl_status": status,
		"updated_at":   time.Now(),
	}).Error; err != nil {
		logger.Error("Fail to finish Update", zap.String("url", url), zap.Any("status", status), zap.Error(err))
		return err
	}
	return nil
}

func (*addressDao) GetByUrl(url string) (*model.Address, error) {
	var result model.Address
	if err := mysqlDB.First(&result, "url = ?", url).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logger.Error("Fail to finish mysqlDB.First", zap.String("url", url), zap.Error(err))
		return nil, err
	}
	return &result, nil
}

func (*addressDao) GetNeedCrawlAddress(count int, crawlSource model.CrawlSource) ([]*model.Address, error) {
	var result []*model.Address
	db := mysqlDB.Limit(count).Find(&result, "crawl_status = ? and crawl_source = ?", model.NotCrawl, crawlSource)
	if err := db.Error; err != nil {
		logger.Error("Fail to finish mysqlDB.Find", zap.Error(err))
		return nil, err
	}
	return result, nil
}
