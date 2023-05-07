package repository

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"web-crawler/internal"
)

func TestSaveOnRoot_ShouldCreateNewFolderWithIndexHtml(t *testing.T) {
	rootDir := "foo.bar"

	page := internal.Page{
		Url: internal.URL{
			Host: rootDir,
			Path: "baz/qux",
		},
		Content: "foo",
	}

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(rootDir)

	repository := NewRepository("")
	err := repository.Save(page)
	assert.Nil(t, err)

	destinationFilePath := "foo.bar/baz/qux/index.html"
	bytes, err := os.ReadFile(destinationFilePath)
	assert.Nil(t, err)
	assert.Equal(t, "foo", string(bytes))
}

func TestIsAlreadySave_ShouldReturnFalseWhenNone(t *testing.T) {
	repository := NewRepository("")
	isAlreadySaved := repository.IsAlreadySaved(internal.URL{
		Host: "foo.bar",
		Path: "baz/qux",
	})
	assert.False(t, isAlreadySaved)
}

func TestIsAlreadySave_ShouldReturnTrueWhenHasAlreadySaved(t *testing.T) {
	rootDir := "foo.bar"
	url := internal.URL{Host: rootDir, Path: "baz/qux"}
	page := internal.Page{Url: url, Content: "foo"}

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(rootDir)

	repository := NewRepository("")
	_ = repository.Save(page)

	isAlreadySaved := repository.IsAlreadySaved(url)
	assert.True(t, isAlreadySaved)
}

func TestGetPage_ShouldReturnEmptyPageWhenNotFound(t *testing.T) {
	repository := NewRepository("")
	page := repository.GetPage(internal.URL{
		Host: "foo.bar",
		Path: "baz/qux",
	})
	assert.Empty(t, page)
}

func TestGetPage_ShouldReturnPageWhenFound(t *testing.T) {
	rootDir := "foo.bar"
	url := internal.URL{Host: rootDir, Path: "baz/qux"}
	page := internal.Page{Url: url, Content: "foo"}

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(rootDir)

	repository := NewRepository("")
	_ = repository.Save(page)

	actualPage := repository.GetPage(url)
	assert.Equal(t, page, actualPage)
}
