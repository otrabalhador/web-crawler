package main

import (
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"web-crawler/internal"
	"web-crawler/internal/extractor"
	"web-crawler/internal/repository"
	"web-crawler/internal/web_client"
)

var Root string
var Destination string
var DummyMode bool

func Run(_ *cobra.Command, _ []string) {
	if DummyMode {

		rootUrl, webClient, repository, extractor := internal.GenerateDummyDependencies()
		crawler := internal.NewCrawler(webClient, repository, extractor)

		log.Printf("Starting on dummy mode")
		crawler.Execute(rootUrl)
	} else {
		rootUrl, err := url.Parse(Root)
		if err != nil {
			log.Fatal(err)
		}

		c := web_client.NewWebClient()
		r := repository.NewRepository(Destination)
		e := extractor.NewExtractor()

		crawler := internal.NewCrawler(c, r, e)

		log.Printf("Starting crawler on %v with destinatino on %v", rootUrl.String(), Destination)
		crawler.Execute(rootUrl)
	}
}

var rootCmd = &cobra.Command{
	Use:   "web-crawler",
	Short: "Web crawler",
	Long:  `Recursive mirroring web crawler`,
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd, args)
	},
}

func main() {
	rootCmd.Flags().BoolVar(&DummyMode, "dummy-mode", false, "Dummy mode: use dummy dependencies for web client, repository and extraction")
	_ = rootCmd.MarkFlagRequired("root")

	rootCmd.Flags().StringVarP(&Root, "root", "r", "", "Url to begin crawl")
	_ = rootCmd.MarkFlagRequired("root")

	rootCmd.Flags().StringVarP(&Destination, "destination", "d", "", "Destination folder")
	_ = rootCmd.MarkFlagRequired("destination")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
