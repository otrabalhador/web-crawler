package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeWebClient struct {
	CalledUrls []string
	CallCount  int
}

func NewFakeWebClient() *FakeWebClient {
	return &FakeWebClient{
		CalledUrls: []string{},
		CallCount:  0,
	}
}

func (f *FakeWebClient) GetPageContent(url string) error {
	f.CallCount++
	f.CalledUrls = append(f.CalledUrls, url)

	return nil
}

func TestShouldCallWebClientWithUrl(t *testing.T) {
	webClient := NewFakeWebClient()
	crawler := NewCrawler(webClient)

	url := "https://google.com"
	crawler.Execute(url)

	assert.Equal(t, 1, webClient.CallCount)
	assert.Equal(t, url, webClient.CalledUrls[0])
}
