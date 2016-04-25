package main

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

type BuildJob struct{}

func Builder() *BuildJob {
	return &BuildJob{}
}

func (e *BuildJob) Build(directory string, builds []Build) {

	log.Info("Starting build process...")

	for _, build := range builds {
		e.Run(directory, build)
	}
}

// https://golang.org/src/os/exec/example_test.go
// http://www.darrencoxall.com/golang/executing-commands-in-go/
func (e *BuildJob) Run(directory string, build Build) error {

	log.Info("Run command on path", build.Command, directory)

	var out bytes.Buffer

	cmd := exec.Command(build.Command)
	cmd.Dir = directory
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return err
	}

	log.Debugf("Result: %q\n", out.String())

	return nil
}
