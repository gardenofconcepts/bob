package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"bob/path"
	"bob/config"
	"bob/parser"
	"bob/reader"
)

var ActionFind = cli.Command{
	Name:  "find",
	Usage: "Find build declaration files in working directory or given path",
	Action: func(c *cli.Context) error {
		app := AppConfig(c)
		app.Configure()

		find(app)

		return nil
	},
}

func find(app config.App) {

	builds := reader.NewReader(app.Path).Read(app.Pattern, app.Include, app.Exclude)

	for _, build := range builds {
		printInfo(build)

		fmt.Print("\n-----------------------------------------\n\n")
	}
}

func printInfo(build parser.BuildFile) {
	fmt.Printf("Name     : %s\n", build.Name)
	fmt.Printf("File     : %s\n", build.File)
	fmt.Printf("Priority : %d\n", build.Priority)
	fmt.Printf("Root     : %s\n", build.Root)
	fmt.Printf("Cwd      : %s\n", build.Cwd)
	fmt.Print("Constant :\n")

	for _, constant := range build.Constant {
		fmt.Printf("           (*) %s: %s (%s)\n", constant.Constant, constant.Result, constant.Command)
	}

	fmt.Print("Constraint:\n")

	for _, constraint := range build.Constraint {
		result := "-"

		if constraint.Result {
			result = "+"
		}

		fmt.Printf("           (%s) %s %s\n", result, constraint.ResultString, constraint.Condition)
	}

	fmt.Print("Verify   :\n")

	for _, path := range build.Verify.Include {
		fmt.Printf("           (+) %s\n", path)
	}

	for _, path := range build.Verify.Exclude {
		fmt.Printf("           (-) %s\n", path)
	}

	fmt.Print("Package  :\n")

	for _, path := range build.Package.Include {
		fmt.Printf("           (+) %s\n", path)
	}

	for _, path := range build.Package.Exclude {
		fmt.Printf("           (-) %s\n", path)
	}

	fmt.Print("Command  :\n")

	for _, cmd := range build.Build {
		command := path.MakePathRelative(build.Root, build.Directory, cmd.Command)
		fmt.Printf("           (o) %s\n", command)
	}
}
