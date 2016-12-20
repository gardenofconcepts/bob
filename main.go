package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "bob"
	cliApp.Usage = "Der Baumeister"
	cliApp.Version = APP_VERSION
	cliApp.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "pattern, p",
			Value: CONFIG_PATTERN,
			Usage: "file pattern for build files",
		},
		cli.StringFlag{
			Name: "include, i",
			Value: CONFIG_INCLUDE,
			Usage: "pattern for directory traversal",
		},
		cli.StringFlag{
			Name: "exclude, e",
			Value: CONFIG_EXCLUDE,
			Usage: "excludes directories with this pattern (e.g. **/node_modules/**,.git)",
		},
		cli.BoolFlag{
			Name: "debug",
			Usage: "enable debug mode (Log level: debug)",
		},
		cli.BoolFlag{
			Name: "verbose",
			Usage: "enable verbose mode (Log level: info)",
		},
		cli.BoolFlag{
			Name: "force, f",
			Usage: "rebuild data without checking remote",
		},
		cli.BoolFlag{
			Name: "skip-download",
			Usage: "don't download builds",
		},
		cli.BoolFlag{
			Name: "skip-upload",
			Usage: "don't upload builds",
		},
		cli.StringFlag{
			Name: "s3-region",
			Usage: "specify S3 region",
		},
		cli.StringFlag{
			Name: "s3-bucket",
			Usage: "specify S3 bucket name",
		},
		cli.StringFlag{
			Name: "cache",
			Value: os.TempDir(),
			Usage: "directory for local (cache) files",
		},
		cli.StringFlag{
			Name: "storage",
			Value: "local",
			Usage: "specify storage engine(s): local, s3",
		},
		cli.StringFlag{
			Name: "config",
			Usage: "path to configuration file",
		},
	}
	cliApp.Commands = []cli.Command{
		{
			Name:    "build",
			Usage:   "find build files to start build process",
			Action:  func(c *cli.Context) error {
				app := App{
					Path:         getPath(c),
					Force:        c.GlobalBool("force"),
					Config:       c.GlobalString("config"),
					Include:      cleanList(strings.Split(c.GlobalString("include"), ",")),
					Exclude:      cleanList(strings.Split(c.GlobalString("exclude"), ",")),
					Pattern:      c.GlobalString("pattern"),
					Debug:        c.GlobalBool("debug"),
					Verbose:      c.GlobalBool("verbose"),
					Download:     c.GlobalBool("skip-download"),
					Upload:       c.GlobalBool("skip-upload"),
					Cache:        c.GlobalString("cache"),
					Storage:      c.GlobalString("storage"),
					S3: S3Config{
						Region: c.GlobalString("region"),
						Bucket: c.GlobalString("bucket"),
					},
				}

				app.configure()
				app.run()

				return nil
			},
		},
	}
	cliApp.Run(os.Args)
}

func getPath(c *cli.Context) string {
	path, _ := os.Getwd()

	if c.Args().Present() {
		path = c.Args().First()
	}

	path, err := filepath.Abs(path)

	if err != nil {
		log.Fatal("Invalid directory", err)
		os.Exit(-1)
	}

	return path
}
