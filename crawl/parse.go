package crawl

type CrawlerAction struct {
	Content []byte
}

type CrawlerResult struct {
	Content []byte
	NewUrls []string
	Error   error
}

// StartRequestWorker starts a goroutine listening on its input channel,
// gets resources from the given url
// then passes the response bodies or errors via its output channel
func StartRequestWorker(crawlerActions chan CrawlerAction) chan CrawlerResult {
	//Create result channel
	resultChan := make(chan CrawlerResult, 100)

	// Start working
	go func() {
		for {
			select {
			case crawlerAction, ok := <-crawlerActions:
				if !ok {
					break
				}
				resultChan <- CrawlerResult{Content: crawlerAction.Content}
			}
		}
	}()

	return resultChan
}
