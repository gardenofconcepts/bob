package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
)

func Analyzer(directory string, include []string, exclude []string) (string, error) {
	log.Info("Analyze directory", directory)

	hashes := read(directory, include, exclude)
	hash, err := hashList(hashes)

	return hash, err
}

func read(baseDir string, includes []string, excludes []string) []string {
	hashes := []string{}

	filepath.Walk(baseDir, func(filePath string, file os.FileInfo, err error) error {
		if err != nil {
			log.Warning(err)
			return nil
		}

		if file.IsDir() {
			return nil
		}

		if matchList(includes, filePath, baseDir) && !matchList(excludes, filePath, baseDir) {
			hash, _ := hashFile(filePath)
			hashes = append(hashes, hash)

			log.Debug("Include file with hash", filePath, hash)
		} else {
			log.Debug("Skip file", filePath)
		}

		return nil
	})

	return hashes
}
