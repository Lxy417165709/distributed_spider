package model

type CrawlResult struct {
	Url          string
	Err          error
	CrawlSource  CrawlSource
	AddressId    int
	ImageUrls    []string
	CrawlUrls    []string
	CrawlNodeNum int
}
