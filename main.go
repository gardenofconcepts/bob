package main

import (
	"flag"
	"fmt"
)

func main() {
	path := flag.String("path", ".", "Path for searching build files")

	flag.Parse()

	fmt.Println("Searching for build files in path:", *path)

	builds := Read(*path, "*.build.yml")

	for _, build := range builds {
		fmt.Println("Found build file", build)

		hash, _	:= Analyzer(build.Directory, build.Verify.Include, build.Verify.Exclude)
		build.Hash	= hash
		build.Archive	= "hash_" + hash + ".tar.gz"

		fmt.Println("Analyzing ends up with hash", hash)

		if Has(build) {
			Get(build)
			Extract(build.Archive, build.Directory)
		} else {
			Builder(build.Directory, build.Build)
			Compress(build.Archive, build.Directory, build.Package.Include, build.Package.Exclude)
			Put(build)
		}
	}

	fmt.Println("Ready!")
}
