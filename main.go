package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tarent/techtalk-reactive-programming-in-go/crawl"
	"github.com/tarent/techtalk-reactive-programming-in-go/persistence"
)

func main() {
	crawlIn, crawlOut := crawl.StartCrawlWorker(2)
	fileIn, fileOut := persistence.StartFileWorker()

	go func() {
		for {
			crawledBody := <-crawlOut
			if crawledBody.Error != nil {
				fmt.Println(crawledBody.Error.Error())
				os.Exit(1)
			}
			fileName := crawledBody.Url[strings.Index(crawledBody.Url, "://")+3:]
			fileIn <- persistence.FileAction{Path: fileName, Content: crawledBody.Data}
		}
	}()

	go func() {
		for {
			fileResult := <-fileOut
			if fileResult.Error != nil {
				fmt.Println(fileResult.Error.Error())
				os.Exit(1)
			}
		}
	}()

	// Scan command line input line by line
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter URLs line by line. Type 'exit' to quit.")
		for scanner.Scan() {
			if scanner.Text() == "exit" {
				os.Exit(0)
			}
			crawlIn <- crawl.CrawlerAction{Url: scanner.Text()}
		}
		if scanner.Err() != nil {
			fmt.Println(scanner.Err().Error())
			os.Exit(1)
		}
	}

}
