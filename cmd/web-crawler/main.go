package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "web-crawler",
	Short: "Web crawler",
	Long:  `Recursive mirroring web crawler`,
	Run: func(cmd *cobra.Command, args []string) {
		Run(cmd, args)
	},
}

func Run(cmd *cobra.Command, _ []string) {
	rootUrl := cmd.Flag("root").Value
	destinationUrl := cmd.Flag("destination").Value
	fmt.Println(rootUrl)
	fmt.Println(destinationUrl)
}

func main() {
	rootCmd.PersistentFlags().StringP("root", "r", "", "Url to begin crawl")
	rootCmd.PersistentFlags().StringP("destination", "d", "", "Destination folder")
	if err := rootCmd.MarkFlagRequired("root"); err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.MarkFlagRequired("destination"); err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
