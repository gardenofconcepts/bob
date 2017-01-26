package main

import (
	"path/filepath"
	"testing"
)

func TestHash(t *testing.T) {
	files := []string{
		"4f172e3e56dca3ba3ec0be72224bfa83",
		"c4ca4238a0b923820dcc509a6f75849b",
	}

	for _, file := range files {
		hash, err := hashFile(filepath.Join("assets/hash", file))

		if err != nil {
			t.Errorf("Error while hashing file %s: %s", file, err)

			continue
		}

		if hash != file {
			t.Errorf("Expect hash %s, instead of %s", file, hash)
		}
	}

	hash := hashList(map[string]string{
		"a": "4f172e3e56dca3ba3ec0be72224bfa83",
		"b": "",
		"c": "c4ca4238a0b923820dcc509a6f75849b",
	})

	if hash != "ae92e8a612e1c24dd913cfbbbb23f924" {
		t.Errorf("Expected hash 'ae92e8a612e1c24dd913cfbbbb23f924', instead of %s", hash)
	}
}
