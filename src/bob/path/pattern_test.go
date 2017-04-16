package path

import (
	"testing"
	"path/filepath"
)

func TestPattern(t *testing.T) {

	if MatchPattern("src/*", "src/blub/test.js", "") != false {
		t.Error("Path doesn't match pattern")
	}

	if MatchPattern("/var/www/node_modules/**", "/var/www/node_modules/file.js", "/var/www") != false {
		t.Error("Path doesn't match pattern")
	}

	//path, _ := filepath.Abs("../../..")
	//
	//if MakePathRelative(path, "src/bob", "archive") != "" {
	//	t.Error("Path is not relative")
	//}

}
