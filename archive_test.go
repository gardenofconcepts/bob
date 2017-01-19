package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"testing"
)

func TestArchive(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	workingDirectory := "build/" + randomString(10)
	sourceDirectory := "assets/archive"
	targetDirectory := workingDirectory
	archive := workingDirectory + "/test.tar.gz"
	checkExistingFiles := []string{
		targetDirectory + "/src/test.js",
	}
	checkNonExistingFiles := []string{
		targetDirectory + "/src/blub.gif",
	}
	checkExistingLinks := []string{
		targetDirectory + "/src/this_is_a_link.js",
		targetDirectory + "/src/.bin/run",
	}
	checkExistingExecutables := []string{
		targetDirectory + "/src/.bin/run",
		targetDirectory + "/src/executables/run.sh",
	}

	include := []string{"**"}
	exclude := []string{"**/**.gif"}

	if err := os.Mkdir(workingDirectory, 0777); err != nil {
		t.Error("Cannot create working dir:", err)
	}

	if err := NewArchive(archive).Compress(sourceDirectory, include, exclude); err != nil {
		t.Error("Error while compressing data:", err)
	}

	if err := NewArchive(archive).Extract(targetDirectory); err != nil {
		t.Error("Error while extracting data:", err)
	}

	for _, path := range checkExistingFiles {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Error("File not found, but should be there:", path)
		}
	}

	for _, path := range checkNonExistingFiles {
		if _, err := os.Stat(path); err == nil {
			t.Error("File found, but should not there:", path)
		}
	}

	for _, path := range checkExistingLinks {
		info, err := os.Lstat(path)

		if err != nil {
			t.Error("Link not found, but should be there:", path)
		}

		if info.Mode()&os.ModeSymlink == 0 {
			t.Error("File found, but should be a link:", path)
		}
	}

	for _, path := range checkExistingExecutables {
		info, err := os.Stat(path)

		if err != nil || os.IsNotExist(err) {
			t.Error("File not found, but should be there:", path)
		} else if mode := info.Mode(); !info.IsDir() && mode&0111 == 0 {
			t.Error("File found, but should be executable:", path)
		}
	}
}
