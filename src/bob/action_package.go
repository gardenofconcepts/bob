package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"bob/archive"
	"bob/config"
	"bob/reader"
)

var ActionPackage = cli.Command{
	Name:  "package",
	Usage: "Shows detailed information about the packaging process",
	Action: func(c *cli.Context) error {
		app := AppConfig(c)
		app.Configure()

		doPackage(app)

		return nil
	},
}

func doPackage(app config.App) {
	builds := reader.NewReader(app.Path).Read(app.Pattern, app.Include, app.Exclude)

	for _, build := range builds {
		printInfo(build)

		if len(build.Package.Include) > 0 {
			fmt.Print("Files    :\n")

			hashes, _ := archive.Compress(build.Root, build.Package.Include, build.Package.Exclude)

			for _, path := range hashes {
				fmt.Printf("           (o) %s\n", path)
			}
		} else {
			fmt.Print("Files    : skipped\n")
		}

		fmt.Print("\n-----------------------------------------\n\n")
	}
}
