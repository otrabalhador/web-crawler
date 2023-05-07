package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldGetPageContent(t *testing.T) {
	page := Page{
		Content: `Hello crawler`,
		Url:     "https://foo.com",
	}
	pageMap := map[string]Page{page.Url: page}

	webClient := NewFakeWebClient(pageMap)
	crawler := NewCrawler(webClient, NewFakeRepository(), NewFakeExtractor(nil))

	crawler.Execute(page.Url)

	assert.Equal(t, []string{page.Url}, webClient.CalledUrls)
}

func TestShouldSavePage(t *testing.T) {
	page := Page{
		Content: `Hello crawler`,
		Url:     "https://foo.com",
	}
	pageMap := map[string]Page{page.Url: page}

	repository := NewFakeRepository()
	crawler := NewCrawler(NewFakeWebClient(pageMap), repository, NewFakeExtractor(nil))
	crawler.Execute(page.Url)

	assert.Equal(t, []Page{page}, repository.SavedPages)
}

func TestShouldExtractUrls(t *testing.T) {
	page := Page{
		Content: `Hello crawler`,
		Url:     "https://foo.com",
	}
	pageMap := map[string]Page{page.Url: page}

	extractor := NewFakeExtractor(nil)
	crawler := NewCrawler(NewFakeWebClient(pageMap), NewFakeRepository(), extractor)
	crawler.Execute(page.Url)

	assert.Equal(t, []Page{page}, extractor.ExtractedPages)
}

func TestShouldCrawlAgainForEachChildUrl(t *testing.T) {
	rootPageUrl := "https://foo.com"
	childPageUrl1 := "https://foo.com/bar"
	childPageUrl2 := "https://foo.com/baz"
	pages := []Page{
		{
			Content: `Hello crawler`,
			Url:     rootPageUrl,
		},
		{
			Content: `I bar`,
			Url:     childPageUrl1,
		},
		{
			Content: `I am baz`,
			Url:     childPageUrl2,
		},
	}

	pageMap := map[string]Page{
		rootPageUrl:   pages[0],
		childPageUrl1: pages[1],
		childPageUrl2: pages[2],
	}

	extractionMap := map[Page][]string{
		pages[0]: {
			childPageUrl1,
			childPageUrl2,
		},
	}

	extractor := NewFakeExtractor(extractionMap)
	webClient := NewFakeWebClient(pageMap)
	repository := NewFakeRepository()
	crawler := NewCrawler(webClient, repository, extractor)

	crawler.Execute(rootPageUrl)

	expectedCalledUrls := []string{rootPageUrl, childPageUrl1, childPageUrl2}
	assert.Equal(t, expectedCalledUrls, webClient.CalledUrls)
	assert.Equal(t, pages, repository.SavedPages)
	assert.Equal(t, pages, extractor.ExtractedPages)
}

// FakeWebClient

type FakeWebClient struct {
	Pages      map[string]Page
	CalledUrls []string
	CallCount  int
}

func NewFakeWebClient(pageMap map[string]Page) *FakeWebClient {
	return &FakeWebClient{
		Pages:      pageMap,
		CalledUrls: []string{},
		CallCount:  0,
	}
}

func (f *FakeWebClient) SetupPage(url string, page Page) {
	f.Pages[url] = page
}

func (f *FakeWebClient) GetPageContent(url string) (Page, error) {
	f.CallCount++
	f.CalledUrls = append(f.CalledUrls, url)

	page := f.Pages[url]
	return page, nil
}

// FakeRepository

type FakeRepository struct {
	CallCount  int
	SavedPages []Page
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		CallCount:  0,
		SavedPages: []Page{},
	}
}

func (f *FakeRepository) Save(page Page) error {
	f.CallCount++
	f.SavedPages = append(f.SavedPages, page)

	return nil
}

// FakeExtractor

type FakeExtractor struct {
	CallCount      int
	ExtractedPages []Page
	UrlMap         map[Page][]string
}

func NewFakeExtractor(extractedUrlMap map[Page][]string) *FakeExtractor {
	return &FakeExtractor{
		UrlMap:         extractedUrlMap,
		CallCount:      0,
		ExtractedPages: []Page{},
	}
}

func (f *FakeExtractor) Extract(page Page) []string {
	f.CallCount++
	f.ExtractedPages = append(f.ExtractedPages, page)

	if f.ExtractedPages == nil {
		return nil
	}

	return f.UrlMap[page]
}
