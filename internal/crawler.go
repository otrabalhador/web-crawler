package internal

type Page struct {
	Url     string
	Content string
}

type Extractor interface {
	Extract(page Page) []string
}

type Repository interface {
	Save(page Page) error
	IsAlreadySaved(url string) bool
	GetPage(url string) Page
}

type WebClient interface {
	GetPageContent(url string) (Page, error)
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

func (c *Crawler) Execute(url string) {
	c.travelledUrls[url] = true

	if c.repository.IsAlreadySaved(url) {
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
