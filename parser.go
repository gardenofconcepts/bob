package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

func Parser() *BuildFile {
	return &BuildFile{
		Priority: 0,
		Name:     "Unknown",
		Verify: Verify{
			Include: []string{},
			Exclude: []string{},
		},
		Package: Package{
			Include: []string{},
			Exclude: []string{},
		},
		Build: []Build{},
	}
}

func (build *BuildFile) load(path string) *BuildFile {
	build.File = path
	build.Directory = filepath.Dir(path)

	data, err := ioutil.ReadFile(build.File)

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(data, build)

	return build.determine()
}

func (build *BuildFile) determine() *BuildFile {
	if len(build.Cwd) == 0 {
		build.Cwd = build.Directory
	} else if !filepath.IsAbs(build.Cwd) {
		build.Cwd = filepath.Join(build.Directory, build.Cwd)
		build.Cwd, _ = filepath.Abs(build.Cwd)
	}

	if len(build.Root) == 0 {
		build.Root = build.Directory
	} else if !filepath.IsAbs(build.Root) {
		build.Root = filepath.Join(build.Directory, build.Root)
		build.Root, _ = filepath.Abs(build.Root)
	}

	build.Verify.Include = buildPaths(build.Root, build.Cwd, build.Verify.Include)
	build.Verify.Exclude = buildPaths(build.Root, build.Cwd, build.Verify.Exclude)

	build.Package.Include = buildPaths(build.Root, build.Cwd, build.Package.Include)
	build.Package.Exclude = buildPaths(build.Root, build.Cwd, build.Package.Exclude)

	return build
}
