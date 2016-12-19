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
	builds := NewReader(app.path).read(app.pattern, app.include, app.exclude)
	storage := S3Storage(app.region, app.bucket)

	for _, build := range builds {
		app.build(build, *storage)
	}

	log.Info("Ready!")
}

func (app *App) build(build BuildFile, storage Storage) {
	log.WithFields(log.Fields{
		"file":      build.File,
		"directory": build.Directory,
		"name":      build.Name,
		"priority":  build.Priority,
	}).Info("Executing build")

	if len(build.Verify.Include) > 0 {
		hash, _ := Analyzer(build.Directory, build.Verify.Include, build.Verify.Exclude)

		build.Hash = hash
		build.Archive = "build/" + hash + ".tar.gz"

		log.WithField("hash", hash).Info("Analyzing ends up with hash")

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
	} else {
		log.Info("No verification steps given, skip verification")

		Builder().Build(build.Directory, build.Build)
	}
}
