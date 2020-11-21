package model

import "time"

type CrawlStatus int

const (
	NotCrawl CrawlStatus = iota
	HadCrawled
	FailCrawled
)

type CrawlSource string

type Address struct {
	Id           int
	URL          string
	CrawlStatus  CrawlStatus
	CrawlSource  CrawlSource
	CrawlNodeNum int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (*Address) TableName() string {
	return "spd_address"
}

type Image struct {
	Id        int
	URL       string
	MD5       string
	AddressId int
	CreatedAt time.Time
}

func (*Image) TableName() string {
	return "spd_image"
}

type CrawlResult struct {
	URL       string
	AddressId int
	Err       error
	ImageUrls []string
	CrawlUrls []string
}
