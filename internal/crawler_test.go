package internal

import (
	"github.com/stretchr/testify/assert"
	netUrl "net/url"
	"testing"
)

func parseUrl(textUrl string) URL {
	rootUrl, _ := netUrl.Parse(textUrl)
	return URL(*rootUrl)
}

func TestShouldGetPageContent(t *testing.T) {
	page := Page{
		Content: `Hello crawler`,
		Url:     parseUrl("https://foo.com"),
	}
	pageMap := map[URL]Page{page.Url: page}

	webClient := NewFakeWebClient(pageMap)
	crawler := NewCrawler(webClient, NewFakeRepository(), NewFakeExtractor(nil))

	crawler.Execute(page.Url)

	assert.Equal(t, []URL{page.Url}, webClient.GetPageContentCalls)
}

func TestShouldSavePage(t *testing.T) {
	page := Page{
		Content: `Hello crawler`,
		Url:     parseUrl("https://foo.com"),
	}
	pageMap := map[URL]Page{page.Url: page}

	repository := NewFakeRepository()
	crawler := NewCrawler(NewFakeWebClient(pageMap), repository, NewFakeExtractor(nil))
	crawler.Execute(page.Url)

	assert.Equal(t, []Page{page}, repository.SaveCalls)
}

func TestShouldExtractUrls(t *testing.T) {
	page := Page{
		Content: `Hello crawler`,
		Url:     parseUrl("https://foo.com"),
	}
	pageMap := map[URL]Page{page.Url: page}

	extractor := NewFakeExtractor(nil)
	crawler := NewCrawler(NewFakeWebClient(pageMap), NewFakeRepository(), extractor)
	crawler.Execute(page.Url)

	assert.Equal(t, []Page{page}, extractor.ExtractCalls)
}

func TestShouldCrawlAgainForEachChildUrl(t *testing.T) {
	rootPageUrl := parseUrl("https://foo.com")
	childPageUrl1 := parseUrl("https://foo.com/bar")
	childPageUrl2 := parseUrl("https://foo.com/baz")
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

	pageMap := map[URL]Page{
		rootPageUrl:   pages[0],
		childPageUrl1: pages[1],
		childPageUrl2: pages[2],
	}

	extractionMap := map[Page][]URL{
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

	expectedCalledUrls := []URL{rootPageUrl, childPageUrl1, childPageUrl2}
	assert.Equal(t, expectedCalledUrls, webClient.GetPageContentCalls)
	assert.Equal(t, pages, repository.SaveCalls)
	assert.Equal(t, pages, extractor.ExtractCalls)
}

func TestShouldIgnoreAlreadyExtractedPages(t *testing.T) {
	url := parseUrl("https://foo.com")
	page := Page{
		Content: `Hello crawler. I have a circular dependency'`,
		Url:     url,
	}

	extractionMap := map[Page][]URL{
		page: {url},
	}

	extractor := NewFakeExtractor(extractionMap)

	pageMap := map[URL]Page{page.Url: page}

	repository := NewFakeRepository()
	webClient := NewFakeWebClient(pageMap)
	crawler := NewCrawler(webClient, repository, extractor)

	crawler.Execute(url)

	assert.Equal(t, 1, webClient.CallCount)
	assert.Equal(t, 1, repository.CallCount)
	assert.Equal(t, 1, extractor.CallCount)
}

func TestShouldResumePreviousWorkAndCrawlOnlyLeafPages(t *testing.T) {
	crawledPage1 := Page{
		Content: `I am root, but i was already crawled`,
		Url:     parseUrl("https://foo.com"),
	}
	crawledPage2 := Page{
		Content: `I am bar, but i was already crawled`,
		Url:     parseUrl("https://foo.com/bar"),
	}
	leafPage1 := Page{
		Content: `I am baz, and I haven't been yet crawled`,
		Url:     parseUrl("https://foo.com/baz"),
	}
	leafPage2 := Page{
		Content: `I am a qux, and I haven't been yet crawled`,
		Url:     parseUrl("https://foo.com/qux"),
	}

	pageMap := map[URL]Page{
		crawledPage1.Url: crawledPage1,
		crawledPage2.Url: crawledPage2,
		leafPage1.Url:    leafPage1,
		leafPage2.Url:    leafPage2,
	}

	extractionMap := map[Page][]URL{
		crawledPage1: {crawledPage1.Url, crawledPage2.Url, leafPage1.Url},
		crawledPage2: {crawledPage1.Url, leafPage2.Url},
		leafPage1:    {crawledPage1.Url},
		leafPage2:    {crawledPage1.Url},
	}

	extractor := NewFakeExtractor(extractionMap)

	repository := NewFakeRepository()
	repository.SetSavedPages([]Page{crawledPage1, crawledPage2})

	webClient := NewFakeWebClient(pageMap)
	crawler := NewCrawler(webClient, repository, extractor)

	crawler.Execute(crawledPage1.Url)

	assert.Equal(t, []URL{leafPage2.Url, leafPage1.Url}, webClient.GetPageContentCalls)
	assert.Equal(t, []Page{leafPage2, leafPage1}, repository.SaveCalls)
}

// FakeWebClient

type FakeWebClient struct {
	Pages               map[URL]Page
	GetPageContentCalls []URL
	CallCount           int
}

func NewFakeWebClient(pageMap map[URL]Page) *FakeWebClient {
	return &FakeWebClient{
		Pages:               pageMap,
		GetPageContentCalls: []URL{},
		CallCount:           0,
	}
}

func (f *FakeWebClient) SetupPage(url URL, page Page) {
	f.Pages[url] = page
}

func (f *FakeWebClient) GetPageContent(url URL) (Page, error) {
	f.CallCount++
	f.GetPageContentCalls = append(f.GetPageContentCalls, url)

	page := f.Pages[url]
	return page, nil
}

// FakeRepository

type FakeRepository struct {
	CallCount  int
	SaveCalls  []Page
	SavedPages []Page
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		CallCount: 0,
		SaveCalls: []Page{},
	}
}

func (f *FakeRepository) SetSavedPage(page Page) {
	f.SavedPages = append(f.SavedPages, page)
}

func (f *FakeRepository) SetSavedPages(pages []Page) {
	f.SavedPages = append(f.SavedPages, pages...)
}

func (f *FakeRepository) Save(page Page) error {
	f.CallCount++
	f.SaveCalls = append(f.SaveCalls, page)
	f.SavedPages = append(f.SavedPages, page)

	return nil
}

func (f *FakeRepository) IsAlreadySaved(url URL) bool {
	for _, page := range f.SavedPages {
		if url == page.Url {
			return true
		}
	}

	return false
}

func (f *FakeRepository) GetPage(url URL) Page {
	for _, page := range f.SavedPages {
		if url == page.Url {
			return page
		}
	}
	return Page{}
}

// FakeExtractor

type FakeExtractor struct {
	CallCount    int
	ExtractCalls []Page
	UrlMap       map[Page][]URL
}

func NewFakeExtractor(extractedUrlMap map[Page][]URL) *FakeExtractor {
	return &FakeExtractor{
		UrlMap:       extractedUrlMap,
		CallCount:    0,
		ExtractCalls: []Page{},
	}
}

func (f *FakeExtractor) Extract(page Page) []URL {
	f.CallCount++
	f.ExtractCalls = append(f.ExtractCalls, page)

	if f.ExtractCalls == nil {
		return nil
	}

	return f.UrlMap[page]
}
