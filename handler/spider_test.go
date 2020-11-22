package handler

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"
	"spider/common/logger"
	"spider/dao"
	"spider/model"
	"strings"
	"testing"
)

func TestSpiderWorker(t *testing.T) {
	f1 := func(doc *goquery.Document, res *http.Response) []string {
		result := make([]string, 0)
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			absoluteUrl := getAbsoluteUrl(s, res.Request, "href")
			if isHttpUrl(absoluteUrl) {
				result = append(result, absoluteUrl)
			}
		})
		return result
	}
	f2 := func(doc *goquery.Document, res *http.Response) []string {
		result := make([]string, 0)
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			absoluteUrl := getAbsoluteUrl(s, res.Request, "src")
			if isImageUrl(absoluteUrl) {
				result = append(result, absoluteUrl)
			}
			absoluteUrl = getAbsoluteUrl(s, res.Request, "data-src")
			if isImageUrl(absoluteUrl) {
				result = append(result, absoluteUrl)
			}
		})
		return result
	}
	f3 := func(doc *goquery.Document, res *http.Response) []string {
		result := make([]string, 0)
		doc.Find("div").Each(func(i int, s *goquery.Selection) {
			absoluteUrl := getAbsoluteUrl(s, res.Request, "data-bkg")
			if isImageUrl(absoluteUrl) {
				result = append(result, absoluteUrl)
			}
		})
		return result
	}
	spiderWorker := NewSpiderWorker(
		[]FilterFunction{f2, f3},
		[]FilterFunction{f1},
		context.Background(),
		nil,
	)
	logger.Info("Crawl test", zap.Any("result", spiderWorker.Crawl("http://www.baidu.com",model.Baidu,0)))
}

func TestSpiderBoss_Run(t *testing.T) {
	dao.InitDB("root:123456@tcp(120.26.162.39:40000)", "spider", 100)
	baiduSpider := baiduSpider()
	baiduSpider.Run()
	select {}
}

func baiduSpider() *SpiderBoss {
	parentCtx := context.Background()
	f1 := func(doc *goquery.Document, res *http.Response) []string {
		result := make([]string, 0)
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			absoluteUrl := getAbsoluteUrl(s, res.Request, "href")
			if isHttpUrl(absoluteUrl) {
				result = append(result, absoluteUrl)
			}
		})
		return result
	}
	f2 := func(doc *goquery.Document, res *http.Response) []string {
		result := make([]string, 0)
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			absoluteUrl := getAbsoluteUrl(s, res.Request, "src")
			if isImageUrl(absoluteUrl) {
				result = append(result, absoluteUrl)
			}
			absoluteUrl = getAbsoluteUrl(s, res.Request, "data-src")
			if isImageUrl(absoluteUrl) {
				result = append(result, absoluteUrl)
			}
		})
		return result
	}
	f3 := func(doc *goquery.Document, res *http.Response) []string {
		result := make([]string, 0)
		doc.Find("div").Each(func(i int, s *goquery.Selection) {
			absoluteUrl := getAbsoluteUrl(s, res.Request, "data-bkg")
			if isImageUrl(absoluteUrl) {
				result = append(result, absoluteUrl)
			}
		})
		return result
	}

	spiderWorkers := make([]*SpiderWorker, 0)
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithCancel(parentCtx)
		spiderWorker := NewSpiderWorker(
			[]FilterFunction{f2, f3},
			[]FilterFunction{f1},
			ctx,
			cancel,
		)
		spiderWorkers = append(spiderWorkers, spiderWorker)
	}

	spiderHandlers := make([]*SpiderHandler, 0)
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithCancel(parentCtx)
		spiderHandler := NewSpiderHandler(NewClient(), ctx, cancel)
		spiderHandlers = append(spiderHandlers, spiderHandler)
	}

	spiderSuppliers := make([]*SpiderSupplier, 0)
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithCancel(parentCtx)
		spiderSupplier := NewSpiderSupplier(10, ctx, cancel)
		spiderSuppliers = append(spiderSuppliers, spiderSupplier)
	}

	spiderBoss := NewSpiderBoss(
		spiderSuppliers,
		spiderWorkers,
		spiderHandlers,
		1000,
		1000,
		model.Baidu,
		0,
	)

	return spiderBoss
}

func isHttpUrl(url string) bool {
	return strings.HasPrefix(url, "http") && strings.Contains(url, "")
}

func isImageUrl(url string) bool {
	imgIdentifies := []string{"jpg", "png", "bmp", "photo"}
	for _, imgIdentify := range imgIdentifies {
		if strings.HasSuffix(url, imgIdentify) {
			return true
		}
	}
	return false
}
