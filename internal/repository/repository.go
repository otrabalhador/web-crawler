package repository

import (
	netUrl "net/url"
	"os"
	"path"
	"web-crawler/internal"
)

type Repository struct {
	directory string
}

func NewRepository(directory string) *Repository {
	return &Repository{
		directory: directory,
	}
}

func (r *Repository) Save(page internal.Page) error {
	destinationFolder := path.Join(r.directory, page.Url.Host, page.Url.Path)
	err := os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		return err
	}

	destinationFilePath := path.Join(destinationFolder, "index.html")

	file, err := os.OpenFile(destinationFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(page.Content))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) IsAlreadySaved(url *netUrl.URL) bool {
	filePath := path.Join(r.directory, url.Host, url.Path, "index.html")

	// TODO: Handle error
	f, _ := os.Stat(filePath)

	if f != nil {
		return true
	} else {
		return false
	}
}

func (r *Repository) GetPage(url *netUrl.URL) internal.Page {
	filePath := path.Join(r.directory, url.Host, url.Path, "index.html")

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return internal.Page{}
	}

	return internal.Page{
		Url:     url,
		Content: string(bytes),
	}
}
