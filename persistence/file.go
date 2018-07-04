package persistence

import "os"

// FileAction describes the data expected in the input channel
type FileAction struct {
	Path    string
	Content []byte
}

// FileResult describes the data in the output channel
type FileResult struct {
	File  *os.File
	Error error
}

// StartFileWorker starts a goroutine listening on its input channel, creates files with given path and content and
// then passes the files or errors via its output channel
func StartFileWorker(fileActions chan FileAction) chan FileResult {
	//Create result channel
	resultChan := make(chan FileResult, 100)

	// Start working
	go func() {
		for {
			select {
			case fileAction, ok := <-fileActions:
				if !ok {
					break
				}
				resultChan <- writeFile(fileAction)
			}
		}
	}()

	return resultChan
}

func writeFile(fileAction FileAction) FileResult {
	// Create new file
	file, err := os.Open(fileAction.Path)
	if err != nil {
		return FileResult{Error: err}
	}

	// Write to file
	_, err = file.Write(fileAction.Content)
	if err != nil {
		return FileResult{Error: err}
	}
	return FileResult{File: file}
}
