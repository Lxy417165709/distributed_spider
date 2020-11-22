package handler

import (
	"context"
	"go.uber.org/zap"
	"spider/cache"
	"spider/common/logger"
	"spider/dao"
	"spider/model"
	"time"
)

type SpiderSupplier struct {
	ctx         context.Context
	cancelFunc context.CancelFunc
	supplyCount int
	id int
}

func NewSpiderSupplier(supplyCount int,id int,ctx context.Context,cancelFunc context.CancelFunc) *SpiderSupplier {
	return &SpiderSupplier{
		ctx:         ctx,
		supplyCount: supplyCount,
		cancelFunc:cancelFunc,
		id:id,
	}
}

func (s *SpiderSupplier) Run(addressChannel chan *model.Address, crawlSource model.CrawlSource) {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			addresses := s.GetAddresses(crawlSource)
			for _, ad := range addresses {
				addressChannel <- ad
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *SpiderSupplier) GetAddresses(source model.CrawlSource) []*model.Address {
	if !cache.Spider.GetSupplierLock() {
		return nil
	}
	logger.Info("Get supplier lock",zap.Int("ID",s.id))
	defer func() {
		logger.Info("Release supplier lock",zap.Int("ID",s.id))
		cache.Spider.ReleaseSupplierLock()
	}()

	addresses, err := dao.AddressDB.GetNeedCrawlAddress(s.supplyCount, source)
	if err != nil {
		logger.Error("Fail to finish AddressDB.GetNeedCrawlAddress", zap.Error(err))
		return nil
	}
	logger.Info("Get addresses",zap.Int("ID",s.id))
	urls := make([]string, 0)
	for _, address := range addresses {
		urls = append(urls, address.Url)
	}
	if err := dao.AddressDB.UpdateStatusBatch(model.Crawling, urls...); err != nil {
		logger.Info("Fail to finish AddressDB.UpdateStatus", zap.Error(err))
		return nil
	}
	return addresses
}
