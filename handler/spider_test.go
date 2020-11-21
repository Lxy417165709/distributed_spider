package handler

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"
	"spider/common/logger"
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
		nil,
	)
	logger.Info("Crawl test", zap.Any("result", spiderWorker.Crawl("http://www.baidu.com")))
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
