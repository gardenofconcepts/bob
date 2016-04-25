package main

type Build struct {
	Command string `json:"command"`
}

type Verify struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type Package struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type BuildFile struct {
	File      string
	Directory string
	Hash      string
	Archive   string
	Name      string  `json:"name"`
	Priority  int     `json:"priority"`
	Verify    Verify  `json:"verify"`
	Package   Package `json:"package"`
	Build     []Build `json:"build"`
}

type AppConfig struct {
	path         *string
	debug        *bool
	verbose      *bool
	force        *bool
	skipDownload *bool
	skipUpload   *bool
	region       *string
	bucket       *string
}
