package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
)

var ActionVerify = cli.Command{
	Name:  "verify",
	Usage: "Shows detailed information about the verification process",
	Action: func(c *cli.Context) error {
		app := AppFactory(c)
		app.configure()
		app.verify()

		return nil
	},
}

func (app App) verify() {
	builds := NewReader(app.Path).read(app.Pattern, app.Include, app.Exclude)

	for _, build := range builds {
		hash, _ := Analyzer(build.Root, build.Verify.Include, build.Verify.Exclude)

		fmt.Printf("%s %s\n", hash, build.File)
	}
}
