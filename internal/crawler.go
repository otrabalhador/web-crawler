package internal

type WebClient interface {
	GetPageContent(url string) error
}

type Crawler struct {
	webClient WebClient
}

func NewCrawler(webClient WebClient) *Crawler {
	return &Crawler{
		webClient: webClient,
	}
}

func (c *Crawler) Execute(url string) {
	c.webClient.GetPageContent(url)
}
