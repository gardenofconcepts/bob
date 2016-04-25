package main

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

func Parser(path string) *BuildFile {
	build 		:= new(BuildFile)
	build.File	= path
	build.Directory	= filepath.Dir(path)
	build.Priority	= 0
	build.Name	= "Unknown"

	parse(build)

	return build
}

func parse(build *BuildFile) {
	data, err := ioutil.ReadFile(build.File)

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(data, build)
}
