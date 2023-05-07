package web_client

import (
	"github.com/stretchr/testify/assert"
	netUrl "net/url"
	"testing"
)

// TODO: Create tests with mocked http client
func TestShouldCallUrl(t *testing.T) {
	webClient := NewWebClient()

	url, _ := netUrl.Parse("https://www.google.com")
	page, err := webClient.GetPageContent(url)
	assert.Nil(t, err)
	assert.NotEmpty(t, page)
}
