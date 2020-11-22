package model

import "time"

type CrawlStatus int

const (
	NotCrawl CrawlStatus = iota
	HadCrawled
	FailCrawled
	Crawling
)

type CrawlSource string

const (
	Baidu CrawlSource = "baidu"
)

type Address struct {
	Id           int
	Url          string
	CrawlStatus  CrawlStatus
	CrawlSource  CrawlSource
	CrawlNodeNum int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (*Address) TableName() string {
	return "spd_address"
}

