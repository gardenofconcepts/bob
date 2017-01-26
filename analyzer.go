package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
)

func Analyzer(rootDir string, include []string, exclude []string) (string, error) {
	hashes := read(rootDir, include, exclude)
	hash := hashList(hashes)

	return hash, nil
}

func read(rootDir string, includes []string, excludes []string) map[string]string {
	hashList := make(map[string]string)

	log.WithFields(log.Fields{
		"cwd":     rootDir,
		"include": includes,
		"exclude": excludes,
	}).Info("Analyzing directory")

	filepath.Walk(rootDir, func(filePath string, file os.FileInfo, err error) error {
		if err != nil {
			log.Warning(err)
			return nil
		}

		if file.IsDir() {
			return nil
		}

		if matchList(includes, filePath, rootDir) && !matchList(excludes, filePath, rootDir) {
			hash, _ := hashFile(filePath)
			hashList[filePath] = hash

			log.WithFields(log.Fields{
				"file": filePath,
				"hash": hash,
			}).Debug("Append file")
		} else {
			log.WithField("file", filePath).Debug("Skip file")
		}

		return nil
	})

	return hashList
}
