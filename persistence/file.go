package persistence

import "os"

const basePath = "sites/"

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
	return FileResult{File: file, Lines: len(fileAction.Content)}
}

func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}
