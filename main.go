package main

import (
	"flag"
	"os"
	"path/filepath"
	"log"
	"fmt"
)

func main() {
	pattern := flag.String("pattern", "*.build.yml", "File pattern for build files")
	debug := flag.Bool("debug", false, "Enable debug mode (Log level: debug)")
	verbose := flag.Bool("verbose", false, "Enable verbose mode (Log level: info)")
	force := flag.Bool("force", false, "Rebuild data without checking remote")
	skipDownload := flag.Bool("skip-download", false, "Don't download builds")
	skipUpload := flag.Bool("skip-upload", false, "Don't upload builds")
	region := flag.String("s3-region", "eu-central-1", "Specify S3 region")
	bucket := flag.String("s3-bucket", "cmsbuild", "Specify S3 bucket name")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	app := App{
		path:         getPath(),
		pattern:      *pattern,
		debug:        *debug,
		verbose:      *verbose,
		force:        *force,
		skipDownload: *skipDownload,
		skipUpload:   *skipUpload,
		region:       *region,
		bucket:       *bucket,
	}

	if (*version) {
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

	if (err != nil) {
		log.Fatal("Invalid directory", err)
		os.Exit(-1)
	}

	return path
}