package internal

import (
	"log"
	netUrl "net/url"
)

func parseUrl(textUrl string) URL {
	rootUrl, _ := netUrl.Parse(textUrl)
	return URL(*rootUrl)
}

func GenerateDummyDependencies() (URL, *DummyWebClient, *DummyRepository, *DummyExtractor) {
	page1 := Page{
		Content: `I am root`,
		Url:     parseUrl("https://foo.com"),
	}
	page2 := Page{
		Content: `I am bar`,
		Url:     parseUrl("https://foo.com/bar"),
	}
	page3 := Page{
		Content: `I am baz`,
		Url:     parseUrl("https://foo.com/baz"),
	}
	page4 := Page{
		Content: `I am a qux`,
		Url:     parseUrl("https://foo.com/qux"),
	}

	pageMap := map[URL]Page{
		page1.Url: page1,
		page2.Url: page2,
		page3.Url: page3,
		page4.Url: page4,
	}

	extractionMap := map[Page][]URL{
		page1: {page1.Url, page2.Url, page3.Url},
		page2: {page1.Url, page4.Url},
		page3: {page1.Url},
		page4: {page1.Url},
	}

	extractor := NewDummyExtractor(extractionMap)

	repository := NewDummyRepository()

	webClient := NewDummyWebClient(pageMap)

	return page1.Url, webClient, repository, extractor
}

// DummyWebClient

type DummyWebClient struct {
	Pages               map[URL]Page
	GetPageContentCalls []URL
	CallCount           int
}

func NewDummyWebClient(pageMap map[URL]Page) *DummyWebClient {
	return &DummyWebClient{
		Pages:               pageMap,
		GetPageContentCalls: []URL{},
		CallCount:           0,
	}
}

func (f *DummyWebClient) SetupPage(url URL, page Page) {
	f.Pages[url] = page
}

func (f *DummyWebClient) GetPageContent(url URL) (Page, error) {
	log.Printf("Getting page content for %v", url)

	f.CallCount++
	f.GetPageContentCalls = append(f.GetPageContentCalls, url)

	page := f.Pages[url]
	return page, nil
}

// DummyRepository

type DummyRepository struct {
	CallCount  int
	SaveCalls  []Page
	SavedPages []Page
}

func NewDummyRepository() *DummyRepository {
	return &DummyRepository{
		CallCount: 0,
		SaveCalls: []Page{},
	}
}

func (f *DummyRepository) SetSavedPage(page Page) {
	f.SavedPages = append(f.SavedPages, page)
}

func (f *DummyRepository) SetSavedPages(pages []Page) {
	f.SavedPages = append(f.SavedPages, pages...)
}

func (f *DummyRepository) Save(page Page) error {
	log.Printf("Saving page for %v", page)

	f.CallCount++
	f.SaveCalls = append(f.SaveCalls, page)
	f.SavedPages = append(f.SavedPages, page)

	return nil
}

func (f *DummyRepository) IsAlreadySaved(url URL) bool {
	log.Printf("Is url %v already saved saved?", url)

	for _, page := range f.SavedPages {
		if url == page.Url {
			return true
		}
	}

	return false
}

func (f *DummyRepository) GetPage(url URL) Page {
	for _, page := range f.SavedPages {
		if url == page.Url {
			return page
		}
	}
	return Page{}
}

// DummyExtractor

type DummyExtractor struct {
	CallCount    int
	ExtractCalls []Page
	UrlMap       map[Page][]URL
}

func NewDummyExtractor(extractedUrlMap map[Page][]URL) *DummyExtractor {
	return &DummyExtractor{
		UrlMap:       extractedUrlMap,
		CallCount:    0,
		ExtractCalls: []Page{},
	}
}

func (f *DummyExtractor) Extract(page Page) []URL {
	f.CallCount++
	f.ExtractCalls = append(f.ExtractCalls, page)

	if f.ExtractCalls == nil {
		return nil
	}

	return f.UrlMap[page]
}
