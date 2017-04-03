package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
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
		app.printInfo(build)

		hashes := read(build.Root, build.Verify.Include, build.Verify.Exclude)
		status := "not found"

		if app.StorageBag.Has(BuildFile{
			Hash: hashList(hashes),
		}) {
			status = "found"
		}

		fmt.Printf("Hash     : %s\n", hashList(hashes))
		fmt.Printf("Status   : %s\n", status)
		fmt.Print("Verified :\n")

		for path, hash := range hashes {
			path, _ = filepath.Rel(build.Root, path)
			fmt.Printf("           %s\t%s\n", hash, path)
		}

		fmt.Print("\n-----------------------------------------\n\n")
	}
}
