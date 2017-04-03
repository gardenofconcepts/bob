package main

import "testing"

func TestPattern(t *testing.T) {

	if matchPattern("src/*", "src/blub/test.js", "") != false {
		t.Error("Path doesn't match pattern")
	}

	if matchPattern("/var/www/node_modules/**", "/var/www/node_modules/file.js", "/var/www") != false {
		t.Error("Path doesn't match pattern")
	}

}
