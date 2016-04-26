package main

import "testing"

func TestParser(t *testing.T) {
	build := Parser("assets/high_priority.build.yml").parse()

	if build.Name != "blubpuuups" {
		t.Error("Expected 'blubpuuups', got ", build.Name)
	}
}
