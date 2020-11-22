package model

type CrawlResult struct {
	Address          *Address
	Err          error
	ImageUrls    []string
	CrawlUrls    []string
}
