package handler

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net"
	"net/http"
	"spider/common/logger"
	"time"
)

type FilterFunction = func(doc *goquery.Document, res *http.Response) []string

func NewClient() *http.Client {
	//proxy := "127.0.0.1:52768"

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					logger.Error("Dial timeout", zap.Error(err))
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   10000,
			ResponseHeaderTimeout: time.Second * 20,
			//Proxy: func(_ *http.Request) (*url.URL, error) {
			//	return url.Parse("http://" + proxy)
			//},
		},
	}
	return client
}

func getAbsoluteUrl(s *goquery.Selection, request *http.Request, attr string) string {
	relativeUrl, _ := s.Attr(attr)
	urlStruct, _ := request.URL.Parse(relativeUrl)
	if urlStruct == nil {
		return ""
	}
	return urlStruct.String()
}

func NewReq(visitUrl string) (*http.Request, error) {
	req, err := http.NewRequest("GET", visitUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authority", "cn.pornhub.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("sec-fetch-site", "none")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Add("cookie", "platform_cookie_reset=pc; platform=pc; bs=u9yaxw53xyljp40ti7bipjxve64dhmyd; ss=362185722324793321; fg_9d12f2b2865de2f8c67706feaa332230=68378.100000; _ga=GA1.2.353548939.1605527397; d_uidb=dbf4167a-295c-42a2-83f7-15559dc67833; d_uid=dbf4167a-295c-42a2-83f7-15559dc67833; ua=080219060735c2a535b621c11879b95a; _gid=GA1.2.853342992.1605758827; fg_7133c455c2e877ecb0adfd7a6ec6d6fe=44953.100000; sm_track=v18mGN9uGvTXqCGo2ZkH52yI_7oENmaIp39X52Z4Hz8gsxNjA1ODU4NTg3NjAwMQ..; RNKEY=2326991*2948479:1110389720:3108730167:1")
	return req, nil

}
