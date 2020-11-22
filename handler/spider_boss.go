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
	addressChannel     chan *model.Address
	crawlNodeNum       int
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
		addressChannel:     make(chan *model.Address, crawlUrlChannelCap),
		crawlSource:        crawlSource,
		crawlNodeNum:       crawlNodeNum,
	}
}

func (s *SpiderBoss) Run() {
	go s.ShowInfo()
	for _, supplier := range s.spiderSuppliers {
		go supplier.Run(s.addressChannel, s.crawlSource)
	}

	for _, worker := range s.spiderWorkers {
		go worker.Run(s.addressChannel, s.crawlResultChannel, s.crawlSource, s.crawlNodeNum)
	}
	for _, handler := range s.spiderHandlers {
		go handler.Run(s.crawlResultChannel)
	}
}

func (s *SpiderBoss) ShowInfo() {
	for {
		logger.Info("Channel show",
			zap.Int("Len of addressChannel", len(s.addressChannel)),
			zap.Int("Len of crawlResultChannel", len(s.crawlResultChannel)))
		time.Sleep(1 * time.Second)
	}

}

