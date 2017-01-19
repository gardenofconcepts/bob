package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	build := Parser().load("assets/parser/test.build.yml")

	if build.Name != "test file" {
		t.Error("Expected 'blubpuuups', got ", build.Name)
	}

	t.Run("Default directory", func(t *testing.T) {
		build := Parser()
		build.Directory = "/var/www"
		build.Verify = Verify{
			Include: []string{"**"},
			Exclude: []string{},
		}

		build.determine()

		if build.Cwd != "/var/www" {
			t.Errorf("Expect %s, got %s", "/var/www", build.Cwd)
		}
		if build.Root != "/var/www" {
			t.Errorf("Expect %s, got %s", "/var/www", build.Root)
		}
		if build.Verify.Include[0] != "**" {
			t.Errorf("Expect %s, got %s", "/var/www/**", build.Verify.Include[0])
		}
	})

	t.Run("Current directory is default directory", func(t *testing.T) {
		build := BuildFile{
			Directory: "/var/www",
			Cwd:       ".",
			Root:      ".",
		}

		build.determine()

		if build.Cwd != "/var/www" {
			t.Errorf("Expect %s, got %s", "/var/www", build.Cwd)
		}
		if build.Root != "/var/www" {
			t.Errorf("Expect %s, got %s", "/var/www", build.Root)
		}
	})

	t.Run("Parent directory", func(t *testing.T) {
		build := BuildFile{
			Directory: "/var/www/test",
			Cwd:       "..",
			Root:      "..",
		}

		build.determine()

		if build.Cwd != "/var/www" {
			t.Errorf("Expect %s, got %s", "/var/www", build.Cwd)
		}
		if build.Root != "/var/www" {
			t.Errorf("Expect %s, got %s", "/var/www", build.Root)
		}
	})

	t.Run("Parent directory", func(t *testing.T) {
		build := BuildFile{
			Directory: "/var/www/test",
			Cwd:       "..",
			Root:      "..",
			Verify: Verify{
				Include: []string{"**"},
			},
		}

		build.determine()

		if build.Cwd != "/var/www" {
			t.Errorf("Expect %s, got %s", "/var/www", build.Cwd)
		}
		if build.Root != "/var/www" {
			t.Errorf("Expect %s, got %s", "/var/www", build.Root)
		}
	})
}
