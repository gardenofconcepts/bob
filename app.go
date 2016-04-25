package main

import (
	log "github.com/Sirupsen/logrus"
)

func (app *App) configure() {
	level := log.WarnLevel

	if app.verbose {
		level = log.InfoLevel
	}

	if app.debug {
		level = log.DebugLevel
	}

	log.SetLevel(level)
}

func (app *App) run() {
	log.Info("Searching for build files in path:", app.path)

	builds := NewReader(app.path).read(app.pattern)
	storage := S3Storage(app.region, app.bucket)

	for _, build := range builds {
		app.build(build, *storage)
	}

	log.Info("Ready!")
}

func (app *App) build(build BuildFile, storage Storage) {
	log.Info("Found build file", build)

	hash, _ := Analyzer(build.Directory, build.Verify.Include, build.Verify.Exclude)

	build.Hash = hash
	build.Archive = "build/" + hash + ".tar.gz"

	log.Info("Analyzing ends up with hash", hash)

	if !app.force && storage.Has(build) {
		if !app.skipDownload {
			storage.Get(build)
		}
		NewArchive(build.Archive).Extract(build.Directory)
	} else {
		Builder().Build(build.Directory, build.Build)
		NewArchive(build.Archive).Compress(build.Directory, build.Package.Include, build.Package.Exclude)

		if !app.skipUpload {
			storage.Put(build)
		}
	}
}
