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
}

type WebClient interface {
	GetPageContent(url string) (Page, error)
}

type Crawler struct {
	webClient  WebClient
	repository Repository
	extractor  Extractor
}

func NewCrawler(webClient WebClient, repository Repository, extractor Extractor) *Crawler {
	return &Crawler{
		webClient:  webClient,
		repository: repository,
		extractor:  extractor,
	}
}

func (c *Crawler) Execute(url string) {
	page, _ := c.webClient.GetPageContent(url)

	_ = c.repository.Save(page)

	_ = c.extractor.Extract(page)
}
