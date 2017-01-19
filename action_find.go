package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
)

var ActionFind = cli.Command{
	Name:  "find",
	Usage: "Find build declaration files in working directory or given path",
	Action: func(c *cli.Context) error {
		app := AppFactory(c)
		app.configure()
		app.find()

		return nil
	},
}

func (app App) find() {

	builds := NewReader(app.Path).read(app.Pattern, app.Include, app.Exclude)

	for _, build := range builds {
		fmt.Printf("Name     : %s\n", build.Name)
		fmt.Printf("File     : %s\n", build.File)
		fmt.Printf("Priority : %d\n", build.Priority)
		fmt.Printf("Root     : %s\n", build.Root)
		fmt.Printf("Cwd      : %s\n", build.Cwd)
		fmt.Print("Verify   :\n")

		for _, path := range buildPaths(app.Path, build.Directory, build.Verify.Include) {
			fmt.Printf("           (+) %s\n", path)
		}

		for _, path := range buildPaths(app.Path, build.Directory, build.Verify.Exclude) {
			fmt.Printf("           (-) %s\n", path)
		}

		fmt.Print("Package  :\n")

		for _, path := range buildPaths(app.Path, build.Directory, build.Package.Include) {
			fmt.Printf("           (+) %s\n", path)
		}

		for _, path := range buildPaths(app.Path, build.Directory, build.Package.Exclude) {
			fmt.Printf("           (-) %s\n", path)
		}

		fmt.Print("Command  :\n")

		for _, cmd := range build.Build {
			fmt.Printf("           (o) %s\n", cmd.Command)
		}

		fmt.Print("\n-----------------------------------------\n\n")
	}
}
