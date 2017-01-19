package main

import (
	log "github.com/Sirupsen/logrus"
	"testing"
)

func TestReader(t *testing.T) {
	log.SetLevel(log.ErrorLevel)

	builds := NewReader("assets/reader").read("*.build.yml", []string{"**"}, []string{"**/exclude"})

	if len(builds) != 2 {
		t.Error("Expected 2, got ", len(builds))
	}

	if builds[0].Name != "sub file" {
		t.Error("Expected 'sub file', got ", builds[0].Name)
	}

	if builds[1].Name != "test file" {
		t.Error("Expected 'test file', got ", builds[1].Name)
	}
}
