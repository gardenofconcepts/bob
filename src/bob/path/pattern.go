package path

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bmatcuk/doublestar"
	"path/filepath"
)

func Match(includes []string, excludes []string, filePath string, baseDir string) bool {
	return MatchList(includes, filePath, baseDir) && !MatchList(excludes, filePath, baseDir)
}

func MatchList(patterns []string, path string, baseDir string) bool {
	for _, pattern := range CleanList(patterns) {
		if MatchPattern(pattern, path, baseDir) {
			return true
		}
	}

	return false
}

func MatchPattern(pattern string, path string, baseDir string) bool {
	path, _ = filepath.Rel(baseDir, path)
	result, _ := doublestar.Match(pattern, path)

	log.WithFields(log.Fields{
		"file":    path,
		"pattern": pattern,
		"result":  result,
	}).Debug("Matching file pattern")

	return result
}

func CleanList(s []string) []string {
	var r []string

	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}

	return r
}

func MakePathsRelative(rootDir string, baseDir string, paths []string) []string {
	result := []string{}

	for _, path := range paths {
		result = append(result, MakePathRelative(rootDir, baseDir, path))
	}

	return result
}

func MakePathRelative(rootDir string, baseDir string, path string) string {
	path = filepath.Join(baseDir, path)
	path, err := filepath.Rel(rootDir, path)

	if err != nil {
		panic(err)
	}

	return path
}
