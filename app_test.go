package main

import (
	"testing"
	"os"
)

func TestApp(t *testing.T) {

	t.Run("Check current working directory", func(t *testing.T) {
		app := App{}
		app.configure();

		path, _ := os.Getwd()

		if app.WorkingDir != path {
			t.Errorf("Expect '%s', instead of %s", path, app.WorkingDir)
		}
	});

	t.Run("Build app without configuration", func(t *testing.T) {
		app := App{
			Path: ".",
		}
		app.configure();

		if app.Path != "." {
			t.Errorf("Expect '.', instead of %s", app.Path)
		}
	});

	t.Run("Build app with given configuration", func(t *testing.T) {
		app := App{
			Path: ".",
			Config: "assets/app/build.yml",
		}
		app.configure();

		if app.Cache != "test" {
			t.Errorf("Expect 'test', instead of %s", app.Cache)
		}
	});

	t.Run("Build app with configuration file in working directory", func(t *testing.T) {
		app := App{
			WorkingDir: "assets/app",
		}
		app.configure();

		if app.Cache != "test" {
			t.Errorf("Expect 'test', instead of %s", app.Cache)
		}
	});

}
