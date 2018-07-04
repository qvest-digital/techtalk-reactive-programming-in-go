package crawl

import "mvdan.cc/xurls"

type CrawlerAction struct {
	Url     string
	Content []byte
}

type CrawlerResult struct {
	Url     string
	Content []byte
	NewUrls []string
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
				urls := xurls.Relaxed().FindAllString(string(crawlerAction.Content), -1)
				outputChan <- CrawlerResult{Content: crawlerAction.Content, NewUrls: urls}
			}
		}
	}()

	return inputChan, outputChan
}
