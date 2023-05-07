package extractor

import (
	"github.com/stretchr/testify/assert"
	netUrl "net/url"
	"testing"
	"web-crawler/internal"
)

func TestShouldReturnEmptyListIfNoUrlsAreFound(t *testing.T) {
	page := internal.Page{
		Url: &netUrl.URL{
			Host: "foo.bar",
			Path: "baz/qux",
		},
		Content: "<html>foo</html>",
	}

	extractor := NewExtractor()
	urls := extractor.Extract(page)

	assert.Empty(t, urls)
}

func TestShouldIgnoreUrlOfOtherHosts(t *testing.T) {
	page := internal.Page{
		Url: &netUrl.URL{
			Host: "foo.bar",
			Path: "baz/qux",
		},
		Content: `
<html>
<body>
<a href=bar.foo>A link to bar.foo</a>
</body>
</html>
`,
	}

	extractor := NewExtractor()
	urls := extractor.Extract(page)

	assert.Empty(t, urls)
}

func TestShouldReturnAllMatchedUrls(t *testing.T) {
	page := internal.Page{
		Url: &netUrl.URL{
			Host: "foo.bar",
			Path: "baz/qux",
		},
		Content: `
<html>
<body>
<a href="https://foo.bar"">A link to foo.bar</a>
<a href="https://foo.bar/baz"">A link to foo.bar/baz</a>
<a href="https://foo.bar/baz/qux"">A link to foo.bar/baz/qux</a>
</body>
</html>
`,
	}

	extractor := NewExtractor()
	urls := extractor.Extract(page)

	expectedUrls := []*netUrl.URL{
		parseUrl("https://foo.bar"),
		parseUrl("https://foo.bar/baz"),
		parseUrl("https://foo.bar/baz/qux"),
	}

	assert.Equal(t, expectedUrls, urls)
}

func TestMatchRelativeUrls(t *testing.T) {
	rootUrl := parseUrl("https://foo.bar/baz/qux")
	page := internal.Page{
		Url: rootUrl,
		Content: `
<html>
<body>
<a href="https://foo.bar"">A link to foo.bar</a>
<a href="/baz"">A link to foo.bar/baz</a>
<a href="/baz/qux"">A link to foo.bar/baz/qux</a>
</body>
</html>
`,
	}

	extractor := NewExtractor()
	urls := extractor.Extract(page)

	expectedUrls := []*netUrl.URL{
		parseUrl("https://foo.bar"),
		parseUrl("https://foo.bar/baz"),
		parseUrl("https://foo.bar/baz/qux"),
	}

	assert.Equal(t, expectedUrls, urls)
}
