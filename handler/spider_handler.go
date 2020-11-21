package handler

import (
	"context"
	"go.uber.org/zap"
	"spider/common/logger"
	"spider/model"
	"time"
)

type SpiderHandler struct {
	crawlResultChannel chan *model.CrawlResult
	ctx                context.Context
}

func (s *SpiderHandler) Run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case result := <-s.crawlResultChannel:
			s.Handle(result)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *SpiderHandler) Handle(result *model.CrawlResult) {
	// 5.1 爬取错误时
	if result.Err != nil {
		logger.Error("Fail to crawl", zap.String("url", result.Url), zap.Error(result.Err))
		if err := dao.AddressDB.UpdateStatus(result.Url, model.FailCrawled); err != nil {
			logger.Error("Fail to finish AddressDB.UpdateStatus", zap.Error(result.Err))
		}
		return
	}

	// 5.2 爬取成功时
	if err := dao.AddressDB.UpdateStatus(result.Url, model.HadCrawled); err != nil {
		logger.Error("Fail to finish AddressDB.UpdateStatus", zap.Error(result.Err))
		return
	}

	// 5.1.2 将该链接下所有未爬取的子链接加入爬取通道
	go func() {
		for _, crawlUrl := range result.CrawlUrls {
			if err := dao.AddressDB.Create(crawlUrl, s.identity); err != nil {
				//logger.Error("Fail to finish AddressDB.Create", zap.String("url", crawlUrl), zap.Error(err))
				continue
			}
		}
	}()

	// 5.1.3 将该链接下所有爬取到、且未存储的图片进行存储 {
	go func(){
		for _, imageUrl := range result.ImageUrls {
			s.storeImage(imageUrl)
		}
	}()
}

func (s *ImageSpider) handleCrawlResult() {
	for {
		go func(){
			redisResult := s.spiderStorage.GetCrawlResult()
			if redisResult == nil {
				return
			}
			result := redisResult.ToCrawlResult()
			//logger.Info("Get result",
			//	zap.String("result.url", result.Url),
			//	zap.String("result.url", result.Url),
			//	zap.Any("result.crawlUrls", result.CrawlUrls),
			//	zap.Any("result.imageUrls", result.ImageUrls),
			//	zap.Any("result.err", result.Err))
			// 5.1 爬取错误时
			if result.Err != nil {
				logger.Error("Fail to crawl", zap.String("url", result.Url), zap.Error(result.Err))
				if err := dao.AddressDB.UpdateStatus(result.Url, model.FailCrawled); err != nil {
					logger.Error("Fail to finish AddressDB.UpdateStatus", zap.Error(result.Err))
				}
				return
			}

			// 5.2 爬取成功时
			if err := dao.AddressDB.UpdateStatus(result.Url, model.HadCrawled); err != nil {
				logger.Error("Fail to finish AddressDB.UpdateStatus", zap.Error(result.Err))
				return
			}

			// 5.1.2 将该链接下所有未爬取的子链接加入爬取通道
			go func() {
				for _, crawlUrl := range result.CrawlUrls {
					if err := dao.AddressDB.Create(crawlUrl, s.identity); err != nil {
						//logger.Error("Fail to finish AddressDB.Create", zap.String("url", crawlUrl), zap.Error(err))
						continue
					}
					s.spiderStorage.PushCrawlUrlIfNotCrawlAndNotFull(crawlUrl)
				}
			}()

			// 5.1.3 将该链接下所有爬取到、且未存储的图片进行存储 {
			go func(){
				for _, imageUrl := range result.ImageUrls {
					s.storeImage(imageUrl)
				}
			}()
		}()
		time.Sleep(1 * time.Second)
	}
}
