package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bradfitz/slice"
	"os"
	"path/filepath"
)

type Result struct {
	path string
}

func NewReader(path string) *Result {
	return &Result{
		path: path,
	}
}

func (reader *Result) read(glob string, includes []string, excludes []string) []BuildFile {
	matches := []BuildFile{}

	log.WithFields(log.Fields{
		"path":    reader.path,
		"pattern": glob,
	}).Info("Searching for build files")

	filepath.Walk(reader.path, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Warning(err)

			return nil
		}

		if file.IsDir() && (!matchList(includes, path, reader.path) || matchList(excludes, path, reader.path)) {
			log.WithFields(log.Fields{
				"path":     path,
				"includes": includes,
				"excludes": excludes,
			}).Info("Skipping directory")

			return filepath.SkipDir
		}

		if file.IsDir() {
			return nil
		}

		log.WithField("path", path).Debug("Search")

		matched, err := filepath.Match(glob, file.Name())

		if err != nil {
			log.Warning(err)

			return err
		}

		if matched {
			build := Parser().load(path)
			matches = append(matches, *build)

			log.WithFields(log.Fields{
				"file":      build.File,
				"directory": build.Directory,
				"name":      build.Name,
				"priority":  build.Priority,
			}).Info("Found build")
		}

		return nil
	})

	log.Debug("Sorting build files...")

	slice.Sort(matches[:], func(i, j int) bool {
		return matches[i].Priority < matches[j].Priority
	})

	return matches
}
