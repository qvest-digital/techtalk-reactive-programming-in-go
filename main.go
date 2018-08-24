package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tarent/techtalk-reactive-programming-in-go/crawl"
	"github.com/tarent/techtalk-reactive-programming-in-go/persistence"
)

func main() {
	crawlIn, crawlOut := crawl.StartCrawlWorker(3)
	fileIn, fileOut := persistence.StartFileWorker()

	go func() {
		for {
			crawlResult := <-crawlOut
			if crawlResult.Error != nil {
				log.Println(crawlResult.Error.Error())
				continue
			}

			fileName := getFileName(crawlResult.Url)
			fileIn <- persistence.FileAction{Path: fileName, Content: crawlResult.Data}
		}
	}()

	go func() {
		for {
			fileResult := <-fileOut
			if fileResult.Error != nil {
				log.Println(fileResult.Error.Error())
				continue
			}
			log.Printf("Wrote %d lines to %s", fileResult.Lines, fileResult.File.Name())
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
			log.Fatalln(scanner.Err().Error())
		}
	}
}

func getFileName(in string) string {
	return in[strings.Index(in, "://")+3:]
}
