package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"path/filepath"
	"bob/storage"
	"bob/hash"
	"bob/config"
	"bob/analyzer"
	"bob/reader"
)

var ActionVerify = cli.Command{
	Name:  "verify",
	Usage: "Shows detailed information about the verification process",
	Action: func(c *cli.Context) error {
		app := AppConfig(c)
		app.Configure()

		verify(app)

		return nil
	},
}

func verify(app config.App) {
	builds := reader.NewReader(app.Path).Read(app.Pattern, app.Include, app.Exclude)

	for _, build := range builds {
		printInfo(build)

		hashes := analyzer.Read(build.Root, build.Verify.Include, build.Verify.Exclude)
		status := "not found"

		for _, constraint := range build.Constraint {
			hashes[constraint.Condition] = constraint.Hash
		}

		hash := hash.List(hashes)

		if app.StorageBag.Has(storage.StorageRequest{
			Hash: hash,
		}) {
			status = "found"
		}

		fmt.Printf("Hash     : %s\n", hash)
		fmt.Printf("Status   : %s\n", status)
		fmt.Print("Verified :\n")

		for path, hash := range hashes {
			displayedPath, err := filepath.Rel(build.Root, path)

			if err != nil {
				displayedPath = "Constraint: " + path
			}

			fmt.Printf("           %s\t%s\n", hash, displayedPath)
		}

		fmt.Print("\n-----------------------------------------\n\n")
	}
}
