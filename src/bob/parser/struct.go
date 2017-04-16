package parser

import "bob/builder"

type Constant struct {
	Command  string `yaml:"command"`
	Constant string `yaml:"constant"` // CONSTANT_OS, CONSTANT_VERSION (bob version)
	Result   string
}

type Constraint struct {
	Condition    string `yaml:"condition"`
	Name         string `yaml:"name"`
	Hash         string
	Result       bool
	ResultString string
}

type Verify struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

type Package struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

type BuildFile struct {
	File       string
	Directory  string
	Hash       string
	Archive    string
	Name       string       `yaml:"name"`
	Cwd        string       `yaml:"cwd"`  // change working directory
	Root       string       `yaml:"root"` // lowest level for finding / verification
	Priority   int          `yaml:"priority"`
	Verify     Verify       `yaml:"verify"`
	Package    Package      `yaml:"package"`
	Build      []builder.Build `yaml:"build"`
	Constant   []Constant   `yaml:"constant"`
	Constraint []Constraint `yaml:"constraint"` // node/npm version, OS identifier // https://github.com/Knetic/govaluate
}
