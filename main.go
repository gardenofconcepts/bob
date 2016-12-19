package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	pattern := flag.String("pattern", "*.build.yml", "File pattern for build files")
	include := flag.String("include", "**", "Pattern for directory traversal")
	exclude := flag.String("exclude", "", "Excludes directories with this pattern (e.g. **/node_modules/**,.git)")
	debug := flag.Bool("debug", false, "Enable debug mode (Log level: debug)")
	verbose := flag.Bool("verbose", false, "Enable verbose mode (Log level: info)")
	force := flag.Bool("force", false, "Rebuild data without checking remote")
	skipDownload := flag.Bool("skip-download", false, "Don't download builds")
	skipUpload := flag.Bool("skip-upload", false, "Don't upload builds")
	region := flag.String("s3-region", "eu-central-1", "Specify S3 region")
	bucket := flag.String("s3-bucket", "cmsbuild", "Specify S3 bucket name")
	version := flag.Bool("version", false, "Show version")
	cache := flag.String("cache", "build", "Directory for local (cache) files")
	storage := flag.String("storage", "local", "Specify storage engine(s): local, s3")

	flag.Parse()

	app := App{
		path:         getPath(),
		pattern:      *pattern,
		include:      strings.Split(*include, ","),
		exclude:      strings.Split(*exclude, ","),
		debug:        *debug,
		verbose:      *verbose,
		force:        *force,
		skipDownload: *skipDownload,
		skipUpload:   *skipUpload,
		region:       *region,
		bucket:       *bucket,
		cache:        *cache,
		storage:      *storage,
	}

	if *version {
		fmt.Printf("Builder v%s (Build %s)\nCopyright (c) 2016 Garden of Concepts GmbH\n", APP_VERSION, APP_BUILD)
		os.Exit(0)
	}

	app.path = getPath()

	app.configure()
	app.run()
}

func getPath() string {
	path, _ := os.Getwd()

	if len(flag.Args()) > 0 && len(flag.Arg(0)) > 0 {
		path = flag.Arg(0)
	}

	path, err := filepath.Abs(path)

	if err != nil {
		log.Fatal("Invalid directory", err)
		os.Exit(-1)
	}

	return path
}
