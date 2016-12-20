package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
)

func (app *App) configure() {
	app.configureLog()

	if len(app.Config) > 0 {
		config := ConfigReader()
		config.Read(app.Config)
		config.Apply(app)
	} else if _, err := os.Stat(CONFIG_FILE); os.IsExist(err) {
		config := ConfigReader()
		config.Read(CONFIG_FILE)
		config.Apply(app)
	}

	fmt.Printf("Apply configuration: %+v\n", app)

	app.configureLog()
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

func (app *App) run() {
	builds := NewReader(app.Path).read(app.Pattern, app.Include, app.Exclude)

	storage := Storage()
	storage.Register(StorageLocal(app.Cache))

	if app.Storage == "s3" {
		storage.Register(StorageS3(app.S3.Region, app.S3.Bucket))
	}

	for _, build := range builds {
		app.build(app.Cache, build, *storage)
	}

	log.Info("Ready!")
}

func (app *App) build(cacheDir string, build BuildFile, storage StorageBag) {
	log.WithFields(log.Fields{
		"file":      build.File,
		"directory": build.Directory,
		"name":      build.Name,
		"priority":  build.Priority,
	}).Info("Executing build")

	if len(build.Verify.Include) > 0 {
		hash, _ := Analyzer(build.Directory, build.Verify.Include, build.Verify.Exclude)

		build.Hash = hash
		build.Archive = path.Join(cacheDir, hash+".tar.gz")

		log.WithField("hash", hash).Info("Analyzing ends up with hash")

		if !app.Force && storage.Has(build) {
			if !app.Download {
				storage.Get(build)
			}
			NewArchive(build.Archive).Extract(build.Directory)
		} else {
			Builder().Build(build.Directory, build.Build)
			NewArchive(build.Archive).Compress(build.Directory, build.Package.Include, build.Package.Exclude)

			if !app.Upload {
				storage.Put(build)
			}
		}
	} else {
		log.Info("No verification steps given, skip verification")

		Builder().Build(build.Directory, build.Build)
	}
}
