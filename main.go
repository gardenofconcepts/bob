package main

import "fmt"
import (
	"flag"
	"app"
)

func main() {
	path := flag.String("path", ".", "Path for searching build files")

	flag.Parse()

	fmt.Println("Searching for build files in path:", *path)

	builds := app.Read(*path, "*.yml")

	for _, build := range builds {
		fmt.Println("Found build file", build)

		hash, _ := app.Analyzer(build.Directory, build.Verify.Include, build.Verify.Exclude)
		archive := "/home/dennis/htdocs/builder/hash.tar.gz"

		fmt.Println("Analyzing ends up with hash", hash)

		if app.Has(hash) {
			app.Get(hash, archive)
			app.Extract(archive, build.Directory)
		} else {
			//app.Builder(build.Directory, build.Build)
			app.Archive(archive, build.Directory, build.Package.Include, build.Package.Exclude)
			app.Put(hash, archive)
		}
	}
}
