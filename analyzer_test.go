package main

import (
	"testing"
)

func TestAnalyzer(t *testing.T) {

	expectedHash := "9c90746368d07aa971e4ddd37e7d5c98"
	hash, _ := Analyzer("assets/analyzer", []string{"**"}, []string{"*.js"})

	if hash != expectedHash {
		t.Fatalf("Expect hash '%s', instead of %s", expectedHash, hash)
	}
}
