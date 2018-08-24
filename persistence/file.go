package persistence

import "os"

const basePath = "sites/"

// FileAction describes the data expected in the input channel
type FileAction struct {
	Path    string
	Content []string
}

// FileResult describes the data in the output channel
type FileResult struct {
	File  *os.File
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

func writeFile(fileAction FileAction) FileResult {
	createDirIfNotExist(basePath)

	// Create new file
	file, err := os.Create(basePath + fileAction.Path)
	if err != nil {
		return FileResult{Error: err}
	}

	// Write to file
	for _, line := range fileAction.Content {
		_, err = file.Write([]byte(line + "\n"))
		if err != nil {
			return FileResult{Error: err}
		}
	}
	return FileResult{File: file}
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
