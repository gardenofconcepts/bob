package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bmatcuk/doublestar"
)

func match(patternList []string, path string) bool {
	for _, pattern := range patternList {
		matched, _ := doublestar.Match(pattern, path)

		log.Debug("Match file against pattern with result", path, pattern, matched)

		if matched {
			return true
		}
	}

	return false
}
