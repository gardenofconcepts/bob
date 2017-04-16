package parser

import (
	"testing"
)

func TestParser(t *testing.T) {
	build := Parser().Load("test-fixtures/test.build.yml")

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

	t.Run("Read and parse constants", func(t *testing.T) {
		build := Parser()
		build.Directory = "/var/www"
		build.Constant = []Constant{
			{
				Constant: "OS_VERSION",
				Command:  "uname -r", // "4.10.8-1-ARCH"
			},
			{ // overwrite ARCH
				Constant: "ARCH",
				Command:  "uname -m", // "x86_64"
			},
		}
		build.Constraint = []Constraint{
			{
				Condition: "OS == 'linux'",
			},
			{
				Condition: "ARCH != 'x86_64'",
			},
			{
				Condition: "version_compare(OS_VERSION, '>= 4.10')",
			},
			{
				Condition: "version_compare(OS_VERSION, '<= 4.10')",
			},
			{
				Condition: "OS",
			},
		}

		build.determine()
		build.loadConstants()
		build.executeConstraints()

		// check constants

		if build.Constant[1].Result != "x86_64" { // overwritten, instead of amd64
			t.Errorf("Expect %s, got %s", "x86_64", build.Constant[1].Result)
		}


		// check constraints

		if build.Constraint[0].Result == false {
			t.Errorf("Expect %s, got %s", "true", build.Constraint[0].Result)
		}

		if build.Constraint[1].Result == true {
			t.Errorf("Expect %s, got %s", "false", build.Constraint[1].Result)
		}

		if build.Constraint[2].Result == false {
			t.Errorf("Expect %s, got %s", "false", build.Constraint[2].Result)
		}

		if build.Constraint[3].Result == true {
			t.Errorf("Expect %s, got %s", "false", build.Constraint[3].Result)
		}

		if build.Constraint[4].ResultString != "linux" {
			t.Errorf("Expect %s, got %s", "linux", build.Constraint[4].ResultString)
		}
	})
}
