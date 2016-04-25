package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

func Parser(path string) *BuildFile {
	return &BuildFile{
		File:      path,
		Directory: filepath.Dir(path),
		Priority:  0,
		Name:      "Unknown",
	}
}

func (build *BuildFile) parse() *BuildFile {
	data, err := ioutil.ReadFile(build.File)

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(data, build)

	return build
}
