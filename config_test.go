package main

import "testing"

func TestConfig(t *testing.T) {
	config := ConfigReader()
	config.Read("assets/config.yml")

	if config.Pattern != "*.build.yml" {
		t.Errorf("Expect '*.build.yml', instead of %s", config.Pattern)
	}

	if config.Debug == false {
		t.Errorf("Expect 'true', instead of %s", config.Debug)
	}

	if config.Verbose == false {
		t.Errorf("Expect 'true', instead of %s", config.Verbose)
	}

	if config.Upload == true {
		t.Errorf("Expect 'true', instead of %s", config.Upload)
	}

	if config.Download == true {
		t.Errorf("Expect 'true', instead of %s", config.Download)
	}
}
