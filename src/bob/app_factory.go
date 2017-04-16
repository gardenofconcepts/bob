package main

import (
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"path/filepath"
	"strings"
	"bob/path"
	"bob/config"
)

func AppConfig(c *cli.Context) config.App {
	return config.App{
		Path:   getPath(c),
		Force:  c.GlobalBool("force"),
		Config: c.GlobalString("config"),
		Defaults: config.Defaults{
			Pattern:      c.GlobalIsSet("pattern"),
			Cache:        c.GlobalIsSet("cache"),
			Storage:      c.GlobalIsSet("storage"),
			Debug:        c.GlobalIsSet("debug"),
			Verbose:      c.GlobalIsSet("verbose"),
			SkipDownload: c.GlobalIsSet("skip-download"),
			SkipUpload:   c.GlobalIsSet("skip-upload"),
			SkipBuild:    c.GlobalIsSet("skip-build"),
			Include:      c.GlobalIsSet("include"),
			Exclude:      c.GlobalIsSet("exclude"),
		},
		AppConfig: config.AppConfig{
			Include:      path.CleanList(strings.Split(c.GlobalString("include"), ",")),
			Exclude:      path.CleanList(strings.Split(c.GlobalString("exclude"), ",")),
			SkipDownload: c.GlobalBool("skip-download"),
			SkipUpload:   c.GlobalBool("skip-upload"),
			SkipBuild:    c.GlobalBool("skip-build"),
			Pattern:      c.GlobalString("pattern"),
			Debug:        c.GlobalBool("debug"),
			Verbose:      c.GlobalBool("verbose"),
			Cache:        c.GlobalString("cache"),
			Storage:      c.GlobalString("storage"),
			S3: config.S3Config{
				Region: c.GlobalString("region"),
				Bucket: c.GlobalString("bucket"),
			},
		},
	}
}

func getPath(c *cli.Context) string {
	path, _ := os.Getwd()

	if c.Args().Present() {
		path = c.Args().First()
	}

	path, err := filepath.Abs(path)

	if err != nil {
		log.Fatal("Invalid directory", err)
		os.Exit(-1)
	}

	return path
}
