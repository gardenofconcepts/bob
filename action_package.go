package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
)

var ActionPackage = cli.Command{
	Name:  "package",
	Usage: "Shows detailed information about the packaging process",
	Action: func(c *cli.Context) error {
		app := AppFactory(c)
		app.configure()
		app.doPackage()

		return nil
	},
}

func (app App) doPackage() {
	builds := NewReader(app.Path).read(app.Pattern, app.Include, app.Exclude)

	for _, build := range builds {
		app.printInfo(build)

		if len(build.Package.Include) > 0 {
			fmt.Print("Files    :\n")

			hashes, _ := compress(build.Root, build.Package.Include, build.Package.Exclude)

			for _, path := range hashes {
				fmt.Printf("           (o) %s\n", path)
			}
		} else {
			fmt.Print("Files    : skipped\n")
		}

		fmt.Print("\n-----------------------------------------\n\n")
	}
}
