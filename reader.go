package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/bradfitz/slice"
)

func Read(path string, glob string) []BuildFile {
	matches := []BuildFile{}

	filepath.Walk(path, func (path string, file os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if file.IsDir() {
			return nil
		}

		matched, err := filepath.Match(glob, file.Name())

		if err != nil {
			fmt.Println(err)
			return err
		}

		if matched {
			build := Parser(path)
			matches = append(matches, *build)
		}

		return nil
	})

	slice.Sort(matches[:], func(i, j int) bool {
	    return matches[i].Priority < matches[j].Priority
	})

	return matches
}