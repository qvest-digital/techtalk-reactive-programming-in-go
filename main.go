package main

import (
	"fmt"
	"os"

	"github.com/tarent/techtalk-reactive-programming-in-go/crawl"
	"github.com/tarent/techtalk-reactive-programming-in-go/persistence"
)

func main() {
	args := os.Args[1:]

	reqIn, reqOut := crawl.StartRequestWorker()
	crawlIn, crawlOut := crawl.StartCrawlWorker()
	fileIn, fileOut := persistence.StartFileWorker()

	// Push in some Urls
	for _, arg := range args {
		reqIn <- crawl.RequestAction{arg}
	}

	// Push Response bodies into the crawler workers
	go func() {
		for {
			select {
			case response, ok := <-reqOut:
				if !ok {
					break
				}
				if response.Error != nil {
					fmt.Println(response.Error.Error())
					continue
				}
				crawlIn <- crawl.CrawlerAction{response.Url, response.Content}
			}
		}
	}()

	// Push response bodies to file workers and print out newly found URLs
	go func() {
		for {
			select {
			case crawledBody, ok := <-crawlOut:
				if !ok {
					break
				}
				fileIn <- persistence.FileAction{Path: crawledBody.Url, Content: crawledBody.Content}
				for _, url := range crawledBody.NewUrls {
					fmt.Println(url)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case fileResult, ok := <-fileOut:
				if !ok {
					break
				}
				if fileResult.Error != nil {
					fmt.Println(fileResult.Error.Error())
				}
			}
		}
	}()
	for {
	}
}
