package web_client

import (
	"io"
	"net/http"
	netUrl "net/url"
	"web-crawler/internal"
)

type WebClient struct {
}

func NewWebClient() *WebClient {
	return &WebClient{}
}

func (w WebClient) GetPageContent(url *netUrl.URL) (internal.Page, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return internal.Page{}, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return internal.Page{}, err
	}

	return internal.Page{
		Url:     url,
		Content: string(data),
	}, nil
}
