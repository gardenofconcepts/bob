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

func read(path string, include []string, exclude []string) []string {
	hashes := []string{}

	filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Warning(err)
			return nil
		}

		if file.IsDir() {
			return nil
		}

		if match(include, path) && !match(exclude, path) {
			hash, _ := hashFile(path)
			hashes = append(hashes, hash)

			log.Debug("Include file with hash", path, hash)
		} else {
			log.Debug("Skip file", path)
		}

		return nil
	})

	return hashes
}
