package config

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func ConfigReader() *AppConfig {
	return &AppConfig{}
}

func (config *AppConfig) Read(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithField("file", path).Fatal("Configuratoin file not found")
	}

	log.WithField("file", path).Info("Read configuration file")

	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.WithError(err).Fatal("Error while reading configuration file")
	}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		log.WithError(err).Fatal("Error while parsing configuration file")
	}
}

func (config *AppConfig) Apply(app *App) error {
	if !app.Defaults.Pattern && len(config.Pattern) > 0 {
		app.Pattern = config.Pattern
	}

	if !app.Defaults.Storage && len(config.Storage) > 0 {
		app.Storage = config.Storage
	}

	if !app.Defaults.Cache && len(config.Cache) > 0 {
		app.Cache = config.Cache
	}

	if !app.Defaults.Include && len(config.Include) > 0 {
		app.Include = config.Include
	}

	if !app.Defaults.Exclude && len(config.Exclude) > 0 {
		app.Exclude = config.Exclude
	}

	if !app.Defaults.SkipDownload {
		app.SkipDownload = config.SkipDownload
	}

	if !app.Defaults.SkipUpload {
		app.SkipUpload = config.SkipUpload
	}

	if !app.Defaults.SkipBuild {
		app.SkipBuild = config.SkipBuild
	}

	if !app.Defaults.Verbose && !app.Verbose {
		app.Verbose = config.Verbose
	}

	if !app.Defaults.Debug && !app.Debug {
		app.Debug = config.Debug
	}

	// s3

	return nil
}
