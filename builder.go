package main

import (
	"fmt"
	"bytes"
	"os/exec"
)

func Builder(directory string, builds []Build) {

	fmt.Println("Starting build process...")

	for _, build := range builds {
		run(directory, build)
	}
}

// https://golang.org/src/os/exec/example_test.go
// http://www.darrencoxall.com/golang/executing-commands-in-go/
func run(directory string, build Build) error {

	fmt.Println("Run command on path", build.Command, directory)

	var out bytes.Buffer

	cmd 		:= exec.Command(build.Command)
	cmd.Dir 	= directory
	cmd.Stdout	= &out

	err := cmd.Run()

	if err != nil {
		return err
	}

	fmt.Printf("Result: %q\n", out.String())

	return nil
}
