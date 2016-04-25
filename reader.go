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

func (reader *Result) read(glob string) []BuildFile {
	matches := []BuildFile{}

	filepath.Walk(reader.path, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Warning(err)

			return nil
		}

		if file.IsDir() {
			return nil
		}

		matched, err := filepath.Match(glob, file.Name())

		if err != nil {
			log.Warning(err)

			return err
		}

		if matched {
			build := Parser(path).parse()
			matches = append(matches, *build)
		}

		return nil
	})

	slice.Sort(matches[:], func(i, j int) bool {
		return matches[i].Priority < matches[j].Priority
	})

	return matches
}
