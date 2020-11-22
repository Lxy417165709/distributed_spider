package handler

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"
	"spider/common/logger"
	"spider/model"
	"time"
)

type SpiderWorker struct {
	httpClient      *http.Client
	imageUrlFilters []FilterFunction
	crawlUrlFilters []FilterFunction
	ctx             context.Context
	cancelFunc      context.CancelFunc
}

func NewSpiderWorker(
	imageUrlFilters []FilterFunction,
	crawlUrlFilters []FilterFunction,
	ctx context.Context,
	cancelFunc context.CancelFunc,

) *SpiderWorker {
	return &SpiderWorker{
		httpClient:      NewClient(),
		imageUrlFilters: imageUrlFilters,
		crawlUrlFilters: crawlUrlFilters,
		ctx:             ctx,
		cancelFunc:      cancelFunc,
	}
}

func (s *SpiderWorker) Run(
	addressChannel chan *model.Address,
	crawlResultChannel chan *model.CrawlResult,
	source model.CrawlSource,
	crawlNodeNum int,
) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case address := <-addressChannel:
			crawlResultChannel <- s.Crawl(address)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *SpiderWorker) Crawl(address *model.Address) *model.CrawlResult {
	// 1. 请求 url，获得响应
	req, err := NewReq(address.Url)
	if err != nil {
		return &model.CrawlResult{
			Address: address,
			Err: err,
		}
	}
	res, err := s.httpClient.Do(req)
	if err != nil {
		logger.Error("Fail to finish http.Get", zap.String("url", address.Url), zap.Error(err))
		return &model.CrawlResult{
			Address: address,
			Err: err,
		}
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			logger.Error("Fail to finish res.Body.Close", zap.Any("body", res.Body), zap.Error(err))
		}
	}()
	if res.StatusCode != http.StatusOK {
		return &model.CrawlResult{
			Address: address,
			Err: fmt.Errorf("status is %d, expected %d", res.StatusCode, http.StatusOK),
		}
	}

	// 2. 发现链接、处理链接
	crawlUrls, imageUrls := s.getCrawlUrlsAndImageUrls(res)
	if len(crawlUrls) == 0 && len(imageUrls) == 0 {
		return &model.CrawlResult{
			Address: address,
			Err: fmt.Errorf("response blank"),
		}
	}
	return &model.CrawlResult{
		Address: address,
		ImageUrls: imageUrls,
		CrawlUrls: crawlUrls,
	}
}

func (s *SpiderWorker) getCrawlUrlsAndImageUrls(res *http.Response) ([]string, []string) {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Error("Fail to finish goquery.NewDocumentFromReader", zap.Any("body", res.Body), zap.Error(err))
		return nil, nil
	}

	crawlUrls := make([]string, 0)
	imageUrls := make([]string, 0)
	for _, filterFunction := range s.crawlUrlFilters {
		crawlUrls = append(crawlUrls, filterFunction(doc, res)...)
	}
	for _, filterFunction := range s.imageUrlFilters {
		imageUrls = append(imageUrls, filterFunction(doc, res)...)
	}
	return crawlUrls, imageUrls
}

