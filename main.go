package main

import (
	"flag"
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
			Usage: "File pattern for build files",
		},
		cli.StringFlag{
			Name: "include",
			Value: CONFIG_INCLUDE,
			Usage: "Pattern for directory traversal",
		},
		cli.StringFlag{
			Name: "exclude",
			Value: CONFIG_EXCLUDE,
			Usage: "Excludes directories with this pattern (e.g. **/node_modules/**,.git)",
		},
		cli.BoolFlag{
			Name: "debug",
			Usage: "Enable debug mode (Log level: debug)",
		},
		cli.BoolFlag{
			Name: "verbose",
			Usage: "Enable verbose mode (Log level: info)",
		},
		cli.BoolFlag{
			Name: "force",
			Usage: "Rebuild data without checking remote",
		},
		cli.BoolFlag{
			Name: "skip-download",
			Usage: "Don't download builds",
		},
		cli.BoolFlag{
			Name: "skip-upload",
			Usage: "Don't upload builds",
		},
		cli.StringFlag{
			Name: "s3-region",
			Usage: "Specify S3 region",
		},
		cli.StringFlag{
			Name: "s3-bucket",
			Usage: "Specify S3 bucket name",
		},
		cli.StringFlag{
			Name: "cache",
			Usage: "Directory for local (cache) files",
		},
		cli.StringFlag{
			Name: "storage",
			Value: "local",
			Usage: "Specify storage engine(s): local, s3",
		},
		cli.StringFlag{
			Name: "config",
			Usage: "Specify a configuration path",
		},
	}
	cliApp.Commands = []cli.Command{
		{
			Name:    "build",
			Usage:   "add a task to the list",
			Action:  func(c *cli.Context) error {
				app := App{
					Path:         getPath(),
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

				app.Path = getPath()

				app.configure()
				app.run()

				return nil
			},
		},
	}
	cliApp.Run(os.Args)
}

func getPath() string {
	path, _ := os.Getwd()

	if len(flag.Args()) > 0 && len(flag.Arg(0)) > 0 {
		path = flag.Arg(0)
	}

	path, err := filepath.Abs(path)

	if err != nil {
		log.Fatal("Invalid directory", err)
		os.Exit(-1)
	}

	return path
}
