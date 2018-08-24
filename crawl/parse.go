package crawl

import (
	"net/url"

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
func StartCrawlWorker(numWorkers int) (chan CrawlerAction, chan CrawlerResult) {
	//Create channels
	inputChan := make(chan CrawlerAction, 100)
	outputChan := make(chan CrawlerResult, 100)

	// Start working
	go func() {
		for i := 0; i < numWorkers; i++ {
			select {
			case crawlerAction := <-inputChan:
				result := crawl(crawlerAction.Url)
				outputChan <- result
			}
		}
	}()
	return inputChan, outputChan
}

// Creates and crawls a document from a given url
func crawl(in string) CrawlerResult {
	url, err := normalizeUrl(in)
	if err != nil {
		return CrawlerResult{Error: err}
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return CrawlerResult{Error: err}
	}
	links := findLinks(doc)
	return CrawlerResult{Url: url, Data: links}
}

// Extracts links from a document
func findLinks(doc *goquery.Document) []string {
	links := make([]string, 0)
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		links = append(links, link)
	})
	return links
}

func normalizeUrl(in string) (out string, err error) {
	url, err := url.Parse(in)
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	out = url.String()
	return
}
