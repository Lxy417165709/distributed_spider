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
	crawlUrlChannel chan string,
	crawlResultChannel chan *model.CrawlResult,
	source model.CrawlSource,
	crawlNodeNum int,
) {
	for {
		select {
		case <-s.ctx.Done():
			return
		case url := <-crawlUrlChannel:
			crawlResultChannel <- s.Crawl(url,source,crawlNodeNum)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *SpiderWorker) Crawl(visitUrl string,source model.CrawlSource,crawlNodeNum int) *model.CrawlResult {
	// 1. 请求 url，获得响应
	req, err := NewReq(visitUrl)
	if err != nil {
		return &model.CrawlResult{
			Url: visitUrl,
			Err: err,
		}
	}
	res, err := s.httpClient.Do(req)
	if err != nil {
		logger.Error("Fail to finish http.Get", zap.String("url", visitUrl), zap.Error(err))
		return &model.CrawlResult{
			Url: visitUrl,
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
			Url: visitUrl,
			Err: fmt.Errorf("status is %d, expected %d", res.StatusCode, http.StatusOK),
		}
	}

	// 2. 发现链接、处理链接
	crawlUrls, imageUrls := s.getCrawlUrlsAndImageUrls(res)
	if len(crawlUrls) == 0 && len(imageUrls) == 0 {
		return &model.CrawlResult{
			Url: visitUrl,
			Err: fmt.Errorf("response blank"),
		}
	}
	return &model.CrawlResult{
		Url:       visitUrl,
		ImageUrls: imageUrls,
		CrawlSource: source,
		CrawlNodeNum: crawlNodeNum,
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

