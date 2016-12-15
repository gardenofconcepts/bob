package main

import (
	log "github.com/Sirupsen/logrus"
	"testing"
)

func TestReader(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	builds := NewReader("assets").read("*.build.yml", []string{"**"}, []string{"**/exclude"})

	if len(builds) != 2 {
		t.Error("Expected 2, got ", len(builds))
	}

	if builds[0].Name != "blubpuuups" {
		t.Error("Expected 'blubpuuups', got ", builds[0].Name)
	}
}
