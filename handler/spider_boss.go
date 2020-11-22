package handler

import "spider/model"

type SpiderBoss struct {
	spiderSuppliers    []*SpiderSupplier
	spiderWorkers      []*SpiderWorker
	spiderHandlers     []*SpiderHandler
	crawlResultChannel chan *model.CrawlResult
	crawlUrlChannel    chan string
	crawlSource        model.CrawlSource
}

func NewSpiderBoss(
	suppliers []*SpiderSupplier,
	workers []*SpiderWorker,
	handlers []*SpiderHandler,
	crawlResultChannelCap int,
	crawlUrlChannelCap int,
	crawlSource model.CrawlSource,
) *SpiderBoss {
	return &SpiderBoss{
		spiderSuppliers:    suppliers,
		spiderWorkers:      workers,
		spiderHandlers:     handlers,
		crawlResultChannel: make(chan *model.CrawlResult, crawlResultChannelCap),
		crawlUrlChannel:    make(chan string, crawlUrlChannelCap),
		crawlSource:        crawlSource,
	}
}

func (s *SpiderBoss) Run() {
	for _, supplier := range s.spiderSuppliers {
		go supplier.Run(s.crawlUrlChannel, s.crawlSource)
	}

	for _, worker := range s.spiderWorkers {
		go worker.Run(s.crawlUrlChannel, s.crawlResultChannel)
	}
	for _, handler := range s.spiderHandlers {
		go handler.Run(s.crawlResultChannel)
	}
}
