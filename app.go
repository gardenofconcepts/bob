package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
)

func (app *App) configure() {
	app.configureLog()

	if app.WorkingDir == "" {
		path, _ := os.Getwd()
		app.WorkingDir = path
	}

	if len(app.Config) > 0 {
		config := ConfigReader()
		config.Read(app.Config)
		config.Apply(app)
	} else if _, err := os.Stat(filepath.Join(app.WorkingDir, CONFIG_FILE)); err == nil {
		config := ConfigReader()
		config.Read(filepath.Join(app.WorkingDir, CONFIG_FILE))
		config.Apply(app)
	}

	app.configureLog()
	app.configureBackend()
}

func (app *App) configureLog() {
	level := log.WarnLevel

	if app.Verbose {
		level = log.InfoLevel
	}

	if app.Debug {
		level = log.DebugLevel
	}

	log.SetLevel(level)
}

func (app *App) configureBackend() {
	storage := Storage()
	storage.Register(StorageLocal(app.Cache))

	if app.Storage == "s3" {
		storage.Register(StorageS3(app.S3.Region, app.S3.Bucket))
	}

	app.StorageBag = storage
}
