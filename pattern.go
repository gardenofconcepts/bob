package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bmatcuk/doublestar"
	"path/filepath"
)

func matchList(patterns []string, path string, baseDir string) bool {
	for _, pattern := range patterns {
		if match(pattern, path, baseDir) {
			return true
		}
	}

	return false
}

func match(pattern string, path string, baseDir string) bool {
	path, _ = filepath.Rel(baseDir, path)
	matched, _ := doublestar.Match(pattern, path)

	log.Debug("Match file against pattern with result", path, pattern, matched)

	return matched
}
