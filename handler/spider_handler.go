package handler

import (
	"context"
	"crypto/md5"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"spider/common/logger"
	"spider/dao"
	"spider/model"
	"time"
)

type SpiderHandler struct {
	httpClient *http.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewSpiderHandler(client *http.Client, ctx context.Context, cancelFunc context.CancelFunc) *SpiderHandler {
	return &SpiderHandler{
		httpClient: client,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

func (s *SpiderHandler) Run(crawlResultChannel chan *model.CrawlResult) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case result := <-crawlResultChannel:
			s.Handle(result)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *SpiderHandler) Handle(result *model.CrawlResult) {
	// 1. 爬取错误时
	if result.Err != nil {
		logger.Error("Fail to crawl", zap.String("url", result.Url), zap.Error(result.Err))
		if err := dao.AddressDB.UpdateStatus(result.Url, model.FailCrawled); err != nil {
			logger.Error("Fail to finish AddressDB.UpdateStatus", zap.Error(result.Err))
		}
		return
	}

	// 2. 爬取成功时
	// 2.1 将该链接下所有未爬取的子链接加入爬取通道
	for _, crawlUrl := range result.CrawlUrls {
		go func(crawlUrl string) {
			s.storeUrl(crawlUrl, result.CrawlSource, result.CrawlNodeNum)
		}(crawlUrl)
	}

	// 2.2 将该链接下所有爬取到、且未存储的图片进行存储
	for _, imageUrl := range result.ImageUrls {
		go func(imageUrl string) {
			s.storeImage(imageUrl, result.AddressId)
		}(imageUrl)
	}

	// 2.3 状态更新
	if err := dao.AddressDB.UpdateStatus(result.Url, model.HadCrawled); err != nil {
		logger.Error("Fail to finish AddressDB.UpdateStatus", zap.Error(result.Err))
		return
	}
}

func (s *SpiderHandler) storeUrl(crawlUrl string, source model.CrawlSource, nodeNum int) {
	if err := dao.AddressDB.Create(crawlUrl, source, nodeNum); err != nil {
		//logger.Error("Fail to finish AddressDB.Create", zap.String("url", crawlUrl), zap.Error(err))
	}
}

func (s *SpiderHandler) storeImage(imageUrl string, addressId int) {
	imageMd5 := s.getImageMd5(imageUrl)
	if imageMd5 == "" {
		logger.Warn("Image md5 is blank", zap.String("imageMd5", imageMd5))
		return
	}
	image, err := dao.ImageDB.GetByMd5(imageMd5)
	if err != nil {
		logger.Error("Fail to finish ImageDB.GetByMd5", zap.String("imageMd5", imageMd5), zap.Error(err))
		return
	}
	if image != nil {
		//logger.Info("Image has exist", zap.Any("image", image))
		return
	}
	if err := dao.ImageDB.Create(imageUrl, imageMd5, addressId); err != nil {
		logger.Error("Fail to finish ImageDB.GetByMd5", zap.Error(err))
		return
	}
}

func (s *SpiderHandler) getImageMd5(imageUrl string) string {
	res, err := s.httpClient.Get(imageUrl)
	if err != nil {
		logger.Error("Fail to finish http.Get", zap.String("imageUrl", imageUrl), zap.Error(err))
		return ""
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error("Fail to finish res.Body.Close", zap.Any("body", res.Body), zap.Error(err))
		}
	}()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Fail to finish ioutil.ReadAll", zap.Any("body", res.Body), zap.Error(err))
		return ""
	}
	return fmt.Sprintf("%x", md5.Sum(data))
}
