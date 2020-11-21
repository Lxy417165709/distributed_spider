package model

type CrawlStatus int

const (
	NotCrawl CrawlStatus = iota
	HadCrawled
	FailCrawled
)

type CrawlSource string

type Address struct {
	Id           int
	Url          string
	CrawlStatus  CrawlStatus
	CrawlSource  CrawlSource
	CrawlNodeNum int
}

type Image struct {
	Id        int
	Url       int
	AddressId int
}

type CrawlResult struct {
	Url       string
	AddressId int
	Err       error
	ImageUrls []string
	CrawlUrls []string
}
