package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"path"
)

var ActionBuild = cli.Command{
	Name:  "build",
	Usage: "***Default*** Build artefacts specified via build declaration files",
	Action: func(c *cli.Context) error {
		app := AppFactory(c)
		app.configure()
		app.build()

		return nil
	},
}

func (app *App) build() {
	builds := NewReader(app.Path).read(app.Pattern, app.Include, app.Exclude)

	storage := Storage()
	storage.Register(StorageLocal(app.Cache))

	if app.Storage == "s3" {
		storage.Register(StorageS3(app.S3.Region, app.S3.Bucket))
	}

	for _, build := range builds {
		app.processVerification(app.Cache, build, *storage)
	}

	log.Info("Ready!")
}

func (app *App) processVerification(cacheDir string, build BuildFile, storage StorageBag) {
	log.WithFields(log.Fields{
		"file":      build.File,
		"directory": build.Directory,
		"name":      build.Name,
		"priority":  build.Priority,
	}).Info("Executing build")

	if len(build.Verify.Include) > 0 {
		hash, _ := Analyzer(build.Root, build.Verify.Include, build.Verify.Exclude)

		build.Hash = hash
		build.Archive = path.Join(cacheDir, hash+".tar.gz")

		log.WithField("hash", hash).Info("Analyzing ends up with hash")

		if !app.Force && storage.Has(build) {
			if !app.SkipDownload {
				storage.Get(build)
			}
			NewArchive(build.Archive).Extract(build.Directory)
		} else {
			app.processBuild(build, storage)
		}
	} else {
		log.Info("No verification steps given, skip verification")

		app.processBuild(build, storage)
	}
}

func (app *App) processBuild(build BuildFile, storage StorageBag) {
	// TODO: if anything goes wrong, there should be some error handling
	Builder().Build(build.Directory, build.Build)

	NewArchive(build.Archive).Compress(build.Directory, build.Package.Include, build.Package.Exclude)

	if !app.SkipUpload {
		storage.Put(build)
	}
}
