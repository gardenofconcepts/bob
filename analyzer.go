package main

import (
	"fmt"
	"path/filepath"
	"os"
)

func Analyzer(directory string, include []string, exclude []string) (string, error) {
	fmt.Println("Analyze directory", directory);

	hashes		:= read(directory, include, exclude)
	hash, err	:= hashList(hashes)

	return hash, err
}

func read(path string, include []string, exclude []string) []string {
	hashes := []string{}

	filepath.Walk(path, func (path string, file os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if file.IsDir() {
			return nil
		}

		if match(include, path) && !match(exclude, path) {
			hash, _ := hashFile(path)
			hashes = append(hashes, hash)

			fmt.Println("Include file with hash", path, hash)
		} else {
			fmt.Println("Skip file", path)
		}

		return nil
	})

	return hashes
}
