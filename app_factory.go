package main

import (
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func AppFactory(c *cli.Context) App {
	return App{
		Path:   getPath(c),
		Force:  c.GlobalBool("force"),
		Config: c.GlobalString("config"),
		Defaults: AppConfigDefaults{
			Pattern:      c.GlobalIsSet("pattern"),
			Cache:        c.GlobalIsSet("cache"),
			Storage:      c.GlobalIsSet("storage"),
			Debug:        c.GlobalIsSet("debug"),
			Verbose:      c.GlobalIsSet("verbose"),
			SkipDownload: c.GlobalIsSet("skip-download"),
			SkipUpload:   c.GlobalIsSet("skip-upload"),
			Include:      c.GlobalIsSet("include"),
			Exclude:      c.GlobalIsSet("exclude"),
		},
		AppConfig: AppConfig{
			Include:      cleanList(strings.Split(c.GlobalString("include"), ",")),
			Exclude:      cleanList(strings.Split(c.GlobalString("exclude"), ",")),
			SkipDownload: c.GlobalBool("skip-download"),
			SkipUpload:   c.GlobalBool("skip-upload"),
			Pattern:      c.GlobalString("pattern"),
			Debug:        c.GlobalBool("debug"),
			Verbose:      c.GlobalBool("verbose"),
			Cache:        c.GlobalString("cache"),
			Storage:      c.GlobalString("storage"),
			S3: S3Config{
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
