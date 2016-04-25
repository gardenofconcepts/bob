package main

import (
	log "github.com/Sirupsen/logrus"
)

func (config *AppConfig) configure() {
	level := log.WarnLevel

	if *config.verbose {
		level = log.InfoLevel
	}

	if *config.debug {
		level = log.DebugLevel
	}

	log.SetLevel(level)
}

func (config *AppConfig) run() {
	log.Info("Searching for build files in path:", *config.path)

	builds := NewReader(*config.path).read("*.build.yml")
	storage := Storage(*config.region, *config.bucket)

	for _, build := range builds {
		log.Info("Found build file", build)

		hash, _ := Analyzer(build.Directory, build.Verify.Include, build.Verify.Exclude)

		build.Hash = hash
		build.Archive = "build/" + hash + ".tar.gz"

		log.Info("Analyzing ends up with hash", hash)

		if !*config.force && storage.Has(build) {
			if !*config.skipDownload {
				storage.Get(build)
			}
			NewArchive(build.Archive).Extract(build.Directory)
		} else {
			Builder().Build(build.Directory, build.Build)
			NewArchive(build.Archive).Compress(build.Directory, build.Package.Include, build.Package.Exclude)

			if !*config.skipUpload {
				storage.Put(build)
			}
		}
	}

	log.Info("Ready!")
}
