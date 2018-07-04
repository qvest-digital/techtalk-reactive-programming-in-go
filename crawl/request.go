package crawl

import (
	"io"
	"net/http"
)

type RequestAction struct {
	Url string
}

type RequestResult struct {
	Content io.ReadCloser
	Error   error
}

// StartRequestWorker starts a goroutine listening on its input channel,
// gets resources from the given url
// then passes the response bodies or errors via its output channel
func StartRequestWorker(requestActions chan RequestAction) chan RequestResult {
	//Create result channel
	resultChan := make(chan RequestResult, 100)

	// Start working
	go func() {
		for {
			select {
			case requestAction, ok := <-requestActions:
				if !ok {
					break
				}
				resp, err := http.Get(requestAction.Url)
				if err != nil {
					resultChan <- RequestResult{Error: err}
				}
				resultChan <- RequestResult{Content: resp.Body}
			}
		}
	}()

	return resultChan
}
