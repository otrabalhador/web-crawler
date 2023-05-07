package main

import (
	"github.com/spf13/cobra"
	"log"
	"web-crawler/internal"
)

var Root string
var Destination string
var DummyMode bool

func Run(_ *cobra.Command, _ []string) {
	log.Printf("Received root url as %v and destination folder as %v", Root, Destination)

	if DummyMode {
		log.Printf("Executing on dummy mode")

		rootUrl, webClient, repository, extractor := internal.GenerateDummyDependencies()
		crawler := internal.NewCrawler(webClient, repository, extractor)

		log.Printf("Starting")
		crawler.Execute(rootUrl)
	} else {
		log.Printf("Not ready for not-dummy mode")
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
