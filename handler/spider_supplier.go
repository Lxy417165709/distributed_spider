package handler

import (
	"context"
	"go.uber.org/zap"
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

func (s *SpiderSupplier) Run(crawlResultChannel chan string, crawlSource model.CrawlSource) {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			urls := s.GetUrls(crawlSource)
			for _, url := range urls {
				crawlResultChannel <- url
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *SpiderSupplier) GetUrls(source model.CrawlSource) []string {
	addresses, err := dao.AddressDB.GetNeedCrawlAddress(s.supplyCount, source)
	if err != nil {
		logger.Error("Fail to finish AddressDB.GetNeedCrawlAddress", zap.Error(err))
		return nil
	}
	urls := make([]string, 0)
	for _, address := range addresses {
		urls = append(urls, address.Url)
	}
	return urls
}
