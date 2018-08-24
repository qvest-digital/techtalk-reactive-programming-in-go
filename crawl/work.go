package crawl

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
			for {
				crawlerAction := <-inputChan
				result := crawl(crawlerAction.Url)
				outputChan <- result
			}
		}
	}()
	return inputChan, outputChan
}
