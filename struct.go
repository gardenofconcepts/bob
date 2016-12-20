package main

const CONFIG_FILE = "bob.yml"
const CONFIG_PATTERN = "*.build.yml"
const CONFIG_INCLUDE = "**"
const CONFIG_EXCLUDE = ""

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

type App struct {
	Path     string
	Config   string
	Cache    string `json:"cache"`
	Pattern  string `json:"pattern"`
	Storage  string `json:"storage"`
	Force    bool
	Debug    bool     `json:"debug"`
	Verbose  bool     `json:"verbose"`
	Download bool     `json:"download"`
	Upload   bool     `json:"upload"`
	Include  []string `json:"include"`
	Exclude  []string `json:"exclude"`
	S3       S3Config `json:"s3"`
}

type AppDefaults struct {

}

type S3Config struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
}
