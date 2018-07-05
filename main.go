package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tarent/techtalk-reactive-programming-in-go/crawl"
	"github.com/tarent/techtalk-reactive-programming-in-go/persistence"
)

func main() {
	args := os.Args[1:]

	crawlIn, crawlOut := crawl.StartCrawlWorker()
	fileIn, fileOut := persistence.StartFileWorker()

	fmt.Println("Push in some Urls")
	for _, arg := range args {
		crawlIn <- crawl.CrawlerAction{arg}
	}

	fmt.Println("Push response bodies to file workers and print out newly found URLs")
	go func() {
		for {
			select {
			case crawledBody, ok := <-crawlOut:
				if !ok {
					break
				}
				if crawledBody.Error != nil {
					fmt.Println(crawledBody.Error.Error())
					break
				}
				fmt.Println("Url: " + crawledBody.Url)
				fileName := crawledBody.Url[strings.Index(crawledBody.Url, "://")+3:]
				fmt.Println("Filename: " + fileName)
				fileIn <- persistence.FileAction{Path: fileName, Content: crawledBody.Data}
				for _, url := range crawledBody.Data {
					fmt.Println(url)
				}
			}
		}
		close(fileIn)
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

	fmt.Println("Wait loop")
	for {
	}
}
