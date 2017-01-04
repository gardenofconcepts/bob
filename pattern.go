package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bmatcuk/doublestar"
	"path/filepath"
)

func matchList(patterns []string, path string, baseDir string) bool {
	for _, pattern := range cleanList(patterns) {
		if match(pattern, path, baseDir) {
			return true
		}
	}

	return false
}

func match(pattern string, path string, baseDir string) bool {
	path, _ = filepath.Rel(baseDir, path)
	result, _ := doublestar.Match(pattern, path)

	log.WithFields(log.Fields{
		"file":    path,
		"pattern": pattern,
		"result":  result,
	}).Debug("Matching file pattern")

	return result
}

func cleanList(s []string) []string {
	var r []string

	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}

	return r
}
