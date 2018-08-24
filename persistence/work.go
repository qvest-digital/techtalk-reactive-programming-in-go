package persistence

import "os"

// FileAction describes the data expected in the input channel
type FileAction struct {
	Path    string
	Content []string
}

// FileResult describes the data in the output channel
type FileResult struct {
	File  *os.File
	Lines int
	Error error
}

// StartFileWorker starts a goroutine listening on its input channel, creates files with given path and content and
// then passes the files or errors via its output channel
func StartFileWorker() (chan FileAction, chan FileResult) {
	//Create result channel
	inputChan := make(chan FileAction, 100)
	outputChan := make(chan FileResult, 100)

	// Start working
	go func() {
		for {
			fileAction := <-inputChan
			outputChan <- writeFile(fileAction)
		}
	}()
	return inputChan, outputChan
}
