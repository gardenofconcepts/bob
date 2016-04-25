package main

import "testing"

func TestReader(t *testing.T) {
	builds := Read("assets", "*.yml")

	if len(builds) != 2 {
		t.Error("Expected 2, got ", len(builds))
	}

	if builds[0].Name != "blubpuuups" {
		t.Error("Expected 'blubpuuups', got ", builds[0].Name)
	}
}
