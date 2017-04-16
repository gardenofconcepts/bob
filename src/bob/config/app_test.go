package config

import (
	"testing"
	"os"
)

func TestApp(t *testing.T) {

	t.Run("Check current working directory", func(t *testing.T) {
		app := App{}
		app.Configure();

		path, _ := os.Getwd()

		if app.WorkingDir != path {
			t.Errorf("Expect '%s', instead of %s", path, app.WorkingDir)
		}
	});

	t.Run("Build app without configuration", func(t *testing.T) {
		app := App{
			Path: ".",
		}
		app.Configure();

		if app.Path != "." {
			t.Errorf("Expect '.', instead of %s", app.Path)
		}
	});

	t.Run("Build app with given configuration", func(t *testing.T) {
		app := App{
			Path: ".",
			Config: "test-fixtures/app/build.yml",
		}
		app.Configure();

		if app.Cache != "test" {
			t.Errorf("Expect 'test', instead of %s", app.Cache)
		}
	});

	t.Run("Build app with configuration file in working directory", func(t *testing.T) {
		app := App{
			WorkingDir: "test-fixtures/app",
		}
		app.Configure();

		if app.Cache != "test" {
			t.Errorf("Expect 'test', instead of %s", app.Cache)
		}
	});

}
