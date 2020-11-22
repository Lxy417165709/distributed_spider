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
}

func NewSpiderSupplier(supplyCount int,ctx context.Context,cancelFunc context.CancelFunc) *SpiderSupplier {
	return &SpiderSupplier{
		ctx:         ctx,
		supplyCount: supplyCount,
		cancelFunc:cancelFunc,
	}
}

func (s *SpiderSupplier) Run(crawlUrlChannel chan string, crawlSource model.CrawlSource) {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			urls := s.GetUrls(crawlSource)
			for _, url := range urls {
				crawlUrlChannel <- url
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *SpiderSupplier) GetUrls(source model.CrawlSource) []string {
	if !cache.Spider.GetSupplierLock() {
		return nil
	}

	defer func() {
		cache.Spider.ReleaseSupplierLock()
	}()

	addresses, err := dao.AddressDB.GetNeedCrawlAddress(s.supplyCount, source)
	if err != nil {
		logger.Error("Fail to finish AddressDB.GetNeedCrawlAddress", zap.Error(err))
		return nil
	}
	urls := make([]string, 0)
	for _, address := range addresses {
		urls = append(urls, address.Url)
	}
	if err := dao.AddressDB.UpdateStatusBatch(model.Crawling, urls...); err != nil {
		logger.Info("Fail to finish AddressDB.UpdateStatus", zap.Error(err))
		return nil
	}

	return urls
}
