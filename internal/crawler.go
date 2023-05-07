package internal

import (
	"log"
	netUrl "net/url"
)

type URL netUrl.URL

type Page struct {
	Url     URL
	Content string
}

type Extractor interface {
	Extract(page Page) []URL
}

type Repository interface {
	Save(page Page) error
	IsAlreadySaved(url URL) bool
	GetPage(url URL) Page
}

type WebClient interface {
	GetPageContent(url URL) (Page, error)
}

type Crawler struct {
	webClient     WebClient
	repository    Repository
	extractor     Extractor
	travelledUrls map[URL]bool
}

func NewCrawler(webClient WebClient, repository Repository, extractor Extractor) *Crawler {
	return &Crawler{
		webClient:     webClient,
		repository:    repository,
		extractor:     extractor,
		travelledUrls: map[URL]bool{},
	}
}

func (c *Crawler) Execute(url URL) {
	log.Printf("Crawling url (host: %v, path: %v)", url.Host, url.Path)
	c.travelledUrls[url] = true

	if c.repository.IsAlreadySaved(url) {
		log.Printf("Url (host: %v, path: %v) is already saved", url.Host, url.Path)
		page := c.repository.GetPage(url)
		urls := c.extractor.Extract(page)
		for _, pageUrl := range urls {
			if _, ok := c.travelledUrls[pageUrl]; ok {
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
		if _, ok := c.travelledUrls[pageUrl]; ok {
			continue
		}

		c.Execute(pageUrl)
	}
}
