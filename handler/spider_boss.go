package handler

import "spider/model"

type SpiderBoss struct {
	spiderWorkers      []*SpiderWorker
	crawlResultChannel chan *model.CrawlResult
	crawlUrlChannel    chan string
}

func NewSpiderBoss(workers []*SpiderWorker, crawlResultChannelCap int, crawlUrlChannelCap int) *SpiderBoss {
	return &SpiderBoss{
		spiderWorkers:      workers,
		crawlResultChannel: make(chan *model.CrawlResult, crawlResultChannelCap),
		crawlUrlChannel:    make(chan string, crawlUrlChannelCap),
	}
}
