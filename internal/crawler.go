package internal

import (
	"log"
	netUrl "net/url"
)

type Page struct {
	Url     *netUrl.URL
	Content string
}

type Extractor interface {
	Extract(page Page) []*netUrl.URL
}

type Repository interface {
	Save(page Page) error
	IsAlreadySaved(url *netUrl.URL) bool
	GetPage(url *netUrl.URL) Page
}

type WebClient interface {
	GetPageContent(url *netUrl.URL) (Page, error)
}

type Crawler struct {
	webClient     WebClient
	repository    Repository
	extractor     Extractor
	travelledUrls map[string]bool
}

func NewCrawler(webClient WebClient, repository Repository, extractor Extractor) *Crawler {
	return &Crawler{
		webClient:     webClient,
		repository:    repository,
		extractor:     extractor,
		travelledUrls: map[string]bool{},
	}
}

func (c *Crawler) Execute(url *netUrl.URL) {
	log.Printf("Crawling url %s", url)
	c.travelledUrls[url.String()] = true

	if c.repository.IsAlreadySaved(url) {
		log.Printf("Url %s is already saved", url)
		page := c.repository.GetPage(url)
		urls := c.extractor.Extract(page)
		for _, pageUrl := range urls {
			if _, ok := c.travelledUrls[pageUrl.String()]; ok {
				continue
			}

			c.Execute(pageUrl)
		}

		return
	}

	page, _ := c.webClient.GetPageContent(url)

	_ = c.repository.Save(page)

	urls := c.extractor.Extract(page)
	for _, pageUrl := range urls {
		if _, ok := c.travelledUrls[pageUrl.String()]; ok {
			continue
		}

		c.Execute(pageUrl)
	}
}
