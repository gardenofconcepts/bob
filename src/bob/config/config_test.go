package config

import "testing"

func TestConfig(t *testing.T) {
	config := ConfigReader()
	config.Read("test-fixtures/test.yml")

	if config.Pattern != "*.build.yml" {
		t.Errorf("Expect '*.build.yml', instead of %s", config.Pattern)
	}

	if config.Debug == false {
		t.Errorf("Expect 'true', instead of %s", config.Debug)
	}

	if config.Verbose == false {
		t.Errorf("Expect 'true', instead of %s", config.Verbose)
	}

	if config.SkipUpload == true {
		t.Errorf("Expect 'true', instead of %s", config.SkipUpload)
	}

	if config.SkipDownload == true {
		t.Errorf("Expect 'true', instead of %s", config.SkipDownload)
	}
}
