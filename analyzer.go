package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
)

func Analyzer(directory string, include []string, exclude []string) (string, error) {
	log.WithFields(log.Fields{
		"cwd":     directory,
		"include": include,
		"exclude": exclude,
	}).Info("Analyzing directory")

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

			log.WithFields(log.Fields{
				"file": filePath,
				"hash": hash,
			}).Debug("Include file with hash")
		} else {
			log.WithField("file", filePath).Debug("Skip file")
		}

		return nil
	})

	return hashes
}
