package extractor

import (
	"golang.org/x/net/html"
	netUrl "net/url"
	"strings"
	"web-crawler/internal"
)

type Extractor struct {
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e Extractor) Extract(page internal.Page) []*netUrl.URL {
	stringReader := strings.NewReader(page.Content)
	z := html.NewTokenizer(stringReader)

	urls := []*netUrl.URL{}
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return urls
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						url := parseUrl(attr.Val)
						if url.Host == page.Url.Host {
							urls = append(urls, url)
						}
					}
				}
			}
		}
	}
}

func parseUrl(textUrl string) *netUrl.URL {
	// TODO: Handle error
	parsedUrl, _ := netUrl.Parse(textUrl)
	return parsedUrl
}
