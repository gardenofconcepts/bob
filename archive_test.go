package main

import (
	"testing"
	"os"
)

func TestArchive(t *testing.T) {
	archive := "build/test.tar.gz"
	sourceDirectory := "assets/src"
	targetDirectory := "build"
	checkExistingFile := "build/assets/src/test.js";
	checkNonExistingFile := "build/assets/src/blub.gif";
	include := []string{"**"}
	exclude := []string{"**.gif"}

	err := NewArchive(archive).Compress(sourceDirectory, include, exclude)

	if err != nil {
		t.Error(err)
	}

	err = NewArchive(archive).Extract(targetDirectory)

	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(checkExistingFile); os.IsNotExist(err) {
		t.Error("File not found: ", checkExistingFile)
	}

	if _, err := os.Stat(checkNonExistingFile); os.IsExist(err) {
		t.Error("File found, but should not there: ", checkNonExistingFile)
	}
}
