package main

const CONFIG_FILE = "bob.yml"
const CONFIG_PATTERN = "*.build.yml"
const CONFIG_INCLUDE = "**"
const CONFIG_EXCLUDE = ""

type Build struct {
	Command string `yaml:"command"`
}

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
	Build      []Build      `yaml:"build"`
	Constant   []Constant   `yaml:"constant"`
	Constraint []Constraint `yaml:"constraint"` // node/npm version, OS identifier // https://github.com/Knetic/govaluate
}

type App struct {
	WorkingDir string
	Path       string
	Config     string
	Force      bool
	Defaults   AppConfigDefaults
	StorageBag StorageBag
	AppConfig  `yaml:",inline"`
}

type AppConfig struct {
	Cache        string   `yaml:"cache"`
	Pattern      string   `yaml:"pattern"`
	Storage      string   `yaml:"storage"`
	Debug        bool     `yaml:"debug"`
	Verbose      bool     `yaml:"verbose"`
	SkipDownload bool     `yaml:"skipDownload"`
	SkipUpload   bool     `yaml:"skipUpload"`
	SkipBuild    bool     `yaml:"skipBuild"`
	Include      []string `yaml:"include"`
	Exclude      []string `yaml:"exclude"`
	S3           S3Config `yaml:"s3"`
}

type AppConfigDefaults struct {
	Cache        bool
	Pattern      bool
	Storage      bool
	Debug        bool
	Verbose      bool
	SkipDownload bool
	SkipUpload   bool
	SkipBuild    bool
	Include      bool
	Exclude      bool
}

type S3Config struct {
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
}
