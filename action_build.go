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

	for _, build := range builds {
		app.processVerification(app.Cache, build, app.StorageBag)
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

	// If there are no verification pattern defined skip the whole process
	// and execute the commands directly
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
			if app.SkipBuild {
				log.WithFields(log.Fields{
					"file": build.File,
				}).Warning("Skip build process")
			} else {
				app.processBuild(build, storage)
			}
		}
	} else {
		log.Info("No verification steps given, skip verification")

		if app.SkipBuild {
			log.WithFields(log.Fields{
				"file": build.File,
			}).Warning("Skip build process")
		} else {
			app.processBuild(build, storage)
		}
	}
}

func (app *App) processBuild(build BuildFile, storage StorageBag) {
	// TODO: if anything goes wrong, there should be some error handling
	Builder().Build(build.Directory, build.Build)

	NewArchive(build.Archive).Compress(build.Root, build.Package.Include, build.Package.Exclude)

	if !app.SkipUpload {
		storage.Put(build)
	}
}
