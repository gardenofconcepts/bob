package main

import (
	"flag"
)

func main() {
	app := AppConfig{
		path:         flag.String("path", ".", "Path for searching build files"),
		debug:        flag.Bool("debug", false, "Enable debug mode (Log level: debug)"),
		verbose:      flag.Bool("verbose", false, "Enable verbose mode (Log level: info)"),
		force:        flag.Bool("force", false, "Rebuild data without remote check"),
		skipDownload: flag.Bool("skip-download", false, "Don't download builds"),
		skipUpload:   flag.Bool("skip-upload", false, "Don't upload builds"),
		region:       flag.String("s3-region", "eu-central-1", "Specify S3 region"),
		bucket:       flag.String("s3-bucket", "cmsbuild", "Specify S3 bucket name"),
	}

	flag.Parse()

	app.configure()
	app.run()
}
