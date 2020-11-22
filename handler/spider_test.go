package handler

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"
	"spider/cache"
	"spider/common/env"
	"spider/common/logger"
	"spider/common/utils"
	"spider/dao"
	"spider/model"
	"strings"
	"testing"
)

func Init() {
	const confFilePath = "C:\\Users\\hasee\\Desktop\\spider\\configure\\alpha.json"
	utils.InitConfigure(confFilePath)
	dao.InitDB(env.Conf.MainDB.Link, env.Conf.MainDB.Name, env.Conf.MainDB.MaxConn)
	dao.CloseLog()
	cache.InitCache("120.26.162.39:20000", 0)
}


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
	logger.Info("Crawl test", zap.Any("result", spiderWorker.Crawl(&model.Address{
		Url: "http://www.baidu.com",
	})))
}

func TestSpiderBoss_Run(t *testing.T) {
	Init()
	cache.Spider.ReleaseSupplierLock()
	baiduSpider := baiduSpider()
	if err := dao.AddressDB.Create("http://baidu.com", model.Baidu, 0); err != nil {
		logger.Error("Fail to finish AddressDB.Create", zap.Error(err))
	}
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
	for i := 0; i < 5; i++ {
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
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithCancel(parentCtx)
		spiderHandler := NewSpiderHandler(NewClient(), ctx, cancel)
		spiderHandlers = append(spiderHandlers, spiderHandler)
	}

	spiderSuppliers := make([]*SpiderSupplier, 0)
	for i := 0; i < 2; i++ {
		ctx, cancel := context.WithCancel(parentCtx)
		spiderSupplier := NewSpiderSupplier(1,i, ctx, cancel)
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
