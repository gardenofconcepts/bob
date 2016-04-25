package main

import (
	"flag"
	"fmt"
)

func main() {
	path := flag.String("path", ".", "Path for searching build files")
	region := flag.String("s3-region", "eu-central-1", "S3 region")
	bucket := flag.String("s3-bucket", "cmsbuild", "S3 bucket name")

	flag.Parse()

	fmt.Println("Searching for build files in path:", *path)

	builds := NewReader(*path).read("*.build.yml")
	storage := Storage(*region, *bucket)

	for _, build := range builds {
		fmt.Println("Found build file", build)

		hash, _ := Analyzer(build.Directory, build.Verify.Include, build.Verify.Exclude)

		build.Hash = hash
		build.Archive = "build/" + hash + ".tar.gz"

		fmt.Println("Analyzing ends up with hash", hash)

		if storage.Has(build) {
			storage.Get(build)
			NewArchive(build.Archive).Extract(build.Directory)
		} else {
			Builder().Build(build.Directory, build.Build)
			NewArchive(build.Archive).Compress(build.Directory, build.Package.Include, build.Package.Exclude)
			storage.Put(build)
		}
	}

	fmt.Println("Ready!")
}
