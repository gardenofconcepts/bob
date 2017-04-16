package config

import "bob/storage"

const CONFIG_FILE = "bob.yml"

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

type Defaults struct {
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

type App struct {
	WorkingDir string
	Path       string
	Config     string
	Force      bool
	Defaults   Defaults
	StorageBag storage.StorageBag
	AppConfig  `yaml:",inline"`
}
