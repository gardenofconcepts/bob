package builder

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"path/filepath"
)

type BuildJob struct{}

func Builder() *BuildJob {
	return &BuildJob{}
}

// TODO: error handling
func (e *BuildJob) Build(directory string, builds []Build) {
	for _, build := range builds {
		directory, _ = filepath.Abs(directory)

		err := e.Run(directory, build)

		if err != nil {
			log.WithFields(log.Fields{
				"cmd": build.Command,
				"cwd": directory,
			}).Fatal("Error while executing build")
		}
	}
}

// https://golang.org/src/os/exec/example_test.go
// http://www.darrencoxall.com/golang/executing-commands-in-go/
func (e *BuildJob) Run(directory string, build Build) error {

	log.WithFields(log.Fields{
		"cmd": build.Command,
		"cwd": directory,
	}).Info("Executing build process")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", build.Command)
	cmd.Dir = directory
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		log.WithFields(log.Fields{
			"cmd":    build.Command,
			"cwd":    directory,
			"err":    err,
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Error("Error while running command")

		return err
	}

	log.Debugf("Result: %q\n", stdout.String())

	return nil
}
