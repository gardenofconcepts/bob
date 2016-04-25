package main

import "testing"

func TestParser(t *testing.T) {
	build := Parser("assets/high_priority.yml")

	if build.Name != "blubpuuups" {
		t.Error("Expected 'blubpuuups', got ", build.Name)
	}
}