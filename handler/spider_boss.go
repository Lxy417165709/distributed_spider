package handler

import (
	"go.uber.org/zap"
	"spider/common/logger"
	"spider/model"
	"time"
)

type SpiderBoss struct {
	spiderSuppliers    []*SpiderSupplier
	spiderWorkers      []*SpiderWorker
	spiderHandlers     []*SpiderHandler
	crawlResultChannel chan *model.CrawlResult
	crawlUrlChannel    chan string
	crawlNodeNum int
	crawlSource        model.CrawlSource
}

func NewSpiderBoss(
	suppliers []*SpiderSupplier,
	workers []*SpiderWorker,
	handlers []*SpiderHandler,
	crawlResultChannelCap int,
	crawlUrlChannelCap int,
	crawlSource model.CrawlSource,
	crawlNodeNum int,
) *SpiderBoss {
	return &SpiderBoss{
		spiderSuppliers:    suppliers,
		spiderWorkers:      workers,
		spiderHandlers:     handlers,
		crawlResultChannel: make(chan *model.CrawlResult, crawlResultChannelCap),
		crawlUrlChannel:    make(chan string, crawlUrlChannelCap),
		crawlSource:        crawlSource,
		crawlNodeNum:crawlNodeNum,
	}
}

func (s *SpiderBoss) Run() {
	go s.ShowInfo()
	for _, supplier := range s.spiderSuppliers {
		go supplier.Run(s.crawlUrlChannel, s.crawlSource)
	}

	for _, worker := range s.spiderWorkers {
		go worker.Run(s.crawlUrlChannel, s.crawlResultChannel, s.crawlSource, s.crawlNodeNum)
	}
	for _, handler := range s.spiderHandlers {
		go handler.Run(s.crawlResultChannel)
	}
}

func (s *SpiderBoss) ShowInfo() {
	for {
		logger.Info("Channel show",
			zap.Int("Len of crawlUrlChannel", len(s.crawlUrlChannel)),
			zap.Int("Len of crawlResultChannel", len(s.crawlResultChannel)))
		time.Sleep(1 * time.Second)
	}

}

