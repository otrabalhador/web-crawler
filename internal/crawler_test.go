package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldGetPageContent(t *testing.T) {
	page := Page{
		Content: `<html><body><h1>Hello crawler</h1></body></html>`,
		Url:     "https://google.com",
	}
	pageMap := map[string]Page{page.Url: page}

	webClient := NewFakeWebClient(pageMap)
	crawler := NewCrawler(webClient, NewFakeRepository())

	crawler.Execute(page.Url)

	assert.Equal(t, 1, webClient.CallCount)
	assert.Equal(t, page.Url, webClient.CalledUrls[0])
}

func TestShouldSavePage(t *testing.T) {
	page := Page{
		Content: `<html><body><h1>Hello crawler</h1></body></html>`,
		Url:     "https://google.com",
	}
	pageMap := map[string]Page{page.Url: page}

	webClient := NewFakeWebClient(pageMap)

	repository := NewFakeRepository()

	crawler := NewCrawler(webClient, repository)
	crawler.Execute(page.Url)

	assert.Equal(t, 1, repository.CallCount)
	assert.Equal(t, page, repository.SavedPages[0])
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
