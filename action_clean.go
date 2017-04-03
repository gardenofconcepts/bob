package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
)

var ActionClean = cli.Command{
	Name:  "clean",
	Usage: "Removes files created by this tool",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force",
			Usage: "Removes files",
		},
	},
	Action: func(c *cli.Context) error {
		app := AppFactory(c)
		app.configure()
		app.clean(c.Bool("force"))

		return nil
	},
}

func (app App) clean(force bool) {

	builds := NewReader(app.Path).read(app.Pattern, app.Include, app.Exclude)

	for _, build := range builds {
		app.printInfo(build)

		fmt.Print("Files    :\n")

		fileList, _ := Clean(build.Root, build.Package.Include, build.Package.Exclude, force)

		for _, path := range fileList {
			fmt.Printf("           (o) %s\n", path)
		}

		fmt.Print("\n-----------------------------------------\n\n")
	}
}

func Clean(rootDir string, includes []string, excludes []string, force bool) ([]string, error) {
	fileList := []string{}

	err := filepath.Walk(rootDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			log.Error("Error reading directory", err)

			return err
		}

		if !matchList(includes, filePath, rootDir) || matchList(excludes, filePath, rootDir) {
			log.WithField("file", filePath).Debug("Skipping file")

			return nil
		}

		log.WithFields(log.Fields{
			"file": filePath,
		}).Debug("Remove file")

		if force {
			os.Remove(filePath)
		}

		filePath, _ = filepath.Rel(rootDir, filePath)
		fileList = append(fileList, filePath)

		return nil
	})

	return fileList, err
}
