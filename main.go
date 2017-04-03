package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	path, _ := os.Getwd()

	cliApp := cli.NewApp()
	cliApp.Name = "bob"
	cliApp.Usage = "Der Baumeister"
	cliApp.ArgsUsage = fmt.Sprintf("Path to directory where build declaration files are located OR specify a list of build declaration files (Default: %s)", path)
	cliApp.Version = APP_VERSION
	cliApp.Action = ActionBuild.Action
	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pattern, p",
			Value: CONFIG_PATTERN,
			Usage: "file pattern for build files",
		},
		cli.StringFlag{
			Name:  "include, i",
			Value: CONFIG_INCLUDE,
			Usage: "pattern for directory traversal",
		},
		cli.StringFlag{
			Name:  "exclude, e",
			Value: CONFIG_EXCLUDE,
			Usage: "excludes directories with this pattern (e.g. **/node_modules/**,.git)",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug mode (Log level: debug)",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable verbose mode (Log level: info)",
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "rebuild data without checking remote",
		},
		cli.BoolFlag{
			Name:  "skip-download",
			Usage: "don't download builds",
		},
		cli.BoolFlag{
			Name:  "skip-upload",
			Usage: "don't upload builds",
		},
		cli.BoolFlag{
			Name:  "skip-build",
			Usage: "don't execute build process",
		},
		cli.StringFlag{
			Name:  "s3-region",
			Usage: "specify S3 region",
		},
		cli.StringFlag{
			Name:  "s3-bucket",
			Usage: "specify S3 bucket name",
		},
		cli.StringFlag{
			Name:  "cache",
			Value: os.TempDir(),
			Usage: "directory for local (cache) files",
		},
		cli.StringFlag{
			Name:  "storage",
			Value: "local",
			Usage: "specify storage engine(s): local, s3",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "path to configuration file",
		},
	}
	cliApp.Commands = []cli.Command{
		ActionBuild,
		ActionFind,
		ActionVerify,
		ActionPackage,
		ActionClean,
	}
	cliApp.Run(os.Args)
}
