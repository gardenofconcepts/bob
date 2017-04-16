package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"path"
	"bob/archive"
	"bob/storage"
	"bob/hash"
	"bob/config"
	"bob/analyzer"
	"bob/builder"
	"bob/parser"
	"bob/reader"
)

var ActionBuild = cli.Command{
	Name:  "build",
	Usage: "***Default*** Build artefacts specified via build declaration files",
	Action: func(c *cli.Context) error {
		app := AppConfig(c)
		app.Configure()

		build(app)

		return nil
	},
}

func build(app config.App) {
	builds := reader.NewReader(app.Path).Read(app.Pattern, app.Include, app.Exclude)

	for _, build := range builds {
		processVerification(app, app.Cache, build, app.StorageBag)
	}

	log.Info("Ready!")
}

func processVerification(app config.App, cacheDir string, build parser.BuildFile, backend storage.StorageBag) {
	log.WithFields(log.Fields{
		"file":      build.File,
		"directory": build.Directory,
		"name":      build.Name,
		"priority":  build.Priority,
	}).Info("Executing build")

	// If there are no verification pattern defined skip the whole process
	// and execute the commands directly
	if len(build.Verify.Include) > 0 {
		hashes := analyzer.Read(build.Root, build.Verify.Include, build.Verify.Exclude)

		for _, constraint := range build.Constraint {
			hashes[constraint.Condition] = constraint.Hash
		}

		hash := hash.List(hashes)

		build.Hash = hash
		build.Archive = path.Join(cacheDir, hash+".tar.gz")

		request := storage.StorageRequest{
			Hash: build.Hash,
			Archive: build.Archive,
		}

		log.WithField("hash", hash).Info("Analyzing ends up with hash")

		if !app.Force && backend.Has(request) {
			if !app.SkipDownload {
				backend.Get(request)
			}
			archive.NewArchive(build.Archive).Extract(build.Directory)
		} else {
			if app.SkipBuild {
				log.WithFields(log.Fields{
					"file": build.File,
				}).Warning("Skip build process")
			} else {
				processBuild(build, backend)

				if !app.SkipUpload {
					backend.Put(request)
				}
			}
		}
	} else {
		log.Info("No verification steps given, skip verification")

		if app.SkipBuild {
			log.WithFields(log.Fields{
				"file": build.File,
			}).Warning("Skip build process")
		} else {
			processBuild(build, backend)
		}
	}
}

func processBuild(build parser.BuildFile, backend storage.StorageBag) {
	// TODO: if anything goes wrong, there should be some error handling
	builder.Builder().Build(build.Directory, build.Build)

	archive.NewArchive(build.Archive).Compress(build.Root, build.Package.Include, build.Package.Exclude)
}
