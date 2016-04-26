package main

import "testing"

func TestPattern(t *testing.T) {

	if match("src/*", "src/blub/test.js", "") != false {
		t.Error()
	}

}
