package main

type Build struct {
	Command string `yaml:"command"`
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
	File      string
	Directory string
	Hash      string
	Archive   string
	Name      string  `yaml:"name"`
	Priority  int     `yaml:"priority"`
	Verify    Verify  `yaml:"verify"`
	Package   Package `yaml:"package"`
	Build     []Build `yaml:"build"`
}

type App struct {
	Path      string
	Config    string
	Cache     string `yaml:"cache"`
	Pattern   string `yaml:"pattern"`
	Storage   string `yaml:"storage"`
	Force     bool
	Debug     bool     `yaml:"debug"`
	Verbose   bool     `yaml:"verbose"`
	S3        S3Config `yaml:"s3"`
	AppConfig `yaml:",inline"`
}

type AppConfig struct {
	SkipDownload bool     `yaml:"skipDownload"`
	SkipUpload   bool     `yaml:"skipUpload"`
	Include      []string `yaml:"include"`
	Exclude      []string `yaml:"exclude"`
}

type S3Config struct {
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}
