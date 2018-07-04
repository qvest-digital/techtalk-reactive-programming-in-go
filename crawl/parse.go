package crawl

import (
	"github.com/PuerkitoBio/goquery"
)

type CrawlerAction struct {
	Url string
}

type CrawlerResult struct {
	Url   string
	Data  []string
	Error error
}

// StartRequestWorker starts a goroutine listening on its input channel,
// gets resources from the given url
// then passes the response bodies or errors via its output channel
func StartCrawlWorker() (chan CrawlerAction, chan CrawlerResult) {
	//Create channels
	inputChan := make(chan CrawlerAction, 100)
	outputChan := make(chan CrawlerResult, 100)

	// Start working
	go func() {
		for {
			select {
			case crawlerAction, ok := <-inputChan:
				if !ok {
					break
				}
				doc, err := goquery.NewDocument(crawlerAction.Url)
				if err != nil {
					outputChan <- CrawlerResult{Error: err}
					continue
				}
				links := make([]string, 0)
				doc.Find("a").Each(func(index int, item *goquery.Selection) {
					link, _ := item.Find("a").Attr("href")
					links = append(links, link)

				})
				outputChan <- CrawlerResult{Data: links}
			}
		}
	}()

	return inputChan, outputChan
}
