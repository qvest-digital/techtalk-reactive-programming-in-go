package crawl

import (
	"io/ioutil"
	"net/http"
)

type RequestAction struct {
	Url string
}

type RequestResult struct {
	Url     string
	Content []byte
	Error   error
}

// StartRequestWorker starts a goroutine listening on its input channel,
// gets resources from the given url
// then passes the response bodies or errors via its output channel
func StartRequestWorker() (chan RequestAction, chan RequestResult) {
	//Create channels
	inputChan := make(chan RequestAction, 100)
	outputChan := make(chan RequestResult, 100)

	// Start working
	go func() {
		for {
			select {
			case requestAction, ok := <-inputChan:
				if !ok {
					break
				}
				resp, err := http.Get(requestAction.Url)
				if err != nil {
					outputChan <- RequestResult{Error: err}
					continue
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					outputChan <- RequestResult{Error: err}
					continue
				}
				outputChan <- RequestResult{Url: requestAction.Url, Content: body}
			}
		}
	}()

	return inputChan, outputChan
}
